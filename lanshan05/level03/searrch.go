package main

import (
	"bufio"
	"fmt"
	"lanshan05/taskpool" // 引入Lv2实现的协程池包
	"os"
	"path/filepath"
)

// FileSearchTask 实现taskpool.Task接口，代表单个文件的搜索任务
type FileSearchTask struct {
	Filepath string              // 待搜索的文件路径
	Keyword  string              // 目标关键词
	ResultCh chan<- SearchResult // 结果通道，用于传递搜索结果
}

// SearchResult 存储单个文件的搜索结果
type SearchResult struct {
	Filepath string   // 文件名
	Lines    []string // 包含关键词的行
}

// Execute 实现Task接口，执行文件搜索逻辑
func (t *FileSearchTask) Execute() {
	file, err := os.Open(t.Filepath)
	if err != nil {
		return // 忽略打开失败的文件
	}
	defer file.Close()

	var resultLines []string
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if contains(line, t.Keyword) {
			// 格式化输出行号和内容
			resultLines = append(resultLines, fmt.Sprintf("第%d行: %s", lineNum, line))
		}
		lineNum++
	}

	// 将结果发送到通道（若有匹配行）
	if len(resultLines) > 0 {
		t.ResultCh <- SearchResult{
			Filepath: t.Filepath,
			Lines:    resultLines,
		}
	}
}

// contains 判断字符串是否包含关键词（简单实现，可扩展为正则）
func contains(s, substr string) bool {
	return len(substr) == 0 || filepath.Base(s) == substr || (len(s) >= len(substr) && s[:len(substr)] == substr)
}

// 递归遍历目录，收集所有文件路径
func walkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	// 解析命令行参数
	if len(os.Args) != 3 {
		fmt.Printf("用法: %s [检索目录] [关键词]\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]
	keyword := os.Args[2]

	// 1. 递归获取目录下所有文件
	files, err := walkDir(dir)
	if err != nil {
		fmt.Printf("遍历目录失败: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Println("目录下无文件")
		return
	}

	// 2. 初始化协程池（这里设置5个工作协程，任务缓冲为100）
	pool := taskpool.New(5, 100)
	defer pool.Close()

	// 3. 初始化结果通道，用于收集所有任务的结果
	resultCh := make(chan SearchResult, len(files))
	defer close(resultCh)

	// 4. 提交所有文件的搜索任务到协程池
	for _, file := range files {
		pool.Submit(&FileSearchTask{
			Filepath: file,
			Keyword:  keyword,
			ResultCh: resultCh,
		})
	}

	// 5. 等待所有任务完成，并统一输出结果
	go func() {
		pool.Close()
		close(resultCh) // 关闭协程池，等待所有任务执行完毕
	}()

	// 读取结果通道，输出所有匹配结果
	for result := range resultCh {
		fmt.Printf("\n=== 匹配文件: %s ===\n", result.Filepath)
		for _, line := range result.Lines {
			fmt.Println(line)
		}
	}
}

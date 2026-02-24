package utils_viper

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 全局配置结构体（加互斥锁保证线程安全）
var (
	config   *AppConfig   // 私有变量，避免直接修改
	configMu sync.RWMutex // 读写锁，保障热更新线程安全
)

// AppConfig 配置总结构体
type AppConfig struct {
	App   AppConfigItem `mapstructure:"app"`
	Mysql MysqlConfig   `mapstructure:"mysql"`
	Redis RedisConfig   `mapstructure:"redis"`
	JWT   JWTConfig     `mapstructure:"jwt"`
	Log   LogConfig     `mapstructure:"log"`
}

// AppConfigItem 项目基础配置
type AppConfigItem struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

// MysqlConfig MySQL配置
type MysqlConfig struct {
	Dsn             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`      // Redis地址
	Password string `mapstructure:"password"`  // Redis密码
	DB       int    `mapstructure:"db"`        // 数据库编号
	PoolSize int    `mapstructure:"pool_size"` // 连接池大小
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackup  int    `mapstructure:"max_backup"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// InitConfig 初始化配置（带热更新）
func InitConfig() error {
	// 1. 配置文件基础设置
	viper.SetConfigName("config")   // 配置文件名（无后缀）
	viper.SetConfigType("yaml")     // 配置格式
	viper.AddConfigPath("./config") // 优先读取config目录
	viper.AddConfigPath(".")        // 备用：项目根目录

	// 2. 读取初始配置
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 3. 解析到全局结构体（加写锁）
	configMu.Lock()
	defer configMu.Unlock()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 4. 配置热更新
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("[配置热更新] 检测到文件变化: %s\n", e.Name)
		configMu.Lock()
		defer configMu.Unlock()
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("[配置热更新] 解析失败: %v\n", err)
			return
		}
		fmt.Printf("[配置热更新] 生效 - 日志级别: %s | JWT过期时间: %d小时 | Redis地址: %s\n",
			config.Log.Level, config.JWT.ExpireHours, config.Redis.Addr)
	})

	fmt.Println("配置初始化成功（已开启热更新）")
	return nil
}

// GetConfig 获取最新配置（对外暴露）
func GetConfig() *AppConfig {
	configMu.RLock()
	defer configMu.RUnlock()
	return config
}

// GetEnv 获取当前环境
func GetEnv() string {
	cfg := GetConfig()
	if cfg == nil {
		return "dev"
	}
	return cfg.App.Env
}

// IsDev 判断是否为开发环境
func IsDev() bool {
	return GetEnv() == "dev"
}

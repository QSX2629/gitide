package main

import "testing"

func TestJisuanji(t *testing.T) {
	tests := []struct {
		name    string
		x       int
		y       int
		op      string
		want    int
		wanterr bool
	}{
		{"加法", 3, 5, "+", 8, false},
		{"减法", 3, 5, "-", -2, false},
		{"乘法", 3, 5, "*", 15, false},
		{"除法", 3, 1, "/", 3, false},
		{"除法（除数为0）", 5, 0, "/", 0, true},
		{"非法运算符", 5, 4, "%", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sum := jisuanji(tt.x, tt.y, tt.op)
			if !tt.wanterr {
				if sum != tt.want {
					t.Errorf("输入 %d %s %d，期望结果 %d，实际结果 %d", tt.x, tt.op, tt.y, tt.want, sum)

				}

			}
			if sum != 0 {
				t.Errorf("输入 %d %s %d，期望触发错误（结果应为0），实际结果 %d", tt.x, tt.op, tt.y, sum)
			}
		})
	}
}

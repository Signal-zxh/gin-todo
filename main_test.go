package main

import "testing"

func TestIsValidTask(t *testing.T) {
	tests := []struct {
		name     string
		task     string
		expected bool
	}{
		{"有效任务", "学习 Go", true},
		{"空字符串", "", false},
		{"只有空格", "   ", false},
		{"带空格的有效任务", "  学习测试  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTask(tt.task)
			if result != tt.expected {
				t.Errorf("isValidTask(%q) = %v, 期望 %v", tt.task, result, tt.expected)
			}
		})
	}
}

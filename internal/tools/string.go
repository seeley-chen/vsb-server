package tools

import (
	"strings"
)

// TrimSpace 去除字符串首尾空白
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// IsEmpty 判断字符串是否为空或仅包含空白字符
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Default 返回第一个非空字符串；若 s 为空则返回 fallback
func Default(s, fallback string) string {
	if !IsEmpty(s) {
		return s
	}
	return fallback
}

// ContainsString 判断字符串切片是否包含指定字符串
func ContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// MaskString 对字符串进行脱敏处理，保留前 start 位和后 end 位，中间用 * 填充
// 若字符串长度 <= start+end，返回全 * 的等长字符串
func MaskString(s string, start, end int) string {
	runes := []rune(s)
	n := len(runes)
	if n == 0 {
		return s
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start+end >= n {
		return strings.Repeat("*", n)
	}
	maskLen := n - start - end
	return string(runes[:start]) + strings.Repeat("*", maskLen) + string(runes[n-end:])
}

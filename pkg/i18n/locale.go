package i18n

import "strings"

// Locale 多语言文本，JSON 格式：{"zh-cn":"中文","en-us":"English"}
type Locale map[string]string

// 支持的语言代码
const (
	ZhCN = "zh-cn"
	EnUS = "en-us"
)

// IsEmpty 判断是否为空（nil 或无任何语言条目）
func (l Locale) IsEmpty() bool {
	return len(l) == 0
}

// TrimSpace 去除所有语言条目的首尾空白，返回新的 Locale
func (l Locale) TrimSpace() Locale {
	if l == nil {
		return nil
	}
	result := make(Locale, len(l))
	for lang, text := range l {
		result[lang] = strings.TrimSpace(text)
	}
	return result
}

// Get 获取指定语言的文本，不存在时返回空字符串
func (l Locale) Get(lang string) string {
	if l == nil {
		return ""
	}
	return l[lang]
}

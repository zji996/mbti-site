package i18n

import (
	"strings"
)

// Manager 国际化管理器
type Manager struct {
	supportedLanguages map[string]bool
	defaultLanguage    string
}

// NewManager 创建新的国际化管理器
func NewManager() *Manager {
	return &Manager{
		supportedLanguages: map[string]bool{
			"zh": true,
			"en": true,
		},
		defaultLanguage: "zh",
	}
}

// IsValidLanguage 检查语言是否有效
func (m *Manager) IsValidLanguage(lang string) bool {
	return m.supportedLanguages[strings.ToLower(lang)]
}

// GetDefaultLanguage 获取默认语言
func (m *Manager) GetDefaultLanguage() string {
	return m.defaultLanguage
}

// GetSupportedLanguages 获取支持的语言列表
func (m *Manager) GetSupportedLanguages() []string {
	var languages []string
	for lang := range m.supportedLanguages {
		languages = append(languages, lang)
	}
	return languages
}

// NormalizeLanguage 标准化语言代码
func (m *Manager) NormalizeLanguage(lang string) string {
	normalized := strings.ToLower(lang)
	if m.IsValidLanguage(normalized) {
		return normalized
	}
	return m.defaultLanguage
}

// GetFallbackLanguage 获取回退语言
func (m *Manager) GetFallbackLanguage(lang string) string {
	if m.IsValidLanguage(lang) {
		return lang
	}
	return m.defaultLanguage
}

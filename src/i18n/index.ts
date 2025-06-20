import { createI18n } from 'vue-i18n'
import type { Language } from '@/types'

// 导入语言文件
import zh from './locales/zh.json'
import en from './locales/en.json'

// 创建 i18n 实例
export const i18n = createI18n({
  legacy: false, // 使用 Composition API
  locale: 'zh', // 默认语言
  fallbackLocale: 'zh', // 回退语言
  messages: {
    zh,
    en
  },
  globalInjection: true, // 全局注入 $t
})

// 设置语言
export function setI18nLanguage(locale: Language) {
  i18n.global.locale.value = locale
  document.querySelector('html')?.setAttribute('lang', locale)
}

// 获取当前语言
export function getCurrentLanguage(): Language {
  return i18n.global.locale.value as Language
}

// 检查是否支持该语言
export function isLanguageSupported(locale: string): locale is Language {
  return ['zh', 'en'].includes(locale)
}

export default i18n

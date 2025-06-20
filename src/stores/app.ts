import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Language } from '@/types'

export const useAppStore = defineStore('app', () => {
  // 状态
  const language = ref<Language>('zh')
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  // 计算属性
  const isLoading = computed(() => loading.value)
  const hasError = computed(() => !!error.value)
  
  // 动作
  const setLanguage = (lang: Language) => {
    language.value = lang
    // 保存到本地存储
    localStorage.setItem('mbti-language', lang)
  }
  
  const setLoading = (state: boolean) => {
    loading.value = state
  }
  
  const setError = (message: string | null) => {
    error.value = message
  }
  
  const clearError = () => {
    error.value = null
  }
  
  // 初始化语言设置
  const initializeLanguage = () => {
    // 从本地存储获取语言设置
    const savedLang = localStorage.getItem('mbti-language') as Language
    if (savedLang && ['zh', 'en'].includes(savedLang)) {
      language.value = savedLang
    } else {
      // 根据浏览器语言设置默认语言
      const browserLang = navigator.language.toLowerCase()
      if (browserLang.startsWith('zh')) {
        language.value = 'zh'
      } else {
        language.value = 'en'
      }
    }
  }
  
  return {
    // 状态
    language,
    loading,
    error,
    
    // 计算属性
    isLoading,
    hasError,
    
    // 动作
    setLanguage,
    setLoading,
    setError,
    clearError,
    initializeLanguage,
  }
})

import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import type { ThemeMode } from '@/types'

export const useThemeStore = defineStore('theme', () => {
  // 状态
  const mode = ref<ThemeMode>('system')
  const isDark = ref(false)
  
  // 计算属性
  const currentTheme = computed(() => {
    if (mode.value === 'system') {
      return isDark.value ? 'dark' : 'light'
    }
    return mode.value
  })
  
  const isSystemMode = computed(() => mode.value === 'system')
  
  // 动作
  const setMode = (newMode: ThemeMode) => {
    mode.value = newMode
    applyTheme()
    saveToStorage()
  }
  
  const toggleMode = () => {
    if (mode.value === 'light') {
      setMode('dark')
    } else if (mode.value === 'dark') {
      setMode('system')
    } else {
      setMode('light')
    }
  }
  
  const applyTheme = () => {
    const html = document.documentElement
    
    if (mode.value === 'dark' || (mode.value === 'system' && isDark.value)) {
      html.classList.add('dark')
    } else {
      html.classList.remove('dark')
    }
  }
  
  const detectSystemTheme = () => {
    if (typeof window !== 'undefined') {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      isDark.value = mediaQuery.matches
      
      // 监听系统主题变化
      mediaQuery.addEventListener('change', (e) => {
        isDark.value = e.matches
        if (mode.value === 'system') {
          applyTheme()
        }
      })
    }
  }
  
  const saveToStorage = () => {
    localStorage.setItem('mbti-theme-mode', mode.value)
  }
  
  const loadFromStorage = () => {
    try {
      const saved = localStorage.getItem('mbti-theme-mode') as ThemeMode
      if (saved && ['light', 'dark', 'system'].includes(saved)) {
        mode.value = saved
      }
    } catch (error) {
      console.error('Failed to load theme from storage:', error)
    }
  }
  
  const initialize = () => {
    loadFromStorage()
    detectSystemTheme()
    applyTheme()
  }
  
  // 监听模式变化
  watch(mode, () => {
    applyTheme()
  })
  
  // 监听系统主题变化
  watch(isDark, () => {
    if (mode.value === 'system') {
      applyTheme()
    }
  })
  
  return {
    // 状态
    mode,
    isDark,
    
    // 计算属性
    currentTheme,
    isSystemMode,
    
    // 动作
    setMode,
    toggleMode,
    initialize,
  }
})

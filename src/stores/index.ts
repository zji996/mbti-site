import { createPinia } from 'pinia'

// 创建 Pinia 实例
export const pinia = createPinia()

// 导出所有 store
export { useAppStore } from './app'
export { usePersonalityStore } from './personality'
export { useThemeStore } from './theme'

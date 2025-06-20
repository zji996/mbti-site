<template>
  <div id="app" class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-200">
    <!-- 导航栏 -->
    <nav class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <router-link
              :to="createLocalizedPath('/types')"
              class="text-xl font-bold text-gray-900 dark:text-gray-100 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
            >
              MBTI Explorer
            </router-link>
          </div>

          <!-- 导航链接 -->
          <div class="hidden md:flex items-center space-x-8">
            <router-link
              :to="createLocalizedPath('/types')"
              class="nav-link"
              :class="{ 'nav-link-active': $route.name === 'TypesIndex' }"
            >
              {{ $t('nav.types') }}
            </router-link>
          </div>

          <!-- 工具栏 -->
          <div class="flex items-center space-x-4">
            <!-- 语言切换 -->
            <div class="relative">
              <select
                v-model="currentLanguage"
                @change="handleLanguageChange"
                class="appearance-none bg-transparent border border-gray-300 dark:border-gray-600 rounded px-3 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="zh">中文</option>
                <option value="en">English</option>
              </select>
            </div>

            <!-- 主题切换 -->
            <button
              @click="themeStore.toggleMode()"
              class="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
              :title="$t('theme.title')"
            >
              <svg v-if="themeStore.currentTheme === 'light'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              <svg v-else-if="themeStore.currentTheme === 'dark'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容区域 -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- 页脚 -->
    <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center text-gray-600 dark:text-gray-400">
          <p>&copy; 2025 MBTI Explorer. Made with ❤️ for personality exploration.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore, useThemeStore, usePersonalityStore } from '@/stores'
import { createLocalizedPath, getLanguageFromRoute } from '@/router'
import { setI18nLanguage } from '@/i18n'
import type { Language } from '@/types'

// 组合式API
const route = useRoute()
const router = useRouter()
const { locale } = useI18n()

// 状态管理
const appStore = useAppStore()
const themeStore = useThemeStore()
const personalityStore = usePersonalityStore()

// 计算属性
const currentLanguage = computed({
  get: () => appStore.language,
  set: (value: Language) => {
    appStore.setLanguage(value)
  }
})

// 方法
const handleLanguageChange = () => {
  // 更新 i18n 语言
  setI18nLanguage(currentLanguage.value)

  // 更新路由
  const newPath = route.fullPath.replace(
    `/${getLanguageFromRoute(route)}`,
    `/${currentLanguage.value}`
  )
  router.push(newPath)
}

// 监听路由变化，同步语言设置
watch(() => route.params.lang, (newLang) => {
  if (newLang && newLang !== currentLanguage.value) {
    appStore.setLanguage(newLang as Language)
    setI18nLanguage(newLang as Language)
  }
}, { immediate: true })

// 组件挂载时初始化
onMounted(() => {
  // 初始化应用状态
  appStore.initializeLanguage()
  themeStore.initialize()
  personalityStore.initialize()

  // 设置初始语言
  const routeLang = getLanguageFromRoute(route)
  if (routeLang !== currentLanguage.value) {
    appStore.setLanguage(routeLang)
  }
  setI18nLanguage(currentLanguage.value)
})
</script>

<style scoped>
.nav-link {
  color: #374151;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-link:hover {
  color: #2563eb;
}

.nav-link-active {
  color: #2563eb;
}

.dark .nav-link {
  color: #d1d5db;
}

.dark .nav-link:hover {
  color: #60a5fa;
}

.dark .nav-link-active {
  color: #60a5fa;
}
</style>

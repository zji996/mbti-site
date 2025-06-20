<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
    <div class="text-center">
      <div class="text-gray-400 dark:text-gray-500 mb-8">
        <svg class="mx-auto h-24 w-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" 
                d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 15c-2.34 0-4.47-.881-6.08-2.33" />
        </svg>
      </div>
      
      <h1 class="text-6xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        404
      </h1>
      
      <h2 class="text-2xl font-semibold text-gray-700 dark:text-gray-300 mb-4">
        {{ $t('errors.notFound') }}
      </h2>
      
      <p class="text-gray-500 dark:text-gray-400 mb-8 max-w-md mx-auto">
        您访问的页面不存在，可能已被删除或移动到其他位置。
      </p>
      
      <div class="space-x-4">
        <button 
          @click="goBack"
          class="btn-secondary"
        >
          {{ $t('common.back') }}
        </button>
        
        <router-link 
          :to="createLocalizedPath('/types')"
          class="btn-primary"
        >
          {{ $t('nav.home') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { createLocalizedPath } from '@/router'

const router = useRouter()

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push(createLocalizedPath('/types'))
  }
}
</script>

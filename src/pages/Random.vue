<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="text-center">
        <h1 class="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          {{ $t('nav.random') }}
        </h1>
        <p class="text-lg text-gray-600 dark:text-gray-300 mb-8">
          随机抽卡功能即将推出...
        </p>
        
        <router-link 
          :to="createLocalizedPath('/types')"
          class="btn-primary"
        >
          返回人格类型
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { createLocalizedPath } from '@/router'
</script>

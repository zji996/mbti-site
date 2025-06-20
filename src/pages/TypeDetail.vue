<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center min-h-screen">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="flex justify-center items-center min-h-screen">
      <div class="text-center">
        <div class="text-red-500 mb-4">
          <svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 15.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          加载失败
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-4">{{ error }}</p>
        <button 
          class="btn-primary"
          @click="loadPersonalityData"
        >
          重试
        </button>
      </div>
    </div>

    <!-- 主要内容 -->
    <div v-else-if="personalityData" class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 返回按钮 -->
      <div class="mb-6">
        <button 
          class="flex items-center text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          @click="goBack"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          返回列表
        </button>
      </div>

      <!-- 头部信息 -->
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 mb-8">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between">
          <div class="flex-1">
            <div class="flex items-center mb-4">
              <h1 
                class="text-4xl font-bold mr-4"
                :style="{ color: temperamentColor }"
              >
                {{ personalityData.code }}
              </h1>
              <span 
                class="px-3 py-1 rounded-full text-sm font-medium text-white"
                :style="{ backgroundColor: temperamentColor }"
              >
                {{ temperamentName }}
              </span>
            </div>
            
            <h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-2">
              {{ personalityData.name }}
            </h2>
            
            <p class="text-lg text-gray-600 dark:text-gray-300 leading-relaxed">
              {{ personalityData.tagline }}
            </p>
          </div>

          <!-- 性别切换和操作按钮 -->
          <div class="mt-6 lg:mt-0 lg:ml-8 flex flex-col space-y-4">
            <!-- 性别切换 -->
            <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
              <button
                v-for="gender in genders"
                :key="gender.value"
                class="px-4 py-2 text-sm rounded transition-colors"
                :class="[
                  currentGender === gender.value
                    ? 'bg-white dark:bg-gray-600 shadow-sm text-gray-900 dark:text-gray-100'
                    : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
                ]"
                @click="switchGender(gender.value)"
              >
                {{ gender.label }}
              </button>
            </div>

            <!-- 操作按钮 -->
            <div class="flex space-x-2">
              <button class="btn-secondary text-sm">
                分享
              </button>
              <button class="btn-primary text-sm">
                查看配对
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 概要信息 -->
      <div v-if="personalityData.summary" class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8 mb-8">
        <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-4">
          概要
        </h3>
        <p class="text-gray-700 dark:text-gray-300 leading-relaxed">
          {{ personalityData.summary }}
        </p>
      </div>

      <!-- 详细内容 -->
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg p-8">
        <div 
          class="prose prose-lg dark:prose-invert max-w-none"
          v-html="personalityData.body"
        ></div>
      </div>

      <!-- 相关推荐 -->
      <div class="mt-8">
        <h3 class="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-6">
          相关人格类型
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div
            v-for="relatedCode in relatedTypes"
            :key="relatedCode"
            class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 hover:shadow-lg transition-shadow cursor-pointer"
            @click="navigateToType(relatedCode, currentGender)"
          >
            <div class="flex items-center justify-between mb-2">
              <h4 
                class="font-semibold"
                :style="{ color: getTemperamentColor(getTemperament(relatedCode)) }"
              >
                {{ relatedCode }}
              </h4>
              <span 
                class="text-xs px-2 py-1 rounded text-white"
                :style="{ backgroundColor: getTemperamentColor(getTemperament(relatedCode)) }"
              >
                {{ getTemperamentName(getTemperament(relatedCode)) }}
              </span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              点击查看详情
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { PersonalityCode, Gender, PersonalityType } from '@/types'
import { 
  getTemperament, 
  getTemperamentColor, 
  getTemperamentName,
  PERSONALITY_BY_TEMPERAMENT 
} from '@/types'
import { loadPersonalityContent } from '@/utils/markdown'

// 路由
const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(true)
const error = ref<string | null>(null)
const personalityData = ref<PersonalityType | null>(null)
const currentGender = ref<Gender>('m')

// 配置数据
const genders = [
  { value: 'm' as Gender, label: '♂ 男性' },
  { value: 'f' as Gender, label: '♀ 女性' },
]

// 计算属性
const personalityCode = computed(() => {
  const code = route.params.code as string
  return code?.toUpperCase() as PersonalityCode
})

const temperament = computed(() => 
  personalityData.value ? getTemperament(personalityData.value.code) : null
)

const temperamentColor = computed(() => 
  temperament.value ? getTemperamentColor(temperament.value) : '#6366f1'
)

const temperamentName = computed(() => 
  temperament.value ? getTemperamentName(temperament.value) : ''
)

// 相关类型（同气质的其他类型）
const relatedTypes = computed(() => {
  if (!temperament.value || !personalityData.value) return []
  
  return PERSONALITY_BY_TEMPERAMENT[temperament.value]
    .filter(code => code !== personalityData.value!.code)
    .slice(0, 3) // 只显示3个相关类型
})

// 方法
const loadPersonalityData = async () => {
  if (!personalityCode.value) {
    error.value = '无效的人格类型代码'
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = null
    
    const data = await loadPersonalityContent(personalityCode.value, currentGender.value)
    if (data) {
      personalityData.value = data
    } else {
      error.value = '未找到该人格类型的数据'
    }
  } catch (err) {
    error.value = '加载数据时发生错误'
    console.error('Failed to load personality data:', err)
  } finally {
    loading.value = false
  }
}

const switchGender = (gender: Gender) => {
  if (currentGender.value !== gender) {
    currentGender.value = gender
    // 更新URL
    router.replace({
      name: 'TypeDetail',
      params: { 
        code: personalityCode.value.toLowerCase(),
        gender: gender
      }
    })
  }
}

const goBack = () => {
  router.push({ name: 'TypesIndex' })
}

const navigateToType = (code: PersonalityCode, gender: Gender) => {
  router.push({
    name: 'TypeDetail',
    params: {
      code: code.toLowerCase(),
      gender: gender
    }
  })
}

// 从路由参数获取性别
const getGenderFromRoute = () => {
  const gender = route.params.gender as string
  if (gender === 'f' || gender === 'm') {
    currentGender.value = gender
  }
}

// 监听路由变化
watch(() => route.params, () => {
  getGenderFromRoute()
  loadPersonalityData()
}, { immediate: true })

// 监听性别变化
watch(currentGender, () => {
  loadPersonalityData()
})

// 组件挂载时
onMounted(() => {
  getGenderFromRoute()
})
</script>

<style scoped>
/* 自定义prose样式以适应深色模式 */
.prose {
  color: #374151;
}

.dark .prose {
  color: #d1d5db;
}

.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6 {
  color: #111827;
}

.dark .prose h1, .dark .prose h2, .dark .prose h3, .dark .prose h4, .dark .prose h5, .dark .prose h6 {
  color: #f9fafb;
}

.prose strong {
  color: #111827;
}

.dark .prose strong {
  color: #f9fafb;
}

.prose a {
  color: #2563eb;
}

.dark .prose a {
  color: #60a5fa;
}
</style>

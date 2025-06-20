<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- 页面头部 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center">
          <h1 class="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
            MBTI 人格类型探索
          </h1>
          <p class="text-lg text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            探索16种不同的人格类型，了解每种类型的独特特征、优势和成长方向。
            点击卡片查看详细信息，切换性别查看不同视角。
          </p>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm p-6 mb-8">
        <div class="flex flex-col lg:flex-row gap-4 items-center">
          <!-- 搜索框 -->
          <div class="flex-1 max-w-md">
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="搜索人格类型..."
                class="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg 
                       bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100
                       focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- 气质筛选 -->
          <div class="flex flex-wrap gap-2">
            <button
              v-for="temperament in temperaments"
              :key="temperament.type"
              class="px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              :class="[
                selectedTemperament === temperament.type
                  ? 'text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              ]"
              :style="selectedTemperament === temperament.type ? { backgroundColor: temperament.color } : {}"
              @click="toggleTemperament(temperament.type)"
            >
              {{ temperament.name }}
            </button>
          </div>

          <!-- 性别切换 -->
          <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
            <button
              v-for="gender in genders"
              :key="gender.value"
              class="px-4 py-2 text-sm rounded transition-colors"
              :class="[
                defaultGender === gender.value
                  ? 'bg-white dark:bg-gray-600 shadow-sm text-gray-900 dark:text-gray-100'
                  : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
              ]"
              @click="defaultGender = gender.value"
            >
              {{ gender.label }}
            </button>
          </div>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="flex justify-between items-center mb-6">
        <div class="text-sm text-gray-600 dark:text-gray-400">
          显示 {{ filteredPersonalities.length }} / {{ allPersonalities.length }} 种人格类型
        </div>
        
        <!-- 视图切换 -->
        <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
          <button
            v-for="view in viewModes"
            :key="view.value"
            class="px-3 py-1 text-sm rounded transition-colors"
            :class="[
              currentView === view.value
                ? 'bg-white dark:bg-gray-600 shadow-sm'
                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            ]"
            @click="currentView = view.value"
          >
            {{ view.label }}
          </button>
        </div>
      </div>

      <!-- 人格卡片网格 -->
      <div 
        class="grid gap-6 animate-fade-in"
        :class="gridClasses"
      >
        <PersonalityCard
          v-for="personality in filteredPersonalities"
          :key="`${personality}-${defaultGender}`"
          :code="personality"
          :gender="defaultGender"
          :auto-flip="currentView === 'hover'"
          @gender-change="handleGenderChange"
          @card-click="handleCardClick"
          @view-details="handleViewDetails"
        />
      </div>

      <!-- 空状态 -->
      <div 
        v-if="filteredPersonalities.length === 0"
        class="text-center py-12"
      >
        <div class="text-gray-400 dark:text-gray-500 mb-4">
          <svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 15c-2.34 0-4.47-.881-6.08-2.33" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
          没有找到匹配的人格类型
        </h3>
        <p class="text-gray-500 dark:text-gray-400">
          尝试调整搜索条件或筛选器
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import PersonalityCard from '@/components/PersonalityCard.vue'
import type { PersonalityCode, Gender, TemperamentType } from '@/types'
import { 
  ALL_PERSONALITY_CODES, 
  PERSONALITY_BY_TEMPERAMENT,
  getTemperament,
  getTemperamentColor,
  getTemperamentName
} from '@/types'

// 响应式数据
const searchQuery = ref('')
const selectedTemperament = ref<TemperamentType | null>(null)
const defaultGender = ref<Gender>('m')
const currentView = ref<'grid' | 'hover'>('grid')

// 配置数据
const temperaments = [
  { type: 'NT' as TemperamentType, name: '理性者', color: '#6366f1' },
  { type: 'NF' as TemperamentType, name: '理想主义者', color: '#10b981' },
  { type: 'SJ' as TemperamentType, name: '守护者', color: '#ef4444' },
  { type: 'SP' as TemperamentType, name: '艺术家', color: '#f59e0b' },
]

const genders = [
  { value: 'm' as Gender, label: '♂ 男性' },
  { value: 'f' as Gender, label: '♀ 女性' },
]

const viewModes = [
  { value: 'grid', label: '网格' },
  { value: 'hover', label: '悬停' },
]

const allPersonalities = ALL_PERSONALITY_CODES

// 计算属性
const filteredPersonalities = computed(() => {
  let filtered = [...allPersonalities]

  // 按搜索查询筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(code => 
      code.toLowerCase().includes(query) ||
      getTemperamentName(getTemperament(code)).includes(query)
    )
  }

  // 按气质筛选
  if (selectedTemperament.value) {
    filtered = PERSONALITY_BY_TEMPERAMENT[selectedTemperament.value]
    
    // 如果有搜索查询，需要进一步筛选
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(code => 
        code.toLowerCase().includes(query) ||
        getTemperamentName(getTemperament(code)).includes(query)
      )
    }
  }

  return filtered
})

const gridClasses = computed(() => {
  const baseClasses = 'grid gap-6'
  
  switch (currentView.value) {
    case 'hover':
      return `${baseClasses} grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4`
    default:
      return `${baseClasses} grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4`
  }
})

// 方法
const toggleTemperament = (temperament: TemperamentType) => {
  selectedTemperament.value = selectedTemperament.value === temperament ? null : temperament
}

const handleGenderChange = (gender: Gender) => {
  defaultGender.value = gender
}

const handleCardClick = (code: PersonalityCode, gender: Gender) => {
  console.log('Card clicked:', code, gender)
  // 这里可以添加卡片点击的逻辑，比如显示详情模态框
}

const handleViewDetails = (code: PersonalityCode, gender: Gender) => {
  console.log('View details:', code, gender)
  // 这里可以导航到详情页面
  // router.push(`/types/${code.toLowerCase()}-${gender}`)
}
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

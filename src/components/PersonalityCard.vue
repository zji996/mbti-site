<template>
  <div 
    class="personality-card group cursor-pointer relative"
    :style="{ borderColor: temperamentColor }"
    @click="handleCardClick"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <!-- 卡片翻转容器 -->
    <div 
      class="card-inner relative w-full h-64 transition-transform duration-600 transform-style-preserve-3d"
      :class="{ 'rotate-y-180': isFlipped }"
    >
      <!-- 正面 -->
      <div class="card-face card-front absolute inset-0 backface-hidden">
        <div class="personality-card-front h-full flex flex-col justify-between">
          <!-- 头部：人格代码和性别切换 -->
          <div class="flex justify-between items-start">
            <div class="flex items-center space-x-2">
              <h3 class="text-2xl font-bold" :style="{ color: temperamentColor }">
                {{ personalityData.code }}
              </h3>
              <span class="text-sm text-gray-500 dark:text-gray-400">
                {{ temperamentName }}
              </span>
            </div>
            
            <!-- 性别切换按钮 -->
            <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
              <button
                v-for="gender in ['m', 'f']"
                :key="gender"
                class="px-2 py-1 text-xs rounded transition-colors"
                :class="[
                  currentGender === gender
                    ? 'bg-white dark:bg-gray-600 shadow-sm'
                    : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
                ]"
                @click.stop="switchGender(gender)"
              >
                {{ gender === 'm' ? '♂' : '♀' }}
              </button>
            </div>
          </div>

          <!-- 中部：名称和标语 -->
          <div class="flex-1 flex flex-col justify-center text-center">
            <h4 class="text-xl font-semibold mb-2 text-gray-900 dark:text-gray-100">
              {{ personalityData.name }}
            </h4>
            <p class="text-sm text-gray-600 dark:text-gray-300 leading-relaxed">
              {{ personalityData.tagline }}
            </p>
          </div>

          <!-- 底部：气质标识 -->
          <div class="flex justify-center">
            <span 
              class="px-3 py-1 rounded-full text-xs font-medium text-white"
              :style="{ backgroundColor: temperamentColor }"
            >
              {{ temperamentName }}
            </span>
          </div>
        </div>
      </div>

      <!-- 背面 -->
      <div class="card-face card-back absolute inset-0 backface-hidden rotate-y-180">
        <div class="personality-card-back h-full">
          <div class="flex flex-col h-full">
            <!-- 背面头部 -->
            <div class="flex justify-between items-center mb-4">
              <h4 class="text-lg font-semibold" :style="{ color: temperamentColor }">
                {{ personalityData.name }}
              </h4>
              <button 
                class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                @click.stop="toggleFlip"
              >
                ✕
              </button>
            </div>

            <!-- 背面内容 -->
            <div class="flex-1 overflow-y-auto">
              <p class="text-sm text-gray-700 dark:text-gray-300 leading-relaxed mb-4">
                {{ personalityData.summary || '暂无简介内容...' }}
              </p>
              
              <!-- 特征标签 -->
              <div class="space-y-2">
                <div class="text-xs text-gray-500 dark:text-gray-400">特征标签:</div>
                <div class="flex flex-wrap gap-1">
                  <span 
                    v-for="trait in traits" 
                    :key="trait"
                    class="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-xs rounded"
                  >
                    {{ trait }}
                  </span>
                </div>
              </div>
            </div>

            <!-- 背面底部：查看详情按钮 -->
            <div class="mt-4">
              <button 
                class="w-full btn-primary text-sm py-2"
                @click.stop="viewDetails"
              >
                查看详情
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 悬停效果指示器 -->
    <div 
      class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
      v-if="!isFlipped"
    >
      <div class="w-2 h-2 bg-gray-400 rounded-full"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { PersonalityCode, Gender, PersonalityType } from '@/types'
import { getTemperament, getTemperamentColor, getTemperamentName } from '@/types'
import { loadPersonalityContent } from '@/utils/markdown'

// Props
interface Props {
  code: PersonalityCode
  gender?: Gender
  autoFlip?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  gender: 'm',
  autoFlip: false,
})

// Emits
const emit = defineEmits<{
  genderChange: [gender: Gender]
  cardClick: [code: PersonalityCode, gender: Gender]
  viewDetails: [code: PersonalityCode, gender: Gender]
}>()

// 响应式数据
const currentGender = ref<Gender>(props.gender)
const isFlipped = ref(false)
const isHovered = ref(false)
const personalityData = ref<PersonalityType | null>(null)

// 计算属性
const temperament = computed(() => getTemperament(props.code))
const temperamentColor = computed(() => getTemperamentColor(temperament.value))
const temperamentName = computed(() => getTemperamentName(temperament.value))

// 模拟特征标签
const traits = computed(() => {
  const traitMap: Record<PersonalityCode, string[]> = {
    'INTJ': ['战略思维', '独立', '理性', '远见'],
    'INTP': ['逻辑', '创新', '理论', '分析'],
    'ENTJ': ['领导力', '决断', '效率', '目标导向'],
    'ENTP': ['创意', '辩论', '适应性', '概念思维'],
    'INFJ': ['洞察力', '理想主义', '同理心', '直觉'],
    'INFP': ['价值观', '创造力', '和谐', '真实'],
    'ENFJ': ['魅力', '激励', '合作', '成长导向'],
    'ENFP': ['热情', '创新', '灵活', '人际关系'],
    'ISTJ': ['可靠', '传统', '细致', '责任感'],
    'ISFJ': ['关怀', '忠诚', '服务', '和谐'],
    'ESTJ': ['组织', '效率', '传统', '领导'],
    'ESFJ': ['社交', '合作', '支持', '和谐'],
    'ISTP': ['实用', '灵活', '技能', '独立'],
    'ISFP': ['艺术', '价值观', '灵活', '和谐'],
    'ESTP': ['行动', '适应', '实用', '社交'],
    'ESFP': ['热情', '社交', '灵活', '乐观'],
  }
  return traitMap[props.code] || []
})

// 方法
const switchGender = (gender: Gender) => {
  if (currentGender.value !== gender) {
    currentGender.value = gender
    emit('genderChange', gender)
    loadPersonalityData()
  }
}

const toggleFlip = () => {
  isFlipped.value = !isFlipped.value
}

const handleCardClick = () => {
  if (!isFlipped.value) {
    toggleFlip()
  }
  emit('cardClick', props.code, currentGender.value)
}

const handleMouseEnter = () => {
  isHovered.value = true
  if (props.autoFlip && !isFlipped.value) {
    setTimeout(() => {
      if (isHovered.value) {
        isFlipped.value = true
      }
    }, 500)
  }
}

const handleMouseLeave = () => {
  isHovered.value = false
  if (props.autoFlip && isFlipped.value) {
    setTimeout(() => {
      if (!isHovered.value) {
        isFlipped.value = false
      }
    }, 300)
  }
}

const viewDetails = () => {
  emit('viewDetails', props.code, currentGender.value)
}

const loadPersonalityData = async () => {
  try {
    const data = await loadPersonalityContent(props.code, currentGender.value)
    if (data) {
      personalityData.value = data
    }
  } catch (error) {
    console.error('Failed to load personality data:', error)
  }
}

// 监听性别变化
watch(() => props.gender, (newGender) => {
  currentGender.value = newGender
}, { immediate: true })

// 监听当前性别变化，重新加载数据
watch(currentGender, () => {
  loadPersonalityData()
}, { immediate: true })
</script>

<style scoped>
.transform-style-preserve-3d {
  transform-style: preserve-3d;
}

.backface-hidden {
  backface-visibility: hidden;
}

.rotate-y-180 {
  transform: rotateY(180deg);
}

.card-face {
  border-radius: 0.75rem;
}

.card-front {
  background: inherit;
}

.card-back {
  background: inherit;
  transform: rotateY(180deg);
}
</style>

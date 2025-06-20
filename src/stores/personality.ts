import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { 
  PersonalityType, 
  PersonalityPairing, 
  PersonalityCode, 
  Gender, 
  Language,
  TemperamentType 
} from '@/types'
import { loadPersonalityContent } from '@/utils/markdown'

export const usePersonalityStore = defineStore('personality', () => {
  // 状态
  const personalityCache = ref<Map<string, PersonalityType>>(new Map())
  const pairingCache = ref<Map<string, PersonalityPairing>>(new Map())
  const favorites = ref<Set<string>>(new Set())
  const recentlyViewed = ref<string[]>([])
  
  // 计算属性
  const getCachedPersonality = computed(() => {
    return (code: PersonalityCode, gender: Gender, lang: Language = 'zh') => {
      const key = `${code}-${gender}-${lang}`
      return personalityCache.value.get(key)
    }
  })
  
  const getCachedPairing = computed(() => {
    return (maleCode: PersonalityCode, femaleCode: PersonalityCode, lang: Language = 'zh') => {
      const key = `${maleCode}-${femaleCode}-${lang}`
      return pairingCache.value.get(key)
    }
  })
  
  const getFavorites = computed(() => Array.from(favorites.value))
  
  const getRecentlyViewed = computed(() => recentlyViewed.value.slice(0, 10))
  
  // 动作
  const loadPersonality = async (
    code: PersonalityCode, 
    gender: Gender, 
    lang: Language = 'zh'
  ): Promise<PersonalityType | null> => {
    const key = `${code}-${gender}-${lang}`
    
    // 检查缓存
    const cached = personalityCache.value.get(key)
    if (cached) {
      addToRecentlyViewed(key)
      return cached
    }
    
    try {
      // 从API或文件加载
      const data = await loadPersonalityContent(code, gender, lang)
      if (data) {
        // 缓存数据
        personalityCache.value.set(key, data)
        addToRecentlyViewed(key)
        return data
      }
    } catch (error) {
      console.error('Failed to load personality:', error)
    }
    
    return null
  }
  
  const loadPairing = async (
    maleCode: PersonalityCode,
    femaleCode: PersonalityCode,
    lang: Language = 'zh'
  ): Promise<PersonalityPairing | null> => {
    const key = `${maleCode}-${femaleCode}-${lang}`
    
    // 检查缓存
    const cached = pairingCache.value.get(key)
    if (cached) {
      return cached
    }
    
    try {
      // 这里应该调用实际的API
      // const data = await loadPairingContent(maleCode, femaleCode, lang)
      // 目前返回模拟数据
      const mockPairing: PersonalityPairing = {
        id: `${maleCode.toLowerCase()}_m__${femaleCode.toLowerCase()}_f`,
        male_code: maleCode,
        female_code: femaleCode,
        compatibility_score: Math.floor(Math.random() * 100) + 1,
        body: `<p>${maleCode} 男性与 ${femaleCode} 女性的配对分析...</p>`,
        lang
      }
      
      // 缓存数据
      pairingCache.value.set(key, mockPairing)
      return mockPairing
    } catch (error) {
      console.error('Failed to load pairing:', error)
    }
    
    return null
  }
  
  const addToFavorites = (code: PersonalityCode, gender: Gender) => {
    const key = `${code}-${gender}`
    favorites.value.add(key)
    saveFavoritesToStorage()
  }
  
  const removeFromFavorites = (code: PersonalityCode, gender: Gender) => {
    const key = `${code}-${gender}`
    favorites.value.delete(key)
    saveFavoritesToStorage()
  }
  
  const isFavorite = (code: PersonalityCode, gender: Gender): boolean => {
    const key = `${code}-${gender}`
    return favorites.value.has(key)
  }
  
  const addToRecentlyViewed = (key: string) => {
    // 移除已存在的项目
    const index = recentlyViewed.value.indexOf(key)
    if (index > -1) {
      recentlyViewed.value.splice(index, 1)
    }
    
    // 添加到开头
    recentlyViewed.value.unshift(key)
    
    // 限制数量
    if (recentlyViewed.value.length > 20) {
      recentlyViewed.value = recentlyViewed.value.slice(0, 20)
    }
    
    saveRecentlyViewedToStorage()
  }
  
  const clearCache = () => {
    personalityCache.value.clear()
    pairingCache.value.clear()
  }
  
  const clearRecentlyViewed = () => {
    recentlyViewed.value = []
    saveRecentlyViewedToStorage()
  }
  
  // 本地存储相关
  const saveFavoritesToStorage = () => {
    localStorage.setItem('mbti-favorites', JSON.stringify(Array.from(favorites.value)))
  }
  
  const loadFavoritesFromStorage = () => {
    try {
      const saved = localStorage.getItem('mbti-favorites')
      if (saved) {
        const parsed = JSON.parse(saved) as string[]
        favorites.value = new Set(parsed)
      }
    } catch (error) {
      console.error('Failed to load favorites from storage:', error)
    }
  }
  
  const saveRecentlyViewedToStorage = () => {
    localStorage.setItem('mbti-recently-viewed', JSON.stringify(recentlyViewed.value))
  }
  
  const loadRecentlyViewedFromStorage = () => {
    try {
      const saved = localStorage.getItem('mbti-recently-viewed')
      if (saved) {
        recentlyViewed.value = JSON.parse(saved)
      }
    } catch (error) {
      console.error('Failed to load recently viewed from storage:', error)
    }
  }
  
  // 初始化
  const initialize = () => {
    loadFavoritesFromStorage()
    loadRecentlyViewedFromStorage()
  }
  
  // 搜索和筛选相关
  const searchPersonalities = (
    query: string,
    temperament?: TemperamentType,
    gender?: Gender
  ): PersonalityCode[] => {
    // 这里应该实现实际的搜索逻辑
    // 目前返回模拟结果
    const allCodes: PersonalityCode[] = [
      'INTJ', 'INTP', 'ENTJ', 'ENTP',
      'INFJ', 'INFP', 'ENFJ', 'ENFP',
      'ISTJ', 'ISFJ', 'ESTJ', 'ESFJ',
      'ISTP', 'ISFP', 'ESTP', 'ESFP'
    ]
    
    return allCodes.filter(code => {
      if (query && !code.toLowerCase().includes(query.toLowerCase())) {
        return false
      }
      // 这里可以添加更多筛选逻辑
      return true
    })
  }
  
  return {
    // 状态
    personalityCache,
    pairingCache,
    favorites,
    recentlyViewed,
    
    // 计算属性
    getCachedPersonality,
    getCachedPairing,
    getFavorites,
    getRecentlyViewed,
    
    // 动作
    loadPersonality,
    loadPairing,
    addToFavorites,
    removeFromFavorites,
    isFavorite,
    addToRecentlyViewed,
    clearCache,
    clearRecentlyViewed,
    initialize,
    searchPersonalities,
  }
})

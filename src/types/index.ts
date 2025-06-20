// 性别类型
export type Gender = 'm' | 'f'

// MBTI 人格类型代码
export type PersonalityCode = 
  | 'INTJ' | 'INTP' | 'ENTJ' | 'ENTP'  // NT - 理性者
  | 'INFJ' | 'INFP' | 'ENFJ' | 'ENFP'  // NF - 理想主义者
  | 'ISTJ' | 'ISFJ' | 'ESTJ' | 'ESFJ'  // SJ - 守护者
  | 'ISTP' | 'ISFP' | 'ESTP' | 'ESFP'  // SP - 艺术家

// 四种气质类型
export type TemperamentType = 'NT' | 'NF' | 'SJ' | 'SP'

// 人格类型接口
export interface PersonalityType {
  code: PersonalityCode
  gender: Gender
  name: string
  tagline: string
  summary: string
  body: string
  lang: string
}

// 人格配对接口
export interface PersonalityPairing {
  id: string // 格式: "intj_m__infp_f"
  male_code: PersonalityCode
  female_code: PersonalityCode
  compatibility_score: number
  body: string
  lang: string
}

// 前端显示用的人格卡片数据
export interface PersonalityCard {
  code: PersonalityCode
  gender: Gender
  name: string
  tagline: string
  summary: string
  temperament: TemperamentType
  color: string // 对应的主题色
}

// API 响应类型
export interface ApiResponse<T> {
  data: T
  success: boolean
  message?: string
}

// 分页响应类型
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  success: boolean
}

// 搜索和筛选参数
export interface TypesSearchParams {
  gender?: Gender
  search?: string
  temperament?: TemperamentType
}

// 语言类型
export type Language = 'zh' | 'en'

// 主题模式
export type ThemeMode = 'light' | 'dark' | 'system'

// 工具函数：获取气质类型
export function getTemperament(code: PersonalityCode): TemperamentType {
  const firstTwo = code.slice(0, 2)
  if (firstTwo === 'NT' || (code[1] === 'N' && code[2] === 'T')) return 'NT'
  if (firstTwo === 'NF' || (code[1] === 'N' && code[2] === 'F')) return 'NF'
  if (firstTwo === 'SJ' || (code[1] === 'S' && code[2] === 'J')) return 'SJ'
  if (firstTwo === 'SP' || (code[1] === 'S' && code[2] === 'P')) return 'SP'
  
  // 备用逻辑
  if (code.includes('N') && code.includes('T')) return 'NT'
  if (code.includes('N') && code.includes('F')) return 'NF'
  if (code.includes('S') && code.includes('J')) return 'SJ'
  return 'SP'
}

// 工具函数：获取气质对应的颜色
export function getTemperamentColor(temperament: TemperamentType): string {
  const colors = {
    'NT': '#6366f1', // 理性者 - 靛蓝
    'NF': '#10b981', // 理想主义者 - 绿色
    'SP': '#f59e0b', // 艺术家 - 橙色
    'SJ': '#ef4444', // 守护者 - 红色
  }
  return colors[temperament]
}

// 工具函数：获取气质中文名称
export function getTemperamentName(temperament: TemperamentType, lang: Language = 'zh'): string {
  const names = {
    zh: {
      'NT': '理性者',
      'NF': '理想主义者',
      'SP': '艺术家',
      'SJ': '守护者',
    },
    en: {
      'NT': 'Rational',
      'NF': 'Idealist',
      'SP': 'Artisan',
      'SJ': 'Guardian',
    }
  }
  return names[lang][temperament]
}

// 所有16种人格类型的列表
export const ALL_PERSONALITY_CODES: PersonalityCode[] = [
  'INTJ', 'INTP', 'ENTJ', 'ENTP',
  'INFJ', 'INFP', 'ENFJ', 'ENFP',
  'ISTJ', 'ISFJ', 'ESTJ', 'ESFJ',
  'ISTP', 'ISFP', 'ESTP', 'ESFP',
]

// 按气质分组的人格类型
export const PERSONALITY_BY_TEMPERAMENT: Record<TemperamentType, PersonalityCode[]> = {
  'NT': ['INTJ', 'INTP', 'ENTJ', 'ENTP'],
  'NF': ['INFJ', 'INFP', 'ENFJ', 'ENFP'],
  'SJ': ['ISTJ', 'ISFJ', 'ESTJ', 'ESFJ'],
  'SP': ['ISTP', 'ISFP', 'ESTP', 'ESFP'],
}

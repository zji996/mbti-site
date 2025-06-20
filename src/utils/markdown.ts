import matter from 'gray-matter'
import { marked } from 'marked'
import type { PersonalityType, PersonalityPairing, PersonalityCode, Gender, Language } from '@/types'

// 配置 marked 选项
marked.setOptions({
  breaks: true,
  gfm: true,
})

// 解析人格类型 Markdown 文件
export function parsePersonalityMarkdown(content: string): PersonalityType {
  const { data, content: body } = matter(content)
  
  return {
    code: data.code as PersonalityCode,
    gender: data.gender as Gender,
    name: data.name || '',
    tagline: data.tagline || '',
    summary: data.summary || '',
    body: marked(body),
    lang: data.lang as Language || 'zh',
  }
}

// 解析配对 Markdown 文件
export function parsePairingMarkdown(content: string): PersonalityPairing {
  const { data, content: body } = matter(content)
  
  return {
    id: data.id || '',
    male_code: data.male_code as PersonalityCode,
    female_code: data.female_code as PersonalityCode,
    compatibility_score: data.compatibility_score || 0,
    body: marked(body),
    lang: data.lang as Language || 'zh',
  }
}

// 从文件路径获取人格类型信息
export function getPersonalityInfoFromPath(filePath: string): { code: PersonalityCode; gender: Gender } {
  // 文件名格式: intj_m.md 或 intj_f.md
  const fileName = filePath.split('/').pop()?.replace('.md', '') || ''
  const [code, gender] = fileName.split('_')
  
  return {
    code: code.toUpperCase() as PersonalityCode,
    gender: gender as Gender,
  }
}

// 生成人格类型的文件路径
export function getPersonalityFilePath(code: PersonalityCode, gender: Gender, lang: Language = 'zh'): string {
  return `content/${lang}/types/${code.toLowerCase()}_${gender}.md`
}

// 生成配对的文件路径
export function getPairingFilePath(maleCode: PersonalityCode, femaleCode: PersonalityCode, lang: Language = 'zh'): string {
  return `content/${lang}/pairings/${maleCode.toLowerCase()}_m__${femaleCode.toLowerCase()}_f.md`
}

// 模拟从文件系统读取文件（在实际应用中，这将通过 API 调用）
export async function loadPersonalityContent(code: PersonalityCode, gender: Gender, lang: Language = 'zh'): Promise<PersonalityType | null> {
  try {
    const filePath = getPersonalityFilePath(code, gender, lang)
    // 在实际应用中，这里会是一个 fetch 调用到后端 API
    // const response = await fetch(`/api/v1/${lang}/types/${code.toLowerCase()}-${gender}`)
    // const data = await response.json()
    // return data
    
    // 目前返回模拟数据
    return createMockPersonalityType(code, gender, lang)
  } catch (error) {
    console.error('Failed to load personality content:', error)
    return null
  }
}

// 创建模拟的人格类型数据
function createMockPersonalityType(code: PersonalityCode, gender: Gender, lang: Language): PersonalityType {
  const names = {
    zh: {
      'INTJ': '建筑师',
      'INTP': '逻辑学家',
      'ENTJ': '指挥官',
      'ENTP': '辩论家',
      'INFJ': '提倡者',
      'INFP': '调停者',
      'ENFJ': '主人公',
      'ENFP': '竞选者',
      'ISTJ': '物流师',
      'ISFJ': '守卫者',
      'ESTJ': '总经理',
      'ESFJ': '执政官',
      'ISTP': '鉴赏家',
      'ISFP': '探险家',
      'ESTP': '企业家',
      'ESFP': '娱乐家',
    },
    en: {
      'INTJ': 'Architect',
      'INTP': 'Logician',
      'ENTJ': 'Commander',
      'ENTP': 'Debater',
      'INFJ': 'Advocate',
      'INFP': 'Mediator',
      'ENFJ': 'Protagonist',
      'ENFP': 'Campaigner',
      'ISTJ': 'Logistician',
      'ISFJ': 'Defender',
      'ESTJ': 'Executive',
      'ESFJ': 'Consul',
      'ISTP': 'Virtuoso',
      'ISFP': 'Adventurer',
      'ESTP': 'Entrepreneur',
      'ESFP': 'Entertainer',
    }
  }

  const taglines = {
    zh: {
      'INTJ': '深谋远虑，构筑未来的思想建筑师',
      'INTP': '理性思辨，探索真理的逻辑学家',
      'ENTJ': '天生领袖，统筹全局的指挥官',
      'ENTP': '思维敏捷，挑战传统的辩论家',
      'INFJ': '理想主义，洞察人心的提倡者',
      'INFP': '价值驱动，追求和谐的调停者',
      'ENFJ': '富有魅力，激励他人的主人公',
      'ENFP': '热情洋溢，充满创意的竞选者',
      'ISTJ': '踏实可靠，维护秩序的物流师',
      'ISFJ': '温暖体贴，保护他人的守卫者',
      'ESTJ': '高效务实，管理有方的总经理',
      'ESFJ': '善于交际，关怀他人的执政官',
      'ISTP': '技艺精湛，灵活应变的鉴赏家',
      'ISFP': '艺术天赋，自由自在的探险家',
      'ESTP': '行动力强，把握当下的企业家',
      'ESFP': '活力四射，感染他人的娱乐家',
    },
    en: {
      'INTJ': 'Imaginative and strategic thinkers, with a plan for everything',
      'INTP': 'Innovative inventors with an unquenchable thirst for knowledge',
      'ENTJ': 'Bold, imaginative and strong-willed leaders',
      'ENTP': 'Smart and curious thinkers who cannot resist an intellectual challenge',
      'INFJ': 'Quiet and mystical, yet very inspiring and tireless idealists',
      'INFP': 'Poetic, kind and altruistic people, always eager to help a good cause',
      'ENFJ': 'Charismatic and inspiring leaders, able to mesmerize their listeners',
      'ENFP': 'Enthusiastic, creative and sociable free spirits',
      'ISTJ': 'Practical and fact-minded, reliable and responsible',
      'ISFJ': 'Very dedicated and warm protectors, always ready to defend their loved ones',
      'ESTJ': 'Excellent administrators, unsurpassed at managing things or people',
      'ESFJ': 'Extraordinarily caring, social and popular people, always eager to help',
      'ISTP': 'Bold and practical experimenters, masters of all kinds of tools',
      'ISFP': 'Flexible and charming artists, always ready to explore new possibilities',
      'ESTP': 'Smart, energetic and very perceptive people, truly enjoy living on the edge',
      'ESFP': 'Spontaneous, energetic and enthusiastic people - life is never boring',
    }
  }

  return {
    code,
    gender,
    name: names[lang][code],
    tagline: taglines[lang][code],
    summary: `${names[lang][code]}类型的简要描述...`,
    body: `<h2>${names[lang][code]}的详细分析</h2><p>这里是${code}${gender === 'm' ? '男性' : '女性'}的详细内容...</p>`,
    lang,
  }
}

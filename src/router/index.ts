import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { Language } from '@/types'

// 页面组件
import TypesIndex from '@/pages/TypesIndex.vue'
import TypeDetail from '@/pages/TypeDetail.vue'

// 路由配置
const routes: RouteRecordRaw[] = [
  // 根路径重定向到中文首页
  {
    path: '/',
    redirect: '/zh/types'
  },
  
  // 多语言路由
  {
    path: '/:lang(zh|en)',
    component: { template: '<router-view />' },
    children: [
      // 人格类型列表页
      {
        path: 'types',
        name: 'TypesIndex',
        component: TypesIndex,
        meta: {
          title: 'MBTI 人格类型',
          description: '探索16种MBTI人格类型，了解每种类型的特征和优势'
        }
      },
      
      // 人格类型详情页
      {
        path: 'types/:code([a-zA-Z]{4})-:gender(m|f)',
        name: 'TypeDetail',
        component: TypeDetail,
        meta: {
          title: '人格类型详情',
          description: '查看具体人格类型的详细分析和特征'
        }
      },
      
      // 配对探索页（预留）
      {
        path: 'pairings',
        name: 'PairingsIndex',
        component: () => import('@/pages/PairingsIndex.vue'),
        meta: {
          title: '人格配对',
          description: '探索不同人格类型之间的配对关系'
        }
      },
      
      // 配对详情页（预留）
      {
        path: 'pairings/:male([a-zA-Z]{4})-:female([a-zA-Z]{4})',
        name: 'PairingDetail',
        component: () => import('@/pages/PairingDetail.vue'),
        meta: {
          title: '配对详情',
          description: '查看具体人格配对的详细分析'
        }
      },
      
      // 探索模式页（预留）
      {
        path: 'explorer',
        name: 'Explorer',
        component: () => import('@/pages/Explorer.vue'),
        meta: {
          title: '探索模式',
          description: '通过交互式图表探索人格类型关系'
        }
      },
      
      // 随机抽卡页（预留）
      {
        path: 'random',
        name: 'Random',
        component: () => import('@/pages/Random.vue'),
        meta: {
          title: '随机探索',
          description: '随机发现人格类型和配对'
        }
      }
    ]
  },
  
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/pages/NotFound.vue'),
    meta: {
      title: '页面未找到',
      description: '您访问的页面不存在'
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // 如果有保存的位置（如浏览器后退），恢复到该位置
    if (savedPosition) {
      return savedPosition
    }
    // 如果有锚点，滚动到锚点
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth'
      }
    }
    // 否则滚动到顶部
    return { top: 0 }
  }
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 验证语言参数
  const lang = to.params.lang as Language
  if (lang && !['zh', 'en'].includes(lang)) {
    // 如果语言参数无效，重定向到中文版本
    next('/zh' + to.path.replace(`/${lang}`, ''))
    return
  }
  
  // 验证人格类型代码
  if (to.name === 'TypeDetail') {
    const code = to.params.code as string
    const validCodes = [
      'intj', 'intp', 'entj', 'entp',
      'infj', 'infp', 'enfj', 'enfp',
      'istj', 'isfj', 'estj', 'esfj',
      'istp', 'isfp', 'estp', 'esfp'
    ]
    
    if (!validCodes.includes(code.toLowerCase())) {
      // 如果人格类型代码无效，重定向到列表页
      next(`/${lang}/types`)
      return
    }
  }
  
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - MBTI Explorer`
  }
  
  next()
})

// 路由错误处理
router.onError((error) => {
  console.error('Router error:', error)
})

export default router

// 导出路由相关的工具函数
export function getLanguageFromRoute(route: any): Language {
  return (route.params.lang as Language) || 'zh'
}

export function createLocalizedPath(path: string, lang: Language = 'zh'): string {
  return `/${lang}${path.startsWith('/') ? path : '/' + path}`
}

export function switchLanguage(currentRoute: any, newLang: Language): string {
  const currentLang = getLanguageFromRoute(currentRoute)
  const currentPath = currentRoute.fullPath
  
  // 替换路径中的语言部分
  return currentPath.replace(`/${currentLang}`, `/${newLang}`)
}

// 人格类型路由生成器
export function createTypeRoute(code: string, gender: 'm' | 'f', lang: Language = 'zh'): string {
  return createLocalizedPath(`/types/${code.toLowerCase()}-${gender}`, lang)
}

// 配对路由生成器
export function createPairingRoute(maleCode: string, femaleCode: string, lang: Language = 'zh'): string {
  return createLocalizedPath(`/pairings/${maleCode.toLowerCase()}-${femaleCode.toLowerCase()}`, lang)
}

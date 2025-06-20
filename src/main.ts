import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// 导入路由
import router from './router'

// 导入状态管理
import { pinia } from './stores'

// 导入国际化
import i18n from './i18n'

// 创建应用实例
const app = createApp(App)

// 使用插件
app.use(router)
app.use(pinia)
app.use(i18n)

// 挂载应用
app.mount('#app')

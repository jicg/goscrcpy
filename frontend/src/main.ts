import {createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// @ts-ignore
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import App from './App.vue'

const app = createApp(App)
app.use(ElementPlus, {locale: zhCn, size: 'small', zIndex: 3000}).mount("#app")
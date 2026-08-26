// native-shim 必须第一个 import:它给 APP 环境打的 fetch/base 补丁,
// 要先于 App.vue 组件树里 data.js 等模块的模块级请求执行
import './native-shim.js'
import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import { isNativeApp } from './native-shim.js'

// 滚动淡入指令:元素进入视口时添加 .revealed
const reveal = {
  mounted(el) {
    el.classList.add('reveal')
    const io = new IntersectionObserver(
      (entries) => {
        entries.forEach((e) => {
          if (e.isIntersecting) {
            el.classList.add('revealed')
            io.unobserve(el)
          }
        })
      },
      { threshold: 0.15 }
    )
    io.observe(el)
  },
}

createApp(App).directive('reveal', reveal).mount('#app')

// PWA Service Worker 注册(APP 内不需要,页面已离线打包)
if (!isNativeApp && 'serviceWorker' in navigator && location.protocol === 'https:') {
  navigator.serviceWorker.register('/sw.js').catch(() => {})
}

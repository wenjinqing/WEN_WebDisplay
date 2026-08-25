import { createApp } from 'vue'
import App from './App.vue'
import './style.css'

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

// PWA Service Worker 注册
if ('serviceWorker' in navigator && location.protocol === 'https:') {
  navigator.serviceWorker.register('/sw.js').catch(() => {})
}

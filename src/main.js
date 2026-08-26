import { createApp } from 'vue'
import { Capacitor } from '@capacitor/core'
import './style.css'

// ============================================================
// APP(Capacitor WebView)内运行时的适配:
// 页面本体打包在 APP 里(https://localhost),但 API / 图片 / 下载文件
// 都在线上服务器,需要把相对路径指到线上
// ============================================================
const isNativeApp = Capacitor.isNativePlatform()
if (isNativeApp) {
  document.body.classList.add('is-app')
  const SITE_ORIGIN = 'https://alicefans.cn'
  // 1. fetch 的相对路径全部指向线上 API
  const origFetch = window.fetch.bind(window)
  window.fetch = (input, init) => {
    if (typeof input === 'string' && input.startsWith('/')) input = SITE_ORIGIN + input
    return origFetch(input, init)
  }
  // 2. DOM 里的绝对路径(<img src="/gallery/...">、<a href="/downloads/...">)解析到线上
  //    module 脚本在 DOM 解析完后执行,此时插 base 先于 Vue 挂载,资源都会按 base 解析
  const base = document.createElement('base')
  base.href = SITE_ORIGIN + '/'
  document.head.prepend(base)
  // 3. 页内锚点(导航 #home 等)不能被 base 带去线上,改为本地平滑滚动
  document.addEventListener('click', (e) => {
    const a = e.target.closest('a[href^="#"]')
    if (!a) return
    const el = document.querySelector(a.getAttribute('href'))
    if (el) {
      e.preventDefault()
      el.scrollIntoView({ behavior: 'smooth' })
    }
  })
}

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

// 动态引入根组件:保证上面的 fetch/base 补丁先于各组件模块(如 data.js 的
// 模块级 fetch)执行,否则 APP 里内容接口会打到 https://localhost 拿不到数据
import('./App.vue').then(({ default: App }) => {
  createApp(App).directive('reveal', reveal).mount('#app')
})

// PWA Service Worker 注册(APP 内不需要,页面已离线打包)
if (!isNativeApp && 'serviceWorker' in navigator && location.protocol === 'https:') {
  navigator.serviceWorker.register('/sw.js').catch(() => {})
}

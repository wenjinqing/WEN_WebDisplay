// ============================================================
// APP(Capacitor WebView)运行补丁 —— 必须在所有组件模块之前执行!
// 页面本体打包在 APP 里(https://localhost),API / 图片 / 下载文件都在线上,
// 需要把相对路径指到线上;且 data.js 等模块在 import 时就会发请求,
// 所以这个文件要在 main.js 里第一个 import
// ============================================================
import { Capacitor } from '@capacitor/core'

export const isNativeApp = Capacitor.isNativePlatform()

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

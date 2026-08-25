// 小涩猫咖啡厅 Service Worker —— 静态资源离线缓存
const CACHE = 'catcafe-v2'
const ASSETS = ['/', '/index.html', '/favicon.svg', '/manifest.webmanifest']

self.addEventListener('install', (e) => {
  e.waitUntil(caches.open(CACHE).then((c) => c.addAll(ASSETS)).then(() => self.skipWaiting()))
})

self.addEventListener('activate', (e) => {
  e.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k)))
    ).then(() => self.clients.claim())
  )
})

self.addEventListener('fetch', (e) => {
  const url = new URL(e.request.url)
  // API 和下载文件不缓存,直接走网络
  if (url.pathname.startsWith('/api/') || url.pathname.startsWith('/downloads/')) return
  e.respondWith(
    caches.match(e.request).then((hit) => hit || fetch(e.request).then((res) => {
      // 缓存静态资源
      if (e.request.method === 'GET' && res.ok) {
        const clone = res.clone()
        caches.open(CACHE).then((c) => c.put(e.request, clone))
      }
      return res
    }).catch(() => caches.match('/index.html')))
  )
})

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'

const open = ref(false)
const links = [
  { href: '#home', label: '首页' },
  { href: '#notice', label: '公告板' },
  { href: '#novels', label: '小说菜单' },
  { href: '#gallery', label: '插画墙' },
  { href: '#interact', label: '互动区' },
  { href: '#pigmi', label: '猪咪聚集地' },
]

// 彩蛋:快速点 logo 会触发猫爪雨(由 PawRain 组件监听)
function onLogoClick() {
  window.dispatchEvent(new CustomEvent('catcafe-logo-click'))
}

// PWA 下载应用:可安装时直接弹系统安装框,
// 否则(iOS / 已安装过 / 浏览器不支持)给出手动安装指引
const installed = ref(false)
const tip = ref('')
let deferredPrompt = null
let tipTimer = null

function onBeforeInstall(e) {
  e.preventDefault()
  deferredPrompt = e
}

function onAppInstalled() {
  installed.value = true
  tip.value = ''
}

async function onInstallClick() {
  open.value = false
  // 安卓浏览器:直接下载官方 APK 安装包
  if (/android/i.test(navigator.userAgent)) {
    tip.value = '开始下载 APK 安装包,下好后点开安装(需允许「未知来源」安装)喵~'
    clearTimeout(tipTimer)
    tipTimer = setTimeout(() => (tip.value = ''), 8000)
    // 用页面跳转方式触发下载(响应头是 attachment,页面不会真的跳走),
    // 比程序创建 <a> 点击在国产浏览器里兼容性更好
    location.href = '/downloads/app/catcafe.apk'
    return
  }
  if (deferredPrompt) {
    deferredPrompt.prompt()
    await deferredPrompt.userChoice.catch(() => null)
    deferredPrompt = null
    return
  }
  const isIOS = /iphone|ipad|ipod/i.test(navigator.userAgent)
  tip.value = isIOS
    ? '用 Safari 打开本站 → 点底部「分享」→「添加到主屏幕」就下好了喵~'
    : '点浏览器菜单(右上角 ⋮ 或底部 ≡)→「安装应用」/「添加到主屏幕」即可下载喵~'
  clearTimeout(tipTimer)
  tipTimer = setTimeout(() => (tip.value = ''), 6000)
}

onMounted(() => {
  if (window.matchMedia('(display-mode: standalone)').matches || navigator.standalone) {
    installed.value = true
  }
  window.addEventListener('beforeinstallprompt', onBeforeInstall)
  window.addEventListener('appinstalled', onAppInstalled)
})

onUnmounted(() => {
  window.removeEventListener('beforeinstallprompt', onBeforeInstall)
  window.removeEventListener('appinstalled', onAppInstalled)
  clearTimeout(tipTimer)
})
</script>

<template>
  <header class="nav">
    <div class="container nav-inner">
      <a href="#home" class="logo font-cute" @click="onLogoClick">
        <PawPrint :size="26" color="#f9718f" />
        小涩猫咖啡厅
      </a>
      <button class="burger" :aria-expanded="open" aria-label="菜单" @click="open = !open">
        <span /><span /><span />
      </button>
      <nav :class="['links', { open }]">
        <a v-for="l in links" :key="l.href" :href="l.href" @click="open = false">{{ l.label }}</a>
        <button v-if="!installed" class="install-link" @click="onInstallClick">下载应用</button>
        <a href="/go/pixiv.html" target="_blank" rel="noopener" class="pixiv-link">作者P站 ↗</a>
      </nav>
    </div>
    <div v-if="tip" class="install-tip">{{ tip }}</div>
  </header>
</template>

<style scoped>
.nav {
  position: sticky;
  top: 0;
  z-index: 50;
  background: rgba(255, 246, 248, 0.85);
  backdrop-filter: blur(10px);
  border-bottom: 2px dashed var(--pink-soft);
}

.nav-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.25rem;
  color: var(--ink);
  text-decoration: none;
}

.links {
  display: flex;
  align-items: center;
  gap: 24px;
}

.links a {
  color: var(--ink);
  text-decoration: none;
  font-size: 0.95rem;
  transition: color 0.2s;
}

.links a:hover {
  color: var(--pink-deep);
}

.pixiv-link {
  padding: 6px 16px;
  border-radius: 999px;
  background: var(--pink-pale);
  color: var(--pink-deep) !important;
  font-weight: 500;
}

.install-link {
  padding: 6px 16px;
  border: none;
  border-radius: 999px;
  background: var(--pink-deep);
  color: #fff;
  font-size: 0.95rem;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: opacity 0.2s;
}

.install-link:hover {
  opacity: 0.85;
}

.install-tip {
  position: fixed;
  left: 50%;
  bottom: 32px;
  transform: translateX(-50%);
  max-width: min(90vw, 420px);
  padding: 12px 20px;
  border-radius: 14px;
  background: var(--ink);
  color: var(--card);
  font-size: 0.9rem;
  line-height: 1.6;
  text-align: center;
  z-index: 60;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.18);
}

.burger {
  display: none;
  flex-direction: column;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
}

.burger span {
  width: 22px;
  height: 2.5px;
  border-radius: 2px;
  background: var(--ink);
}

@media (max-width: 720px) {
  .burger {
    display: flex;
  }
  .links {
    position: absolute;
    top: 64px;
    left: 0;
    right: 0;
    flex-direction: column;
    align-items: flex-start;
    gap: 0;
    background: var(--card);
    border-bottom: 2px dashed var(--pink-soft);
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease;
  }
  .links.open {
    max-height: 380px;
  }
  .links a {
    width: 100%;
    padding: 14px 24px;
  }
  .install-link {
    margin: 8px 24px 4px;
    padding: 10px 16px;
    align-self: flex-start;
  }
  .pixiv-link {
    margin: 8px 24px 16px;
    width: auto !important;
  }
}
</style>

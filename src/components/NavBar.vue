<script setup>
import { ref } from 'vue'
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
        <a href="/go/pixiv.html" target="_blank" rel="noopener" class="pixiv-link">作者P站 ↗</a>
      </nav>
    </div>
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
    max-height: 320px;
  }
  .links a {
    width: 100%;
    padding: 14px 24px;
  }
  .pixiv-link {
    margin: 8px 24px 16px;
    width: auto !important;
  }
}
</style>

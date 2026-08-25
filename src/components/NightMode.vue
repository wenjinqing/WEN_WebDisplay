<script setup>
// 深夜暗色模式:晚8点~早7点自动开启,也可手动切换
import { onMounted, onUnmounted } from 'vue'

function applyTheme() {
  const manual = localStorage.getItem('catcafe_theme')
  let night
  if (manual === 'night') night = true
  else if (manual === 'day') night = false
  else {
    const h = new Date().getHours()
    night = h >= 20 || h < 7 // 自动:20:00-07:00
  }
  document.body.classList.toggle('night', night)
  window.dispatchEvent(new CustomEvent('menu-state', { detail: { theme: night ? 'night' : 'day' } }))
}

function onMenuTheme() {
  const manual = localStorage.getItem('catcafe_theme')
  if (manual === 'night') localStorage.setItem('catcafe_theme', 'day')
  else localStorage.setItem('catcafe_theme', 'night')
  applyTheme()
}

let timer = null

onMounted(() => {
  applyTheme()
  window.addEventListener('menu-theme', onMenuTheme)
  timer = setInterval(applyTheme, 60000) // 每分钟检查自动切换
})

onUnmounted(() => {
  window.removeEventListener('menu-theme', onMenuTheme)
  clearInterval(timer)
})
</script>

<template>
  <!-- 无 UI,仅逻辑 -->
</template>

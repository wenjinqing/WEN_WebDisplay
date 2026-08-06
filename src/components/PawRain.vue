<script setup>
// 隐藏彩蛋:猫爪雨。快速点 5 次导航 logo 触发
import { ref, onMounted, onUnmounted } from 'vue'
import PawPrint from './PawPrint.vue'

const raining = ref(false)
const drops = ref([])
let clicks = []
let seq = 0

function onLogoClick() {
  const now = Date.now()
  clicks = clicks.filter((t) => now - t < 1500)
  clicks.push(now)
  if (clicks.length >= 5 && !raining.value) {
    clicks = []
    rain()
  }
}

function rain() {
  raining.value = true
  drops.value = Array.from({ length: 40 }, () => ({
    id: ++seq,
    left: Math.random() * 100,
    delay: Math.random() * 1.5,
    duration: 2.5 + Math.random() * 2,
    size: 18 + Math.random() * 26,
    rotate: Math.random() * 360,
  }))
  setTimeout(() => {
    raining.value = false
    drops.value = []
  }, 6000)
}

onMounted(() => window.addEventListener('catcafe-logo-click', onLogoClick))
onUnmounted(() => window.removeEventListener('catcafe-logo-click', onLogoClick))
</script>

<template>
  <div v-if="raining" class="paw-rain" aria-hidden="true">
    <PawPrint
      v-for="d in drops"
      :key="d.id"
      class="drop"
      :size="d.size"
      color="#f9a8bf"
      :style="{
        left: d.left + '%',
        animationDelay: d.delay + 's',
        animationDuration: d.duration + 's',
        transform: `rotate(${d.rotate}deg)`,
      }"
    />
    <p class="rain-text font-cute">喵呜~猫爪雨!</p>
  </div>
</template>

<style scoped>
.paw-rain {
  position: fixed;
  inset: 0;
  z-index: 200;
  pointer-events: none;
  overflow: hidden;
}

.drop {
  position: absolute;
  top: -8%;
  animation: fall linear forwards;
}

@keyframes fall {
  to {
    top: 110%;
  }
}

.rain-text {
  position: absolute;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 2rem;
  color: var(--pink-deep);
  background: rgba(255, 255, 255, 0.9);
  border: 3px dashed var(--pink-soft);
  border-radius: 24px;
  padding: 16px 36px;
}
</style>

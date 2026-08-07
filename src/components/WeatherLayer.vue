<script setup>
// 猫咖天气:晴 / 下雨 / 飘花,10分钟自动轮换,也可手动切
import { ref, computed, onMounted, onUnmounted } from 'vue'

const ORDER = ['sunny', 'rain', 'petals']
const ICONS = { sunny: '☀️', rain: '🌧️', petals: '🌸' }
const weather = ref('sunny')
let cycleTimer = null

const drops = computed(() => {
  if (weather.value === 'rain') {
    return Array.from({ length: 36 }, (_, i) => ({
      id: i,
      left: Math.random() * 100,
      delay: Math.random() * 1.2,
      duration: 0.7 + Math.random() * 0.5,
    }))
  }
  if (weather.value === 'petals') {
    return Array.from({ length: 16 }, (_, i) => ({
      id: i,
      left: Math.random() * 100,
      delay: Math.random() * 6,
      duration: 6 + Math.random() * 5,
      sway: 20 + Math.random() * 40,
    }))
  }
  return []
})

function setWeather(w) {
  weather.value = w
  window.dispatchEvent(new CustomEvent('cafe-weather', { detail: w }))
}

function cycle() {
  setWeather(ORDER[(ORDER.indexOf(weather.value) + 1) % ORDER.length])
}

onMounted(() => {
  // 随机开局,10分钟一换
  setWeather(ORDER[Math.floor(Math.random() * ORDER.length)])
  cycleTimer = setInterval(cycle, 10 * 60 * 1000)
})

onUnmounted(() => clearInterval(cycleTimer))
</script>

<template>
  <button class="weather-btn" :title="'当前:' + weather" @click="cycle">
    {{ ICONS[weather] }} 天气
  </button>

  <div v-if="weather !== 'sunny'" class="weather-layer" aria-hidden="true">
    <template v-if="weather === 'rain'">
      <i
        v-for="d in drops"
        :key="d.id"
        class="rain-drop"
        :style="{ left: d.left + '%', animationDelay: d.delay + 's', animationDuration: d.duration + 's' }"
      />
    </template>
    <template v-else>
      <i
        v-for="d in drops"
        :key="d.id"
        class="petal"
        :style="{
          left: d.left + '%',
          animationDelay: d.delay + 's',
          animationDuration: d.duration + 's',
          '--sway': d.sway + 'px',
        }"
      />
    </template>
  </div>
</template>

<style scoped>
.weather-btn {
  position: fixed;
  left: 20px;
  bottom: 214px;
  z-index: 60;
  border: 2px solid var(--pink-soft);
  background: rgba(255, 255, 255, 0.92);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 8px 16px;
  font-size: 0.85rem;
  cursor: pointer;
  box-shadow: var(--shadow);
}

.weather-btn:hover {
  background: var(--pink-pale);
}

.weather-layer {
  position: fixed;
  inset: 0;
  z-index: 38;
  pointer-events: none;
  overflow: hidden;
}

.rain-drop {
  position: absolute;
  top: -20px;
  width: 2px;
  height: 16px;
  background: linear-gradient(to bottom, transparent, rgba(140, 160, 190, 0.55));
  border-radius: 2px;
  animation: rainFall linear infinite;
}

@keyframes rainFall {
  to {
    transform: translateY(110vh);
  }
}

.petal {
  position: absolute;
  top: -20px;
  width: 12px;
  height: 9px;
  background: #ffc9d8;
  border-radius: 60% 40% 60% 40%;
  animation: petalFall linear infinite;
}

@keyframes petalFall {
  0% { transform: translate(0, 0) rotate(0deg); }
  25% { transform: translate(var(--sway), 28vh) rotate(120deg); }
  50% { transform: translate(calc(var(--sway) * -0.6), 55vh) rotate(220deg); }
  75% { transform: translate(var(--sway), 82vh) rotate(320deg); }
  100% { transform: translate(0, 110vh) rotate(420deg); }
}

@media (prefers-reduced-motion: reduce) {
  .weather-layer {
    display: none;
  }
}
</style>

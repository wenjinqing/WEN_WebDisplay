<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 弹幕:留言板内容飘过屏幕。reduced-motion 用户默认关闭
const enabled = ref(
  localStorage.getItem('catcafe_danmaku') ??
    (window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'off' : 'on')
) === 'on'

const lanes = ref([]) // {id, text, top, duration, delay}
let timer = null
let seq = 0
let pool = []

function spawn() {
  if (!pool.length || !enabled.value) return
  const text = pool[Math.floor(Math.random() * pool.length)]
  const id = ++seq
  lanes.value.push({
    id,
    text,
    top: 12 + Math.random() * 55, // 屏幕上方 12%~67%
    duration: 10 + Math.random() * 8,
  })
  setTimeout(() => (lanes.value = lanes.value.filter((l) => l.id !== id)), 20000)
}

function toggle() {
  enabled.value = !enabled.value
  localStorage.setItem('catcafe_danmaku', enabled.value ? 'on' : 'off')
  if (!enabled.value) lanes.value = []
}

onMounted(async () => {
  try {
    const res = await fetch('/api/messages')
    const msgs = await res.json()
    pool = msgs.map((m) => `${m.nick}:${m.content}`).filter((t) => t.length <= 40)
  } catch {
    /* 静默 */
  }
  timer = setInterval(spawn, 2500)
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div v-if="enabled" class="danmaku-layer" aria-hidden="true">
    <span
      v-for="l in lanes"
      :key="l.id"
      class="dm"
      :style="{ top: l.top + '%', animationDuration: l.duration + 's' }"
    >
      {{ l.text }}
    </span>
  </div>
  <button class="dm-toggle" :aria-pressed="enabled" @click="toggle">
    {{ enabled ? '💬 弹幕开' : '💭 弹幕关' }}
  </button>
</template>

<style scoped>
.danmaku-layer {
  position: fixed;
  inset: 0;
  z-index: 40;
  pointer-events: none;
  overflow: hidden;
}

.dm {
  position: absolute;
  left: 100%;
  white-space: nowrap;
  font-size: 0.9rem;
  color: var(--pink-deep);
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid var(--pink-pale);
  border-radius: 999px;
  padding: 3px 14px;
  animation: fly linear forwards;
  box-shadow: 0 2px 8px rgba(233, 93, 127, 0.12);
}

@keyframes fly {
  to {
    transform: translateX(calc(-100vw - 100%));
  }
}

.dm-toggle {
  position: fixed;
  left: 20px;
  bottom: 24px;
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

.dm-toggle:hover {
  background: var(--pink-pale);
}
</style>

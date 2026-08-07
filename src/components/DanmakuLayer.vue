<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 弹幕:手绘云朵风。留言板内容飘过屏幕;没留言时飘店内氛围语
const enabled = ref(
  (localStorage.getItem('catcafe_danmaku') ??
    (window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'off' : 'on')) === 'on'
)

const lanes = ref([]) // {id, text, top, duration}
let timer = null
let seq = 0
let pool = []

// 实在没留言时的氛围弹幕
const FALLBACK = [
  '欢迎光临小涩猫咖啡厅~', '今天也在等更新', '猪咪路过~',
  '猫猫酱加油!', '前排兜售小鱼干', '这家店好温暖', '蹲一个新坑',
]

function spawn() {
  if (!enabled.value || !pool.length) return
  const text = pool[Math.floor(Math.random() * pool.length)]
  const id = ++seq
  lanes.value.push({
    id,
    text,
    top: 10 + Math.random() * 55,
    duration: 11 + Math.random() * 8,
  })
  setTimeout(() => (lanes.value = lanes.value.filter((l) => l.id !== id)), 22000)
}

function toggle() {
  enabled.value = !enabled.value
  localStorage.setItem('catcafe_danmaku', enabled.value ? 'on' : 'off')
  if (!enabled.value) {
    lanes.value = []
  } else {
    spawn()
    spawn() // 打开时立刻来两条,不用等
  }
}

onMounted(async () => {
  try {
    const res = await fetch('/api/messages')
    const msgs = await res.json()
    pool = msgs.map((m) => `${m.nick}:${m.content}`).filter((t) => t.length <= 60)
  } catch {
    /* 静默 */
  }
  if (!pool.length) pool = FALLBACK
  timer = setInterval(spawn, 2500)
  if (enabled.value) spawn() // 首屏立刻来一条
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

/* 手绘云朵弹幕:SVG 云朵拉伸做底,字在云中 */
.dm {
  position: absolute;
  left: 100%;
  white-space: nowrap;
  font-size: 0.88rem;
  color: var(--pink-deep);
  padding: 10px 26px 14px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 80 50' preserveAspectRatio='none'%3E%3Cpath d='M20 42 Q 8 42 8 33 Q 8 25 18 24 Q 20 13 32 15 Q 40 6 50 15 Q 61 13 62 24 Q 72 25 71 34 Q 70 42 58 42 Z' fill='%23ffffff' fill-opacity='0.94' stroke='%23f4a9c0' stroke-width='3' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-size: 100% 100%;
  animation: fly linear forwards;
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

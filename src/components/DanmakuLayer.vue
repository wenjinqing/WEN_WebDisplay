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

/* 手绘云朵弹幕:胖云朵铺满底框,文字稳坐云中 */
.dm {
  position: absolute;
  left: 100%;
  white-space: nowrap;
  font-size: 0.88rem;
  line-height: 1.2;
  color: var(--pink-deep);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 14px 34px 16px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 60' preserveAspectRatio='none'%3E%3Cpath vector-effect='non-scaling-stroke' d='M18 52 Q 6 52 6 43 Q 6 35 15 34 Q 15 23 27 22 Q 30 10 44 12 Q 50 4 62 8 Q 74 5 78 15 Q 90 15 90 25 Q 97 28 96 37 Q 95 47 85 47 Q 83 54 72 53 L 22 53 Q 18 53 18 52 Z' fill='%23ffffff' fill-opacity='0.95' stroke='%23f4a9c0' stroke-width='3' stroke-linejoin='round'/%3E%3C/svg%3E");
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

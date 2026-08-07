<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 弹幕:云朵以文字为主体(CSS 胶囊+鼓包,永远包住文字),分航道避免重叠
const enabled = ref(
  (localStorage.getItem('catcafe_danmaku') ??
    (window.matchMedia('(prefers-reduced-motion: reduce)').matches ? 'off' : 'on')) === 'on'
)

const LANES = 6 // 航道数
const lanes = ref([]) // {id, text, lane, duration}
const laneBusy = Array(LANES).fill(0) // 每条航道的解禁时间戳
let timer = null
let seq = 0
let pool = []

const FALLBACK = [
  '欢迎光临小涩猫咖啡厅~', '今天也在等更新', '猪咪路过~',
  '猫猫酱加油!', '前排兜售小鱼干', '这家店好温暖', '蹲一个新坑',
]

function spawn() {
  if (!enabled.value || !pool.length) return
  const text = pool[Math.floor(Math.random() * pool.length)]

  // 字数 → 云大小与速度:字越多云越肥、飘得越慢
  const len = Math.min(text.length, 60)
  const duration = 12 + len * 0.35 // 长文慢慢飘
  const estW = Math.min(len, 20) * 15 + 90 // 折行后宽度有上限
  const vw = window.innerWidth
  const entryMs = (estW / (vw + estW)) * duration * 1000 + 1200

  // 找一条空闲航道;都忙就这波不发了(避免叠云)
  const now = Date.now()
  const free = []
  for (let i = 0; i < LANES; i++) {
    if (laneBusy[i] <= now) free.push(i)
  }
  if (!free.length) return
  const lane = free[Math.floor(Math.random() * free.length)]
  laneBusy[lane] = now + entryMs

  const id = ++seq
  lanes.value.push({ id, text, lane, duration })
  setTimeout(() => (lanes.value = lanes.value.filter((l) => l.id !== id)), (duration + 2) * 1000)
}

function toggle() {
  enabled.value = !enabled.value
  localStorage.setItem('catcafe_danmaku', enabled.value ? 'on' : 'off')
  if (!enabled.value) {
    lanes.value = []
  } else {
    spawn()
    spawn()
  }
  window.dispatchEvent(new CustomEvent('menu-state', { detail: { danmaku: enabled.value } }))
}

function onMenuToggle() {
  toggle()
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
  if (enabled.value) spawn()
  window.addEventListener('menu-danmaku', onMenuToggle)
  window.dispatchEvent(new CustomEvent('menu-state', { detail: { danmaku: enabled.value } }))
})

onUnmounted(() => {
  clearInterval(timer)
  window.removeEventListener('menu-danmaku', onMenuToggle)
})
</script>

<template>
  <div v-if="enabled" class="danmaku-layer" aria-hidden="true">
    <span
      v-for="l in lanes"
      :key="l.id"
      class="dm"
      :style="{ top: 8 + l.lane * 9 + '%', animationDuration: l.duration + 's' }"
    >
      {{ l.text }}
    </span>
  </div>
</template>

<style scoped>
.danmaku-layer {
  position: fixed;
  inset: 0;
  z-index: 40;
  pointer-events: none;
  overflow: hidden;
}

/* 手绘云朵弹幕:云朵随内容伸缩,长文折行云变肥 */
.dm {
  position: absolute;
  left: 100%;
  width: max-content; /* 先按文字自然宽度,超上限才折行 */
  max-width: 230px;
  white-space: normal;
  word-break: break-all;
  text-align: center;
  font-size: 0.88rem;
  line-height: 1.45;
  color: var(--pink-deep);
  padding: 16px 36px 18px; /* 文字锁定在云朵中段安全区 */
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 80 50' preserveAspectRatio='none'%3E%3Cpath d='M20 42 Q 8 42 8 33 Q 8 25 18 24 Q 20 13 32 15 Q 40 6 50 15 Q 61 13 62 24 Q 72 25 71 34 Q 70 42 58 42 Z' fill='%23ffffff' fill-opacity='0.94' stroke='%23f4a9c0' stroke-width='2.5' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-size: 100% 100%;
  animation: fly linear forwards;
}

@keyframes fly {
  to {
    transform: translateX(calc(-100vw - 100%));
  }
}
</style>

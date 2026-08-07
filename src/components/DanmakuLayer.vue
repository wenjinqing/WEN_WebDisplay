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

  // 估算云朵宽度(px):字数×字宽 + 内边距
  const estW = Math.min(text.length, 60) * 15 + 80
  const vw = window.innerWidth
  const duration = 11 + Math.random() * 8
  // 云朵完全进入屏幕所需时间 = 占比 × 总时长;大云朵占航道更久
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
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <div v-if="enabled" class="danmaku-layer" aria-hidden="true">
    <span
      v-for="l in lanes"
      :key="l.id"
      class="dm"
      :style="{ top: 8 + l.lane * 9 + '%', animationDuration: l.duration + 's' }"
    >
      <i class="bump b1" /><i class="bump b2" /><i class="bump b3" />
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

/* 云朵 = 文字胶囊主体 + 顶部鼓包,文字多长云就多长 */
.dm {
  position: absolute;
  left: 100%;
  white-space: nowrap;
  font-size: 0.88rem;
  line-height: 1.2;
  color: var(--pink-deep);
  background: rgba(255, 255, 255, 0.96);
  border: 2.5px solid #f4a9c0;
  border-radius: 999px;
  padding: 9px 24px;
  animation: fly linear forwards;
}

/* 三个鼓包躲在胶囊后面,只露出上半圆 = 云朵轮廓 */
.bump {
  position: absolute;
  background: inherit;
  border: 2.5px solid #f4a9c0;
  border-radius: 50%;
  z-index: -1;
}

.b1 { width: 30px; height: 30px; top: -14px; left: 18%; }
.b2 { width: 38px; height: 38px; top: -19px; left: 42%; }
.b3 { width: 26px; height: 26px; top: -12px; left: 68%; }

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

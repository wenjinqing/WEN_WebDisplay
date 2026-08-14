<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// 弹幕:量体裁衣的手绘云朵 —— 先量文字尺寸,再按尺寸生成云朵路径,永远包住文字
const isMobile = window.innerWidth < 720
const enabled = ref(
  (localStorage.getItem('catcafe_danmaku') ??
    (window.matchMedia('(prefers-reduced-motion: reduce)').matches || isMobile ? 'off' : 'on')) === 'on'
)

const LANES = isMobile ? 3 : 6
const LANE_TOP = isMobile ? 5 : 8
const LANE_GAP = isMobile ? 8 : 9
const lanes = ref([])
const laneBusy = Array(LANES).fill(0)
let timer = null
let seq = 0
let pool = []

const FALLBACK = [
  '欢迎光临小涩猫咖啡厅~', '今天也在等更新', '猪咪路过~',
  '猫猫酱加油!', '前排兜售小鱼干', '这家店好温暖', '蹲一个新坑',
]

// 文字测量器(隐藏元素,样式与弹幕文字一致)
let measurer = null
function measure(text) {
  if (!measurer) {
    measurer = document.createElement('div')
    measurer.style.cssText =
      'position:fixed;left:-9999px;top:0;visibility:hidden;font-size:14px;line-height:1.45;' +
      'max-width:230px;width:max-content;white-space:normal;word-break:break-all;text-align:center;' +
      'font-family:"Noto Sans SC",sans-serif;padding:0;'
    document.body.appendChild(measurer)
  }
  measurer.textContent = text
  return { w: measurer.offsetWidth, h: measurer.offsetHeight }
}

// 按文字尺寸生成"猫咪云":蓬松多层鼓包 + 猫耳朵,文字严格居中
function cloudPath(tw, th, seed) {
  const padX = 34, padY = 24
  const w = tw + padX * 2, h = th + padY * 2
  const TOP = 46
  const rnd = (i) => {
    const x = Math.sin(seed * 97 + i * 131) * 10000
    return x - Math.floor(x)
  }
  const cr = 14 // 圆角

  // 单条轮廓:sweep=1 = 向行进方向左侧鼓 = 顺时针描边时全朝外
  let d = `M ${cr} 0`
  const n = Math.max(3, Math.min(5, Math.round(w / 58)))
  const seg = (w - cr * 2) / n
  for (let i = 0; i < n; i++) {
    const middle = 1 - Math.abs(i - (n - 1) / 2) / ((n - 1) / 2 + 0.01)
    const br = 16 + middle * 12 + rnd(i) * 5
    d += ` A ${br} ${br} 0 0 1 ${cr + seg * (i + 1)} 0`
  }
  d += ` A ${cr} ${cr} 0 0 1 ${w} ${cr}`
  const sr = 12 + rnd(40) * 5
  d += ` A ${sr} ${sr} 0 0 1 ${w} ${h / 2}`
  d += ` L ${w} ${h - cr}`
  d += ` A ${cr} ${cr} 0 0 1 ${w - cr} ${h}`
  const b1 = 14 + rnd(41) * 5
  d += ` A ${b1} ${b1} 0 0 1 ${w * 0.62} ${h}`
  const b2 = 12 + rnd(42) * 5
  d += ` A ${b2} ${b2} 0 0 1 ${w * 0.32} ${h}`
  d += ` L ${cr} ${h}`
  d += ` A ${cr} ${cr} 0 0 1 0 ${h - cr}`
  const sr2 = 11 + rnd(43) * 5
  d += ` A ${sr2} ${sr2} 0 0 1 0 ${h / 2 - 10}`
  d += ` L 0 ${cr}`
  d += ` A ${cr} ${cr} 0 0 1 ${cr} 0 Z`

  // 猫耳朵(圆尖顶,骑在顶部鼓包上)
  const ears = []
  for (const [i, fx] of [[0, 0.28], [1, 0.7]]) {
    const cx = w * (fx + (rnd(i + 30) - 0.5) * 0.05)
    const r = 15 + rnd(i + 31) * 5
    const baseY = 2
    const tipY = baseY - r * 1.5
    ears.push(
      `M ${cx - r} ${baseY} Q ${cx - r * 0.2} ${tipY} ${cx} ${tipY} Q ${cx + r * 0.2} ${tipY} ${cx + r} ${baseY} Z`
    )
  }

  return {
    w: w + 24,
    h: h + TOP,
    vb: `-12 ${-TOP + 6} ${w + 24} ${h + TOP}`,
    offX: -12,
    offY: -TOP + 6,
    d,
    ears,
    padX,
    padY,
  }
}

function spawn() {
  if (!enabled.value || !pool.length) return
  const text = pool[Math.floor(Math.random() * pool.length)]
  const len = Math.min(text.length, 60)
  const duration = 12 + len * 0.35

  const { w } = measure(text)
  const vw = window.innerWidth
  const entryMs = (w / (vw + w)) * duration * 1000 + 1200

  const now = Date.now()
  const free = []
  for (let i = 0; i < LANES; i++) {
    if (laneBusy[i] <= now) free.push(i)
  }
  if (!free.length) return
  const lane = free[Math.floor(Math.random() * free.length)]
  laneBusy[lane] = now + entryMs

  const id = ++seq
  const size = measure(text)
  lanes.value.push({ id, text, lane, duration, cloud: cloudPath(size.w, size.h, id) })
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
  if (measurer) measurer.remove()
})
</script>

<template>
  <div v-if="enabled" class="danmaku-layer" aria-hidden="true">
    <span
      v-for="l in lanes"
      :key="l.id"
      class="dm"
      :style="{ top: LANE_TOP + l.lane * LANE_GAP + '%', animationDuration: l.duration + 's' }"
    >
      <svg
        class="dm-cloud"
        :width="l.cloud.w"
        :height="l.cloud.h"
        :viewBox="l.cloud.vb"
        :style="{ left: l.cloud.offX + 'px', top: l.cloud.offY + 'px' }"
        aria-hidden="true"
      >
        <!-- 单条轮廓的云朵主体 -->
        <path
          :d="l.cloud.d"
          fill="rgba(255,255,255,0.96)" stroke="#f4a9c0" stroke-width="2.5" stroke-linejoin="round"
        />
        <!-- 猫耳朵 -->
        <path
          v-for="(d, i) in l.cloud.ears"
          :key="'e' + i"
          :d="d"
          fill="rgba(255,255,255,0.96)" stroke="#f4a9c0" stroke-width="2.5" stroke-linejoin="round"
        />
      </svg>
      <span class="dm-text" :style="{ padding: l.cloud.padY + 'px ' + l.cloud.padX + 'px' }">{{ l.text }}</span>
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

.dm {
  position: absolute;
  left: 100%;
  display: inline-block;
  width: max-content; /* 先按文字自然宽度撑开 */
  max-width: 300px; /* 再交给文字折行,云变肥 */
  animation: fly linear forwards;
}

.dm-cloud {
  position: absolute;
  overflow: visible; /* 鼓包可以探出本体 */
}

.dm-text {
  position: relative;
  display: inline-block;
  max-width: 230px; /* 与测量器一致 */
  white-space: normal;
  word-break: break-all;
  text-align: center;
  font-size: 0.88rem;
  line-height: 1.45;
  color: var(--pink-deep);
}

@keyframes fly {
  to {
    transform: translateX(calc(-100vw - 100%));
  }
}

@media (max-width: 720px) {
  .dm-text {
    font-size: 0.78rem;
  }
}
</style>

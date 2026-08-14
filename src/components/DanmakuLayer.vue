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

// 按文字尺寸生成"猫咪云":圆角主体(必包文字)+ 两只猫耳朵,鼓包只是装饰
function cloudPath(tw, th, seed) {
  const padX = 32, padY = 22
  const w = tw + padX * 2, h = th + padY * 2
  const TOP = 34 // 耳朵探出空间
  const rnd = (i) => {
    const x = Math.sin(seed * 97 + i * 131) * 10000
    return x - Math.floor(x)
  }
  // 两只猫耳朵(尖顶圆弧,位置和大小轻微随机)
  const ears = []
  for (const [i, fx] of [[0, 0.26], [1, 0.68]]) {
    const cx = w * (fx + (rnd(i) - 0.5) * 0.06)
    const r = 15 + rnd(i + 5) * 6
    const tipY = TOP - r * 1.55
    ears.push(
      `M ${cx - r} ${TOP} Q ${cx - r * 0.25} ${tipY} ${cx} ${tipY} Q ${cx + r * 0.25} ${tipY} ${cx + r} ${TOP}`
    )
  }
  // 两耳之间一个小圆鼓包
  const dr = 9 + rnd(9) * 5
  const domes = [`M ${w / 2 - dr} ${TOP} A ${dr} ${dr} 0 0 0 ${w / 2 + dr} ${TOP}`]
  return {
    w: w + 20,
    h: h + TOP + 8,
    vb: `-10 ${-TOP + 4} ${w + 20} ${h + TOP + 8}`,
    offX: -10,
    offY: -TOP + 4,
    body: { x: 0, y: TOP, w, h, r: Math.min(24, h / 2) },
    ears,
    domes,
    padX,
    padY: padY + 6,
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
        <!-- 主体圆角矩形(包住文字) -->
        <rect
          :x="l.cloud.body.x" :y="l.cloud.body.y"
          :width="l.cloud.body.w" :height="l.cloud.body.h"
          :rx="l.cloud.body.r"
          fill="rgba(255,255,255,0.95)" stroke="#f4a9c0" stroke-width="2.5"
        />
        <!-- 猫耳朵 -->
        <path
          v-for="(d, i) in l.cloud.ears"
          :key="'e' + i"
          :d="d"
          fill="rgba(255,255,255,0.95)" stroke="#f4a9c0" stroke-width="2.5" stroke-linejoin="round"
        />
        <!-- 耳间小鼓包 -->
        <path
          v-for="(d, i) in l.cloud.domes"
          :key="'d' + i"
          :d="d"
          fill="rgba(255,255,255,0.95)" stroke="#f4a9c0" stroke-width="2.5" stroke-linecap="round"
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

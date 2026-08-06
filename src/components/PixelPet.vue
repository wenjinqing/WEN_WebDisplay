<script setup>
// 爱丽丝看板娘:像素小人在页面底部生活
// 行为:散步 / 休息 / 睡觉 / 追鼠标 / 随机说话 / 点她有反应
import { ref, onMounted, onUnmounted } from 'vue'

const x = ref(15) // 位置(vw)
const dir = ref(1) // 朝向:1 右 -1 左
const state = ref('walk') // walk | idle | sleep
const bubble = ref('')
const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')

const BUBBLES = [
  '赶稿中……', '想喝奶茶', '猪咪们好呀~', '在写新坑!',
  '今天也要加油喵', '催更?在写了在写了', '逛逛自己的店',
  '刚才那句写得不错', '想吃草莓蛋糕', '被发现了?',
]

let tick = null
let stateTimer = null
let bubbleTimer = null
let mouseX = 50
let bubbleClear = null

function step() {
  if (hidden.value || state.value !== 'walk') return
  x.value += dir.value * 0.32
  if (x.value >= 90) { x.value = 90; dir.value = -1 }
  if (x.value <= 1) { x.value = 1; dir.value = 1 }
}

function maybeChange() {
  if (hidden.value) return
  const r = Math.random()
  if (state.value === 'walk') {
    if (r < 0.22) state.value = 'idle'
    else if (r < 0.32) state.value = 'sleep'
    else if (r < 0.6) dir.value = mouseX > x.value ? 1 : -1 // 追鼠标方向
  } else if (r < 0.65) {
    state.value = 'walk'
  }
}

function maybeBubble() {
  if (hidden.value || Math.random() > 0.55) return
  bubble.value = BUBBLES[Math.floor(Math.random() * BUBBLES.length)]
  clearTimeout(bubbleClear)
  bubbleClear = setTimeout(() => (bubble.value = ''), 3200)
}

function poke() {
  bubble.value = ['喵?叫我?', '在呢在呢!', '怎么啦~', '摸我就不用赶稿了?'][Math.floor(Math.random() * 4)]
  clearTimeout(bubbleClear)
  bubbleClear = setTimeout(() => (bubble.value = ''), 2500)
  state.value = 'idle'
  setTimeout(() => (state.value = 'walk'), 2000)
}

function hide() {
  hidden.value = true
  localStorage.setItem('catcafe_pet_hide', '1')
}

function show() {
  hidden.value = false
  localStorage.removeItem('catcafe_pet_hide')
}

onMounted(() => {
  window.addEventListener('mousemove', onMouse, { passive: true })
  tick = setInterval(step, 60)
  stateTimer = setInterval(maybeChange, 4000)
  bubbleTimer = setInterval(maybeBubble, 9000)
})

function onMouse(e) {
  mouseX = (e.clientX / window.innerWidth) * 100
}

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouse)
  clearInterval(tick)
  clearInterval(stateTimer)
  clearInterval(bubbleTimer)
  clearTimeout(bubbleClear)
})
</script>

<template>
  <button v-if="hidden" class="pet-restore" aria-label="召回看板娘" @click="show">🐾</button>

  <div v-else class="pet" :style="{ left: x + 'vw' }" :data-state="state">
    <transition name="bub">
      <span v-if="bubble" class="pet-bubble">{{ bubble }}</span>
    </transition>
    <button class="pet-hide" aria-label="让她去休息" @click="hide">×</button>
    <span v-if="state === 'sleep'" class="zzz">💤</span>
    <span class="flipper" :class="{ flip: dir < 0 }">
      <img
        src="/alice-pixel.png"
        alt="爱丽丝看板娘"
        class="pet-img"
        :class="state"
        draggable="false"
        @click="poke"
      />
    </span>
  </div>
</template>

<style scoped>
.pet {
  position: fixed;
  bottom: 0;
  z-index: 55;
  pointer-events: none;
  transition: left 0.06s linear;
}

.flipper {
  display: inline-block;
}

.flipper.flip {
  transform: scaleX(-1);
}

.pet-img {
  width: 76px;
  image-rendering: pixelated;
  pointer-events: auto;
  cursor: pointer;
  user-select: none;
}

/* 散步:上下小跑 */
.pet-img.walk {
  animation: walkBob 0.45s ease-in-out infinite;
}

@keyframes walkBob {
  0%, 100% { transform: translateY(0) rotate(-2deg); }
  50% { transform: translateY(-5px) rotate(2deg); }
}

/* 休息:呼吸 */
.pet-img.idle {
  animation: breathe 2.4s ease-in-out infinite;
}

@keyframes breathe {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.94); }
}

/* 睡觉:变暗静止 */
.pet-img.sleep {
  filter: brightness(0.85) saturate(0.8);
}

.zzz {
  position: absolute;
  top: -18px;
  right: 4px;
  font-size: 18px;
  animation: zfloat 2s ease-in-out infinite;
}

@keyframes zfloat {
  0%, 100% { opacity: 0.4; transform: translateY(0); }
  50% { opacity: 1; transform: translateY(-6px); }
}

.pet-bubble {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 6px;
  background: #fff;
  border: 2px solid var(--pink-soft);
  color: var(--ink);
  border-radius: 14px 14px 14px 4px;
  padding: 4px 12px;
  font-size: 0.8rem;
  white-space: nowrap;
  box-shadow: var(--shadow);
}

.bub-enter-active { transition: all 0.2s ease; }
.bub-leave-active { transition: all 0.3s ease; }
.bub-enter-from { opacity: 0; transform: translateX(-50%) translateY(6px); }
.bub-leave-to { opacity: 0; }

.pet-hide {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: none;
  background: rgba(91, 58, 71, 0.5);
  color: #fff;
  font-size: 12px;
  cursor: pointer;
  pointer-events: auto;
  opacity: 0;
  transition: opacity 0.2s;
}

.pet:hover .pet-hide {
  opacity: 1;
}

.pet-restore {
  position: fixed;
  left: 20px;
  bottom: 76px;
  z-index: 55;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 2px solid var(--pink-soft);
  background: rgba(255, 255, 255, 0.92);
  font-size: 18px;
  cursor: pointer;
  box-shadow: var(--shadow);
}

@media (max-width: 720px) {
  .pet-img {
    width: 58px;
  }
  .pet-hide {
    opacity: 1; /* 手机没有 hover,常显 */
  }
}

@media (prefers-reduced-motion: reduce) {
  .pet-img.walk, .pet-img.idle, .zzz {
    animation: none;
  }
}
</style>

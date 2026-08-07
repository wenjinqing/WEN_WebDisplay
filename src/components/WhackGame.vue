<script setup>
// 拍猪咪小游戏:30秒限时,猪咪随机探头,拍到得分,结算换鱼干
import { ref, onUnmounted } from 'vue'

const playing = ref(false)
const score = ref(0)
const timeLeft = ref(30)
const target = ref(null) // {x, y, size, key}
const best = ref(Number(localStorage.getItem('catcafe_whack_best') || 0))
const result = ref('')
let timer = null
let popTimer = null

function start() {
  playing.value = true
  score.value = 0
  timeLeft.value = 30
  result.value = ''
  window.dispatchEvent(new CustomEvent('whack', { detail: 'start' }))
  timer = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) end()
  }, 1000)
  pop()
}

function pop() {
  target.value = {
    key: Date.now(),
    x: 8 + Math.random() * 78,
    y: 14 + Math.random() * 62,
    size: 62 + Math.random() * 26,
  }
  popTimer = setTimeout(pop, 700 + Math.random() * 400)
}

function boop() {
  score.value++
  target.value = null
  clearTimeout(popTimer)
  popTimer = setTimeout(pop, 220) // 拍中后快速出下一只
}

async function end() {
  clearInterval(timer)
  clearTimeout(popTimer)
  playing.value = false
  target.value = null
  window.dispatchEvent(new CustomEvent('whack', { detail: 'end' }))
  if (score.value > best.value) {
    best.value = score.value
    localStorage.setItem('catcafe_whack_best', String(best.value))
  }
  // 结算鱼干:1分=1鱼干,上限10
  const n = Math.min(score.value, 10)
  const nick = localStorage.getItem('catcafe_nick')
  if (nick && n > 0) {
    try {
      const res = await fetch('/api/points/add', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nick, n }),
      })
      const data = await res.json()
      if (res.ok) {
        result.value = `拍到 ${score.value} 只!+${n} 鱼干 → 现有 ${data.points} 鱼干「${data.title}」`
        return
      }
    } catch { /* 落到通用提示 */ }
  }
  result.value = `拍到 ${score.value} 只!${nick ? '' : '(输昵称可查头衔,鱼干自动入账)'}`
}

onUnmounted(() => {
  clearInterval(timer)
  clearTimeout(popTimer)
})
</script>

<template>
  <button class="whack-btn" @click="start">🎯 拍猪咪</button>

  <div v-if="playing" class="whack-hud">
    <span>⏱️ {{ timeLeft }}s</span>
    <span>🐷 {{ score }} 只</span>
  </div>

  <button
    v-if="target"
    :key="target.key"
    class="whack-target"
    :style="{ left: target.x + 'vw', top: target.y + 'vh', width: target.size + 'px' }"
    @pointerdown="boop"
  >
    <img src="/pets/pigmi.png" alt="拍我" draggable="false" />
  </button>

  <transition name="bub">
    <div v-if="result" class="whack-result" @click="result = ''">
      <p class="font-cute">{{ result }}</p>
      <p class="best">历史最佳:{{ best }} 只 · 点我关闭</p>
    </div>
  </transition>
</template>

<style scoped>
.whack-btn {
  position: fixed;
  left: 20px;
  bottom: 166px;
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

.whack-btn:hover {
  background: var(--pink-pale);
}

.whack-hud {
  position: fixed;
  top: 76px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 90;
  display: flex;
  gap: 18px;
  background: rgba(255, 255, 255, 0.94);
  border: 2px solid var(--pink-soft);
  border-radius: 999px;
  padding: 8px 24px;
  font-size: 1rem;
  color: var(--ink);
  box-shadow: var(--shadow);
}

.whack-target {
  position: fixed;
  z-index: 91;
  background: none;
  border: none;
  cursor: pointer;
  transform: translate(-50%, -100%);
  animation: popIn 0.18s ease-out;
  touch-action: none;
}

.whack-target img {
  width: 100%;
  pointer-events: none;
}

.whack-target:active {
  transform: translate(-50%, -100%) scale(0.85);
}

@keyframes popIn {
  from { transform: translate(-50%, -100%) scale(0); }
  to { transform: translate(-50%, -100%) scale(1); }
}

.whack-result {
  position: fixed;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 95;
  background: #fff;
  border: 3px dashed var(--pink-soft);
  border-radius: 20px;
  padding: 24px 36px;
  text-align: center;
  box-shadow: var(--shadow);
  cursor: pointer;
}

.whack-result p {
  color: var(--pink-deep);
  font-size: 1.05rem;
}

.whack-result .best {
  color: var(--muted);
  font-size: 0.82rem;
  margin-top: 8px;
}

.bub-enter-active { transition: all 0.25s ease; }
.bub-leave-active { transition: all 0.2s ease; }
.bub-enter-from, .bub-leave-to { opacity: 0; transform: translate(-50%, -50%) scale(0.85); }
.bub-enter-to, .bub-leave-from { transform: translate(-50%, -50%) scale(1); }
</style>

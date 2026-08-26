<script setup>
// 拍猪咪小游戏:30秒限时,猪咪随机探头,拍到得分,结算换鱼干
import { ref, onMounted, onUnmounted } from 'vue'

const playing = ref(false)
const score = ref(0)
const timeLeft = ref(30)
const target = ref(null) // {x, y, size, key}
const best = ref(Number(localStorage.getItem('catcafe_whack_best') || 0))
const result = ref('')
let timer = null
let popTimer = null

function start() {
  localStorage.setItem('catcafe_ach_whack', '1')
  playing.value = true
  score.value = 0
  timeLeft.value = 30
  result.value = ''
  document.body.style.overflow = 'hidden' // 冻结页面,防止滚动/误触背后的按钮
  window.dispatchEvent(new CustomEvent('whack', { detail: 'start' }))
  timer = setInterval(() => {
    timeLeft.value--
    if (timeLeft.value <= 0) end()
  }, 1000)
  pop()
}

function quit() {
  clearInterval(timer)
  clearTimeout(popTimer)
  playing.value = false
  target.value = null
  document.body.style.overflow = ''
  window.dispatchEvent(new CustomEvent('whack', { detail: 'end' }))
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
  document.body.style.overflow = ''
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

function onMenuWhack() {
  if (!playing.value) start()
}

onMounted(() => window.addEventListener('menu-whack', onMenuWhack))

onUnmounted(() => {
  clearInterval(timer)
  clearTimeout(popTimer)
  document.body.style.overflow = ''
  window.removeEventListener('menu-whack', onMenuWhack)
})
</script>

<template>
  <!-- 游戏护盾:全屏透明层,吞掉所有点击/滑动,防止点到背后的按钮 -->
  <div v-if="playing" class="whack-shield" @pointerdown.prevent @touchmove.prevent></div>

  <div v-if="playing" class="whack-hud">
    <span>⏱️ {{ timeLeft }}s</span>
    <span>{{ score }} 只</span>
    <button class="whack-quit" @pointerdown.stop="quit">退出</button>
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
/* 全屏护盾:在游戏 HUD/目标之下,页面内容之上 */
.whack-shield {
  position: fixed;
  inset: 0;
  z-index: 89;
  background: rgba(255, 246, 248, 0.45);
  backdrop-filter: blur(1px);
  touch-action: none;
}

.whack-hud {
  position: fixed;
  top: 76px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 90;
  display: flex;
  align-items: center;
  gap: 18px;
  background: rgba(255, 255, 255, 0.94);
  border: 2px solid var(--pink-soft);
  border-radius: 999px;
  padding: 8px 24px;
  font-size: 1rem;
  color: var(--ink);
  box-shadow: var(--shadow);
}

.whack-quit {
  border: none;
  background: var(--pink-pale);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 4px 12px;
  font-size: 0.85rem;
  font-family: inherit;
  cursor: pointer;
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
  background: var(--surface-2);
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

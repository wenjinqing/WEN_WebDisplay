<script setup>
// 摇一摇猫粮盲盒:每天一盒 —— APP 里摇手机开盒(加速度传感器+振动),网页点按钮开盒
import { ref, computed, onMounted, onUnmounted } from 'vue'
import SectionTitle from './SectionTitle.vue'
import { isApp, profile } from '../app.js'

const lastKey = 'catcafe_shake_last'
const today = new Date().toISOString().slice(0, 10)

const opened = ref(localStorage.getItem(lastKey) === today)
const shaking = ref(false)
const reward = ref(null) // { type: 'fish'|'fortune', text }

const canShake = computed(() => isApp && typeof DeviceMotionEvent !== 'undefined')

// 奖池:鱼干靠昵称入账,稀有签图个吉利
const FORTUNES = [
  '【稀有签 · 糖霜流星】今晚的更新会比糖还甜',
  '【稀有签 · 猫爪招财】本周宜催更,忌熬夜,大吉',
  '【稀有签 · 猪咪附体】你嗑的 CP 下一章必发糖',
  '【稀有签 · 月光浴缸】梦里会有猫耳女仆给你端咖啡',
]

function pickReward() {
  const roll = Math.random()
  if (roll < 0.05) {
    return { type: 'fortune', text: FORTUNES[Math.floor(Math.random() * FORTUNES.length)] }
  }
  const n = roll < 0.6 ? 1 + Math.floor(Math.random() * 3) : roll < 0.9 ? 5 : 10
  return { type: 'fish', n, text: '' }
}

async function openBox() {
  if (opened.value || shaking.value) return
  shaking.value = true
  if (navigator.vibrate) navigator.vibrate([40, 60, 40])

  const r = pickReward()
  if (r.type === 'fish') {
    if (profile.nick) {
      try {
        const res = await fetch('/api/points/add', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ nick: profile.nick, n: r.n }),
        })
        const data = await res.json()
        r.text = res.ok
          ? `开出 ${r.n} 根小鱼干!现有 ${data.points} 根「${data.title}」`
          : `开出 ${r.n} 根小鱼干!`
      } catch {
        r.text = `开出 ${r.n} 根小鱼干!(网络开小差,鱼干没入账)`
      }
    } else {
      r.text = `开出 ${r.n} 根小鱼干!可惜没绑昵称,鱼干存不进罐子`
    }
  }

  // 盒子晃一晃再揭晓
  setTimeout(() => {
    reward.value = r
    opened.value = true
    shaking.value = false
    localStorage.setItem(lastKey, today)
    if (navigator.vibrate) navigator.vibrate(80)
  }, 900)
}

// 摇手机检测:加速度突变超过阈值视为摇了一下
let lastMag = 0
let coolDown = 0
function onMotion(e) {
  if (opened.value || shaking.value) return
  const a = e.accelerationIncludingGravity
  if (!a) return
  const mag = Math.sqrt((a.x || 0) ** 2 + (a.y || 0) ** 2 + (a.z || 0) ** 2)
  const delta = Math.abs(mag - lastMag)
  lastMag = mag
  const now = Date.now()
  if (delta > 14 && now > coolDown) {
    coolDown = now + 1200
    openBox()
  }
}

onMounted(() => {
  if (canShake.value) window.addEventListener('devicemotion', onMotion)
})

onUnmounted(() => {
  if (canShake.value) window.removeEventListener('devicemotion', onMotion)
})
</script>

<template>
  <section id="shakebox">
    <div class="container">
      <SectionTitle title="猫粮盲盒" sub="每天一盒,摇出今日份好运" />

      <div class="box-panel" v-reveal>
        <button
          class="box"
          :class="{ shaking, opened }"
          :disabled="opened"
          @click="openBox"
          aria-label="开盲盒"
        >
          <span class="box-emoji">{{ opened ? '🐟' : '🎁' }}</span>
        </button>
        <p v-if="!opened" class="hint">
          {{ canShake ? '摇一摇手机,或者点一下盒子' : '点一下盒子,开出今日好运' }}
        </p>
        <p v-else-if="reward" class="reward">{{ reward.text }}</p>
        <p v-if="opened" class="hint">明天再来一盒~</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.box-panel {
  max-width: 420px;
  margin: 0 auto;
  text-align: center;
  background: var(--card);
  border: 3px dashed var(--pink-soft);
  border-radius: var(--radius);
  padding: 32px 20px;
  box-shadow: var(--shadow);
}

.box {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 4.5rem;
  line-height: 1;
  transition: transform 0.2s;
}

.box:active {
  transform: scale(0.9);
}

.box.shaking .box-emoji {
  display: inline-block;
  animation: shake 0.45s ease-in-out infinite;
}

.box.opened {
  cursor: default;
}

.box-emoji {
  display: inline-block;
}

@keyframes shake {
  0%, 100% { transform: rotate(-12deg) translateY(0); }
  25% { transform: rotate(10deg) translateY(-6px); }
  50% { transform: rotate(-8deg) translateY(0); }
  75% { transform: rotate(12deg) translateY(-4px); }
}

.hint {
  margin-top: 14px;
  color: var(--muted);
  font-size: 0.88rem;
}

.reward {
  margin-top: 14px;
  color: var(--pink-deep);
  font-size: 1.02rem;
  font-weight: 600;
  line-height: 1.7;
  animation: pop 0.35s ease;
}

@keyframes pop {
  from { transform: scale(0.7); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
</style>

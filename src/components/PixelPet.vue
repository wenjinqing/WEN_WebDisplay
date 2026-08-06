<script setup>
// 爱丽丝看板娘 + 猪咪跟屁虫
// 猫:散步 / 休息 / 睡觉 / 追鼠标 / 说话
// 猪:追着猫跑(落后会加速追) / 偶尔撒欢乱跑 / 也会说话
import { ref, reactive, onMounted, onUnmounted } from 'vue'

const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')

const cat = reactive({ x: 15, dir: 1, state: 'walk', bubble: '', anim: '' })
const pig = reactive({ x: 8, dir: 1, state: 'walk', mode: 'follow', bubble: '', anim: '' })
// anim: 临时动作(被戳),优先于 state 显示
// pig.mode: follow 追猫 | sprint 撒欢乱跑 | idle 休息

const CAT_SAYS = [
  '赶稿中……', '想喝奶茶', '猪咪们好呀~', '在写新坑!',
  '今天也要加油喵', '催更?在写了在写了', '逛逛自己的店',
  '刚才那句写得不错', '想吃草莓蛋糕', '被发现了?',
]
const PIG_SAYS = [
  '等等我!', '猪咪来咯~', '拱拱', '今天也很乖',
  '猫猫酱慢点!', '哼唧哼唧', '贴贴!',
]
const CAT_POKED = ['喵?叫我?', '在呢在呢!', '怎么啦~', '摸我就不用赶稿了?']
const PIG_POKED = ['哼唧!', '拱你一下!', '猪咪超乖的!', '干嘛啦~']

let tick = null
let stateTimer = null
let bubbleTimer = null
let mouseX = 50
const clears = { cat: null, pig: null }

function say(who, text, ms = 3000) {
  who.bubble = text
  clearTimeout(clears[who === cat ? 'cat' : 'pig'])
  clears[who === cat ? 'cat' : 'pig'] = setTimeout(() => (who.bubble = ''), ms)
}

function step() {
  if (hidden.value) return
  // ===== 猫移动 =====
  if (cat.state === 'walk') {
    cat.x += cat.dir * 0.3
    if (cat.x >= 90) { cat.x = 90; cat.dir = -1 }
    if (cat.x <= 1) { cat.x = 1; cat.dir = 1 }
  }
  // ===== 猪移动 =====
  if (pig.mode === 'follow') {
    const target = cat.x - cat.dir * 7 // 跟在猫身后 7vw
    const dist = target - pig.x
    if (Math.abs(dist) > 1) {
      pig.dir = dist > 0 ? 1 : -1
      // 落后越多跑得越急,追上立刻减速
      const speed = Math.abs(dist) > 18 ? 0.5 : 0.24
      pig.x += pig.dir * Math.min(speed, Math.abs(dist))
      pig.state = 'walk'
    } else if (cat.state !== 'walk') {
      pig.state = 'idle' // 猫停猪也停
    }
  } else if (pig.mode === 'sprint') {
    pig.x += pig.dir * 0.55
    pig.state = 'walk'
    if (pig.x >= 92) { pig.x = 92; pig.dir = -1 }
    if (pig.x <= 1) { pig.x = 1; pig.dir = 1 }
  }
  pig.x = Math.max(1, Math.min(92, pig.x))
}

function maybeChange() {
  if (hidden.value) return
  const r = Math.random()
  // ===== 猫的状态 =====
  if (cat.state === 'walk') {
    if (r < 0.2) cat.state = 'idle'
    else if (r < 0.28) cat.state = 'sleep'
    else if (r < 0.55) cat.dir = mouseX > cat.x ? 1 : -1
  } else if (r < 0.65) {
    cat.state = 'walk'
  }
  // ===== 猪的模式:主要追猫,偶尔撒欢 =====
  const r2 = Math.random()
  if (pig.mode === 'follow' && r2 < 0.14) {
    pig.mode = 'sprint'
    pig.dir = Math.random() < 0.5 ? 1 : -1
    say(pig, '撒欢啦——!', 2000)
    setTimeout(() => (pig.mode = 'follow'), 2500 + Math.random() * 2500)
  } else if (pig.mode === 'follow' && r2 < 0.22) {
    pig.state = pig.state === 'idle' ? 'walk' : 'idle'
  }
}

function maybeBubble() {
  if (hidden.value) return
  const r = Math.random()
  if (r < 0.3) say(cat, CAT_SAYS[Math.floor(Math.random() * CAT_SAYS.length)], 3200)
  else if (r < 0.5) say(pig, PIG_SAYS[Math.floor(Math.random() * PIG_SAYS.length)], 2800)
}

function pokeCat() {
  say(cat, CAT_POKED[Math.floor(Math.random() * CAT_POKED.length)], 2500)
  cat.anim = 'poke'
  cat.state = 'idle'
  setTimeout(() => (cat.anim = ''), 900)
  setTimeout(() => (cat.state = 'walk'), 2000)
}

function pokePig() {
  say(pig, PIG_POKED[Math.floor(Math.random() * PIG_POKED.length)], 2200)
  pig.anim = 'poke'
  setTimeout(() => (pig.anim = ''), 900)
  // 被点了会吓得跑开一小段
  pig.mode = 'sprint'
  pig.dir = pig.x > cat.x ? 1 : -1
  setTimeout(() => (pig.mode = 'follow'), 1500)
}

function hide() {
  hidden.value = true
  localStorage.setItem('catcafe_pet_hide', '1')
}

function show() {
  hidden.value = false
  localStorage.removeItem('catcafe_pet_hide')
}

function onMouse(e) {
  mouseX = (e.clientX / window.innerWidth) * 100
}

onMounted(() => {
  window.addEventListener('mousemove', onMouse, { passive: true })
  tick = setInterval(step, 60)
  stateTimer = setInterval(maybeChange, 4000)
  bubbleTimer = setInterval(maybeBubble, 8000)
})

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouse)
  clearInterval(tick)
  clearInterval(stateTimer)
  clearInterval(bubbleTimer)
  clearTimeout(clears.cat)
  clearTimeout(clears.pig)
})
</script>

<template>
  <button v-if="hidden" class="pet-restore" aria-label="召回看板娘" @click="show">🐾</button>

  <template v-else>
    <!-- 爱丽丝 -->
    <div class="pet" :style="{ left: cat.x + 'vw' }">
      <transition name="bub">
        <span v-if="cat.bubble" class="pet-bubble">{{ cat.bubble }}</span>
      </transition>
      <button class="pet-hide" aria-label="让她们去休息" @click="hide">×</button>
      <span v-if="cat.state === 'sleep'" class="zzz">💤</span>
      <span class="flipper" :class="{ flip: cat.dir < 0 }">
        <div
          class="sprite alice"
          :class="cat.anim || cat.state"
          role="button"
          tabindex="0"
          aria-label="戳一下爱丽丝"
          @click="pokeCat"
          @keyup.enter="pokeCat"
        />
      </span>
    </div>

    <!-- 猪咪跟屁虫 -->
    <div class="pet pig" :style="{ left: pig.x + 'vw' }">
      <transition name="bub">
        <span v-if="pig.bubble" class="pet-bubble pig-bubble">{{ pig.bubble }}</span>
      </transition>
      <span class="flipper" :class="{ flip: pig.dir < 0 }">
        <div
          class="sprite pigmi"
          :class="pig.anim || pig.state"
          role="button"
          tabindex="0"
          aria-label="戳一下猪咪"
          @click="pokePig"
          @keyup.enter="pokePig"
        />
      </span>
    </div>
  </template>
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

/* ===== 精灵图帧动画 ===== */
.sprite {
  image-rendering: pixelated;
  pointer-events: auto;
  cursor: pointer;
  user-select: none;
  background-repeat: no-repeat;
}

/* 爱丽丝 Evil 全身像素小人(帧宽64,高91) */
.sprite.alice { width: 64px; height: 91px; }
.sprite.pigmi { width: 56px; height: 56px; }

/* 爱丽丝(帧宽64) */
.sprite.alice.walk {
  background-image: url('/pets/evil-walk.png');
  background-size: 256px 91px;
  animation: fr4a 0.55s steps(1) infinite;
}
.sprite.alice.idle {
  background-image: url('/pets/evil-idle.png');
  background-size: 128px 91px;
  animation: fr2a 1.4s steps(1) infinite;
}
.sprite.alice.sleep {
  background-image: url('/pets/evil-sleep.png');
  background-size: 128px 91px;
  animation: fr2a 2.2s steps(1) infinite;
}
.sprite.alice.poke {
  background-image: url('/pets/evil-poke.png');
  background-size: 128px 91px;
  animation: fr2a 0.3s steps(1) infinite;
}

/* 猪咪(帧宽56,腿短倒腾快) */
.sprite.pigmi.walk {
  background-image: url('/pets/pigmi-walk.png');
  background-size: 224px 56px;
  animation: fr4p 0.4s steps(1) infinite;
}
.sprite.pigmi.idle {
  background-image: url('/pets/pigmi-idle.png');
  background-size: 112px 56px;
  animation: fr2p 1.4s steps(1) infinite;
}
.sprite.pigmi.sleep {
  background-image: url('/pets/pigmi-sleep.png');
  background-size: 112px 56px;
  animation: fr2p 2.2s steps(1) infinite;
}
.sprite.pigmi.poke {
  background-image: url('/pets/pigmi-poke.png');
  background-size: 112px 56px;
  animation: fr2p 0.3s steps(1) infinite;
}

@keyframes fr4a {
  0% { background-position-x: 0; }
  25% { background-position-x: -64px; }
  50% { background-position-x: -128px; }
  75% { background-position-x: -192px; }
}
@keyframes fr2a {
  0% { background-position-x: 0; }
  50% { background-position-x: -64px; }
}
@keyframes fr4p {
  0% { background-position-x: 0; }
  25% { background-position-x: -56px; }
  50% { background-position-x: -112px; }
  75% { background-position-x: -168px; }
}
@keyframes fr2p {
  0% { background-position-x: 0; }
  50% { background-position-x: -56px; }
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

.pig-bubble {
  border-color: #c9d8f5;
  border-radius: 14px 14px 4px 14px;
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
  .pet-hide {
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .sprite {
    animation: none !important;
  }
}
</style>

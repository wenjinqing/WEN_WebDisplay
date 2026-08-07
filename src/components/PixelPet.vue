<script setup>
// 爱丽丝猪 & 猪咪跟屁虫 —— 插画版看板娘(irasutoya 免费素材)
// 爱丽丝猪:跑步(散步) / 站立(待机) / 用电脑(赶稿) / 惊讶(被戳)
// 猪咪:梦幻小猪,追爱丽丝、撒欢、被戳会跑
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'

const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')

const cat = reactive({ x: 15, dir: 1, state: 'walk', bubble: '', anim: '' })
const pig = reactive({ x: 8, dir: 1, state: 'walk', mode: 'follow', bubble: '', anim: '' })
// cat.state: walk 散步 | idle 待机 | sleep 赶稿(电脑猪)
// pig.mode: follow 追 | sprint 撒欢

const ALICE_IMG = {
  walk: '/pets/pig-run.png',
  idle: '/pets/pig-stand.png',
  sleep: '/pets/pig-computer.png', // 睡觉?不,她在赶稿
  poke: '/pets/pig-shock.png',
}
const catSrc = computed(() => ALICE_IMG[cat.anim || cat.state] || ALICE_IMG.idle)

const CAT_SAYS = [
  '赶稿中……', '想喝奶茶', '猪咪们好呀~', '在写新坑!',
  '今天也要加油哼', '催更?在写了在写了', '逛逛自己的店',
  '刚才那句写得不错', '想吃草莓蛋糕', '被发现了?',
]
const PIG_SAYS = [
  '等等我!', '猪咪来咯~', '拱拱', '今天也很乖',
  '爱丽丝慢点!', '哼唧哼唧', '贴贴!',
]
const CAT_POKED = ['哼?叫我?', '在写了在写了!', '怎么啦~', '戳我就不用赶稿了?']
const PIG_POKED = ['哼唧!', '拱你一下!', '猪咪超乖的!', '干嘛啦~']

let tick = null
let stateTimer = null
let bubbleTimer = null
let mouseX = 50
const clears = { cat: null, pig: null }

function say(who, text, ms = 3000) {
  who.bubble = text
  const key = who === cat ? 'cat' : 'pig'
  clearTimeout(clears[key])
  clears[key] = setTimeout(() => (who.bubble = ''), ms)
}

function step() {
  if (hidden.value) return
  // ===== 爱丽丝猪移动 =====
  if (cat.state === 'walk') {
    cat.x += cat.dir * 0.3
    if (cat.x >= 90) { cat.x = 90; cat.dir = -1 }
    if (cat.x <= 1) { cat.x = 1; cat.dir = 1 }
  }
  // ===== 猪咪移动 =====
  if (pig.mode === 'follow') {
    const target = cat.x - cat.dir * 7
    const dist = target - pig.x
    if (Math.abs(dist) > 1) {
      pig.dir = dist > 0 ? 1 : -1
      const speed = Math.abs(dist) > 18 ? 0.5 : 0.24
      pig.x += pig.dir * Math.min(speed, Math.abs(dist))
      pig.state = 'walk'
    } else if (cat.state !== 'walk') {
      pig.state = 'idle'
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
  if (cat.state === 'walk') {
    if (r < 0.2) cat.state = 'idle'
    else if (r < 0.32) cat.state = 'sleep' // 坐下来赶稿
    else if (r < 0.55) cat.dir = mouseX > cat.x ? 1 : -1
  } else if (r < 0.6) {
    cat.state = 'walk'
  }
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
  setTimeout(() => (cat.anim = ''), 900)
}

function pokePig() {
  say(pig, PIG_POKED[Math.floor(Math.random() * PIG_POKED.length)], 2200)
  pig.anim = 'poke'
  setTimeout(() => (pig.anim = ''), 900)
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
    <!-- 爱丽丝猪 -->
    <div class="pet" :style="{ left: cat.x + 'vw' }">
      <transition name="bub">
        <span v-if="cat.bubble" class="pet-bubble">{{ cat.bubble }}</span>
      </transition>
      <button class="pet-hide" aria-label="让她们去休息" @click="hide">×</button>
      <span v-if="cat.state === 'sleep'" class="work-tag">💻 赶稿中</span>
      <span class="flipper" :class="{ flip: cat.dir < 0 && cat.state === 'walk' }">
        <img
          :src="catSrc"
          alt="爱丽丝猪"
          class="pet-img alice-img"
          :class="cat.anim || cat.state"
          draggable="false"
          role="button"
          tabindex="0"
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
      <span v-if="pig.state === 'idle' && pig.mode === 'follow'" class="zzz">💤</span>
      <span class="flipper" :class="{ flip: pig.dir < 0 }">
        <img
          src="/pets/pigmi.png"
          alt="猪咪"
          class="pet-img pigmi-img"
          :class="pig.anim || pig.state"
          draggable="false"
          role="button"
          tabindex="0"
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

.pet-img {
  pointer-events: auto;
  cursor: pointer;
  user-select: none;
}

.alice-img {
  height: 88px;
}

.pigmi-img {
  height: 60px;
}

/* 走路:上下小跑 */
.pet-img.walk {
  animation: walkBob 0.4s ease-in-out infinite;
}

@keyframes walkBob {
  0%, 100% { transform: translateY(0) rotate(-3deg); }
  50% { transform: translateY(-6px) rotate(3deg); }
}

/* 待机:呼吸 */
.pet-img.idle {
  animation: breathe 2.4s ease-in-out infinite;
}

@keyframes breathe {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.95); }
}

/* 赶稿:轻微前倾敲键盘感 */
.pet-img.sleep {
  animation: typing 0.6s ease-in-out infinite;
}

@keyframes typing {
  0%, 100% { transform: rotate(0deg); }
  50% { transform: rotate(1.5deg) translateY(1px); }
}

/* 被戳:左右晃 */
.pet-img.poke {
  animation: wiggle 0.3s ease-in-out infinite;
}

@keyframes wiggle {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}

.work-tag {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid var(--pink-pale);
  border-radius: 999px;
  padding: 1px 10px;
  white-space: nowrap;
  color: var(--muted);
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
  .alice-img {
    height: 70px;
  }
  .pigmi-img {
    height: 50px;
  }
  .pet-hide {
    opacity: 1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .pet-img {
    animation: none !important;
  }
}
</style>

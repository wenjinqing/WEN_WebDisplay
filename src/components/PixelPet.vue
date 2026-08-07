<script setup>
// 猫咖看板娘引擎 v3 —— 全屏漫游 / 投喂 / 小剧场 / 追鼠标 / 捡鱼干 / 深夜作息
import { ref, reactive, onMounted, onUnmounted } from 'vue'

const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')
const isNight = ref(false)

// 坐标:x 用 vw,y 用 vh(全屏自由移动)
const cat = reactive({ x: 20, y: 90, dir: 1, state: 'walk', bubble: '', anim: '', tx: 50, ty: 90, food: null })
const pig = reactive({ x: 12, y: 90, dir: 1, state: 'walk', mode: 'follow', bubble: '', anim: '', tx: 40, ty: 90, food: null })
// state: walk | idle | sleep | chase(猫)
// pig.mode: follow | sprint

const CAT_SAYS = ['赶稿中……', '想喝奶茶', '猪咪们好呀~', '在写新坑!', '今天也要加油喵', '催更?在写了在写了', '巡视店铺中', '刚才那句写得不错', '想吃草莓蛋糕', '被发现了?']
const PIG_SAYS = ['等等我!', '猪咪来咯~', '拱拱', '今天也很乖', '猫猫酱慢点!', '哼唧哼唧', '贴贴!']
const CAT_POKED = ['喵?叫我?', '在呢在呢!', '怎么啦~', '摸我就不用赶稿了?']
const PIG_POKED = ['哼唧!', '拱你一下!', '猪咪超乖的!', '干嘛啦~']
const DIALOGUES = [
  ['新坑写不出来……', '摸摸,不急不急'],
  ['猪咪,我饿了', '我也饿了!', '？'],
  ['今天店里好热闹', '都是来看你的!'],
  ['你说新坑写刀还是写糖', '糖!必须是糖!'],
  ['别跟着我啦', '就要跟!'],
]

const foods = ref([]) // {id, x, y, type:'🐟'|'🍰'}
const cracker = ref(null) // 猪咪拱出来的鱼干 {x, y}
const feeds = ref(0)
const toast = ref('')
let foodSeq = 0
let mouse = { x: 50, y: 90 }
let mouseTrail = []
let chaseUntil = 0

let tick, stateTimer, bubbleTimer, nightTimer, bubbleClear = { cat: null, pig: null }, toastTimer, crackerTimer

function say(who, text, ms = 3000) {
  who.bubble = text
  const key = who === cat ? 'cat' : 'pig'
  clearTimeout(bubbleClear[key])
  bubbleClear[key] = setTimeout(() => (who.bubble = ''), ms)
}

function showToast(text) {
  toast.value = text
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 2600)
}

// ===== 目标选择:底部溜达 / 中场巡视 / 趴卡片 =====
function pickTarget(p) {
  const r = Math.random()
  if (r < 0.6) {
    p.tx = 5 + Math.random() * 85
    p.ty = 86 + Math.random() * 8
  } else if (r < 0.8) {
    p.tx = 10 + Math.random() * 75
    p.ty = 35 + Math.random() * 40 // 中场巡视
  } else {
    // 趴卡片:找一个视口内的大卡片,坐在它顶边上
    const els = document.querySelectorAll('.menu-board, .board, .club-card, .send-card, .fortune')
    const seats = []
    els.forEach((el) => {
      const rect = el.getBoundingClientRect()
      if (rect.top > 120 && rect.top < window.innerHeight - 160 && rect.width > 300) {
        seats.push(rect)
      }
    })
    if (seats.length) {
      const s = seats[Math.floor(Math.random() * seats.length)]
      const petH = (100 / window.innerHeight) * 100
      p.tx = ((s.left + Math.random() * s.width * 0.8) / window.innerWidth) * 100
      p.ty = (s.top / window.innerHeight) * 100 - petH + 2
      return
    }
    p.tx = 10 + Math.random() * 75
    p.ty = 86 + Math.random() * 8
  }
}

function moveToward(p, speed) {
  const dx = p.tx - p.x
  const dy = p.ty - p.y
  const dist = Math.hypot(dx, dy)
  if (dist < 1) return true
  p.dir = dx < -0.5 ? -1 : dx > 0.5 ? 1 : p.dir
  p.x += (dx / dist) * speed
  p.y += (dy / dist) * speed
  return false
}

function step() {
  if (hidden.value) return
  if (isNight.value) return // 夜里都睡着

  // ===== 猫 =====
  if (cat.state === 'chase') {
    cat.tx = mouse.x
    cat.ty = Math.min(mouse.y, 92)
    if (moveToward(cat, 0.7)) {
      say(cat, '抓到啦!', 1500)
      cat.state = 'idle'
      setTimeout(() => (cat.state = 'walk'), 1500)
    } else if (Date.now() > chaseUntil) {
      cat.state = 'walk'
      pickTarget(cat)
    }
  } else if (cat.state === 'walk') {
    const arrived = moveToward(cat, cat.food ? 0.55 : 0.28)
    if (cat.food && arrived) eatFood(cat)
    else if (arrived) cat.state = Math.random() < 0.5 ? 'idle' : 'walk'
  }

  // ===== 猪 =====
  if (pig.food) {
    if (moveToward(pig, 0.5)) eatFood(pig)
  } else if (pig.mode === 'follow') {
    pig.tx = cat.x - cat.dir * 8
    pig.ty = cat.y + 3
    const dist = Math.hypot(pig.tx - pig.x, pig.ty - pig.y)
    if (dist > 2) {
      pig.state = 'walk'
      moveToward(pig, dist > 20 ? 0.5 : 0.24)
    } else if (cat.state !== 'walk' && cat.state !== 'chase') {
      pig.state = 'idle'
    }
  } else if (pig.mode === 'sprint') {
    pig.state = 'walk'
    if (moveToward(pig, 0.55)) pig.mode = 'follow'
  }
}

function maybeChange() {
  if (hidden.value || isNight.value) return
  const r = Math.random()
  // 猫的状态机
  if (cat.state === 'walk' && !cat.food) {
    if (r < 0.18) cat.state = 'idle'
    else if (r < 0.26) cat.state = 'sleep'
    else if (r < 0.7) pickTarget(cat)
  } else if (cat.state !== 'chase' && !cat.food && r < 0.55) {
    cat.state = 'walk'
    pickTarget(cat)
  }
  // 猪的模式
  const r2 = Math.random()
  if (pig.mode === 'follow' && !pig.food && r2 < 0.12) {
    pig.mode = 'sprint'
    pig.tx = 5 + Math.random() * 85
    pig.ty = 50 + Math.random() * 44
    say(pig, '撒欢啦——!', 2000)
    setTimeout(() => (pig.mode = 'follow'), 2600)
  }
  // 猪咪拱鱼干
  if (!cracker.value && pig.state === 'walk' && r2 > 0.94) {
    cracker.value = { x: pig.x, y: pig.y - 4 }
    say(pig, '咦?有鱼干!', 2000)
    crackerTimer = setTimeout(() => (cracker.value = null), 5000)
  }
}

function maybeBubble() {
  if (hidden.value || isNight.value) return
  const close = Math.hypot(cat.x - pig.x, cat.y - pig.y) < 12
  const r = Math.random()
  if (close && !cat.bubble && !pig.bubble && r < 0.3) {
    // 小剧场
    const d = DIALOGUES[Math.floor(Math.random() * DIALOGUES.length)]
    say(cat, d[0], 2600)
    setTimeout(() => say(pig, d[1], 2600), 1400)
    if (d[2]) setTimeout(() => say(cat, d[2], 2000), 3000)
  } else if (r < 0.45) {
    say(cat, CAT_SAYS[Math.floor(Math.random() * CAT_SAYS.length)], 3200)
  } else if (r < 0.6) {
    say(pig, PIG_SAYS[Math.floor(Math.random() * PIG_SAYS.length)], 2800)
  }
}

// ===== 投喂 =====
function dropFood() {
  const type = Math.random() < 0.5 ? '🐟' : '🍰'
  const food = {
    id: ++foodSeq,
    type,
    x: 20 + Math.random() * 55,
    y: 76 + Math.random() * 14,
  }
  foods.value.push(food)
  const target = type === '🐟' ? cat : pig // 猫抢鱼,猪抢蛋糕
  target.food = food
  target.state = 'walk'
  say(target, type === '🐟' ? '是鱼!冲!' : '蛋糕!我的!', 2000)
}

function eatFood(p) {
  foods.value = foods.value.filter((f) => f.id !== p.food.id)
  p.food = null
  p.anim = 'eat'
  say(p, '啊呜啊呜……好吃!', 2200)
  setTimeout(() => (p.anim = ''), 1200)
  fetch('/api/feed', { method: 'POST' })
    .then((r) => r.json())
    .then((d) => (feeds.value = d.feeds))
    .catch(() => {})
}

// ===== 捡鱼干 =====
async function pickCracker() {
  cracker.value = null
  clearTimeout(crackerTimer)
  const nick = localStorage.getItem('catcafe_nick')
  if (!nick) {
    showToast('先去「猪咪聚集地」输入昵称查询过头衔,鱼干才能记在你名下哦')
    return
  }
  try {
    const res = await fetch('/api/points/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nick, n: 5 }),
    })
    const data = await res.json()
    if (res.ok) showToast(`🐟 +5 鱼干!${nick} 现有 ${data.points} 鱼干「${data.title}」`)
    else showToast(data.error || '领取失败')
  } catch {
    showToast('网络打了个盹')
  }
}

// ===== 戳一戳 =====
function pokeCat() {
  if (cat.state === 'sleep') {
    say(cat, '唔……再睡五分钟', 2200)
    return
  }
  say(cat, CAT_POKED[Math.floor(Math.random() * CAT_POKED.length)], 2500)
  cat.anim = 'poke'
  setTimeout(() => (cat.anim = ''), 900)
}

function pokePig() {
  say(pig, PIG_POKED[Math.floor(Math.random() * PIG_POKED.length)], 2200)
  pig.anim = 'poke'
  setTimeout(() => (pig.anim = ''), 900)
  pig.mode = 'sprint'
  pig.tx = pig.x > cat.x ? Math.min(pig.x + 20, 90) : Math.max(pig.x - 20, 5)
  pig.ty = pig.y
  setTimeout(() => (pig.mode = 'follow'), 1500)
}

// ===== 鼠标追踪(快速晃动触发猫扑) =====
function onMouse(e) {
  mouse = {
    x: (e.clientX / window.innerWidth) * 100,
    y: (e.clientY / window.innerHeight) * 100,
  }
  const now = Date.now()
  mouseTrail.push({ x: mouse.x, t: now, dir: 0 })
  mouseTrail = mouseTrail.filter((p) => now - p.t < 1200)
  // 统计方向变化
  let turns = 0
  for (let i = 2; i < mouseTrail.length; i++) {
    const d1 = mouseTrail[i - 1].x - mouseTrail[i - 2].x
    const d2 = mouseTrail[i].x - mouseTrail[i - 1].x
    if (d1 * d2 < 0) turns++
  }
  if (turns >= 6 && cat.state === 'walk' && !cat.food && !isNight.value) {
    cat.state = 'chase'
    chaseUntil = now + 3000
    say(cat, '别跑!', 1500)
    mouseTrail = []
  }
}

// 滚动时:不在底部的回去底部
function onScroll() {
  ;[cat, pig].forEach((p) => {
    if (p.y < 80 && !p.food) {
      p.ty = 86 + Math.random() * 8
      if (p.state !== 'chase') p.state = 'walk'
    }
  })
}

function checkNight() {
  const h = new Date().getHours()
  const night = h >= 0 && h < 7
  if (night && !isNight.value) {
    // 一起趴到左下角睡觉
    Object.assign(cat, { x: 8, y: 90, tx: 8, ty: 90, state: 'sleep', food: null })
    Object.assign(pig, { x: 15, y: 92, tx: 15, ty: 92, state: 'sleep', food: null, mode: 'follow' })
    foods.value = []
    cracker.value = null
  } else if (!night && isNight.value) {
    cat.state = 'idle'
    pig.state = 'idle'
    say(cat, '早安,开店啦!', 3000)
  }
  isNight.value = night
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
  window.addEventListener('scroll', onScroll, { passive: true })
  tick = setInterval(step, 60)
  stateTimer = setInterval(maybeChange, 4000)
  bubbleTimer = setInterval(maybeBubble, 8000)
  nightTimer = setInterval(checkNight, 30000)
  checkNight()
  fetch('/api/pet').then((r) => r.json()).catch(() => {})
})

onUnmounted(() => {
  window.removeEventListener('mousemove', onMouse)
  window.removeEventListener('scroll', onScroll)
  clearInterval(tick)
  clearInterval(stateTimer)
  clearInterval(bubbleTimer)
  clearInterval(nightTimer)
  clearTimeout(crackerTimer)
})
</script>

<template>
  <button v-if="hidden" class="pet-restore" aria-label="召回看板娘" @click="show">🐾</button>

  <template v-else>
    <!-- 投喂按钮 -->
    <button class="feed-btn" @click="dropFood">
      🍰 投喂{{ feeds ? ` · 已被喂${feeds}次` : '' }}
    </button>

    <!-- 地上的食物 -->
    <span
      v-for="f in foods"
      :key="f.id"
      class="food"
      :style="{ left: f.x + 'vw', top: f.y + 'vh' }"
    >
      {{ f.type }}
    </span>

    <!-- 猪咪拱出的鱼干(5秒消失,点它+5鱼干) -->
    <button
      v-if="cracker"
      class="cracker"
      :style="{ left: cracker.x + 'vw', top: cracker.y + 'vh' }"
      @click="pickCracker"
    >
      🐟✨
    </button>

    <!-- 看板猫 -->
    <div class="pet" :style="{ left: cat.x + 'vw', top: cat.y + 'vh' }">
      <transition name="bub">
        <span v-if="cat.bubble" class="pet-bubble">{{ cat.bubble }}</span>
      </transition>
      <button class="pet-hide" aria-label="让她们去休息" @click="hide">×</button>
      <span v-if="cat.state === 'sleep'" class="zzz">💤</span>
      <span class="flipper" :class="{ flip: cat.dir < 0 }">
        <img
          src="/pets/cat.png"
          alt="看板猫"
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
    <div class="pet pig" :style="{ left: pig.x + 'vw', top: pig.y + 'vh' }">
      <transition name="bub">
        <span v-if="pig.bubble" class="pet-bubble pig-bubble">{{ pig.bubble }}</span>
      </transition>
      <span v-if="pig.state === 'sleep' || (pig.state === 'idle' && pig.mode === 'follow')" class="zzz">💤</span>
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

    <transition name="bub">
      <p v-if="toast" class="pet-toast">{{ toast }}</p>
    </transition>
  </template>
</template>

<style scoped>
.pet {
  position: fixed;
  z-index: 55;
  pointer-events: none;
  transition: left 0.06s linear, top 0.06s linear;
  transform: translateY(-100%); /* 坐标=落脚点,身子在上方 */
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

.alice-img { height: 88px; }
.pigmi-img { height: 60px; }

/* 走路 */
.pet-img.walk {
  animation: walkBob 0.4s ease-in-out infinite;
}
@keyframes walkBob {
  0%, 100% { transform: translateY(0) rotate(-3deg); }
  50% { transform: translateY(-6px) rotate(3deg); }
}

/* 追鼠标:更急的碎步 */
.pet-img.chase {
  animation: walkBob 0.25s ease-in-out infinite;
}

/* 待机呼吸 */
.pet-img.idle {
  animation: breathe 2.4s ease-in-out infinite;
}
@keyframes breathe {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.95); }
}

/* 睡觉 */
.pet-img.sleep {
  filter: brightness(0.85) saturate(0.85);
  transform: scaleY(0.85);
  transform-origin: bottom;
}

/* 被戳 */
.pet-img.poke {
  animation: wiggle 0.3s ease-in-out infinite;
}
@keyframes wiggle {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}

/* 吃东西 */
.pet-img.eat {
  animation: munch 0.3s ease-in-out infinite;
}
@keyframes munch {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.88); }
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
  bottom: 118px;
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

.feed-btn {
  position: fixed;
  left: 20px;
  bottom: 70px;
  z-index: 60;
  border: 2px solid var(--pink-soft);
  background: rgba(255, 255, 255, 0.92);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 8px 16px;
  font-size: 0.85rem;
  cursor: pointer;
  box-shadow: var(--shadow);
  transition: transform 0.15s;
}

.feed-btn:hover {
  transform: scale(1.06);
  background: var(--pink-pale);
}

.food {
  position: fixed;
  z-index: 54;
  font-size: 26px;
  transform: translate(-50%, -100%);
  animation: foodDrop 0.3s ease-out;
  pointer-events: none;
}

@keyframes foodDrop {
  from { transform: translate(-50%, -160%); opacity: 0; }
  to { transform: translate(-50%, -100%); opacity: 1; }
}

.cracker {
  position: fixed;
  z-index: 56;
  font-size: 24px;
  background: none;
  border: none;
  cursor: pointer;
  transform: translate(-50%, -100%);
  animation: sparkle 0.6s ease-in-out infinite;
  pointer-events: auto;
}

@keyframes sparkle {
  0%, 100% { transform: translate(-50%, -100%) scale(1); }
  50% { transform: translate(-50%, -110%) scale(1.15); }
}

.pet-toast {
  position: fixed;
  bottom: 120px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 70;
  background: var(--ink);
  color: #fff;
  padding: 10px 24px;
  border-radius: 999px;
  font-size: 0.88rem;
  box-shadow: var(--shadow);
  max-width: 80vw;
}

@media (max-width: 720px) {
  .alice-img { height: 70px; }
  .pigmi-img { height: 50px; }
  .pet-hide { opacity: 1; }
}

@media (prefers-reduced-motion: reduce) {
  .pet-img, .cracker {
    animation: none !important;
  }
}
</style>

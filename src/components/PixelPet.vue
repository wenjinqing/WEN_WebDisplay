<script setup>
// 猫咖看板娘引擎 v4 —— 曲线漫游 / 拖拽投喂 / 可抓取小人 / 小剧场 / 追鼠标 / 捡鱼干 / 深夜作息
import { ref, reactive, onMounted, onUnmounted } from 'vue'

const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')
const isNight = ref(false)

const cat = reactive({ x: 20, y: 92, dir: 1, state: 'walk', bubble: '', anim: '', path: null, food: null, held: false })
const pig = reactive({ x: 12, y: 92, dir: 1, state: 'walk', mode: 'follow', bubble: '', anim: '', food: null, held: false, wobble: 0 })
// state: walk | idle | sleep | chase(猫) | held(被抓)

const CAT_SAYS = ['赶稿中……', '想喝奶茶', '猪咪们好呀~', '在写新坑!', '今天也要加油喵', '催更?在写了在写了', '巡视店铺中', '刚才那句写得不错', '想吃草莓蛋糕', '被发现了?']
const PIG_SAYS = ['等等我!', '猪咪来咯~', '拱拱', '今天也很乖', '猫猫酱慢点!', '哼唧哼唧', '贴贴!']
const CAT_POKED = ['喵?叫我?', '在呢在呢!', '怎么啦~', '摸我就不用赶稿了?']
const PIG_POKED = ['哼唧!', '拱你一下!', '猪咪超乖的!', '干嘛啦~']
const CAT_HELD = ['哇啊!放我下来!', '脚够不着地啦!', '你要带我去哪!']
const PIG_HELD = ['哼唧哼唧!', '飞起来了?!', '猪咪害怕!']
const DIALOGUES = [
  ['新坑写不出来……', '摸摸,不急不急'],
  ['猪咪,我饿了', '我也饿了!', '？'],
  ['今天店里好热闹', '都是来看你的!'],
  ['你说新坑写刀还是写糖', '糖!必须是糖!'],
  ['别跟着我啦', '就要跟!'],
]

const foods = ref([])
const cracker = ref(null)
const feeds = ref(0)
const toast = ref('')
const dragFood = ref(null) // 正在拖拽的食物 {type, x, y}
let foodSeq = 0
let mouse = { x: 50, y: 90 }
let mouseTrail = []
let chaseUntil = 0
let heldPet = null
let pendingGrab = null // {pet, x, y} 按下但还没拖动
let grabOffset = { x: 0, y: 0 }
let warnedPig = false

let tick, stateTimer, bubbleTimer, nightTimer, crackerTimer, toastTimer
let bubbleClear = { cat: null, pig: null }

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

// ===== 曲线路径(二次贝塞尔) =====
function setPath(p, tx, ty, curve = true) {
  const x0 = p.x, y0 = p.y
  const dist = Math.hypot(tx - x0, ty - y0)
  let cx = (x0 + tx) / 2, cy = (y0 + ty) / 2
  if (curve && dist > 8) {
    // 垂直于路径方向的随机偏移 → 弧线
    const nx = -(ty - y0) / dist
    const ny = (tx - x0) / dist
    const bend = (Math.random() - 0.5) * dist * 0.9
    cx += nx * bend
    cy += ny * bend * 0.4 // 纵向弯一点就行
    cy = Math.max(15, Math.min(94, cy))
  }
  p.path = { x0, y0, cx, cy, x1: tx, y1: ty, t: 0, len: dist * 1.25 }
}

function stepPath(p, speed) {
  if (!p.path) return true
  p.path.t += speed / Math.max(p.path.len, 1)
  if (p.path.t >= 1) {
    p.x = p.path.x1
    p.y = p.path.y1
    p.path = null
    return true
  }
  const t = p.path.t
  const { x0, y0, cx, cy, x1, y1 } = p.path
  const nx = (1 - t) * (1 - t) * x0 + 2 * (1 - t) * t * cx + t * t * x1
  const ny = (1 - t) * (1 - t) * y0 + 2 * (1 - t) * t * cy + t * t * y1
  p.dir = nx < p.x - 0.05 ? -1 : nx > p.x + 0.05 ? 1 : p.dir
  p.x = nx
  p.y = ny
  return false
}

function pickTarget(p) {
  const r = Math.random()
  if (r < 0.55) {
    setPath(p, 5 + Math.random() * 85, 86 + Math.random() * 8)
  } else if (r < 0.8) {
    setPath(p, 10 + Math.random() * 75, 35 + Math.random() * 40) // 中场巡视
  } else {
    // 趴卡片
    const els = document.querySelectorAll('.menu-board, .board, .club-card, .send-card, .fortune')
    const seats = []
    els.forEach((el) => {
      const rect = el.getBoundingClientRect()
      if (rect.top > 120 && rect.top < window.innerHeight - 160 && rect.width > 300) seats.push(rect)
    })
    if (seats.length) {
      const s = seats[Math.floor(Math.random() * seats.length)]
      const petH = (100 / window.innerHeight) * 100
      setPath(p, ((s.left + Math.random() * s.width * 0.8) / window.innerWidth) * 100, (s.top / window.innerHeight) * 100 - petH + 2)
      return
    }
    setPath(p, 10 + Math.random() * 75, 86 + Math.random() * 8)
  }
}

function moveToward(p, speed) {
  const dx = p.tx2 - p.x
  const dy = p.ty2 - p.y
  const dist = Math.hypot(dx, dy)
  if (dist < 1) return true
  p.dir = dx < -0.5 ? -1 : dx > 0.5 ? 1 : p.dir
  p.x += (dx / dist) * speed
  p.y += (dy / dist) * speed
  return false
}

function step() {
  if (hidden.value || isNight.value) return

  // ===== 拖拽中的食物:对应的小吃货眼巴巴跟着跑 =====
  if (dragFood.value) {
    const hungry = dragFood.value.type === '🐟' ? cat : pig
    if (!hungry.held) {
      hungry.tx2 = dragFood.value.x
      hungry.ty2 = Math.min(dragFood.value.y + 6, 92)
      const dist = Math.hypot(hungry.tx2 - hungry.x, hungry.ty2 - hungry.y)
      if (dist > 2) {
        hungry.state = 'walk'
        moveToward(hungry, 0.5)
      } else {
        hungry.state = 'idle' // 蹲在食物下面等
      }
    }
  }

  // ===== 猫 =====
  if (!cat.held) {
    if (cat.state === 'chase') {
      cat.tx2 = mouse.x
      cat.ty2 = Math.min(mouse.y, 92)
      if (moveToward(cat, 0.7)) {
        say(cat, '抓到啦!', 1500)
        cat.state = 'idle'
        setTimeout(() => (cat.state = 'walk'), 1500)
      } else if (Date.now() > chaseUntil) {
        cat.state = 'walk'
        pickTarget(cat)
      }
    } else if (cat.state === 'walk') {
      const arrived = stepPath(cat, cat.food ? 0.55 : 0.28)
      if (cat.food && arrived) eatFood(cat)
      else if (arrived) cat.state = Math.random() < 0.5 ? 'idle' : 'walk'
    }
  }

  // ===== 猪(跟随带摆动,像撒欢的小跟班) =====
  if (!pig.held) {
    pig.wobble += 0.15
    if (pig.food) {
      if (stepPath(pig, 0.5)) eatFood(pig)
    } else if (pig.mode === 'follow') {
      pig.tx2 = cat.x - cat.dir * 8
      pig.ty2 = Math.min(cat.y + 3 + Math.sin(pig.wobble) * 1.5, 94)
      const dist = Math.hypot(pig.tx2 - pig.x, pig.ty2 - pig.y)
      if (dist > 2) {
        pig.state = 'walk'
        moveToward(pig, dist > 20 ? 0.5 : 0.24)
      } else if (cat.state !== 'walk' && cat.state !== 'chase') {
        pig.state = 'idle'
      }
    } else if (pig.mode === 'sprint') {
      pig.state = 'walk'
      if (stepPath(pig, 0.55)) pig.mode = 'follow'
    }
  }
}

function maybeChange() {
  if (hidden.value || isNight.value) return
  const r = Math.random()
  if (cat.state === 'walk' && !cat.food && !cat.held) {
    if (r < 0.18) cat.state = 'idle'
    else if (r < 0.26) cat.state = 'sleep'
    else if (r < 0.7) pickTarget(cat)
  } else if (cat.state !== 'chase' && !cat.food && !cat.held && r < 0.55) {
    cat.state = 'walk'
    pickTarget(cat)
  }
  const r2 = Math.random()
  if (pig.mode === 'follow' && !pig.food && !pig.held && r2 < 0.12) {
    pig.mode = 'sprint'
    setPath(pig, 5 + Math.random() * 85, 50 + Math.random() * 44)
    say(pig, '撒欢啦——!', 2000)
    setTimeout(() => (pig.mode = 'follow'), 2600)
  }
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

// ===== 拖拽投喂 =====
function startFoodDrag(e) {
  const type = Math.random() < 0.5 ? '🐟' : '🍰'
  dragFood.value = { type, x: px2vw(e.clientX), y: px2vh(e.clientY), moved: false }
  // 对应的小吃货立刻注意到食物
  const target = type === '🐟' ? cat : pig
  if (!target.held && target.state !== 'sleep') {
    target.state = 'walk'
    say(target, type === '🐟' ? '鱼!是鱼!' : '蛋糕!给我给我!', 2000)
  }
}

function px2vw(px) { return (px / window.innerWidth) * 100 }
function px2vh(px) { return (px / window.innerHeight) * 100 }

function onPointerMove(e) {
  const x = px2vw(e.clientX), y = px2vh(e.clientY)
  if (dragFood.value) {
    dragFood.value.x = x
    dragFood.value.y = y
    dragFood.value.moved = true
  }
  // 按下后移动超过阈值才算"抓起来"
  if (pendingGrab && !heldPet) {
    if (Math.hypot(x - pendingGrab.x, y - pendingGrab.y) > 1.5) {
      const p = pendingGrab.pet
      heldPet = p
      p.held = true
      p.state = 'held'
      p.path = null
      grabOffset = { x: p.x - pendingGrab.x, y: p.y - pendingGrab.y }
      say(p, (p === cat ? CAT_HELD : PIG_HELD)[Math.floor(Math.random() * 3)], 2000)
    }
  }
  if (heldPet) {
    heldPet.x = x + grabOffset.x
    heldPet.y = Math.max(12, y + grabOffset.y)
    // 猫被抓时,猪咪会着急
    if (heldPet === cat && !warnedPig) {
      warnedPig = true
      say(pig, '放下她!!', 2000)
    }
    onMouse(e)
    return
  }
  onMouse(e)
}

function onPointerUp(e) {
  if (dragFood.value) {
    const f = dragFood.value
    dragFood.value = null
    const food = {
      id: ++foodSeq,
      type: f.type,
      x: f.moved ? Math.max(3, Math.min(95, f.x)) : 20 + Math.random() * 55,
      y: f.moved ? Math.max(15, Math.min(93, f.y)) : 76 + Math.random() * 14,
    }
    foods.value.push(food)
    const target = f.type === '🐟' ? cat : pig
    target.food = food
    target.state = 'walk'
    setPath(target, food.x, food.y)
    say(target, f.type === '🐟' ? '是鱼!冲!' : '蛋糕!我的!', 2000)
  }
  if (heldPet) {
    const p = heldPet
    heldPet = null
    warnedPig = false
    p.held = false
    // 松手:就地停在当前位置,不再移动
    p.state = 'idle'
    p.path = null
    p.anim = 'poke'
    say(p, p === cat ? '下次轻点放!' : '落地啦!', 1800)
    setTimeout(() => (p.anim = ''), 500)
  }
  pendingGrab = null
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

// ===== 抓取小人(按下先记位置,拖动才生效) =====
function grabPet(p, e) {
  if (isNight.value) return
  pendingGrab = { pet: p, x: px2vw(e.clientX), y: px2vh(e.clientY) }
}

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

function pokeCat() {
  if (cat.held) return
  if (cat.state === 'sleep') {
    say(cat, '唔……再睡五分钟', 2200)
    return
  }
  say(cat, CAT_POKED[Math.floor(Math.random() * CAT_POKED.length)], 2500)
  cat.anim = 'poke'
  setTimeout(() => (cat.anim = ''), 900)
}

function pokePig() {
  if (pig.held) return
  say(pig, PIG_POKED[Math.floor(Math.random() * PIG_POKED.length)], 2200)
  pig.anim = 'poke'
  setTimeout(() => (pig.anim = ''), 900)
  pig.mode = 'sprint'
  setPath(pig, pig.x > cat.x ? Math.min(pig.x + 20, 90) : Math.max(pig.x - 20, 5), pig.y, false)
  setTimeout(() => (pig.mode = 'follow'), 1500)
}

function onMouse(e) {
  mouse = { x: px2vw(e.clientX), y: px2vh(e.clientY) }
  const now = Date.now()
  mouseTrail.push({ x: mouse.x, t: now })
  mouseTrail = mouseTrail.filter((p) => now - p.t < 1200)
  let turns = 0
  for (let i = 2; i < mouseTrail.length; i++) {
    const d1 = mouseTrail[i - 1].x - mouseTrail[i - 2].x
    const d2 = mouseTrail[i].x - mouseTrail[i - 1].x
    if (d1 * d2 < 0) turns++
  }
  if (turns >= 6 && cat.state === 'walk' && !cat.food && !cat.held && !isNight.value) {
    cat.state = 'chase'
    chaseUntil = now + 3000
    say(cat, '别跑!', 1500)
    mouseTrail = []
  }
}

function onScroll() {
  ;[cat, pig].forEach((p) => {
    if (p.y < 80 && !p.food && !p.held) {
      setPath(p, p.x, 86 + Math.random() * 8)
      if (p.state !== 'chase') p.state = 'walk'
    }
  })
}

function checkNight() {
  const h = new Date().getHours()
  const night = h >= 0 && h < 7
  if (night && !isNight.value) {
    Object.assign(cat, { x: 8, y: 92, state: 'sleep', food: null, path: null })
    Object.assign(pig, { x: 15, y: 94, state: 'sleep', food: null, path: null, mode: 'follow' })
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
  window.addEventListener('pointermove', onPointerMove, { passive: true })
  window.addEventListener('pointerup', onPointerUp)
  window.addEventListener('scroll', onScroll, { passive: true })
  tick = setInterval(step, 50)
  stateTimer = setInterval(maybeChange, 4000)
  bubbleTimer = setInterval(maybeBubble, 8000)
  nightTimer = setInterval(checkNight, 30000)
  checkNight()
  pickTarget(cat)
})

onUnmounted(() => {
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
  window.removeEventListener('scroll', onScroll)
  clearInterval(tick)
  clearInterval(stateTimer)
  clearInterval(bubbleTimer)
  clearInterval(nightTimer)
  clearTimeout(crackerTimer)
})
</script>

<template>
  <!-- 左侧按钮堆:弹幕 / 投喂 / 看板娘开关 -->
  <button class="pet-toggle" @click="hidden ? show() : hide()">
    {{ hidden ? '🐾 召回看板娘' : '🐾 看板娘休息' }}
  </button>

  <template v-if="!hidden">
    <!-- 投喂按钮:拖出去丢食物 -->
    <button class="feed-btn" @pointerdown.prevent="startFoodDrag">
      🍰 按住拖出去投喂{{ feeds ? ` · 已被喂${feeds}次` : '' }}
    </button>

    <!-- 拖拽中的食物 -->
    <span v-if="dragFood" class="food ghost" :style="{ left: dragFood.x + 'vw', top: dragFood.y + 'vh' }">
      {{ dragFood.type }}
    </span>

    <!-- 落地的食物 -->
    <span v-for="f in foods" :key="f.id" class="food" :style="{ left: f.x + 'vw', top: f.y + 'vh' }">
      {{ f.type }}
    </span>

    <!-- 猪咪拱出的鱼干 -->
    <button
      v-if="cracker"
      class="cracker"
      :style="{ left: cracker.x + 'vw', top: cracker.y + 'vh' }"
      @click="pickCracker"
    >
      🐟✨
    </button>

    <!-- 看板猫(可拖拽) -->
    <div class="pet" :class="{ held: cat.held }" :style="{ left: cat.x + 'vw', top: cat.y + 'vh' }">
      <transition name="bub">
        <span v-if="cat.bubble" class="pet-bubble">{{ cat.bubble }}</span>
      </transition>
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
          @pointerdown.prevent="grabPet(cat, $event)"
          @click="pokeCat"
          @keyup.enter="pokeCat"
        />
      </span>
    </div>

    <!-- 猪咪跟屁虫(可拖拽) -->
    <div class="pet pig" :class="{ held: pig.held }" :style="{ left: pig.x + 'vw', top: pig.y + 'vh' }">
      <transition name="bub">
        <span v-if="pig.bubble" class="pet-bubble pig-bubble">{{ pig.bubble }}</span>
      </transition>
      <span v-if="pig.state === 'sleep'" class="zzz">💤</span>
      <span class="flipper" :class="{ flip: pig.dir < 0 }">
        <img
          src="/pets/pigmi.png"
          alt="猪咪"
          class="pet-img pigmi-img"
          :class="pig.anim || pig.state"
          draggable="false"
          role="button"
          tabindex="0"
          @pointerdown.prevent="grabPet(pig, $event)"
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
  transition: left 0.05s linear, top 0.05s linear;
  transform: translateY(-100%);
}

.pet.held {
  transition: none; /* 拖拽时跟手 */
  z-index: 58;
}

.pet.held .pet-img {
  cursor: grabbing;
  animation: heldShake 0.25s ease-in-out infinite;
}

@keyframes heldShake {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}

.flipper {
  display: inline-block;
}

.flipper.flip {
  transform: scaleX(-1);
}

.pet-img {
  pointer-events: auto;
  cursor: grab;
  user-select: none;
  touch-action: none; /* 手机上可拖不滚动 */
}

.alice-img { height: 88px; }
.pigmi-img { height: 60px; }

.pet-img.walk {
  animation: walkBob 0.4s ease-in-out infinite;
}
@keyframes walkBob {
  0%, 100% { transform: translateY(0) rotate(-3deg); }
  50% { transform: translateY(-6px) rotate(3deg); }
}

.pet-img.chase {
  animation: walkBob 0.25s ease-in-out infinite;
}

.pet-img.idle {
  animation: breathe 2.4s ease-in-out infinite;
}
@keyframes breathe {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.95); }
}

.pet-img.sleep {
  filter: brightness(0.85) saturate(0.85);
  transform: scaleY(0.85);
  transform-origin: bottom;
}

.pet-img.poke {
  animation: wiggle 0.3s ease-in-out infinite;
}
@keyframes wiggle {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}

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

.pet-toggle {
  position: fixed;
  left: 20px;
  bottom: 118px;
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

.pet-toggle:hover {
  background: var(--pink-pale);
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
  cursor: grab;
  box-shadow: var(--shadow);
  transition: transform 0.15s;
  touch-action: none;
  user-select: none;
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

.food.ghost {
  z-index: 65;
  font-size: 32px;
  animation: none;
  opacity: 0.9;
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
}

@media (prefers-reduced-motion: reduce) {
  .pet-img, .cracker {
    animation: none !important;
  }
}
</style>

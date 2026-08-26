<script setup>
// 猫咖看板娘引擎 v4 —— 曲线漫游 / 拖拽投喂 / 可抓取小人 / 小剧场 / 追鼠标 / 捡鱼干 / 深夜作息
import { ref, reactive, onMounted, onUnmounted } from 'vue'

const hidden = ref(localStorage.getItem('catcafe_pet_hide') === '1')
const isNight = ref(false)
const laserOn = ref(false)   // 激光笔模式
const portal = reactive({ active: false, x: 50, y: 50 }) // 赛博裂缝
const pigGone = ref(false)   // 猪咪被吞进裂缝
let catCommented = false     // 猫已评论过裂缝
let portalTimer = null
let suckTimers = []
const gameMode = ref(false)  // 拍猪咪游戏中,看板娘回避
const loveHearts = ref([])   // 贴贴爱心 {id, x, y}
let togetherUntil = 0        // 贴贴截止时间

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

let tick, stateTimer, bubbleTimer, nightTimer, crackerTimer, toastTimer, heartSeq = 0
let bubbleClear = { cat: null, pig: null }

function spawnHeart(x, y) {
  const id = ++heartSeq
  loveHearts.value.push({ id, x, y })
  setTimeout(() => (loveHearts.value = loveHearts.value.filter((h) => h.id !== id)), 1600)
}

// 贴贴检测:拖动任意一只松手时,若两只靠得近 → 贴贴 8 秒
function tryTogether() {
  const dist = Math.hypot(cat.x - pig.x, cat.y - pig.y)
  if (dist < 9 && pig.mode !== 'together' && !isNight.value) {
    localStorage.setItem('catcafe_ach_together', '1')
    pig.mode = 'together'
    togetherUntil = Date.now() + 8000
    say(cat, '……靠这么近干嘛啦', 2200)
    setTimeout(() => say(pig, '就要贴贴!', 2000), 1300)
    spawnHeart((cat.x + pig.x) / 2, Math.min(cat.y, pig.y) - 10)
  }
}

// ===== 赛博裂缝彩蛋 =====
function schedulePortal() {
  clearTimeout(portalTimer)
  portalTimer = setTimeout(spawnPortal, 90000 + Math.random() * 180000) // 1.5~4.5分钟
}

function spawnPortal() {
  if (hidden.value || isNight.value || gameMode.value || portal.active) {
    schedulePortal()
    return
  }
  portal.x = 18 + Math.random() * 64
  portal.y = 22 + Math.random() * 56
  portal.active = true
  catCommented = false
  setTimeout(closePortal, 60000) // 60秒后消失
}

function closePortal() {
  portal.active = false
  if (pigGone.value) {
    // 猪咪从裂缝里爬回来
    pigGone.value = false
    pig.x = portal.x
    pig.y = portal.y
    pig.state = 'idle'
    say(pig, '我刚才……去了哪里?', 3200)
    say(cat, '猪咪!你回来了!', 3200)
  }
  schedulePortal()
}

function enterGlitch() {
  localStorage.setItem('catcafe_ach_glitch', '1')
  // 用 origin 拼接:网页是 https://alicefans.cn,APP 内是 https://localhost(glitch.html 已打包进 APP),
  // 不能用 '/glitch.html' —— 会被 APP 里的 <base> 带去线上
  location.href = location.origin + '/glitch.html'
}

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
  if (hidden.value || isNight.value || gameMode.value) return

  // ===== 裂缝互动 =====
  if (portal.active) {
    // 猪咪路过被吞
    if (!pigGone.value && !pig.held) {
      const d = Math.hypot(pig.x - portal.x, pig.y - portal.y)
      if (d < 8) {
        pigGone.value = true
        pig.state = 'sucked'
        say(pig, '呀啊啊——!', 1500)
        // 猫的反应时间线
        suckTimers.push(setTimeout(() => say(cat, '咦?猪咪呢?', 2400), 2000))
        suckTimers.push(setTimeout(() => say(cat, '猪咪不见了!你看到它了吗?快救救猪咪!', 4000), 4600))
        suckTimers.push(setTimeout(() => {
          // 猫跑到裂缝边守着
          cat.state = 'walk'
          setPath(cat, portal.x + 8, Math.min(portal.y + 6, 92))
        }, 4800))
      }
    }
    // 猫路过吐槽(只一次)
    if (!pigGone.value && !catCommented) {
      const d = Math.hypot(cat.x - portal.x, cat.y - portal.y)
      if (d < 13) {
        catCommented = true
        say(cat, '这是什么??怪怪的……', 2600)
      }
    }
  }

  // ===== 激光笔:猫全速追红点 =====
  if (laserOn.value && !cat.held) {
    cat.state = 'chase'
    cat.tx2 = mouse.x
    cat.ty2 = Math.min(mouse.y, 92)
    const dist = Math.hypot(cat.tx2 - cat.x, cat.ty2 - cat.y)
    if (dist > 3) {
      moveToward(cat, 0.8)
      cat.pounced = false
    } else if (!cat.pounced) {
      cat.pounced = true
      cat.anim = 'poke'
      say(cat, '按住啦!', 1500)
      setTimeout(() => (cat.anim = ''), 600)
    }
  } else if (cat.state === 'chase' && !laserOn.value) {
    cat.state = 'walk'
    pickTarget(cat)
  }

  // ===== 食物最高优先级:场上有食物而吃货没锁定时,立刻锁定冲过去 =====
  for (const f of foods.value) {
    const hungry = f.type === '🐟' ? cat : pig
    if (!hungry.held && !hungry.food) {
      hungry.food = f
      hungry.state = 'walk'
      setPath(hungry, f.x, f.y, false)
    }
  }

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
  if (!cat.held && !laserOn.value) {
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
      const arrived = stepPath(cat, cat.food ? 0.8 : 0.28) // 干饭速度拉满
      if (cat.food && arrived) eatFood(cat)
      else if (arrived) cat.state = Math.random() < 0.5 ? 'idle' : 'walk'
    }
  }

  // ===== 猪(跟随带摆动;贴贴模式贴紧猫) =====
  if (!pig.held) {
    pig.wobble += 0.15
    if (pig.food) {
      if (stepPath(pig, 0.75)) eatFood(pig) // 干饭速度拉满
    } else if (pig.mode === 'together') {
      // 贴贴:紧贴猫身边同步走
      pig.tx2 = cat.x - cat.dir * 6
      pig.ty2 = cat.y
      pig.state = cat.state === 'walk' || cat.state === 'chase' ? 'walk' : 'idle'
      moveToward(pig, 0.3)
      if (Date.now() > togetherUntil) {
        pig.mode = 'follow'
        say(pig, '嘿嘿,贴贴~', 2000)
      } else if (Math.random() < 0.02) {
        spawnHeart((cat.x + pig.x) / 2, Math.min(cat.y, pig.y) - 10)
      }
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
    if (heldPet === cat && !warnedPig) {
      warnedPig = true
      say(pig, '放下她!!', 2000)
    }
    return
  }
  // 触屏滑动(滚动页面)不参与"晃鼠标逗猫"判定
  if (e.pointerType === 'touch') {
    mouse = { x, y } // 但更新位置,激光笔在触屏可用
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
    p.path = null
    if (p.food) {
      // 场上有它的食物:放下立刻直奔食物
      p.state = 'walk'
      setPath(p, p.food.x, p.food.y, false)
      say(p, p === cat ? '鱼干等我!' : '蛋糕等我!', 1500)
    } else {
      // 没有食物:就地停住
      p.state = 'idle'
      p.anim = 'poke'
      say(p, p === cat ? '下次轻点放!' : '落地啦!', 1800)
      setTimeout(() => (p.anim = ''), 500)
    }
    tryTogether() // 松手检测贴贴
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

// 菜单事件:投喂(进入放置模式) / 看板娘开关 / 激光笔开关
function onMenuFeed() {
  if (hidden.value) return
  const type = Math.random() < 0.5 ? '🐟' : '🍰'
  dragFood.value = { type, x: mouse.x, y: mouse.y, moved: true }
  showToast('点击任意位置放下食物')
  // 下一次点击 = 放置
  setTimeout(() => {
    window.addEventListener('pointerdown', placeFoodOnce, { once: true })
  }, 50)
}

function placeFoodOnce(e) {
  if (!dragFood.value) return
  const f = dragFood.value
  dragFood.value = null
  const food = {
    id: ++foodSeq,
    type: f.type,
    x: Math.max(3, Math.min(95, px2vw(e.clientX))),
    y: Math.max(15, Math.min(93, px2vh(e.clientY))),
  }
  foods.value.push(food)
  const target = f.type === '🐟' ? cat : pig
  target.food = food
  target.state = 'walk'
  setPath(target, food.x, food.y)
  say(target, f.type === '🐟' ? '是鱼!冲!' : '蛋糕!我的!', 2000)
}

function onMenuPet() {
  hidden.value ? show() : hide()
}

function onMenuLaser() {
  laserOn.value = !laserOn.value
  window.dispatchEvent(new CustomEvent('menu-state', { detail: { laser: laserOn.value } }))
}

function onWhack(e) {
  gameMode.value = e.detail === 'start'
  if (gameMode.value) {
    // 游戏开始:两只回避到左下角
    cat.path = null
    pig.path = null
    cat.state = 'idle'
    pig.state = 'idle'
  } else {
    cat.state = 'walk'
    pickTarget(cat)
  }
}

function onWeather(e) {
  if (isNight.value || hidden.value) return
  if (e.detail === 'rain') {
    // 下雨:躲到公告板下面
    const board = document.querySelector('.board')
    if (board) {
      const r = board.getBoundingClientRect()
      cat.state = 'walk'
      setPath(cat, (r.left / window.innerWidth) * 100 + 10, Math.min((r.bottom / window.innerHeight) * 100 + 6, 92))
      say(cat, '下雨了,躲一躲', 2500)
    }
  } else if (e.detail === 'petals') {
    say(cat, '花瓣雨,好漂亮', 2500)
    setTimeout(() => say(pig, '哇——', 1800), 1400)
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
  window.addEventListener('whack', onWhack)
  window.addEventListener('cafe-weather', onWeather)
  window.addEventListener('menu-feed', onMenuFeed)
  window.addEventListener('menu-pet', onMenuPet)
  window.addEventListener('menu-laser', onMenuLaser)
  tick = setInterval(step, 50)
  stateTimer = setInterval(maybeChange, 4000)
  bubbleTimer = setInterval(maybeBubble, 8000)
  nightTimer = setInterval(checkNight, 30000)
  checkNight()
  pickTarget(cat)
  schedulePortal() // 启动裂缝彩蛋调度
})

onUnmounted(() => {
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('whack', onWhack)
  window.removeEventListener('cafe-weather', onWeather)
  window.removeEventListener('menu-feed', onMenuFeed)
  window.removeEventListener('menu-pet', onMenuPet)
  window.removeEventListener('menu-laser', onMenuLaser)
  clearInterval(tick)
  clearInterval(stateTimer)
  clearInterval(bubbleTimer)
  clearInterval(nightTimer)
  clearTimeout(crackerTimer)
  clearTimeout(portalTimer)
  suckTimers.forEach(clearTimeout)
})
</script>

<template>
  <template v-if="!hidden">
    <!-- 激光红点 -->
    <span v-if="laserOn" class="laser-dot" :style="{ left: mouse.x + 'vw', top: mouse.y + 'vh' }" />

    <!-- 贴贴爱心 -->
    <span v-for="h in loveHearts" :key="h.id" class="love-heart" :style="{ left: h.x + 'vw', top: h.y + 'vh' }">💗</span>

    <!-- 拖拽/放置中的食物 -->
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

    <!-- 看板猫(可拖拽);拍猪咪游戏时隐藏 -->
    <div v-show="!gameMode" class="pet" :class="{ held: cat.held }" :style="{ left: cat.x + 'vw', top: cat.y + 'vh' }">
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

    <!-- 赛博裂缝(点了进去) -->
    <div
      v-if="portal.active"
      class="portal"
      :style="{ left: portal.x + 'vw', top: portal.y + 'vh' }"
      role="button"
      aria-label="神秘裂缝"
      @click="enterGlitch"
    >
      <div class="portal-ring" />
      <div class="portal-ring r2" />
      <div class="portal-core" />
      <span class="portal-label">???</span>
    </div>

    <!-- 猪咪跟屁虫(可拖拽) -->
    <div v-show="!gameMode && !pigGone" class="pet pig" :class="{ held: pig.held, sucked: pig.state === 'sucked' }" :style="{ left: pig.x + 'vw', top: pig.y + 'vh' }">
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





.laser-dot {
  position: fixed;
  z-index: 90;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #ff2d55;
  box-shadow: 0 0 12px 4px rgba(255, 45, 85, 0.55);
  transform: translate(-50%, -50%);
  pointer-events: none;
  transition: left 0.03s linear, top 0.03s linear;
}

.love-heart {
  position: fixed;
  z-index: 57;
  font-size: 20px;
  transform: translate(-50%, -100%);
  pointer-events: none;
  animation: loveUp 1.6s ease-out forwards;
}

@keyframes loveUp {
  0% { opacity: 0; transform: translate(-50%, -90%) scale(0.5); }
  20% { opacity: 1; }
  100% { opacity: 0; transform: translate(-50%, -260%) scale(1.2); }
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

/* ===== 赛博裂缝 ===== */
.portal {
  position: fixed;
  z-index: 53;
  width: 96px;
  height: 96px;
  transform: translate(-50%, -50%);
  cursor: pointer;
  animation: portalIn 0.5s ease;
}

@keyframes portalIn {
  from { transform: translate(-50%, -50%) scale(0); }
  to { transform: translate(-50%, -50%) scale(1); }
}

.portal-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  background: conic-gradient(#ff5da2, #c5a3ff, #1a0b1e, #ff5da2);
  animation: spin 2.4s linear infinite;
  filter: blur(1px);
}

.portal-ring.r2 {
  inset: 12px;
  background: conic-gradient(#1a0b1e, #ff9ecf, #c5a3ff, #1a0b1e);
  animation: spinR 1.7s linear infinite;
}

.portal-core {
  position: absolute;
  inset: 30px;
  border-radius: 50%;
  background: #07070c;
  box-shadow: 0 0 24px 6px rgba(255, 93, 162, 0.65), inset 0 0 12px #000;
  animation: corePulse 1.2s ease-in-out infinite;
}

.portal-label {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  color: #ff9ecf;
  font-size: 12px;
  font-family: 'JetBrains Mono', monospace;
  animation: flick 1.6s steps(1) infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }
@keyframes spinR { to { transform: rotate(-360deg); } }
@keyframes corePulse {
  0%, 100% { box-shadow: 0 0 24px 6px rgba(255, 93, 162, 0.65), inset 0 0 12px #000; }
  50% { box-shadow: 0 0 36px 12px rgba(197, 163, 255, 0.5), inset 0 0 16px #000; }
}
@keyframes flick {
  0%, 80% { opacity: 1; }
  90% { opacity: 0.2; }
  100% { opacity: 1; }
}

/* 猪被吞:旋转缩小进洞 */
.pet.sucked .pet-img {
  animation: suckedIn 0.8s ease-in forwards !important;
}

@keyframes suckedIn {
  to { transform: scale(0) rotate(720deg); }
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
  .alice-img { height: 56px; }
  .pigmi-img { height: 42px; }
}

@media (prefers-reduced-motion: reduce) {
  .pet-img, .cracker {
    animation: none !important;
  }
}
</style>

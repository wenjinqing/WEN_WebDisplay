<script setup>
import { ref, computed, onMounted } from 'vue'
import CatFace from './CatFace.vue'
import PigmiFace from './PigmiFace.vue'
import PawPrint from './PawPrint.vue'
import HandDoodle from './HandDoodle.vue'
import { site } from '../data.js'

// 季节限定装饰:冬雪/春樱/夏瓜/秋叶
const seasonal = computed(() => {
  const m = new Date().getMonth() + 1
  if (m === 12 || m <= 2) return '❄️'
  if (m <= 5) return '🌸'
  if (m <= 8) return '🍉'
  return '🍂'
})

// ===== 撸猫 =====
const pets = ref(0)
const meow = ref('')
const hearts = ref([])
let heartSeq = 0
let meowTimer = null

const meows = ['喵!', '呼噜呼噜…', '再摸一下嘛~', '喵呜♪', '蹭蹭你', '尾巴翘起来了!']

onMounted(async () => {
  try {
    const res = await fetch('/api/pet')
    pets.value = (await res.json()).pets || 0
  } catch {
    /* 静默 */
  }
})

async function petCat() {
  meow.value = meows[Math.floor(Math.random() * meows.length)]
  clearTimeout(meowTimer)
  meowTimer = setTimeout(() => (meow.value = ''), 1500)
  const id = ++heartSeq
  hearts.value.push(id)
  setTimeout(() => (hearts.value = hearts.value.filter((h) => h !== id)), 1200)
  pets.value++
  try {
    const res = await fetch('/api/pet', { method: 'POST' })
    pets.value = (await res.json()).pets
  } catch {
    /* 本地计数即可 */
  }
}
</script>

<template>
  <section id="home" class="hero">
    <!-- 梦幻渐变天幕 + 漂浮涂鸦 -->
    <div class="paws" aria-hidden="true">
      <HandDoodle kind="cloud" :size="110" class="paw d1" />
      <HandDoodle kind="cloud" :size="72" class="paw d2" />
      <HandDoodle kind="cloud" :size="56" class="paw d6" />
      <HandDoodle kind="star" :size="26" class="paw d3" />
      <HandDoodle kind="star" :size="18" class="paw d4" />
      <HandDoodle kind="heart" :size="30" class="paw d5" />
      <PawPrint :size="30" class="paw p2" :rotate="18" />
      <span class="paw seasonal s1">{{ seasonal }}</span>
      <span class="paw seasonal s2">{{ seasonal }}</span>
    </div>

    <div class="container hero-inner">
      <div class="hero-text">
        <div class="badges">
          <span class="clay-chip font-cute">☕ 营业中 · OPEN</span>
          <span class="clay-chip alt font-cute">店主现在:{{ site.authorStatus || '赶稿中 ✍️' }}</span>
        </div>

        <h1 class="font-cute title">
          <span class="title-line">爱丽丝的</span>
          <span class="title-line pop">小涩猫咖啡厅</span>
        </h1>

        <p class="slogan">{{ site.slogan }}</p>
        <p class="sub">
          这里是 <b>{{ site.author }}</b> 的非官方粉丝后援会
        </p>
        <div class="actions">
          <a href="#novels" class="btn-clay primary">📖 去看小说菜单</a>
          <a href="/go/pixiv.html" target="_blank" rel="noopener" class="btn-clay ghost">
            作者主页 ↗
          </a>
          <HandDoodle kind="arrow" :size="56" class="cta-arrow" />
        </div>
      </div>

      <div class="hero-mascots">
        <div class="blob" aria-hidden="true" />
        <div class="cat-zone" role="button" tabindex="0" aria-label="撸猫" @click="petCat" @keyup.enter="petCat">
          <transition name="meow-pop">
            <span v-if="meow" class="meow-bubble font-cute">{{ meow }}</span>
          </transition>
          <span v-for="h in hearts" :key="h" class="fly-heart">💗</span>
          <CatFace :size="220" class="cat" />
        </div>
        <PigmiFace :size="110" class="pigmi" />
        <p class="pet-count">🫳 猫猫已被撸 {{ pets }} 次,你也来一下?</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
/* ===== Claymorphism × 手绘 × 梦幻少女风 ===== */
.hero {
  position: relative;
  padding: 88px 0 96px;
  /* 梦幻粉紫渐变天幕 */
  background:
    radial-gradient(ellipse 60% 45% at 85% 20%, rgba(230, 230, 250, 0.8), transparent 70%),
    radial-gradient(ellipse 50% 40% at 10% 80%, rgba(253, 188, 180, 0.45), transparent 70%),
    linear-gradient(165deg, #fff6f8 0%, #ffe9f0 55%, #f6efff 100%);
  overflow: hidden;
}

.hero-inner {
  display: flex;
  align-items: center;
  gap: 48px;
}

.hero-text {
  flex: 1.2;
}

/* 黏土徽章:厚边 + 内外双阴影 */
.clay-chip {
  display: inline-block;
  padding: 7px 20px;
  border-radius: 18px;
  background: #fff;
  border: 3px solid var(--pink-soft);
  color: var(--pink-deep);
  font-size: 0.9rem;
  box-shadow:
    inset -2px -2px 6px rgba(255, 194, 212, 0.6),
    4px 4px 10px rgba(233, 93, 127, 0.15);
}

.clay-chip.alt {
  background: var(--pink-pale);
  color: var(--ink);
}

.badges {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 26px;
}

/* 标题:两行错位 + 描边投影层 */
.title {
  margin-bottom: 18px;
  line-height: 1.2;
}

.title-line {
  display: block;
  font-size: clamp(1.9rem, 5vw, 3rem);
  color: var(--ink);
  letter-spacing: 2px;
}

.title-line.pop {
  font-size: clamp(2.4rem, 7.2vw, 4.2rem);
  color: var(--pink-deep);
  text-shadow:
    3px 3px 0 #ffc2d4,
    6px 6px 0 rgba(233, 93, 127, 0.18);
  letter-spacing: 4px;
  margin-top: 4px;
}

.slogan {
  font-size: 1.15rem;
  color: var(--pink-deep);
  margin-bottom: 8px;
}

.sub {
  color: var(--muted);
  margin-bottom: 34px;
}

.sub b {
  color: var(--pink-deep);
}

/* 黏土按钮:厚边框 + 双层阴影 + 按压回弹 */
.actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  position: relative;
}

.btn-clay {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 14px 30px;
  border-radius: 22px;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  border: 3px solid transparent;
  transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.2s ease;
}

.btn-clay.primary {
  background: var(--pink);
  color: #fff;
  border-color: #d14a6e;
  box-shadow:
    inset -3px -3px 8px rgba(209, 74, 110, 0.55),
    5px 6px 0 rgba(209, 74, 110, 0.9);
}

.btn-clay.primary:hover {
  transform: translate(-2px, -3px);
  box-shadow:
    inset -3px -3px 8px rgba(209, 74, 110, 0.55),
    7px 9px 0 rgba(209, 74, 110, 0.9);
}

.btn-clay.primary:active {
  transform: translate(3px, 4px);
  box-shadow:
    inset -3px -3px 8px rgba(209, 74, 110, 0.55),
    1px 2px 0 rgba(209, 74, 110, 0.9);
}

.btn-clay.ghost {
  background: #fff;
  color: var(--pink-deep);
  border-color: var(--pink-soft);
  box-shadow:
    inset -2px -2px 6px rgba(255, 194, 212, 0.6),
    5px 6px 0 var(--pink-pale);
}

.btn-clay.ghost:hover {
  transform: translate(-2px, -3px);
  background: var(--pink-pale);
  box-shadow:
    inset -2px -2px 6px rgba(255, 194, 212, 0.6),
    7px 9px 0 var(--pink-soft);
}

.cta-arrow {
  position: absolute;
  right: -64px;
  top: -34px;
  transform: rotate(8deg);
}

/* ===== 吉祥物区:梦幻色团 + 软影 ===== */
.hero-mascots {
  position: relative;
  flex-shrink: 0;
  text-align: center;
}

.blob {
  position: absolute;
  inset: -30px -40px;
  background: radial-gradient(ellipse at 45% 40%, rgba(255, 255, 255, 0.9), rgba(255, 228, 236, 0.5) 70%, transparent);
  border-radius: 46% 54% 58% 42% / 52% 44% 56% 48%;
  z-index: -1;
  animation: blobMorph 9s ease-in-out infinite;
}

@keyframes blobMorph {
  0%, 100% { border-radius: 46% 54% 58% 42% / 52% 44% 56% 48%; }
  50% { border-radius: 54% 46% 42% 58% / 46% 56% 44% 54%; }
}

.cat-zone {
  position: relative;
  cursor: pointer;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.cat-zone:focus-visible {
  outline: 3px solid var(--pink);
  border-radius: 20px;
}

.cat {
  display: block;
  filter: drop-shadow(0 10px 18px rgba(233, 93, 127, 0.25));
  animation: bob 3.5s ease-in-out infinite;
}

.pigmi {
  position: absolute;
  right: -24px;
  bottom: -28px;
  filter: drop-shadow(0 6px 12px rgba(233, 93, 127, 0.22));
  animation: bob 3.5s ease-in-out 0.8s infinite;
}

@keyframes bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.meow-bubble {
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  border: 2px solid var(--pink-soft);
  color: var(--pink-deep);
  border-radius: 16px 16px 16px 4px;
  padding: 4px 14px;
  font-size: 0.95rem;
  white-space: nowrap;
  z-index: 2;
  box-shadow: var(--shadow);
}

.meow-pop-enter-active { transition: all 0.2s ease; }
.meow-pop-leave-active { transition: all 0.4s ease; }
.meow-pop-enter-from { opacity: 0; transform: translateX(-50%) translateY(8px) scale(0.8); }
.meow-pop-leave-to { opacity: 0; transform: translateX(-50%) translateY(-8px); }

.fly-heart {
  position: absolute;
  left: 50%;
  top: 40%;
  font-size: 22px;
  pointer-events: none;
  animation: flyUp 1.2s ease-out forwards;
}

@keyframes flyUp {
  0% { opacity: 1; transform: translate(-50%, 0) scale(0.6); }
  100% { opacity: 0; transform: translate(-50%, -90px) scale(1.3); }
}

.pet-count {
  position: relative;
  z-index: 2;
  display: inline-block;
  margin-top: 14px;
  font-size: 0.85rem;
  color: var(--muted);
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 4px 16px;
}

/* 涂鸦漂浮 */
.paws .paw {
  position: absolute;
  opacity: 0.85;
}

.d1 { top: 4%; left: 46%; animation: drift 9s ease-in-out infinite; }
.d2 { bottom: 16%; left: 6%; animation: drift 11s ease-in-out 1.5s infinite; }
.d3 { top: 18%; right: 8%; animation: drift 7s ease-in-out 0.5s infinite; }
.d4 { top: 44%; left: 24%; animation: drift 8s ease-in-out 2s infinite; }
.d5 { bottom: 26%; right: 30%; animation: drift 10s ease-in-out 1s infinite; }
.d6 { top: 10%; right: 34%; animation: drift 12s ease-in-out 0.3s infinite; }
.p2 { top: 62%; left: 44%; opacity: 0.5; }

.seasonal {
  position: absolute;
  font-size: 26px;
  opacity: 0.75;
}

.s1 { bottom: 30%; left: 8%; animation: drift 8s ease-in-out 0.8s infinite; }
.s2 { top: 5%; right: 44%; font-size: 20px; animation: drift 12s ease-in-out 2.2s infinite; }

@keyframes drift {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-12px) rotate(4deg); }
}

@media (max-width: 720px) {
  .hero {
    padding: 56px 0 64px;
  }
  .hero-inner {
    flex-direction: column-reverse;
    text-align: center;
    gap: 28px;
  }
  .badges, .actions {
    justify-content: center;
  }
  .hero-mascots {
    margin: 0 auto;
  }
  .paws .paw {
    display: none;
  }
  .cta-arrow {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .cat, .pigmi, .paws .paw, .blob { animation: none; }
  .btn-clay { transition: none; }
}
</style>

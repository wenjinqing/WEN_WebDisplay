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
const hearts = ref([]) // 飘出的爱心 id 列表
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
  // 随机猫叫气泡
  meow.value = meows[Math.floor(Math.random() * meows.length)]
  clearTimeout(meowTimer)
  meowTimer = setTimeout(() => (meow.value = ''), 1500)
  // 飘爱心
  const id = ++heartSeq
  hearts.value.push(id)
  setTimeout(() => (hearts.value = hearts.value.filter((h) => h !== id)), 1200)
  // 上报计数
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
    <!-- 漂浮的猫爪 + 手绘涂鸦背景 -->
    <div class="paws" aria-hidden="true">
      <PawPrint :size="44" class="paw p1" :rotate="-24" />
      <PawPrint :size="30" class="paw p2" :rotate="18" />
      <PawPrint :size="56" class="paw p3" :rotate="40" />
      <PawPrint :size="26" class="paw p4" :rotate="-10" />
      <HandDoodle kind="cloud" :size="90" class="paw d1" />
      <HandDoodle kind="cloud" :size="64" class="paw d2" />
      <HandDoodle kind="star" :size="26" class="paw d3" />
      <HandDoodle kind="star" :size="18" class="paw d4" />
      <HandDoodle kind="heart" :size="30" class="paw d5" />
      <span class="paw seasonal s1" aria-hidden="true">{{ seasonal }}</span>
      <span class="paw seasonal s2" aria-hidden="true">{{ seasonal }}</span>
    </div>

    <div class="container hero-inner">
      <div class="hero-text">
        <div class="badges">
          <span class="badge font-cute">☕ 营业中 · OPEN</span>
          <span class="badge status font-cute">店主现在:{{ site.authorStatus || '赶稿中 ✍️' }}</span>
        </div>
        <h1 class="font-cute">{{ site.title }}</h1>
        <p class="slogan">{{ site.slogan }}</p>
        <p class="sub">
          这里是 <b>{{ site.author }}</b> 的非官方粉丝后援会
        </p>
        <div class="actions">
          <a href="#novels" class="btn btn-primary">📖 去看小说菜单</a>
          <a :href="site.authorPixiv" target="_blank" rel="noopener" class="btn btn-ghost">
            作者主页 ↗
          </a>
          <HandDoodle kind="arrow" :size="56" class="cta-arrow" />
        </div>
      </div>

      <div class="hero-mascots">
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
.hero {
  position: relative;
  padding: 96px 0 80px;
  background: linear-gradient(160deg, var(--bg) 0%, var(--bg-deep) 100%);
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

.badges {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.badge {
  display: inline-block;
  padding: 4px 16px;
  border-radius: 999px;
  background: #fff;
  border: 2px solid var(--pink-soft);
  color: var(--pink-deep);
  font-size: 0.9rem;
}

.badge.status {
  background: var(--pink-pale);
  border-color: var(--pink-pale);
  color: var(--ink);
}

h1 {
  font-size: clamp(2rem, 6vw, 3.4rem);
  line-height: 1.25;
  color: var(--ink);
  margin-bottom: 16px;
}

.slogan {
  font-size: 1.15rem;
  color: var(--pink-deep);
  margin-bottom: 8px;
}

.sub {
  color: var(--muted);
  margin-bottom: 32px;
}

.sub b {
  color: var(--pink-deep);
}

.actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.hero-mascots {
  position: relative;
  flex-shrink: 0;
  text-align: center;
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
  z-index: 2; /* 盖过旁边的小猪咪 */
  display: inline-block;
  margin-top: 10px;
  font-size: 0.85rem;
  color: var(--muted);
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 4px 16px;
}

.pigmi {
  z-index: 1;
}

.cat {
  display: block;
  animation: bob 3.5s ease-in-out infinite;
}

.pigmi {
  position: absolute;
  right: -24px;
  bottom: -28px;
  animation: bob 3.5s ease-in-out 0.8s infinite;
}

@keyframes bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.paws .paw {
  position: absolute;
  opacity: 0.5;
}

.p1 { top: 12%; left: 4%; }
.p2 { top: 64%; left: 46%; }
.p3 { bottom: 8%; right: 6%; }
.p4 { top: 8%; right: 28%; }

/* 手绘涂鸦:云和星星慢悠悠地漂 */
.d1 { top: 6%; left: 34%; animation: drift 9s ease-in-out infinite; }
.d2 { bottom: 14%; left: 8%; animation: drift 11s ease-in-out 1.5s infinite; }
.d3 { top: 20%; right: 10%; animation: drift 7s ease-in-out 0.5s infinite; }
.d4 { top: 46%; left: 28%; animation: drift 8s ease-in-out 2s infinite; }
.d5 { bottom: 24%; right: 32%; animation: drift 10s ease-in-out 1s infinite; }

.seasonal {
  position: absolute;
  font-size: 26px;
  opacity: 0.75;
}

.s1 { top: 30%; left: 14%; animation: drift 8s ease-in-out 0.8s infinite; }
.s2 { top: 8%; right: 42%; font-size: 20px; animation: drift 12s ease-in-out 2.2s infinite; }

@keyframes drift {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-12px) rotate(4deg); }
}

.actions {
  position: relative;
}

.cta-arrow {
  position: absolute;
  right: -64px;
  top: -34px;
  transform: rotate(8deg);
}

@media (max-width: 720px) {
  .hero-inner {
    flex-direction: column-reverse;
    text-align: center;
    gap: 32px;
  }
  .actions {
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
  .cat, .pigmi, .paws .paw { animation: none; }
}
</style>

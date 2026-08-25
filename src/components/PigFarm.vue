<script setup>
// 全店养猪:全站共养一只猪,喂食/摸头攒经验进化
import { ref, computed, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'

const pig = ref(null)
const busy = ref(false)
const bubble = ref('')
const evolved = ref('')
let bubbleTimer = null

const STAGE_IMG = {
  猪崽: '/pets/pigmi.png',
  小猪咪: '/pets/farm-teen.png',
  圆润猪咪: '/pets/farm-teen.png',
  猪王: '/pets/farm-king.png',
}

const stageImg = computed(() => STAGE_IMG[pig.value?.stage] || STAGE_IMG['猪崽'])
const stageSize = computed(() => {
  const s = pig.value?.stage
  return s === '猪崽' ? 90 : s === '小猪咪' ? 110 : s === '圆润猪咪' ? 130 : 140
})

// 猪的心情语录
const moodText = computed(() => {
  if (!pig.value) return ''
  const { hunger, mood } = pig.value
  if (hunger < 30) return '肚肚好饿……谁来喂喂我'
  if (mood < 30) return '有点无聊,摸摸我嘛'
  if (hunger > 80 && mood > 80) return '好幸福!今天也是被爱着的一天!'
  return '哼唧哼唧~'
})

onMounted(load)

async function load() {
  try {
    const res = await fetch('/api/pig')
    pig.value = await res.json()
  } catch {
    /* 静默 */
  }
}

function say(text, ms = 2600) {
  bubble.value = text
  clearTimeout(bubbleTimer)
  bubbleTimer = setTimeout(() => (bubble.value = ''), ms)
}

async function act(action) {
  if (action === 'feed') localStorage.setItem('catcafe_ach_feedpig', '1')
  if (busy.value) return
  busy.value = true
  try {
    const res = await fetch(`/api/pig/${action}`, { method: 'POST' })
    const data = await res.json()
    if (res.ok) {
      pig.value = { ...pig.value, ...data }
      if (data.evolved) {
        evolved.value = data.evolved
        say(`🎉 进化成「${data.evolved}」啦!!`, 4000)
        setTimeout(() => (evolved.value = ''), 4000)
      } else {
        say(action === 'feed' ? '啊呜啊呜……好吃!' : '呼噜呼噜,好舒服~')
      }
    } else {
      say(data.error || '稍后再试')
    }
  } catch {
    say('网络打了个盹')
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <section id="farm" class="farm">
    <div class="container">
      <SectionTitle title="店里的猪" sub="全站共养一只猪 · 喂它摸它,看着它长大" />

      <div class="pen" v-reveal>
        <!-- 猪圈 -->
        <div class="pig-area">
          <transition name="bub">
            <span v-if="bubble" class="bubble font-cute">{{ bubble }}</span>
          </transition>
          <transition name="evo">
            <div v-if="evolved" class="evo-banner font-cute">✨ 进化!{{ evolved }} ✨</div>
          </transition>
          <img
            :src="stageImg"
            :alt="pig?.stage"
            class="pig-img"
            :style="{ height: stageSize + 'px' }"
            draggable="false"
          />
          <div class="fence" aria-hidden="true">🌱 🌾 🌱 🌾 🌱 🌾 🌱</div>
        </div>

        <!-- 状态面板 -->
        <div class="status">
          <div class="stage-line">
            <span class="stage-badge font-cute">{{ pig?.stage || '猪崽' }}</span>
            <span class="xp" v-if="pig">
              经验 {{ pig.xp }}{{ pig.next ? ` / ${pig.next}` : ' · 已满级!' }}
            </span>
          </div>

          <div class="bar-row">
            <span class="bar-label">饱食</span>
            <div class="bar"><div class="fill hunger" :style="{ width: (pig?.hunger ?? 0) + '%' }" /></div>
            <span class="bar-num">{{ pig?.hunger ?? 0 }}</span>
          </div>
          <div class="bar-row">
            <span class="bar-label">心情</span>
            <div class="bar"><div class="fill mood" :style="{ width: (pig?.mood ?? 0) + '%' }" /></div>
            <span class="bar-num">{{ pig?.mood ?? 0 }}</span>
          </div>

          <div class="actions">
            <button class="btn btn-primary" :disabled="busy" @click="act('feed')">🍰 喂它</button>
            <button class="btn btn-ghost" :disabled="busy" @click="act('pet')">🫳 摸头</button>
          </div>

          <p class="stats" v-if="pig">
            全站已投喂 {{ pig.fed }} 次 · 摸头 {{ pig.pats }} 次
          </p>
          <p class="mood-line font-cute">"{{ moodText }}"</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.pen {
  max-width: 860px;
  margin: 0 auto;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 32px;
  display: flex;
  gap: 36px;
  align-items: center;
}

/* 猪圈:草地感 */
.pig-area {
  flex: 1;
  position: relative;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  background: linear-gradient(180deg, var(--bg-deep) 0%, var(--card) 70%, var(--surface-2) 100%);
  border: 2px dashed var(--pink-soft);
  border-radius: 18px;
  padding: 20px 12px 0;
  overflow: hidden;
}

.pig-img {
  transition: height 0.5s ease;
  animation: pigBob 2.2s ease-in-out infinite;
  z-index: 1;
}

@keyframes pigBob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.fence {
  width: 100%;
  text-align: center;
  font-size: 1.1rem;
  padding: 4px 0 6px;
  letter-spacing: 6px;
}

.bubble {
  position: absolute;
  top: 12px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--surface-2);
  border: 2px solid var(--pink-soft);
  color: var(--pink-deep);
  border-radius: 14px 14px 14px 4px;
  padding: 5px 16px;
  font-size: 0.9rem;
  white-space: nowrap;
  z-index: 2;
  box-shadow: var(--shadow);
}

.evo-banner {
  position: absolute;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: linear-gradient(135deg, #ffd76e, #f9718f);
  color: #fff;
  font-size: 1.2rem;
  border-radius: 16px;
  padding: 12px 28px;
  z-index: 3;
  box-shadow: 0 10px 30px rgba(249, 113, 143, 0.5);
}

.bub-enter-active, .evo-enter-active { transition: all 0.25s ease; }
.bub-leave-active, .evo-leave-active { transition: all 0.3s ease; }
.bub-enter-from, .evo-enter-from { opacity: 0; transform: translateX(-50%) translateY(8px) scale(0.85); }
.bub-leave-to, .evo-leave-to { opacity: 0; }
.evo-enter-from, .evo-leave-to { transform: translate(-50%, -50%) scale(0.8); }

.status {
  flex: 1;
  min-width: 260px;
}

.stage-line {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-bottom: 18px;
}

.stage-badge {
  background: linear-gradient(160deg, var(--pink), var(--pink-deep));
  color: #fff;
  border-radius: 12px;
  padding: 6px 18px;
  font-size: 1.2rem;
}

.xp {
  color: var(--muted);
  font-size: 0.85rem;
}

.bar-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.bar-label {
  font-size: 0.85rem;
  color: var(--muted);
  width: 32px;
  flex-shrink: 0;
}

.bar {
  flex: 1;
  height: 14px;
  background: var(--pink-pale);
  border-radius: 999px;
  overflow: hidden;
}

.fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.6s ease;
}

.fill.hunger { background: linear-gradient(90deg, #ffc98a, #ff9d5c); }
.fill.mood { background: linear-gradient(90deg, var(--pink-soft), var(--pink)); }

.bar-num {
  font-size: 0.8rem;
  color: var(--muted);
  width: 28px;
  text-align: right;
}

.actions {
  display: flex;
  gap: 12px;
  margin: 20px 0 12px;
}

.stats {
  color: var(--muted);
  font-size: 0.8rem;
}

.mood-line {
  margin-top: 8px;
  color: var(--pink-deep);
  font-size: 0.95rem;
}

@media (max-width: 720px) {
  .pen {
    flex-direction: column;
    padding: 22px;
  }
  .pig-area {
    width: 100%;
  }
}

@media (prefers-reduced-motion: reduce) {
  .pig-img { animation: none; }
}
</style>

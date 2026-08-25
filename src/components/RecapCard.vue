<script setup>
// 营业报告:今日 / 本月 / 总账 三档小票
import { ref, computed, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'

const d = ref(null)
const tab = ref('day') // day | month | all

const dateStr = new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })
const monthStr = new Date().toLocaleDateString('zh-CN', { month: 'long' })

const tabs = [
  { key: 'day', label: '今日' },
  { key: 'month', label: '本月' },
  { key: 'all', label: '总账' },
]

const sub = computed(() => {
  if (tab.value === 'day') return `${dateStr} · 猫咖日报`
  if (tab.value === 'month') return `${monthStr} · 猫咖月报`
  return '开店以来 · 总账'
})

// 按当前标签取数
const lines = computed(() => {
  if (!d.value) return []
  if (tab.value === 'all') {
    const t = d.value
    return [
      ['累计到店猪咪', `${t.visits} 只`],
      ['猫猫被撸', `${t.pets} 次`],
      ['被投喂', `${t.feeds} 次`],
      ['催更', `${t.urges} 次`],
      ['留言', `${t.messages} 条`],
      ['明信片', `${t.postcards} 张 · 获赞 ${t.likes}`],
      ['攒鱼干的猪咪', `${t.pigmis} 只`],
    ]
  }
  const src = tab.value === 'day' ? d.value.day : d.value.month
  if (!src || !src.date) return []
  const p = tab.value === 'day' ? '今日' : '本月'
  return [
    [`${p}到店猪咪`, `${src.visits} 只`],
    [`猫猫被撸`, `${src.pets} 次`],
    [`被投喂`, `${src.feeds} 次`],
    [`催更`, `${src.urges} 次`],
    [`新留言`, `${src.messages} 条`],
    [`新明信片`, `${src.postcards} 张 · 获赞 ${src.likes}`],
  ]
})

onMounted(async () => {
  try {
    const res = await fetch('/api/recap')
    d.value = await res.json()
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <section id="recap" v-if="d">
    <div class="container">
      <SectionTitle title="营业报告" :sub="sub" />

      <div class="mini-tabs" v-reveal>
        <button
          v-for="t in tabs"
          :key="t.key"
          :class="{ active: tab === t.key }"
          @click="tab = t.key"
        >
          {{ t.label }}
        </button>
      </div>

      <div class="receipt" v-reveal>
        <template v-if="lines.length">
          <div v-for="(l, i) in lines" :key="i" class="rline">
            <span>{{ l[0] }}</span><b>{{ l[1] }}</b>
          </div>
        </template>
        <p v-else class="empty">还没有开账,等你来第一单~</p>
        <div class="rsep" />
        <div v-if="d.topPost" class="rline hi">
          <span>👑 人气明信片</span><b>{{ d.topPost.nick }} · {{ d.topPost.likes }} 赞</b>
        </div>
        <div v-if="d.topPigmi" class="rline hi">
          <span>🐟 鱼干首富</span><b>{{ d.topPigmi.nick }} · {{ d.topPigmi.points }}</b>
        </div>
        <div class="rfoot">—— 今日营业中,欢迎常来 ——</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.mini-tabs {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 24px;
}

.mini-tabs button {
  border: 2px solid var(--pink-soft);
  background: var(--surface-2);
  color: var(--ink);
  border-radius: 999px;
  padding: 6px 22px;
  font-size: 0.88rem;
  cursor: pointer;
  transition: all 0.2s;
}

.mini-tabs button.active {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
}

.receipt {
  max-width: 420px;
  margin: 0 auto;
  background: var(--surface-2);
  padding: 28px 28px 20px;
  font-family: 'JetBrains Mono', 'Noto Sans SC', monospace;
  box-shadow: var(--shadow);
  /* 小票锯齿边 */
  clip-path: polygon(
    0 0, 100% 0,
    100% calc(100% - 8px), 97% 100%, 94% calc(100% - 8px), 91% 100%,
    88% calc(100% - 8px), 85% 100%, 82% calc(100% - 8px), 79% 100%,
    76% calc(100% - 8px), 73% 100%, 70% calc(100% - 8px), 67% 100%,
    64% calc(100% - 8px), 61% 100%, 58% calc(100% - 8px), 55% 100%,
    52% calc(100% - 8px), 49% 100%, 46% calc(100% - 8px), 43% 100%,
    40% calc(100% - 8px), 37% 100%, 34% calc(100% - 8px), 31% 100%,
    28% calc(100% - 8px), 25% 100%, 22% calc(100% - 8px), 19% 100%,
    16% calc(100% - 8px), 13% 100%, 10% calc(100% - 8px), 7% 100%,
    4% calc(100% - 8px), 0 100%
  );
}

.rline {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 7px 0;
  font-size: 0.88rem;
  color: var(--ink);
  border-bottom: 1px dashed var(--pink-pale);
}

.rline span {
  color: var(--muted);
}

.rline b {
  text-align: right;
}

.rline.hi b {
  color: var(--pink-deep);
}

.rsep {
  border-bottom: 2px dashed var(--pink-soft);
  margin: 10px 0;
}

.empty {
  text-align: center;
  color: var(--muted);
  font-size: 0.9rem;
  padding: 16px 0;
}

.rfoot {
  text-align: center;
  color: var(--muted);
  font-size: 0.78rem;
  margin-top: 14px;
}
</style>

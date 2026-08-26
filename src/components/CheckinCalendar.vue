<script setup>
// 猫爪签到日历:每天盖一枚爪印,连续签到鱼干加倍,当季签到点亮季节限定成就
import { ref, computed, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'
import PawPrint from './PawPrint.vue'
import { profile } from '../app.js'
import { currentSeason, unlockSeasonAch } from '../seasons.js'

const todayDone = ref(false)
const streak = ref(0)
const total = ref(0)
const dates = ref([])
const doing = ref(false)
const tip = ref('')
const season = currentSeason()

// 当月日历格子
const now = new Date()
const year = now.getFullYear()
const month = now.getMonth()
const days = computed(() => {
  const first = new Date(year, month, 1).getDay() // 周日开头
  const count = new Date(year, month + 1, 0).getDate()
  const cells = []
  for (let i = 0; i < first; i++) cells.push(null)
  for (let d = 1; d <= count; d++) {
    const key = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    cells.push({ day: d, stamped: dates.value.includes(key), isToday: d === now.getDate() })
  }
  return cells
})

async function loadStatus() {
  if (!profile.nick) return
  try {
    const res = await fetch('/api/checkin?nick=' + encodeURIComponent(profile.nick))
    if (!res.ok) return
    const data = await res.json()
    todayDone.value = !!data.todayDone
    streak.value = data.streak || 0
    total.value = data.total || 0
    dates.value = data.dates || []
  } catch { /* 离线就不显示状态 */ }
}

async function doCheckin() {
  if (!profile.nick) {
    tip.value = '先在「我的」绑个昵称,爪印才能记在你名下喵~'
    return
  }
  doing.value = true
  tip.value = ''
  try {
    const res = await fetch('/api/checkin', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nick: profile.nick }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || '签到失败,再试试')
    todayDone.value = true
    streak.value = data.streak
    if (!data.already) {
      total.value++
      const key = `${year}-${String(month + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
      dates.value = [...dates.value, key]
      unlockSeasonAch(season) // 当季签到 → 点亮季节限定成就
      let msg = `盖爪成功!+${data.reward} 鱼干`
      if (data.bonus) msg += `,连击奖励再 +${data.bonus}`
      msg += `(现有 ${data.points},「${data.title}」)`
      if (season) msg += ` · ${season.icon}点亮「${season.achName}」!`
      tip.value = msg
      if (navigator.vibrate) navigator.vibrate(60)
    } else {
      tip.value = '今天已经盖过爪啦,明天再来~'
    }
  } catch (err) {
    tip.value = err.message || '网络被猫踩了,稍后再试'
  } finally {
    doing.value = false
  }
}

onMounted(loadStatus)
</script>

<template>
  <section id="checkin">
    <div class="container">
      <SectionTitle
        title="猫爪签到"
        :sub="season ? `${season.icon} ${season.name}进行中 · 当季签到点亮限定成就「${season.achName}」` : '每天来盖一枚爪印'"
      />

      <div class="panel" v-reveal>
        <div class="stats">
          <div class="stat">
            <b>{{ streak }}</b>
            <span>连续盖爪</span>
          </div>
          <div class="stat">
            <b>{{ total }}</b>
            <span>累计盖爪</span>
          </div>
          <button class="paw-btn" :class="{ done: todayDone }" :disabled="doing || todayDone" @click="doCheckin">
            <PawPrint :size="30" :color="todayDone ? '#c9b8be' : '#fff'" />
            {{ todayDone ? '今日已盖' : doing ? '盖爪中…' : '盖爪签到' }}
          </button>
        </div>

        <div class="calendar">
          <span v-for="w in ['日', '一', '二', '三', '四', '五', '六']" :key="w" class="week">{{ w }}</span>
          <span v-for="(c, i) in days" :key="i" class="cell" :class="{ stamped: c && c.stamped, today: c && c.isToday }">
            <template v-if="c">
              <PawPrint v-if="c.stamped" :size="18" color="#f9718f" />
              <template v-else>{{ c.day }}</template>
            </template>
          </span>
        </div>

        <p class="rule">每日盖爪 +2 鱼干 · 连续 7 天 +5 · 连续 30 天 +15</p>
        <p v-if="tip" class="tip">{{ tip }}</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.panel {
  max-width: 560px;
  margin: 0 auto;
  background: var(--card);
  border: 3px solid var(--pink-soft);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 24px;
}

.stats {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
}

.stat {
  text-align: center;
}

.stat b {
  display: block;
  font-size: 1.6rem;
  color: var(--pink-deep);
}

.stat span {
  font-size: 0.8rem;
  color: var(--muted);
}

.paw-btn {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
  border: none;
  border-radius: 999px;
  padding: 12px 22px;
  background: var(--pink-deep);
  color: #fff;
  font-size: 1rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  transition: transform 0.15s;
}

.paw-btn:active {
  transform: scale(0.92);
}

.paw-btn.done {
  background: var(--pink-pale);
  color: var(--muted);
  cursor: default;
}

.calendar {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
  text-align: center;
}

.week {
  font-size: 0.75rem;
  color: var(--muted);
  padding-bottom: 4px;
}

.cell {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  color: var(--ink);
  border-radius: 8px;
  background: var(--surface-2);
}

.cell.stamped {
  background: var(--pink-pale);
}

.cell.today {
  outline: 2px dashed var(--pink-deep);
  outline-offset: -2px;
}

.rule {
  margin-top: 16px;
  text-align: center;
  font-size: 0.8rem;
  color: var(--muted);
}

.tip {
  margin-top: 10px;
  text-align: center;
  font-size: 0.9rem;
  color: var(--pink-deep);
  line-height: 1.6;
}
</style>

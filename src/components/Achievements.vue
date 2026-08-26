<script setup>
// 成就图鉴:点亮你在猫咖的足迹(含季节限定章,过季不候)
import { ref, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'
import { seasons, currentSeason } from '../seasons.js'

const ACHIEVEMENTS = [
  { key: 'visit', icon: '🐾', name: '初入猫咖', desc: '第一次来到咖啡厅' },
  { key: 'petcat', icon: '🫳', name: '撸猫达人', desc: '撸了一次猫' },
  { key: 'feedpig', icon: '🍰', name: '饲养员见习', desc: '喂过一次店里的猪' },
  { key: 'postcard', icon: '💌', name: '明信片上手', desc: '寄出过一张明信片' },
  { key: 'urge', icon: '📣', name: '催更小能手', desc: '拍爪催更过' },
  { key: 'gate', icon: '🔓', name: '真猪咪认证', desc: '通过入群答题验证' },
  { key: 'whack', icon: '🎯', name: '拍拍乐', desc: '玩过拍猪咪小游戏' },
  { key: 'together', icon: '💕', name: '贴贴大师', desc: '让两只小家伙贴贴了' },
  { key: 'glitch', icon: '🌀', name: '裂缝探险家', desc: '点开过神秘裂缝' },
]

// 季节限定章:当季签到点亮,过季变灰显示「已过季」
const now = new Date()
const curYear = now.getFullYear()
const cur = currentSeason(now)
const SEASONAL = seasons.map((s) => ({
  key: 'season_' + s.id,
  icon: s.icon,
  name: s.achName,
  desc: s.desc,
  active: cur && cur.id === s.id, // 是否当季(可点亮)
}))

const unlocked = ref(new Set())
const total = ref(0)
const totalCount = ACHIEVEMENTS.length + SEASONAL.length

onMounted(() => {
  // 首次访问自动点亮"初入猫咖"
  localStorage.setItem('catcafe_ach_visit', '1')
  const set = new Set()
  for (const a of ACHIEVEMENTS) {
    if (localStorage.getItem('catcafe_ach_' + a.key) === '1') set.add(a.key)
  }
  for (const s of SEASONAL) {
    if (localStorage.getItem('catcafe_ach_' + s.key) === '1') set.add(s.key)
  }
  unlocked.value = set
  total.value = set.size
})
</script>

<template>
  <section id="achievements" v-if="total > 0">
    <div class="container">
      <SectionTitle title="成就图鉴" :sub="`已点亮 ${total} / ${totalCount} 枚`" />

      <div class="grid" v-reveal>
        <div
          v-for="a in ACHIEVEMENTS"
          :key="a.key"
          class="ach"
          :class="{ on: unlocked.has(a.key) }"
        >
          <span class="icon">{{ a.icon }}</span>
          <b>{{ a.name }}</b>
          <p>{{ a.desc }}</p>
        </div>
      </div>

      <p class="season-title">季节限定 · {{ curYear }}</p>
      <div class="grid" v-reveal>
        <div
          v-for="s in SEASONAL"
          :key="s.key"
          class="ach seasonal"
          :class="{ on: unlocked.has(s.key), active: s.active }"
        >
          <span class="icon">{{ s.icon }}</span>
          <b>{{ s.name }}</b>
          <p>{{ s.desc }}</p>
          <span v-if="!unlocked.has(s.key)" class="season-tag" :class="{ now: s.active }">
            {{ s.active ? '进行中' : '已过季' }}
          </span>
          <span v-else class="season-tag got">已点亮</span>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  max-width: 720px;
  margin: 0 auto;
}

.ach {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: 16px;
  padding: 18px 12px;
  text-align: center;
  opacity: 0.35;
  filter: grayscale(0.8);
  transition: all 0.3s ease;
}

.ach.on {
  opacity: 1;
  filter: none;
  border-color: var(--pink-soft);
  box-shadow: var(--shadow);
}

.ach .icon {
  font-size: 1.8rem;
  display: block;
  margin-bottom: 6px;
}

.ach b {
  display: block;
  color: var(--ink);
  font-size: 0.92rem;
  margin-bottom: 4px;
}

.ach p {
  color: var(--muted);
  font-size: 0.75rem;
  line-height: 1.5;
}

.ach.on b {
  color: var(--pink-deep);
}

.season-title {
  text-align: center;
  color: var(--pink-deep);
  font-size: 0.9rem;
  font-weight: 600;
  margin: 28px 0 16px;
}

.ach.seasonal {
  position: relative;
  border-style: dashed;
}

.ach.seasonal.active:not(.on) {
  opacity: 0.7;
  filter: none;
  border-color: var(--pink-soft);
  animation: breathe 2.2s ease-in-out infinite;
}

@keyframes breathe {
  0%, 100% { box-shadow: 0 0 0 0 rgba(249, 113, 143, 0.25); }
  50% { box-shadow: 0 0 0 6px rgba(249, 113, 143, 0); }
}

.season-tag {
  display: inline-block;
  margin-top: 6px;
  font-size: 0.68rem;
  padding: 1px 10px;
  border-radius: 999px;
  background: var(--surface-2);
  color: var(--muted);
}

.season-tag.now {
  background: var(--pink-pale);
  color: var(--pink-deep);
}

.season-tag.got {
  background: var(--pink);
  color: #fff;
}

@media (max-width: 720px) {
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>

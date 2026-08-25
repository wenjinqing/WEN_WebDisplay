<script setup>
// 成就图鉴:点亮你在猫咖的足迹
import { ref, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'

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

const unlocked = ref(new Set())
const total = ref(0)

onMounted(() => {
  // 首次访问自动点亮"初入猫咖"
  localStorage.setItem('catcafe_ach_visit', '1')
  const set = new Set()
  for (const a of ACHIEVEMENTS) {
    if (localStorage.getItem('catcafe_ach_' + a.key) === '1') set.add(a.key)
  }
  unlocked.value = set
  total.value = set.size
})
</script>

<template>
  <section id="achievements" v-if="total > 0">
    <div class="container">
      <SectionTitle title="成就图鉴" :sub="`已点亮 ${total} / ${ACHIEVEMENTS.length} 枚`" />

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

@media (max-width: 720px) {
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>

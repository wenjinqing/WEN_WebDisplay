<script setup>
import PawPrint from './PawPrint.vue'
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'

const emit = defineEmits(['read'])
</script>

<template>
  <section id="novels" class="menu-section">
    <div class="container">
      <SectionTitle title="今日特供 · 小说菜单" sub="点单即下载 · 本区作品含轻度成人向内容,未成年猪咪请自觉绕行喵~" />

      <div class="menu-board" v-reveal>
        <div v-for="n in site.novels" :key="n.file" class="menu-item">
          <div class="item-info">
            <div class="item-head">
              <h3 class="font-cute">{{ n.title }}</h3>
              <span class="cup">{{ n.cup }}</span>
            </div>
            <p class="desc">{{ n.desc }}</p>
          </div>
          <div class="item-actions">
            <button class="btn btn-ghost read" @click="emit('read', n)">在线阅读</button>
            <a class="btn btn-primary download" :href="`/downloads/${n.file}`" download>
              <PawPrint :size="18" color="#fff" />
              下载
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.menu-section {
  background: var(--bg-deep);
}

.menu-board {
  max-width: 860px;
  margin: 0 auto;
  background: var(--card);
  border: 3px solid var(--pink-soft);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 12px 32px;
  position: relative;
}

/* 菜单板顶部挂绳孔装饰 */
.menu-board::before {
  content: '';
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 14px;
  border-radius: 8px 8px 0 0;
  background: var(--pink-soft);
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px 0;
  border-bottom: 2px dashed var(--pink-pale);
}

.menu-item:last-child {
  border-bottom: none;
}

.item-info {
  flex: 1;
}

.item-head {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

h3 {
  font-size: 1.15rem;
  color: var(--ink);
}

.cup {
  font-size: 0.8rem;
  color: var(--pink-deep);
  background: var(--pink-pale);
  border-radius: 999px;
  padding: 2px 12px;
  white-space: nowrap;
}

.desc {
  color: var(--muted);
  font-size: 0.92rem;
  margin-top: 6px;
}

.item-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.read {
  padding: 12px 20px;
}

@media (max-width: 720px) {
  .menu-board {
    padding: 8px 20px;
  }
  .menu-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 14px;
  }
  .item-actions {
    align-self: stretch;
  }
  .item-actions .btn {
    flex: 1;
    justify-content: center;
  }
}
</style>

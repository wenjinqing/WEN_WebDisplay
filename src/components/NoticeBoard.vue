<script setup>
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'

const tagColor = {
  公告: '#f9718f',
  新坑: '#e85d7f',
  更新: '#c98aa0',
  活动: '#d4758e',
}
</script>

<template>
  <section id="notice">
    <div class="container">
      <SectionTitle title="咖啡厅公告板" sub="新坑上架 · 活动通知 · 店主的碎碎念" />

      <div class="board" v-reveal>
        <div v-for="(n, i) in site.notices" :key="i" class="notice">
          <div class="meta">
            <span class="tag" :style="{ background: tagColor[n.tag] || '#f9718f' }">{{ n.tag }}</span>
            <time>{{ n.date }}</time>
          </div>
          <p>{{ n.text }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.board {
  max-width: 760px;
  margin: 0 auto;
  background: #6d4a55; /* 深色木黑板 */
  border-radius: var(--radius);
  border: 6px solid #f0d7de; /* 浅色木框 */
  box-shadow: var(--shadow);
  padding: 28px;
}

.notice {
  background: rgba(255, 249, 250, 0.08);
  border: 1px dashed rgba(255, 214, 226, 0.4);
  /* 手绘便利贴:歪扭圆角 */
  border-radius: 255px 14px 225px 14px / 14px 225px 14px 255px;
  padding: 16px 22px;
  margin-bottom: 14px;
  color: #fff6f8;
}

.notice:last-child {
  margin-bottom: 0;
}

.meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.tag {
  font-size: 0.75rem;
  color: #fff;
  border-radius: 999px;
  padding: 2px 12px;
}

time {
  font-size: 0.8rem;
  color: rgba(255, 214, 226, 0.7);
}

.notice p {
  font-size: 0.95rem;
  line-height: 1.8;
}

@media (max-width: 720px) {
  .board {
    padding: 18px;
  }
}
</style>

<script setup>
import { ref } from 'vue'
import UrgeWall from './UrgeWall.vue'
import MessageForm from './MessageForm.vue'
import MessageList from './MessageList.vue'
import SectionTitle from './SectionTitle.vue'

const tab = ref('wall') // wall 留言墙 | write 写留言 | urge 催更墙

const tabs = [
  { key: 'wall', label: '💬 留言墙' },
  { key: 'write', label: '✍️ 写留言' },
  { key: 'urge', label: '🐾 催更墙' },
]
</script>

<template>
  <section id="interact" class="interact">
    <div class="container">
      <SectionTitle title="互动区" sub="催更和留言,都是对猫猫酱的爱" />

      <!-- 三个子页面切换 -->
      <div class="tab-bar" v-reveal>
        <button
          v-for="t in tabs"
          :key="t.key"
          :class="{ active: tab === t.key }"
          @click="tab = t.key"
        >
          {{ t.label }}
        </button>
      </div>

      <div class="panel hd-card" v-reveal>
        <MessageList v-if="tab === 'wall'" />
        <MessageForm v-else-if="tab === 'write'" />
        <div v-else class="urge-wrap">
          <UrgeWall />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.interact {
  background: var(--bg-deep);
}

.tab-bar {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 28px;
  flex-wrap: wrap;
}

.tab-bar button {
  border: 2px solid var(--pink-soft);
  background: #fff;
  color: var(--ink);
  border-radius: 999px;
  padding: 9px 26px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-bar button.active {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
  box-shadow: 0 4px 12px rgba(249, 113, 143, 0.35);
}

.panel {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  box-shadow: var(--shadow);
  padding: 32px 30px;
  min-height: 320px;
}

.urge-wrap {
  max-width: 460px;
  margin: 0 auto;
}

.urge-wrap :deep(.urge-wall) {
  border: none;
  box-shadow: none;
  padding: 0;
}

@media (max-width: 720px) {
  .panel {
    padding: 22px 18px;
  }
}
</style>

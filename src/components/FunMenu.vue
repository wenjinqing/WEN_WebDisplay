<script setup>
// 玩乐菜单:把散落在屏幕上的功能按钮收进一个悬浮球
// 各功能组件监听这里派发的 window 事件
import { ref, onMounted, onUnmounted } from 'vue'

const open = ref(false)
const active = ref({ danmaku: false, laser: false }) // 由组件回同步状态

const items = [
  { key: 'feed', icon: '🍰', label: '投喂', event: 'menu-feed' },
  { key: 'danmaku', icon: '💬', label: '弹幕', event: 'menu-danmaku', toggle: 'danmaku' },
  { key: 'weather', icon: '🌦️', label: '天气', event: 'menu-weather' },
  { key: 'laser', icon: '🔴', label: '激光笔', event: 'menu-laser', toggle: 'laser' },
  { key: 'whack', icon: '🎯', label: '拍猪咪', event: 'menu-whack' },
  { key: 'pet', icon: '🐾', label: '看板娘', event: 'menu-pet' },
]

function fire(item) {
  window.dispatchEvent(new CustomEvent(item.event))
  if (!item.toggle) open.value = false // 动作类按钮点了就收起
}

// 组件回同步开关状态(高亮用)
function onState(e) {
  active.value = { ...active.value, ...e.detail }
}

function onDocClick(e) {
  if (open.value && !e.target.closest('.fun-menu')) open.value = false
}

onMounted(() => {
  window.addEventListener('menu-state', onState)
  document.addEventListener('click', onDocClick)
})

onUnmounted(() => {
  window.removeEventListener('menu-state', onState)
  document.removeEventListener('click', onDocClick)
})
</script>

<template>
  <div class="fun-menu">
    <transition-group name="pop" tag="div" class="menu-items" v-if="open">
      <button
        v-for="item in items"
        :key="item.key"
        class="menu-item"
        :class="{ on: item.toggle && active[item.toggle] }"
        @click="fire(item)"
      >
        <span class="icon">{{ item.icon }}</span>
        <span class="label">{{ item.label }}</span>
      </button>
    </transition-group>
    <button class="fab" :class="{ open }" aria-label="玩乐菜单" @click="open = !open">
      {{ open ? '✕' : '🎀' }}
    </button>
  </div>
</template>

<style scoped>
.fun-menu {
  position: fixed;
  left: 20px;
  bottom: 24px;
  z-index: 65;
  display: flex;
  flex-direction: column-reverse;
  align-items: flex-start;
  gap: 10px;
}

.fab {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  border: 2px solid var(--pink-soft);
  background: rgba(255, 255, 255, 0.95);
  font-size: 22px;
  cursor: pointer;
  box-shadow: var(--shadow);
  transition: transform 0.2s ease, background 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fab:hover {
  transform: scale(1.08);
}

.fab.open {
  background: var(--pink);
  transform: rotate(90deg);
}

.menu-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 4px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 2px solid var(--pink-soft);
  background: rgba(255, 255, 255, 0.95);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 7px 16px 7px 10px;
  font-size: 0.85rem;
  cursor: pointer;
  box-shadow: var(--shadow);
  transition: transform 0.15s, background 0.15s;
}

.menu-item:hover {
  transform: translateX(4px);
  background: var(--pink-pale);
}

.menu-item.on {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
}

.menu-item .icon {
  font-size: 1rem;
}

.pop-enter-active { transition: all 0.18s ease; }
.pop-leave-active { transition: all 0.12s ease; }
.pop-enter-from, .pop-leave-to { opacity: 0; transform: translateY(10px) scale(0.9); }

@media (max-width: 720px) {
  .fun-menu {
    left: 14px;
    bottom: 16px;
  }
  .fab {
    width: 46px;
    height: 46px;
    font-size: 19px;
  }
  .menu-item {
    padding: 6px 14px 6px 9px;
    font-size: 0.8rem;
  }
}
</style>

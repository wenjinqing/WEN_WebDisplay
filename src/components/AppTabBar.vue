<script setup>
// APP 底部切换栏:首页 / 小说 / 插画 / 互动 / 我的
const props = defineProps({ modelValue: { type: String, default: 'home' } })
const emit = defineEmits(['update:modelValue'])

const tabs = [
  { key: 'home', label: '首页', icon: 'paw' },
  { key: 'novels', label: '小说', icon: 'book' },
  { key: 'gallery', label: '插画', icon: 'image' },
  { key: 'interact', label: '互动', icon: 'chat' },
  { key: 'me', label: '我的', icon: 'user' },
]
</script>

<template>
  <nav class="tabbar">
    <button
      v-for="t in tabs"
      :key="t.key"
      :class="['tab', { active: modelValue === t.key }]"
      @click="emit('update:modelValue', t.key)"
    >
      <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor"
           stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
        <!-- 猫爪 -->
        <template v-if="t.icon === 'paw'">
          <ellipse cx="12" cy="15" rx="4" ry="3.2" />
          <ellipse cx="6.2" cy="9.5" rx="1.7" ry="2.1" />
          <ellipse cx="10" cy="7.2" rx="1.7" ry="2.1" />
          <ellipse cx="14" cy="7.2" rx="1.7" ry="2.1" />
          <ellipse cx="17.8" cy="9.5" rx="1.7" ry="2.1" />
        </template>
        <!-- 书 -->
        <template v-else-if="t.icon === 'book'">
          <path d="M4 5.5A2.5 2.5 0 0 1 6.5 3H20v15H6.5A2.5 2.5 0 0 0 4 20.5z" />
          <path d="M4 20.5V5.5" />
          <path d="M8 7.5h8" />
        </template>
        <!-- 图片 -->
        <template v-else-if="t.icon === 'image'">
          <rect x="3.5" y="4.5" width="17" height="15" rx="2.5" />
          <circle cx="9" cy="9.5" r="1.6" />
          <path d="M4.5 17.5 10 12l3.5 3.5L17 12l3.5 3.5" />
        </template>
        <!-- 对话 -->
        <template v-else-if="t.icon === 'chat'">
          <path d="M21 11.5a7.5 7.5 0 0 1-7.5 7.5c-1.2 0-2.4-.3-3.4-.8L4 20l1.6-4.2A7.5 7.5 0 1 1 21 11.5z" />
        </template>
        <!-- 我的 -->
        <template v-else>
          <circle cx="12" cy="8" r="3.6" />
          <path d="M4.5 20c1.4-3.6 4.2-5.2 7.5-5.2s6.1 1.6 7.5 5.2" />
        </template>
      </svg>
      <span>{{ t.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.tabbar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  display: flex;
  background: var(--card);
  border-top: 2px dashed var(--pink-soft);
  padding: 6px 4px calc(6px + env(safe-area-inset-bottom));
}

.tab {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 4px 0;
  background: none;
  border: none;
  color: var(--ink-soft, #9b8a90);
  font-size: 0.72rem;
  font-family: inherit;
  cursor: pointer;
  transition: color 0.2s;
}

.tab.active {
  color: var(--pink-deep);
  font-weight: 600;
}
</style>

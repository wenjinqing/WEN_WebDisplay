<script setup>
// 留言展示(独占整页,大空间)
import { ref, onMounted, onUnmounted } from 'vue'
import PawPrint from './PawPrint.vue'

const messages = ref([])
const loaded = ref(false)
const expanded = ref(new Set())

async function load() {
  try {
    const res = await fetch('/api/messages')
    messages.value = await res.json()
  } catch {
    /* 静默降级 */
  } finally {
    loaded.value = true
  }
}

function isLong(text) {
  return text.length > 120 || (text.match(/\n/g) || []).length > 3
}

function toggleExpand(i) {
  if (expanded.value.has(i)) expanded.value.delete(i)
  else expanded.value.add(i)
  expanded.value = new Set(expanded.value)
}

onMounted(() => {
  load()
  window.addEventListener('catcafe-msg-posted', load)
})
onUnmounted(() => window.removeEventListener('catcafe-msg-posted', load))
</script>

<template>
  <div class="list-wrap">
    <ul v-if="messages.length" class="msgs">
      <li v-for="(m, i) in messages" :key="i" class="msg">
        <div class="avatar"><PawPrint :size="18" color="#fff" /></div>
        <div class="bubble-wrap">
          <div class="msg-head">
            <b>
              {{ m.nick }}
              <span v-if="m.by === 'agent'" class="keeper-tag">🐷 饲养员</span>
            </b>
            <time>{{ m.time }}</time>
          </div>
          <div class="bubble">
            <p :class="{ clamped: isLong(m.content) && !expanded.has(i) }">{{ m.content }}</p>
            <button v-if="isLong(m.content)" class="expand-btn" @click="toggleExpand(i)">
              {{ expanded.has(i) ? '收起 ▲' : '展开全文 ▼' }}
            </button>
          </div>
          <div v-if="m.reply" class="reply" :class="{ keeper: m.replyBy === '猪咪君君' }">
            <span class="owner-badge font-cute">
              {{ m.replyBy === '猪咪君君' ? '🐷 猪咪饲养员 · 猪咪君君' : '🐾 店长回复' }}
            </span>
            <p>{{ m.reply }}</p>
          </div>
        </div>
      </li>
    </ul>
    <p v-else-if="loaded" class="empty">墙上还空空的,来贴第一张便利贴喵~</p>
  </div>
</template>

<style scoped>
/* 双列网格,大空间展示 */
.msgs {
  list-style: none;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  align-items: start;
}

.msg {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  background: var(--surface-2);
  border: 1px solid var(--pink-pale);
  border-radius: 18px;
  padding: 14px 16px;
}

.avatar {
  flex-shrink: 0;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: linear-gradient(160deg, var(--pink), var(--pink-deep));
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 2px;
}

.bubble-wrap {
  flex: 1;
  min-width: 0;
}

.msg-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 4px;
}

.msg-head b {
  color: var(--pink-deep);
  font-size: 0.88rem;
}

.msg-head time {
  color: var(--muted);
  font-size: 0.75rem;
  flex-shrink: 0;
}

.keeper-tag {
  display: inline-block;
  background: #9b7ede;
  color: #fff;
  font-size: 0.68rem;
  border-radius: 999px;
  padding: 1px 8px;
  margin-left: 6px;
  vertical-align: 1px;
}

.bubble p {
  font-size: 0.92rem;
  line-height: 1.7;
  word-break: break-word;
  white-space: pre-wrap;
  color: var(--ink);
}

.bubble p.clamped {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.expand-btn {
  border: none;
  background: none;
  color: var(--pink-deep);
  font-size: 0.8rem;
  cursor: pointer;
  padding: 4px 0 0;
}

.expand-btn:hover {
  text-decoration: underline;
}

.reply {
  margin-top: 8px;
  background: var(--pink-pale);
  border-left: 3px solid var(--pink);
  border-radius: 0 14px 14px 14px;
  padding: 8px 14px;
}

.owner-badge {
  display: inline-block;
  background: var(--pink);
  color: #fff;
  font-size: 0.72rem;
  border-radius: 999px;
  padding: 2px 10px;
  margin-bottom: 4px;
}

.reply p {
  font-size: 0.86rem;
  color: var(--ink);
  line-height: 1.6;
  word-break: break-word;
}

/* 饲养员回复:淡紫区分店长粉 */
.reply.keeper {
  background: #f0eafd;
  border-left-color: #9b7ede;
}

.reply.keeper .owner-badge {
  background: #9b7ede;
}
.empty {
  text-align: center;
  color: var(--muted);
  padding: 40px 0;
}

@media (max-width: 720px) {
  .msgs {
    grid-template-columns: 1fr;
  }
}
</style>

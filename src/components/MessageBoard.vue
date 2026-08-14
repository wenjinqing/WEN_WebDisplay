<script setup>
import { ref, onMounted } from 'vue'
import PawPrint from './PawPrint.vue'

const messages = ref([])
const nick = ref('')
const content = ref('')
const sending = ref(false)
const error = ref('')
const loaded = ref(false)
const expanded = ref(new Set()) // 展开的长留言

onMounted(load)

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

async function submit() {
  if (!content.value.trim() || sending.value) return
  sending.value = true
  error.value = ''
  try {
    const res = await fetch('/api/messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nick: nick.value, content: content.value }),
    })
    const data = await res.json()
    if (res.ok) {
      messages.value = [data, ...messages.value]
      content.value = ''
      if (nick.value.trim()) localStorage.setItem('catcafe_nick', nick.value.trim())
    } else {
      error.value = data.error || '留言失败,稍后再试'
    }
  } catch {
    error.value = '网络打了个盹,稍后再试'
  } finally {
    sending.value = false
  }
}

function isLong(text) {
  return text.length > 120 || (text.match(/\n/g) || []).length > 3
}

function toggleExpand(i) {
  if (expanded.value.has(i)) expanded.value.delete(i)
  else expanded.value.add(i)
  expanded.value = new Set(expanded.value) // 触发响应
}
</script>

<template>
  <div class="board">
    <h3 class="font-cute">留言板</h3>
    <p class="hint">给猫猫酱留句话吧,催更、夸夸、聊剧情都可以~</p>

    <form class="form" @submit.prevent="submit">
      <input v-model="nick" maxlength="20" placeholder="昵称(不填就是匿名猪咪)" />
      <textarea
        v-model="content"
        maxlength="300"
        rows="3"
        placeholder="写下你的留言……(300字以内)"
        required
      />
      <div class="form-foot">
        <span v-if="error" class="err">{{ error }}</span>
        <span v-else class="count">{{ content.length }}/300</span>
        <button type="submit" class="btn btn-primary" :disabled="sending || !content.trim()">
          {{ sending ? '投递中…' : '贴到墙上' }}
        </button>
      </div>
    </form>

    <ul class="msgs">
      <li v-for="(m, i) in messages" :key="i" class="msg">
        <div class="avatar"><PawPrint :size="18" color="#fff" /></div>
        <div class="bubble-wrap">
          <div class="msg-head">
            <b>{{ m.nick }}</b>
            <time>{{ m.time }}</time>
          </div>
          <div class="bubble">
            <p :class="{ clamped: isLong(m.content) && !expanded.has(i) }">{{ m.content }}</p>
            <button v-if="isLong(m.content)" class="expand-btn" @click="toggleExpand(i)">
              {{ expanded.has(i) ? '收起 ▲' : '展开全文 ▼' }}
            </button>
          </div>
          <div v-if="m.reply" class="reply">
            <span class="owner-badge font-cute">🐾 店长回复</span>
            <p>{{ m.reply }}</p>
          </div>
        </div>
      </li>
      <li v-if="loaded && !messages.length" class="empty">墙上还空空的,来贴第一张便利贴喵~</li>
    </ul>
  </div>
</template>

<style scoped>
.board {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 32px 26px;
  height: 100%;
}

h3 {
  font-size: 1.3rem;
  color: var(--pink-deep);
  margin-bottom: 4px;
  text-align: center;
}

.hint {
  color: var(--muted);
  font-size: 0.9rem;
  text-align: center;
  margin-bottom: 20px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 22px;
}

input,
textarea {
  border: 2px solid var(--pink-pale);
  border-radius: 14px;
  padding: 10px 16px;
  font-size: 0.95rem;
  font-family: inherit;
  color: var(--ink);
  background: #fff;
  outline: none;
  transition: border-color 0.2s;
  resize: vertical;
}

input:focus,
textarea:focus {
  border-color: var(--pink);
}

.form-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.count {
  color: var(--muted);
  font-size: 0.8rem;
}

.err {
  color: var(--pink-deep);
  font-size: 0.85rem;
}

/* ===== 留言气泡 ===== */
.msgs {
  list-style: none;
  max-height: 420px;
  overflow-y: auto;
  padding-right: 4px;
}

.msg {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
  align-items: flex-start;
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
  padding: 0 4px;
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

.bubble {
  background: #fff9f4;
  border: 1px solid var(--pink-pale);
  border-radius: 4px 18px 18px 18px; /* 气泡小尾巴角 */
  padding: 10px 16px;
}

.bubble p {
  font-size: 0.92rem;
  line-height: 1.7;
  word-break: break-word;
  white-space: pre-wrap;
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

/* 店长回复:高亮嵌套气泡 */
.reply {
  margin: 8px 0 2px 14px;
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

.empty {
  text-align: center;
  color: var(--muted);
  font-size: 0.9rem;
  padding: 24px 0;
}
</style>

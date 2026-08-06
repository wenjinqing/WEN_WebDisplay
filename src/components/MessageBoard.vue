<script setup>
import { ref, onMounted } from 'vue'

const messages = ref([])
const nick = ref('')
const content = ref('')
const sending = ref(false)
const error = ref('')
const loaded = ref(false)

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
    } else {
      error.value = data.error || '留言失败,稍后再试'
    }
  } catch {
    error.value = '网络打了个盹,稍后再试'
  } finally {
    sending.value = false
  }
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
        <div class="msg-head">
          <b>{{ m.nick }}</b>
          <time>{{ m.time }}</time>
        </div>
        <p>{{ m.content }}</p>
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
  padding: 36px 28px;
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
  margin-bottom: 24px;
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

.msgs {
  list-style: none;
  max-height: 380px;
  overflow-y: auto;
}

.msg {
  background: #fff9f4;
  border: 1px solid var(--pink-pale);
  border-radius: 14px;
  padding: 12px 16px;
  margin-bottom: 10px;
}

.msg-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 4px;
}

.msg-head b {
  color: var(--pink-deep);
  font-size: 0.9rem;
}

.msg-head time {
  color: var(--muted);
  font-size: 0.78rem;
  flex-shrink: 0;
}

.msg p {
  font-size: 0.92rem;
  line-height: 1.7;
  word-break: break-word;
  white-space: pre-wrap;
}

.empty {
  text-align: center;
  color: var(--muted);
  font-size: 0.9rem;
  padding: 24px 0;
}
</style>

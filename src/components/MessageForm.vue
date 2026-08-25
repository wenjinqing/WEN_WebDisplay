<script setup>
// 写留言(独立页面,宽敞的表单)
import { ref } from 'vue'

const nick = ref(localStorage.getItem('catcafe_nick') || '')
const content = ref('')
const sending = ref(false)
const error = ref('')
const done = ref(false)

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
      content.value = ''
      done.value = true
      setTimeout(() => (done.value = false), 3000)
      if (nick.value.trim()) localStorage.setItem('catcafe_nick', nick.value.trim())
      window.dispatchEvent(new CustomEvent('catcafe-msg-posted'))
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
  <form class="form" @submit.prevent="submit">
    <label class="field">
      <span class="label">你的昵称</span>
      <input v-model="nick" maxlength="20" placeholder="不填就是匿名猪咪" />
    </label>
    <label class="field">
      <span class="label">留言内容</span>
      <textarea
        v-model="content"
        maxlength="300"
        rows="6"
        placeholder="催更、夸夸、聊剧情都可以……(300字以内)"
        required
      />
    </label>
    <div class="form-foot">
      <span v-if="error" class="err">{{ error }}</span>
      <span v-else-if="done" class="ok">已贴到墙上啦 ✅ 去「留言墙」看看</span>
      <span v-else class="count">{{ content.length }}/300</span>
      <button type="submit" class="btn btn-primary" :disabled="sending || !content.trim()">
        {{ sending ? '投递中…' : '贴到墙上' }}
      </button>
    </div>
  </form>
</template>

<style scoped>
.form {
  max-width: 640px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.label {
  font-size: 0.88rem;
  color: var(--muted);
  padding-left: 6px;
}

input,
textarea {
  border: 2px solid var(--pink-pale);
  border-radius: 16px;
  padding: 12px 18px;
  font-size: 1rem;
  font-family: inherit;
  color: var(--ink);
  background: var(--surface-2);
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

.ok {
  color: #4caf7d;
  font-size: 0.88rem;
}
</style>

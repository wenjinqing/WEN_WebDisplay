<script setup>
import { ref, computed, onMounted } from 'vue'
import PawPrint from './PawPrint.vue'

const props = defineProps({
  novel: { type: Object, required: true },
})
const emit = defineEmits(['close'])

const list = ref([])
const nick = ref('')
const content = ref('')
const score = ref(5)
const sending = ref(false)
const error = ref('')

const avg = computed(() => {
  if (!list.value.length) return 0
  return (list.value.reduce((s, c) => s + c.score, 0) / list.value.length).toFixed(1)
})

onMounted(load)

async function load() {
  try {
    const res = await fetch(`/api/comments?file=${encodeURIComponent(props.novel.file)}`)
    list.value = await res.json()
  } catch {
    /* 静默 */
  }
}

async function submit() {
  if (!content.value.trim() || sending.value) return
  sending.value = true
  error.value = ''
  try {
    const res = await fetch('/api/comments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        file: props.novel.file,
        nick: nick.value,
        content: content.value,
        score: score.value,
      }),
    })
    const data = await res.json()
    if (res.ok) {
      list.value = [data, ...list.value]
      content.value = ''
    } else {
      error.value = data.error || '评论失败,稍后再试'
    }
  } catch {
    error.value = '网络打了个盹,稍后再试'
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <div class="mask" @click.self="emit('close')">
    <div class="box" role="dialog" :aria-label="novel.title + '的评论'">
      <header class="head">
        <div>
          <h3 class="font-cute">{{ novel.title }}</h3>
          <p class="avg">
            <template v-if="list.length">综合 {{ avg }} 分 · {{ list.length }} 条评论</template>
            <template v-else>还没有评论,来抢沙发~</template>
          </p>
        </div>
        <button class="close" aria-label="关闭" @click="emit('close')">✕</button>
      </header>

      <form class="form" @submit.prevent="submit">
        <div class="row">
          <input v-model="nick" maxlength="20" placeholder="昵称(可空)" class="nick" />
          <div class="paws-picker">
            <button
              v-for="n in 5"
              :key="n"
              type="button"
              :class="{ on: score >= n }"
              :aria-label="`${n}分`"
              @click="score = n"
            >
              <PawPrint :size="20" :color="score >= n ? '#f9718f' : '#e8d5dc'" />
            </button>
          </div>
        </div>
        <textarea v-model="content" maxlength="300" rows="2" placeholder="说说你的读后感……" required />
        <div class="form-foot">
          <span class="err">{{ error }}</span>
          <button type="submit" class="btn btn-primary" :disabled="sending || !content.trim()">
            {{ sending ? '发表中…' : '发表评论' }}
          </button>
        </div>
      </form>

      <ul class="clist">
        <li v-for="(cm, i) in list" :key="i" class="cmt">
          <div class="cmt-head">
            <b>{{ cm.nick }}</b>
            <span class="cmt-paws">
              <PawPrint v-for="n in cm.score" :key="n" :size="12" color="#f9718f" />
            </span>
            <time>{{ cm.time }}</time>
          </div>
          <p>{{ cm.content }}</p>
        </li>
      </ul>
    </div>
  </div>
</template>

<style scoped>
.mask {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(91, 58, 71, 0.45);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.box {
  width: min(640px, 100%);
  max-height: 86vh;
  background: var(--card);
  border-radius: var(--radius);
  border: 2px solid var(--pink-soft);
  box-shadow: 0 24px 64px rgba(91, 58, 71, 0.3);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  padding: 18px 22px 12px;
  border-bottom: 2px dashed var(--pink-pale);
}

.head h3 {
  color: var(--ink);
  font-size: 1.05rem;
}

.avg {
  color: var(--pink-deep);
  font-size: 0.85rem;
  margin-top: 2px;
}

.close {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: none;
  background: var(--pink-pale);
  color: var(--pink-deep);
  cursor: pointer;
  flex-shrink: 0;
}

.form {
  padding: 14px 22px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  border-bottom: 2px dashed var(--pink-pale);
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.nick {
  width: 160px;
}

input, textarea {
  border: 2px solid var(--pink-pale);
  border-radius: 12px;
  padding: 9px 14px;
  font-size: 0.92rem;
  font-family: inherit;
  color: var(--ink);
  outline: none;
  resize: none;
}

input:focus, textarea:focus {
  border-color: var(--pink);
}

.paws-picker button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 2px;
  transition: transform 0.15s;
}

.paws-picker button:hover {
  transform: scale(1.2);
}

.form-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.err {
  color: var(--pink-deep);
  font-size: 0.82rem;
}

.clist {
  list-style: none;
  overflow-y: auto;
  padding: 14px 22px 20px;
}

.cmt {
  background: #fff9f4;
  border: 1px solid var(--pink-pale);
  border-radius: 12px;
  padding: 10px 14px;
  margin-bottom: 10px;
}

.cmt-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 3px;
}

.cmt-head b {
  color: var(--pink-deep);
  font-size: 0.88rem;
}

.cmt-paws {
  display: inline-flex;
}

.cmt-head time {
  margin-left: auto;
  color: var(--muted);
  font-size: 0.75rem;
}

.cmt p {
  font-size: 0.9rem;
  line-height: 1.6;
  word-break: break-word;
}

@media (max-width: 720px) {
  .box {
    max-height: 92vh;
  }
  .row {
    flex-direction: column;
    align-items: stretch;
  }
  .nick {
    width: 100%;
  }
  .paws-picker {
    text-align: center;
  }
}
</style>

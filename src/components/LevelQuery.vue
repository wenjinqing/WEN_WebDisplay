<script setup>
import { ref } from 'vue'

const nick = ref('')
const result = ref(null)
const error = ref('')
const querying = ref(false)

async function query() {
  if (!nick.value.trim() || querying.value) return
  querying.value = true
  error.value = ''
  result.value = null
  try {
    const res = await fetch(`/api/points?nick=${encodeURIComponent(nick.value.trim())}`)
    const data = await res.json()
    if (res.ok) {
      result.value = data
    } else {
      error.value = data.error || '查询失败'
    }
  } catch {
    error.value = '网络打了个盹'
  } finally {
    querying.value = false
  }
}
</script>

<template>
  <div class="level-card">
    <h4 class="font-cute">🐟 我的猪咪头衔</h4>
    <p class="tip">留言+5 · 评论+5 · 寄明信片+10 · 催更+2(填昵称才累计哦)</p>
    <div class="query-row">
      <input v-model="nick" maxlength="20" placeholder="输入你常用的昵称" @keyup.enter="query" />
      <button class="btn btn-primary" :disabled="querying" @click="query">查询</button>
    </div>
    <p v-if="error" class="err">{{ error }}</p>

    <div v-if="result" class="result">
      <div class="title-line">
        <span class="my-title font-cute">{{ result.title }}</span>
        <span class="my-points">{{ result.points }} 鱼干</span>
      </div>
      <div v-if="result.nextAt" class="progress">
        <div class="bar" :style="{ width: Math.min(100, (result.points / result.nextAt) * 100) + '%' }" />
      </div>
      <p class="next-tip">
        {{ result.nextAt ? `再攒 ${result.nextAt - result.points} 鱼干升级!` : '你已经是本店最高头衔啦,呱唧呱唧!' }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.level-card {
  max-width: 760px;
  margin: 24px auto 0;
  background: var(--card);
  border: 2px dashed var(--pink-soft);
  border-radius: 255px 16px 225px 16px / 16px 225px 16px 255px;
  box-shadow: var(--shadow);
  padding: 24px 32px;
  text-align: center;
}

h4 {
  color: var(--pink-deep);
  font-size: 1.1rem;
  margin-bottom: 4px;
}

.tip {
  color: var(--muted);
  font-size: 0.8rem;
  margin-bottom: 14px;
}

.query-row {
  display: flex;
  gap: 10px;
  max-width: 420px;
  margin: 0 auto;
}

.query-row input {
  flex: 1;
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 9px 18px;
  font-size: 0.95rem;
  font-family: inherit;
  color: var(--ink);
  outline: none;
}

.query-row input:focus {
  border-color: var(--pink);
}

.err {
  color: var(--pink-deep);
  font-size: 0.85rem;
  margin-top: 10px;
}

.result {
  margin-top: 18px;
}

.title-line {
  display: flex;
  justify-content: center;
  align-items: baseline;
  gap: 14px;
}

.my-title {
  font-size: 1.6rem;
  color: var(--pink-deep);
}

.my-points {
  color: var(--muted);
  font-size: 0.9rem;
}

.progress {
  max-width: 360px;
  height: 12px;
  margin: 12px auto 6px;
  background: var(--pink-pale);
  border-radius: 999px;
  overflow: hidden;
}

.bar {
  height: 100%;
  background: linear-gradient(90deg, var(--pink-soft), var(--pink));
  border-radius: 999px;
  transition: width 0.6s ease;
}

.next-tip {
  color: var(--muted);
  font-size: 0.82rem;
}
</style>

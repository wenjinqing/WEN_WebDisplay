<script setup>
import { ref } from 'vue'

const nick = ref('')
const result = ref(null)
const special = ref(null) // 店长专属头衔彩蛋
const error = ref('')
const querying = ref(false)

// 彩蛋:店长专属昵称 → 专属头衔
const SPECIAL_TITLES = {
  '爱丽丝': { title: '店长本人?', desc: '这个味道……是本尊来了吗', icon: '👑' },
  '爱丽丝猫猫酱': { title: '猫咖店长', desc: '欢迎店长巡视自家店铺!', icon: '🎀' },
  '小涩猫爱丽丝': { title: '小涩猫本猫', desc: '喵呜,本体出现,全体注目!', icon: '✨' },
  '猪咪爱丽丝': { title: '猪咪团团长', desc: '全体猪咪听令,列队欢迎!', icon: '🐷' },
}

async function query() {
  if (!nick.value.trim() || querying.value) return
  querying.value = true
  error.value = ''
  result.value = null
  const name = nick.value.trim()
  special.value = SPECIAL_TITLES[name] || null
  try {
    const res = await fetch(`/api/points?nick=${encodeURIComponent(name)}`)
    const data = await res.json()
    if (res.ok) {
      result.value = data
      localStorage.setItem('catcafe_nick', name) // 记住昵称,捡鱼干用
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
    <h4 class="font-cute">我的猪咪头衔</h4>
    <p class="tip">留言+5 · 评论+5 · 寄明信片+10 · 催更+2(填昵称才累计哦)</p>
    <div class="query-row">
      <input v-model="nick" maxlength="20" placeholder="输入你常用的昵称" @keyup.enter="query" />
      <button class="btn btn-primary" :disabled="querying" @click="query">查询</button>
    </div>
    <p v-if="error" class="err">{{ error }}</p>

    <div v-if="result && special" class="result special-card">
      <div class="crown-icon">{{ special.icon }}</div>
      <div class="title-line">
        <span class="my-title special-title font-cute">{{ special.title }}</span>
      </div>
      <p class="special-desc">{{ special.desc }}</p>
      <p class="next-tip">店长的鱼干罐:∞(店长有吃不完的鱼干)</p>
    </div>

    <div v-else-if="result" class="result">
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

/* 店长专属头衔彩蛋卡 */
.special-card {
  background: linear-gradient(135deg, #fff8e7, #ffeef4);
  border: 2px solid #f0c86e;
  border-radius: 16px;
  padding: 18px 20px;
  animation: goldGlow 1.8s ease-in-out infinite;
}

.crown-icon {
  font-size: 1.8rem;
  margin-bottom: 4px;
}

.special-title {
  background: linear-gradient(90deg, #d4a017, #e85d7f);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.special-desc {
  color: var(--ink);
  font-size: 0.95rem;
  margin: 6px 0;
}

@keyframes goldGlow {
  0%, 100% { box-shadow: 0 0 8px rgba(240, 200, 110, 0.4); }
  50% { box-shadow: 0 0 20px rgba(240, 200, 110, 0.8); }
}
</style>

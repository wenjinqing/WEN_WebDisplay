<script setup>
import { ref } from 'vue'
import PigmiFace from './PigmiFace.vue'
import SectionTitle from './SectionTitle.vue'
import LevelQuery from './LevelQuery.vue'
import { site } from '../data.js'

// 入群验证:答对任意一部作品名或主角名即可看到群号
// 校验本地存档:必须含数字(防止打码值被误存的脏数据)
const savedQQ = localStorage.getItem('catcafe_gate_qq') || ''
const unlocked = ref(/\d{5,}/.test(savedQQ) ? savedQQ : '')
if (savedQQ && !/\d{5,}/.test(savedQQ)) localStorage.removeItem('catcafe_gate_qq')
const answer = ref('')
const checking = ref(false)
const wrong = ref(false)

async function verify() {
  if (!answer.value.trim() || checking.value) return
  checking.value = true
  wrong.value = false
  try {
    const res = await fetch('/api/gate/verify', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ answer: answer.value }),
    })
    const data = await res.json()
    if (data.pass) {
      unlocked.value = data.qq
      localStorage.setItem('catcafe_ach_gate', '1')
      localStorage.setItem('catcafe_gate_qq', data.qq)
    } else {
      wrong.value = true
    }
  } catch {
    wrong.value = false
  } finally {
    checking.value = false
  }
}
</script>

<template>
  <section id="pigmi">
    <div class="container">
      <SectionTitle title="猪咪聚集地" sub="本店常客的专属小窝" />

      <div class="club-card hd-card" v-reveal>
        <div class="pigmis">
          <PigmiFace :size="90" />
          <PigmiFace :size="70" style="margin-top: 24px" />
          <PigmiFace :size="80" style="margin-top: 8px" />
        </div>
        <div class="club-info">
          <h3 class="font-cute">{{ site.fanClub.name }}</h3>
          <p>{{ site.fanClub.desc }}</p>

          <!-- 已解锁:显示群号 + 聚集地入口 -->
          <template v-if="unlocked">
            <p class="qq font-cute">🐾 {{ unlocked }}</p>
            <a href="/hub" class="btn btn-primary hub-link">🏠 进入猪咪聚集地 →</a>
          </template>

          <!-- 未解锁:答题验证 -->
          <div v-else class="gate">
            <p class="gate-tip">🔒 群号藏起来了!答对<strong>任意一部作品名或主角名</strong>即可解锁</p>
            <div class="gate-row">
              <input
                v-model="answer"
                placeholder="例如:《…》或主角名字"
                @keyup.enter="verify"
              />
              <button class="btn btn-primary" :disabled="checking || !answer.trim()" @click="verify">
                {{ checking ? '验证中…' : '解锁' }}
              </button>
            </div>
            <p v-if="wrong" class="wrong">不对哦,真猪咪再想想~</p>
          </div>
        </div>
      </div>

      <LevelQuery v-reveal />
    </div>
  </section>
</template>

<style scoped>
.club-card {
  display: flex;
  align-items: center;
  gap: 40px;
  max-width: 760px;
  margin: 0 auto;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 40px;
}

.pigmis {
  display: flex;
  align-items: flex-end;
  gap: 6px;
  flex-shrink: 0;
}

.club-info h3 {
  font-size: 1.4rem;
  color: var(--pink-deep);
  margin-bottom: 10px;
}

.club-info p {
  color: var(--muted);
  margin-bottom: 10px;
}

.qq {
  color: var(--ink) !important;
  font-size: 1.05rem;
}

.gate {
  margin-top: 10px;
}

.gate-tip {
  color: var(--muted);
  font-size: 0.85rem;
  margin-bottom: 10px;
}

.gate-tip strong {
  color: var(--pink-deep);
}

.gate-row {
  display: flex;
  gap: 10px;
}

.gate-row input {
  flex: 1;
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 10px 18px;
  font-size: 0.92rem;
  font-family: inherit;
  color: var(--ink);
  outline: none;
}

.gate-row input:focus {
  border-color: var(--pink);
}

.wrong {
  color: var(--pink-deep);
  font-size: 0.82rem;
  margin-top: 8px;
}

.hub-link {
  margin-top: 12px;
  display: inline-flex;
}

@media (max-width: 720px) {
  .club-card {
    flex-direction: column;
    padding: 28px;
    text-align: center;
  }
  .pigmis {
    justify-content: center;
  }
}
</style>

<script setup>
import { ref, onMounted } from 'vue'
import PawPrint from './PawPrint.vue'

const total = ref(0)
const recent = ref([])
const sending = ref(false)
const pressed = ref(false)
const toast = ref('')

onMounted(async () => {
  try {
    const res = await fetch('/api/urge')
    const data = await res.json()
    total.value = data.total || 0
    recent.value = data.recent || []
  } catch {
    /* API 未就绪时静默降级,只显示静态墙 */
  }
})

async function urge() {
  if (sending.value) return
  sending.value = true
  pressed.value = true
  setTimeout(() => (pressed.value = false), 300)
  try {
    const res = await fetch('/api/urge', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}' })
    const data = await res.json()
    if (res.ok) {
      total.value = data.total
      recent.value = [data.urge, ...recent.value].slice(0, 20)
      showToast('催更成功!猫猫酱收到啦 🐾')
    } else {
      showToast(data.error || '催更失败,稍后再试')
    }
  } catch {
    showToast('网络打了个盹,稍后再试')
  } finally {
    sending.value = false
  }
}

function showToast(text) {
  toast.value = text
  setTimeout(() => (toast.value = ''), 2500)
}

// 相对时间:刚刚 / N小时前 / N天前
function ago(timeStr) {
  const t = new Date(timeStr.replace(' ', 'T'))
  const diff = Date.now() - t.getTime()
  const min = Math.floor(diff / 60000)
  if (min < 1) return '刚刚'
  if (min < 60) return `${min}分钟前`
  const hours = Math.floor(min / 60)
  if (hours < 24) return `${hours}小时前`
  return `${Math.floor(hours / 24)}天前`
}
</script>

<template>
  <div class="urge-wall">
    <h3 class="font-cute">催更墙</h3>
    <p class="hint">猫猫酱已被催更</p>

    <div class="counter-wrap">
      <div class="counter font-cute">{{ total }}<span>次</span></div>
      <div class="paw-trail" aria-hidden="true">
        <PawPrint :size="16" :rotate="-20" /><PawPrint :size="20" :rotate="10" /><PawPrint :size="16" :rotate="25" />
      </div>
    </div>

    <button class="urge-btn" :class="{ pressed }" :disabled="sending" @click="urge">
      <PawPrint :size="36" color="#fff" />
      拍爪催更!
    </button>

    <transition name="fade">
      <p v-if="toast" class="toast">{{ toast }}</p>
    </transition>

    <div v-if="recent.length" class="recent">
      <span v-for="(u, i) in recent.slice(0, 4)" :key="i" class="chip">
        <b>{{ u.nick }}</b> · {{ ago(u.time) }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.urge-wall {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 32px 26px;
  text-align: center;
  height: 100%;
  display: flex;
  flex-direction: column;
}

h3 {
  font-size: 1.3rem;
  color: var(--pink-deep);
  margin-bottom: 4px;
}

.hint {
  color: var(--muted);
  font-size: 0.9rem;
}

.counter-wrap {
  margin: 10px 0 20px;
}

.counter {
  font-size: 3.4rem;
  color: var(--ink);
  line-height: 1.1;
}

.counter span {
  font-size: 1.1rem;
  color: var(--muted);
  margin-left: 4px;
}

.paw-trail {
  display: flex;
  justify-content: center;
  gap: 6px;
  margin-top: 6px;
  color: var(--pink-soft);
}

.urge-btn {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  align-self: center;
  border: none;
  cursor: pointer;
  background: linear-gradient(160deg, var(--pink), var(--pink-deep));
  color: #fff;
  font-size: 1.05rem;
  font-weight: 500;
  padding: 20px 34px;
  border-radius: 26px;
  box-shadow: 0 10px 24px rgba(249, 113, 143, 0.4);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.urge-btn:hover {
  transform: translateY(-2px) scale(1.02);
}

.urge-btn.pressed {
  transform: scale(0.94);
  box-shadow: 0 4px 12px rgba(249, 113, 143, 0.4);
}

.toast {
  margin-top: 14px;
  color: var(--pink-deep);
  font-size: 0.9rem;
}

/* 最近催更:小胶囊 */
.recent {
  margin-top: auto;
  padding-top: 18px;
  border-top: 2px dashed var(--pink-pale);
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
}

.chip {
  background: var(--pink-pale);
  color: var(--muted);
  border-radius: 999px;
  padding: 3px 12px;
  font-size: 0.75rem;
}

.chip b {
  color: var(--pink-deep);
  font-weight: 500;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

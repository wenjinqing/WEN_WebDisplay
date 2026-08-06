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
</script>

<template>
  <div class="urge-wall">
    <h3 class="font-cute">催更墙</h3>
    <p class="hint">猫猫酱已被催更</p>
    <div class="counter font-cute">{{ total }}<span>次</span></div>

    <button class="urge-btn" :class="{ pressed }" :disabled="sending" @click="urge">
      <PawPrint :size="40" color="#fff" />
      拍爪催更!
    </button>

    <transition name="fade">
      <p v-if="toast" class="toast">{{ toast }}</p>
    </transition>

    <ul v-if="recent.length" class="recent">
      <li v-for="(u, i) in recent.slice(0, 5)" :key="i">
        <b>{{ u.nick }}</b> 拍了拍猫猫酱 · {{ u.time }}
      </li>
    </ul>
  </div>
</template>

<style scoped>
.urge-wall {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 36px 28px;
  text-align: center;
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

.counter {
  font-size: 3rem;
  color: var(--ink);
  line-height: 1.2;
  margin: 8px 0 20px;
}

.counter span {
  font-size: 1.1rem;
  color: var(--muted);
  margin-left: 4px;
}

.urge-btn {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  border: none;
  cursor: pointer;
  background: linear-gradient(160deg, var(--pink), var(--pink-deep));
  color: #fff;
  font-size: 1.05rem;
  font-weight: 500;
  padding: 22px 36px;
  border-radius: 28px;
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

.recent {
  list-style: none;
  margin-top: 20px;
  border-top: 2px dashed var(--pink-pale);
  padding-top: 14px;
}

.recent li {
  color: var(--muted);
  font-size: 0.85rem;
  padding: 3px 0;
}

.recent b {
  color: var(--ink);
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

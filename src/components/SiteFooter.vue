<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'

const visits = ref(0)
const online = ref(0)
let timer = null

async function heartbeat() {
  try {
    const res = await fetch('/api/online')
    online.value = (await res.json()).online || 0
  } catch {
    /* 静默 */
  }
}

onMounted(async () => {
  try {
    const res = await fetch('/api/hit')
    const data = await res.json()
    visits.value = data.total || 0
  } catch {
    /* 统计失败不影响页面 */
  }
  heartbeat()
  timer = setInterval(heartbeat, 30000) // 每 30 秒报到一次
})

onUnmounted(() => clearInterval(timer))
</script>

<template>
  <footer class="footer">
    <div class="container footer-inner">
      <PawPrint :size="28" color="#ffc2d4" />
      <p class="font-cute">{{ site.title }}</p>
      <p class="note">
        本站为「{{ site.author }}」的非官方粉丝后援页 · 作品版权归原作者所有 · 用爱发电 💗
      </p>
      <p class="note">
        插画素材:<a class="credit" href="https://www.irasutoya.com/" target="_blank" rel="noopener">いらすとや</a>(免费授权使用)
      </p>
      <p v-if="visits" class="note">☕ 已有 {{ visits }} 只猪咪来喝过咖啡</p>
      <p v-if="online" class="note online-note">
        <span class="dot" />现在有 {{ online }} 只猪咪在店里
      </p>
      <a href="/go/pixiv.html" target="_blank" rel="noopener" class="pixiv">
        作者 P 站主页 ↗
      </a>
      <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener" class="beian">
        闽ICP备2026031516号
      </a>
      <a href="/admin" class="pixiv admin-link">店主通道</a>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  background: var(--bg-deep);
  border-top: 2px dashed var(--pink-soft);
  padding: 40px 0;
}

.footer-inner {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.footer-inner > p:first-of-type {
  font-size: 1.1rem;
  color: var(--ink);
}

.note {
  font-size: 0.85rem;
  color: var(--muted);
}

.pixiv {
  color: var(--pink-deep);
  text-decoration: none;
  font-size: 0.9rem;
  margin-top: 4px;
}

.pixiv:hover {
  text-decoration: underline;
}

.admin-link {
  color: var(--muted);
  font-size: 0.8rem;
  opacity: 0.7;
}

.beian {
  color: var(--muted);
  font-size: 0.8rem;
  text-decoration: none;
  margin-top: 6px;
}

.beian:hover {
  color: var(--pink-deep);
  text-decoration: underline;
}

.online-note {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #4caf7d;
  animation: blink 1.6s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.credit {
  color: var(--muted);
  text-decoration: none;
}

.credit:hover {
  color: var(--pink-deep);
  text-decoration: underline;
}
</style>

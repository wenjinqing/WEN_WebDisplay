<script setup>
import { ref, onMounted } from 'vue'
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'

const visits = ref(0)

onMounted(async () => {
  try {
    const res = await fetch('/api/hit')
    const data = await res.json()
    visits.value = data.total || 0
  } catch {
    /* 统计失败不影响页面 */
  }
})
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
      <a :href="site.authorPixiv" target="_blank" rel="noopener" class="pixiv">
        作者 P 站主页 ↗
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

.credit {
  color: var(--muted);
  text-decoration: none;
}

.credit:hover {
  color: var(--pink-deep);
  text-decoration: underline;
}
</style>

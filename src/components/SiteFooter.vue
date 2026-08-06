<script setup>
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'
</script>

<template>
  <footer class="footer">
    <div class="container footer-inner">
      <PawPrint :size="28" color="#ffc2d4" />
      <p class="font-cute">{{ site.title }}</p>
      <p class="note">
        本站为「{{ site.author }}」的非官方粉丝后援页 · 作品版权归原作者所有 · 用爱发电 💗
      </p>
      <a :href="site.authorPixiv" target="_blank" rel="noopener" class="pixiv">
        作者 P 站主页 ↗
      </a>
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
</style>

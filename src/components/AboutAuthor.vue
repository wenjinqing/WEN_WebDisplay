<script setup>
import CatFace from './CatFace.vue'
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'
</script>

<template>
  <section id="author">
    <div class="container">
      <SectionTitle title="关于猫猫酱" sub="本店店长 · 百合小说製造机" />

      <div class="card hd-card" v-reveal>
        <div class="avatar">
          <CatFace :size="140" />
          <span class="name font-cute">{{ site.author }}</span>
        </div>
        <div class="bio">
          <p v-for="(p, i) in site.aboutAuthor" :key="i">{{ p }}</p>
          <a :href="site.authorPixiv" target="_blank" rel="noopener" class="btn btn-primary">
            🎨 去 P 站看她的主页 ↗
          </a>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.card {
  display: flex;
  gap: 40px;
  align-items: center;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 40px;
  max-width: 860px;
  margin: 0 auto;
}

.avatar {
  flex-shrink: 0;
  text-align: center;
}

.name {
  display: block;
  margin-top: 8px;
  font-size: 1.1rem;
  color: var(--pink-deep);
}

.bio p {
  margin-bottom: 14px;
  color: var(--ink);
}

.bio .btn {
  margin-top: 10px;
}

@media (max-width: 720px) {
  .card {
    flex-direction: column;
    padding: 28px;
    text-align: center;
  }
}
</style>

<script setup>
import PigmiFace from './PigmiFace.vue'
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'
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
          <p class="qq font-cute">🐾 {{ site.fanClub.qq }}</p>
        </div>
      </div>
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

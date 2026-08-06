<script setup>
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'
</script>

<template>
  <section id="gallery">
    <div class="container">
      <h2 class="section-title" v-reveal>插画墙</h2>
      <p class="section-sub" v-reveal>猫猫酱的画作展示区 · 上架前先取得本人授权</p>

      <div class="wall" v-reveal>
        <figure v-for="(g, i) in site.gallery" :key="i" class="polaroid" :class="`tilt-${i % 3}`">
          <div class="photo">
            <img v-if="g.img" :src="`/gallery/${g.img}`" :alt="g.title" loading="lazy" />
            <div v-else class="placeholder">
              <PawPrint :size="40" color="#f0c9d6" />
              <span>照片冲洗中…</span>
            </div>
          </div>
          <figcaption>
            <b>{{ g.title }}</b>
            <span>{{ g.note }}</span>
          </figcaption>
        </figure>
      </div>
    </div>
  </section>
</template>

<style scoped>
#gallery {
  background: var(--bg-deep);
}

.wall {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  max-width: 900px;
  margin: 0 auto;
}

.polaroid {
  background: #fff;
  border-radius: 8px;
  padding: 12px 12px 14px;
  box-shadow: 0 8px 20px rgba(233, 93, 127, 0.15);
  transition: transform 0.25s ease;
}

.polaroid:hover {
  transform: rotate(0deg) scale(1.03) !important;
}

.tilt-0 { transform: rotate(-2deg); }
.tilt-1 { transform: rotate(1.5deg); }
.tilt-2 { transform: rotate(-1deg); }

.photo {
  aspect-ratio: 4 / 3;
  border-radius: 6px;
  overflow: hidden;
  background: var(--pink-pale);
}

.photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #c98aa0;
  font-size: 0.85rem;
}

figcaption {
  padding-top: 10px;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

figcaption b {
  font-size: 0.95rem;
  color: var(--ink);
}

figcaption span {
  font-size: 0.8rem;
  color: var(--muted);
}

@media (max-width: 720px) {
  .wall {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
}
</style>

<script setup>
import { ref } from 'vue'
import PawPrint from './PawPrint.vue'
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'

const zoomImg = ref(null) // 双击放大的图片
</script>

<template>
  <section id="gallery">
    <div class="container">
      <SectionTitle title="插画墙" sub="猫猫酱的画作展示区 · 上架前先取得本人授权" />

      <div class="wall" v-reveal>
        <figure v-for="(g, i) in site.gallery" :key="i" class="polaroid" :class="`tilt-${i % 3}`">
          <div class="photo">
            <img
              v-if="g.img"
              :src="`/gallery/${g.img}`"
              :alt="g.title"
              loading="lazy"
              class="zoomable"
              @dblclick.prevent="zoomImg = `/gallery/${g.img}`"
              @contextmenu.prevent
            />
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

    <!-- 双击放大灯箱 -->
    <div v-if="zoomImg" class="lightbox" @click="zoomImg = null">
      <img :src="zoomImg" alt="放大预览" @contextmenu.prevent @dragstart.prevent />
      <p class="tip">双击打开 · 点任意处关闭</p>
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
  position: relative;
  background: #fff;
  border-radius: 8px;
  padding: 12px 12px 14px;
  box-shadow: 0 8px 20px rgba(233, 93, 127, 0.15);
  transition: transform 0.25s ease;
}

/* 和纸胶带 */
.polaroid::before {
  content: '';
  position: absolute;
  top: -11px;
  left: 50%;
  width: 78px;
  height: 22px;
  background: rgba(249, 113, 143, 0.35);
  border-left: 2px dashed rgba(255, 255, 255, 0.7);
  border-right: 2px dashed rgba(255, 255, 255, 0.7);
  z-index: 1;
}

.tilt-0 { transform: rotate(-2deg); }
.tilt-0::before { transform: translateX(-50%) rotate(-4deg); }
.tilt-1 { transform: rotate(1.5deg); }
.tilt-1::before { transform: translateX(-50%) rotate(3deg); }
.tilt-2 { transform: rotate(-1deg); }
.tilt-2::before { transform: translateX(-50%) rotate(-2deg); }

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

.zoomable {
  cursor: zoom-in;
}

/* 双击放大灯箱 */
.lightbox {
  position: fixed;
  inset: 0;
  z-index: 110;
  background: rgba(91, 58, 71, 0.75);
  backdrop-filter: blur(6px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  cursor: zoom-out;
  animation: lbIn 0.2s ease;
}

.lightbox img {
  max-width: 92vw;
  max-height: 82vh;
  object-fit: contain;
  border-radius: 12px;
  border: 4px solid #fff;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.35);
  user-select: none;
  -webkit-user-drag: none;
}

.lightbox .tip {
  color: rgba(255, 255, 255, 0.85);
  font-size: 0.8rem;
}

@keyframes lbIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@media (max-width: 720px) {
  .wall {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
}
</style>

<script setup>
// 猪咪聚集地 —— 图文小窝(agent 猪咪君君供稿,带饲养员标识)
import { ref, onMounted } from 'vue'

const posts = ref([])
const loaded = ref(false)
const zoomImg = ref(null)
let zoomAt = 0

onMounted(async () => {
  try {
    const res = await fetch('/api/hub')
    posts.value = await res.json()
  } catch {
    /* 静默 */
  } finally {
    loaded.value = true
  }
})

function openZoom(src) {
  zoomImg.value = src
  zoomAt = Date.now()
}

function closeZoom() {
  if (Date.now() - zoomAt > 350) zoomImg.value = null
}
</script>

<template>
  <div class="hub">
    <header class="hub-head">
      <a href="/" class="back">← 回猫咖</a>
      <h1 class="font-cute">猪咪聚集地</h1>
      <p class="sub">饲养员猪咪君君的小仓库 · 图与话的收容所</p>
    </header>

    <main class="feed">
      <p v-if="loaded && !posts.length" class="empty">还空着,等君君来投喂内容~</p>

      <article v-for="p in posts" :key="p.id" class="post" :class="p.type">
        <span v-if="p.by === 'agent'" class="keeper-badge font-cute">🐷 猪咪饲养员</span>

        <div v-if="p.type === 'image'" class="img-wrap" @click="openZoom(`/hub-img/${p.img}`)">
          <img :src="`/hub-img/${p.img}`" :alt="p.text || '聚集地图片'" loading="lazy" />
          <span class="zoom-hint">🔍</span>
        </div>
        <p v-if="p.text" class="text">{{ p.text }}</p>

        <time>{{ p.time }}</time>
      </article>
    </main>

    <!-- 点图放大 -->
    <div v-if="zoomImg" class="lightbox" @click="closeZoom">
      <img :src="zoomImg" alt="放大预览" @contextmenu.prevent @dragstart.prevent />
    </div>
  </div>
</template>

<style scoped>
.hub {
  min-height: 100vh;
  background: var(--bg);
  padding-bottom: 60px;
}

.hub-head {
  text-align: center;
  padding: 40px 20px 24px;
  position: relative;
}

.back {
  position: absolute;
  left: 20px;
  top: 24px;
  color: var(--muted);
  text-decoration: none;
  font-size: 0.9rem;
}

.back:hover {
  color: var(--pink-deep);
}

.hub-head h1 {
  font-size: clamp(1.6rem, 5vw, 2.2rem);
  color: var(--ink);
}

.sub {
  color: var(--muted);
  font-size: 0.9rem;
  margin-top: 6px;
}

.feed {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.post {
  position: relative;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: 255px 16px 225px 16px / 16px 225px 16px 255px; /* 手绘歪扭 */
  box-shadow: var(--shadow);
  padding: 16px 20px;
}

.keeper-badge {
  display: inline-block;
  background: #9b7ede;
  color: #fff;
  font-size: 0.72rem;
  border-radius: 999px;
  padding: 2px 12px;
  margin-bottom: 8px;
}

.img-wrap {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  cursor: zoom-in;
  margin-bottom: 10px;
}

.img-wrap img {
  width: 100%;
  display: block;
}

.zoom-hint {
  position: absolute;
  right: 8px;
  bottom: 8px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
}

.post .text {
  color: var(--ink);
  font-size: 0.95rem;
  line-height: 1.8;
  word-break: break-word;
  white-space: pre-wrap;
}

.post time {
  display: block;
  margin-top: 8px;
  color: var(--muted);
  font-size: 0.75rem;
  text-align: right;
}

.empty {
  text-align: center;
  color: var(--muted);
  padding: 60px 0;
}

.lightbox {
  position: fixed;
  inset: 0;
  z-index: 110;
  background: rgba(91, 58, 71, 0.75);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: zoom-out;
}

.lightbox img {
  max-width: 92vw;
  max-height: 86vh;
  object-fit: contain;
  border-radius: 12px;
  border: 4px solid #fff;
}
</style>

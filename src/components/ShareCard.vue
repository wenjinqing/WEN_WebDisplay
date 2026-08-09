<script setup>
// 分享卡片:网站风格 + 二维码,截图即可分享
import { ref, onMounted, onUnmounted } from 'vue'
import { site } from '../data.js'

const open = ref(false)
const url = ref('')

function onMenuShare() {
  open.value = true
}

onMounted(() => {
  url.value = window.location.origin
  window.addEventListener('menu-share', onMenuShare)
})

onUnmounted(() => window.removeEventListener('menu-share', onMenuShare))
</script>

<template>
  <div v-if="open" class="share-mask" @click.self="open = false">
    <div class="share-card" role="dialog" aria-label="分享本站">
      <button class="close" aria-label="关闭" @click="open = false">✕</button>

      <div class="head">
        <img src="/og-cover.png" alt="吉祥物" class="mascot" />
        <div>
          <h3 class="font-cute">{{ site.title }}</h3>
          <p class="slogan">{{ site.slogan }}</p>
        </div>
      </div>

      <div class="qr-box">
        <img :src="'/api/qr'" alt="本站二维码" class="qr" />
        <p class="scan-tip">扫码进店坐坐</p>
      </div>

      <p class="url">{{ url }}</p>
      <div class="divider">🐾 · 🐾 · 🐾</div>
      <p class="hint">截图这张卡片,发给群里的猪咪吧</p>
    </div>
  </div>
</template>

<style scoped>
.share-mask {
  position: fixed;
  inset: 0;
  z-index: 120;
  background: rgba(91, 58, 71, 0.5);
  backdrop-filter: blur(5px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.share-card {
  position: relative;
  width: min(360px, 100%);
  background: var(--bg);
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='140' height='140'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='2'/%3E%3C/filter%3E%3Crect width='140' height='140' filter='url(%23n)' opacity='0.04'/%3E%3C/svg%3E");
  border: 3px solid var(--pink-soft);
  border-radius: 255px 20px 225px 20px / 20px 225px 20px 255px; /* 手绘歪扭框 */
  box-shadow: 0 24px 64px rgba(91, 58, 71, 0.35);
  padding: 28px 24px 22px;
  text-align: center;
  animation: cardIn 0.25s ease;
}

@keyframes cardIn {
  from { opacity: 0; transform: translateY(16px) scale(0.96); }
  to { opacity: 1; transform: none; }
}

.close {
  position: absolute;
  top: 10px;
  right: 12px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: none;
  background: var(--pink-pale);
  color: var(--pink-deep);
  cursor: pointer;
}

.head {
  display: flex;
  align-items: center;
  gap: 14px;
  text-align: left;
  margin-bottom: 18px;
}

.mascot {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  border: 3px solid var(--pink-soft);
  object-fit: cover;
  background: #fff;
}

.head h3 {
  color: var(--ink);
  font-size: 1.15rem;
}

.slogan {
  color: var(--pink-deep);
  font-size: 0.82rem;
  margin-top: 2px;
}

.qr-box {
  background: #fff;
  border: 2px dashed var(--pink-soft);
  border-radius: 18px;
  padding: 16px;
  display: inline-block;
}

.qr {
  width: 160px;
  height: 160px;
  display: block;
}

.scan-tip {
  color: var(--muted);
  font-size: 0.8rem;
  margin-top: 8px;
}

.url {
  margin-top: 12px;
  color: var(--pink-deep);
  font-size: 0.85rem;
  font-family: 'JetBrains Mono', monospace;
}

.divider {
  color: var(--pink-soft);
  font-size: 0.8rem;
  margin: 10px 0 6px;
  letter-spacing: 6px;
}

.hint {
  color: var(--muted);
  font-size: 0.78rem;
}
</style>

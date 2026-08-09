<script setup>
// 分享卡片:网站风格卡片 + 实时数据 + 二维码,可保存为图片
import { ref, onMounted, onUnmounted } from 'vue'
import { site } from '../data.js'

const open = ref(false)
const url = ref('')
const stats = ref(null)
const saving = ref(false)

function onMenuShare() {
  open.value = true
  if (!stats.value) {
    fetch('/api/recap').then((r) => r.json()).then((d) => (stats.value = d)).catch(() => {})
  }
}

onMounted(() => {
  url.value = window.location.origin
  window.addEventListener('menu-share', onMenuShare)
})

onUnmounted(() => window.removeEventListener('menu-share', onMenuShare))

// ===== canvas 绘制分享图 =====
function roundRect(ctx, x, y, w, h, r) {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.arcTo(x + w, y, x + w, y + h, r)
  ctx.arcTo(x + w, y + h, x, y + h, r)
  ctx.arcTo(x, y + h, x, y, r)
  ctx.arcTo(x, y, x + w, y, r)
  ctx.closePath()
}

function drawPaw(ctx, x, y, s, color) {
  ctx.fillStyle = color
  ctx.beginPath(); ctx.ellipse(x, y + s * 0.35, s * 0.42, s * 0.32, 0, 0, 7); ctx.fill()
  for (const [dx, dy] of [[-0.55, -0.15], [-0.2, -0.45], [0.2, -0.45], [0.55, -0.15]]) {
    ctx.beginPath(); ctx.ellipse(x + dx * s, y + dy * s, s * 0.16, s * 0.2, 0, 0, 7); ctx.fill()
  }
}

async function saveImage() {
  if (saving.value) return
  saving.value = true
  try {
    const W = 750, H = 1060
    const cv = document.createElement('canvas')
    cv.width = W; cv.height = H
    const ctx = cv.getContext('2d')

    // 粉底 + 散落猫爪
    ctx.fillStyle = '#FFF6F8'
    ctx.fillRect(0, 0, W, H)
    for (const [x, y, s] of [[80, 90, 44], [660, 70, 36], [690, 950, 46], [70, 980, 34], [620, 560, 30], [110, 540, 26]]) {
      drawPaw(ctx, x, y, s, 'rgba(255,194,212,0.5)')
    }

    // 白卡
    ctx.fillStyle = '#FFFDFC'
    roundRect(ctx, 45, 45, W - 90, H - 90, 28)
    ctx.fill()
    ctx.strokeStyle = '#FFC2D4'
    ctx.lineWidth = 5
    ctx.stroke()

    // 吉祥物(圆形裁切)
    const mascot = await loadImg('/og-cover.png')
    ctx.save()
    ctx.beginPath(); ctx.arc(W / 2, 195, 95, 0, 7); ctx.clip()
    ctx.drawImage(mascot, W / 2 - 95, 100, 190, 190)
    ctx.restore()
    ctx.strokeStyle = '#F9718F'
    ctx.lineWidth = 6
    ctx.beginPath(); ctx.arc(W / 2, 195, 95, 0, 7); ctx.stroke()

    // 标题 & 标语
    ctx.textAlign = 'center'
    ctx.fillStyle = '#5B3A47'
    ctx.font = '52px "ZCOOL KuaiLe", sans-serif'
    ctx.fillText(site.title, W / 2, 370)
    ctx.fillStyle = '#E85D7F'
    ctx.font = '26px "Noto Sans SC", sans-serif'
    ctx.fillText(site.slogan, W / 2, 420)

    // 波浪线
    ctx.strokeStyle = '#F9718F'
    ctx.lineWidth = 4
    ctx.beginPath()
    for (let x = W / 2 - 90; x <= W / 2 + 90; x += 4) {
      ctx.lineTo(x, 450 + Math.sin((x - W / 2 + 90) / 14) * 5)
    }
    ctx.stroke()

    // 实时数据胶囊
    if (stats.value) {
      const s = stats.value
      const items = [
        `${s.visits} 只猪咪到访`, `${s.postcards} 张明信片`, `${s.urges} 次催更`,
      ]
      ctx.font = '24px "Noto Sans SC", sans-serif'
      let cx = W / 2 - 280
      for (const t of items) {
        const tw = ctx.measureText(t).width + 40
        ctx.fillStyle = '#FFE4EC'
        roundRect(ctx, cx, 485, tw, 52, 26)
        ctx.fill()
        ctx.fillStyle = '#E85D7F'
        ctx.fillText(t, cx + tw / 2, 520)
        cx += tw + 18
      }
    }

    // 二维码
    const qr = await loadImg('/api/qr')
    ctx.fillStyle = '#FFFFFF'
    roundRect(ctx, W / 2 - 140, 570, 280, 300, 20)
    ctx.fill()
    ctx.setLineDash([8, 6])
    ctx.strokeStyle = '#FFC2D4'
    ctx.lineWidth = 3
    ctx.stroke()
    ctx.setLineDash([])
    ctx.drawImage(qr, W / 2 - 115, 590, 230, 230)
    ctx.fillStyle = '#A07886'
    ctx.font = '22px "Noto Sans SC", sans-serif'
    ctx.fillText('扫码进店坐坐', W / 2, 850)

    // 网址 & 底部
    ctx.fillStyle = '#E85D7F'
    ctx.font = '26px "JetBrains Mono", monospace'
    ctx.fillText(url.value, W / 2, 920)
    ctx.fillStyle = '#A07886'
    ctx.font = '20px "Noto Sans SC", sans-serif'
    ctx.fillText(`「${site.author}」粉丝后援会 · 用爱发电`, W / 2, 985)

    const a = document.createElement('a')
    a.href = cv.toDataURL('image/png')
    a.download = '小涩猫咖啡厅-分享卡.png'
    a.click()
  } finally {
    saving.value = false
  }
}

function loadImg(src) {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => resolve(img)
    img.onerror = reject
    img.src = src
  })
}
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

      <!-- 实时数据胶囊 -->
      <div v-if="stats" class="stats-row">
        <span>{{ stats.visits }} 只猪咪到访</span>
        <span>{{ stats.postcards }} 张明信片</span>
        <span>{{ stats.urges }} 次催更</span>
      </div>

      <div class="qr-box">
        <img :src="'/api/qr'" alt="本站二维码" class="qr" />
        <p class="scan-tip">扫码进店坐坐</p>
      </div>

      <p class="url">{{ url }}</p>
      <div class="divider">🐾 · 🐾 · 🐾</div>
      <p class="hint">「{{ site.author }}」粉丝后援会 · 用爱发电</p>

      <button class="btn btn-primary save-btn" :disabled="saving" @click="saveImage">
        {{ saving ? '绘制中…' : '💾 保存分享卡图片' }}
      </button>
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
  width: min(380px, 100%);
  max-height: 92vh;
  overflow-y: auto;
  background: var(--bg);
  border: 3px solid var(--pink-soft);
  border-radius: 255px 20px 225px 20px / 20px 225px 20px 255px;
  box-shadow: 0 24px 64px rgba(91, 58, 71, 0.35);
  padding: 26px 24px 22px;
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
  margin-bottom: 14px;
}

.mascot {
  width: 60px;
  height: 60px;
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

.stats-row {
  display: flex;
  justify-content: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 14px;
}

.stats-row span {
  background: var(--pink-pale);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 3px 12px;
  font-size: 0.75rem;
}

.qr-box {
  background: #fff;
  border: 2px dashed var(--pink-soft);
  border-radius: 18px;
  padding: 14px;
  display: inline-block;
}

.qr {
  width: 150px;
  height: 150px;
  display: block;
}

.scan-tip {
  color: var(--muted);
  font-size: 0.78rem;
  margin-top: 6px;
}

.url {
  margin-top: 10px;
  color: var(--pink-deep);
  font-size: 0.85rem;
  font-family: 'JetBrains Mono', monospace;
}

.divider {
  color: var(--pink-soft);
  font-size: 0.8rem;
  margin: 8px 0 4px;
  letter-spacing: 6px;
}

.hint {
  color: var(--muted);
  font-size: 0.78rem;
}

.save-btn {
  margin-top: 14px;
  width: 100%;
  justify-content: center;
}
</style>

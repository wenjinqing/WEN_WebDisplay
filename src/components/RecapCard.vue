<script setup>
// 年度回忆:猫咖年度账本(小票风格)
import { ref, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'

const d = ref(null)
const year = new Date().getFullYear()

onMounted(async () => {
  try {
    const res = await fetch('/api/recap')
    d.value = await res.json()
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <section id="recap" v-if="d">
    <div class="container">
      <SectionTitle title="年度账本" :sub="`${year} 年 · 猫咖营业报告`" />

      <div class="receipt" v-reveal>
        <div class="rline"><span>累计到店猪咪</span><b>{{ d.visits }} 只</b></div>
        <div class="rline"><span>猫猫被撸</span><b>{{ d.pets }} 次</b></div>
        <div class="rline"><span>被投喂</span><b>{{ d.feeds }} 次</b></div>
        <div class="rline"><span>催更</span><b>{{ d.urges }} 次</b></div>
        <div class="rline"><span>留言</span><b>{{ d.messages }} 条</b></div>
        <div class="rline"><span>明信片</span><b>{{ d.postcards }} 张 · 获赞 {{ d.likes }}</b></div>
        <div class="rline"><span>攒鱼干的猪咪</span><b>{{ d.pigmis }} 只</b></div>
        <div class="rsep" />
        <div v-if="d.topPost" class="rline hi">
          <span>👑 人气明信片</span><b>{{ d.topPost.nick }} · {{ d.topPost.likes }} 赞</b>
        </div>
        <div v-if="d.topPigmi" class="rline hi">
          <span>🐟 鱼干首富</span><b>{{ d.topPigmi.nick }} · {{ d.topPigmi.points }}</b>
        </div>
        <div class="rfoot">—— 感谢惠顾,明年也请多指教 ——</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.receipt {
  max-width: 420px;
  margin: 0 auto;
  background: #fff;
  padding: 28px 28px 20px;
  font-family: 'JetBrains Mono', 'Noto Sans SC', monospace;
  box-shadow: var(--shadow);
  /* 小票锯齿边 */
  clip-path: polygon(
    0 0, 100% 0,
    100% calc(100% - 8px), 97% 100%, 94% calc(100% - 8px), 91% 100%,
    88% calc(100% - 8px), 85% 100%, 82% calc(100% - 8px), 79% 100%,
    76% calc(100% - 8px), 73% 100%, 70% calc(100% - 8px), 67% 100%,
    64% calc(100% - 8px), 61% 100%, 58% calc(100% - 8px), 55% 100%,
    52% calc(100% - 8px), 49% 100%, 46% calc(100% - 8px), 43% 100%,
    40% calc(100% - 8px), 37% 100%, 34% calc(100% - 8px), 31% 100%,
    28% calc(100% - 8px), 25% 100%, 22% calc(100% - 8px), 19% 100%,
    16% calc(100% - 8px), 13% 100%, 10% calc(100% - 8px), 7% 100%,
    4% calc(100% - 8px), 0 100%
  );
}

.rline {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 7px 0;
  font-size: 0.88rem;
  color: var(--ink);
  border-bottom: 1px dashed var(--pink-pale);
}

.rline span {
  color: var(--muted);
}

.rline b {
  text-align: right;
}

.rline.hi b {
  color: var(--pink-deep);
}

.rsep {
  border-bottom: 2px dashed var(--pink-soft);
  margin: 10px 0;
}

.rfoot {
  text-align: center;
  color: var(--muted);
  font-size: 0.78rem;
  margin-top: 14px;
}
</style>

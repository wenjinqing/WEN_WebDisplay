<script setup>
import { ref } from 'vue'

// 今日运势签:以日期为种子,同一天所有猪咪摇到同一签(猫猫酱の天意)
const fortunes = [
  { level: '大吉', text: '今天适合催更,一催一个准,猫猫酱会回你表情', yi: '宜催更 宜二刷', ji: '忌囤稿' },
  { level: '大吉', text: '灵感之神附体,给猫猫酱的留言会被秒回', yi: '宜表白 宜寄明信片', ji: '忌潜水' },
  { level: '中吉', text: '今天重读旧作会有新发现,比如某处伏笔', yi: '宜复习 宜写长评', ji: '忌催更过头' },
  { level: '中吉', text: '宜安利,把咖啡厅分享给一只新猪咪吧', yi: '宜分享 宜安利', ji: '忌吃独食' },
  { level: '小吉', text: '平平淡淡才是真,去留言板说句晚安也好', yi: '宜签到 宜点赞', ji: '忌熬夜' },
  { level: '小吉', text: '今天的猪咪也是可爱的猪咪,这就够了', yi: '宜摸鱼 宜贴贴', ji: '忌 emo' },
  { level: '末吉', text: '小心!今天点开小说会一口气看到凌晨三点', yi: '宜早睡 宜喝水', ji: '忌开新坑' },
]

const result = ref(null)
const shaking = ref(false)

function draw() {
  if (shaking.value) return
  shaking.value = true
  result.value = null
  setTimeout(() => {
    // 日期做种子,同一天全群同一签
    const d = new Date()
    const seed = d.getFullYear() * 10000 + (d.getMonth() + 1) * 100 + d.getDate()
    result.value = fortunes[seed % fortunes.length]
    shaking.value = false
  }, 800)
}
</script>

<template>
  <div class="fortune hd-card-s" v-reveal>
    <h3 class="font-cute">🥠 今日猪咪运势</h3>
    <p class="tip">每天一签,全群的猪咪共享同一份天意</p>

    <button class="tube" :class="{ shaking }" :disabled="shaking" @click="draw">
      {{ shaking ? '摇签中…' : result ? '已领签,明天再来' : '摇一卦' }}
    </button>

    <transition name="pop">
      <div v-if="result" class="fortune-result">
        <span class="level font-cute">{{ result.level }}</span>
        <p class="text">{{ result.text }}</p>
        <p class="yiji">{{ result.yi }} · {{ result.ji }}</p>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.fortune {
  max-width: 460px;
  margin: 32px auto 0;
  background: var(--card);
  border: 2px solid var(--pink-soft);
  box-shadow: var(--shadow);
  padding: 24px 28px;
  text-align: center;
}

h3 {
  color: var(--pink-deep);
  font-size: 1.2rem;
}

.tip {
  color: var(--muted);
  font-size: 0.8rem;
  margin: 4px 0 16px;
}

.tube {
  border: none;
  background: linear-gradient(160deg, var(--pink), var(--pink-deep));
  color: #fff;
  font-size: 1rem;
  font-family: 'ZCOOL KuaiLe', sans-serif;
  padding: 12px 36px;
  border-radius: 999px;
  cursor: pointer;
  box-shadow: 3px 4px 0 #d14a6e;
  transition: transform 0.15s;
}

.tube:hover:not(:disabled) {
  transform: translate(-1px, -2px);
}

.tube.shaking {
  animation: shake 0.15s linear infinite;
}

@keyframes shake {
  0%, 100% { transform: rotate(-6deg); }
  50% { transform: rotate(6deg); }
}

.fortune-result {
  margin-top: 18px;
  background: #fff9f4;
  border: 2px dashed var(--pink-soft);
  border-radius: 14px;
  padding: 16px;
}

.level {
  display: inline-block;
  background: var(--pink);
  color: #fff;
  border-radius: 999px;
  padding: 3px 20px;
  font-size: 1.1rem;
  margin-bottom: 8px;
}

.text {
  color: var(--ink);
  font-size: 0.95rem;
  line-height: 1.7;
}

.yiji {
  color: var(--muted);
  font-size: 0.8rem;
  margin-top: 6px;
}

.pop-enter-active { transition: all 0.35s ease; }
.pop-enter-from { opacity: 0; transform: translateY(12px) scale(0.95); }
</style>

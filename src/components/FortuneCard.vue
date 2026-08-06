<script setup>
import { ref } from 'vue'

// 今日运势签:以日期为种子,同一天所有猪咪摇到同一签(猫猫酱の天意)
const fortunes = [
  // ===== 大吉(10) =====
  { level: '大吉', text: '今天适合催更,一催一个准,猫猫酱会回你表情', yi: '宜催更 宜二刷', ji: '忌囤稿' },
  { level: '大吉', text: '灵感之神附体,给猫猫酱的留言会被秒回', yi: '宜表白 宜寄明信片', ji: '忌潜水' },
  { level: '大吉', text: '今天写下的长评,会被猫猫酱截图珍藏', yi: '宜写长评 宜走心', ji: '忌敷衍' },
  { level: '大吉', text: '桃花运喵喵叫,你嗑的 CP 今天必发糖', yi: '宜嗑糖 宜分享', ji: '忌泼冷水' },
  { level: '大吉', text: '今天是猫猫酱的灵感爆棚日,蹲一蹲说不定有更新', yi: '宜蹲更新 宜刷新', ji: '忌早睡' },
  { level: '大吉', text: '随手一抽就是隐藏款,今天做什么都顺', yi: '宜抽卡 宜尝试', ji: '忌犹豫' },
  { level: '大吉', text: '你的彩虹屁会精准命中猫猫酱的心巴', yi: '宜夸夸 宜输出', ji: '忌含蓄' },
  { level: '大吉', text: '今天寄出的明信片,会登上人气王宝座', yi: '宜寄明信片 宜露脸', ji: '忌害羞' },
  { level: '大吉', text: '梦里会有百合花开,醒来记得讲给群里听', yi: '宜记录 宜分享梦', ji: '忌赖床' },
  { level: '大吉', text: '今天的你就是本店的幸运小猫,撸猫必出喵叫', yi: '宜撸猫 宜摇签', ji: '忌低调' },
  // ===== 中吉(12) =====
  { level: '中吉', text: '今天重读旧作会有新发现,比如某处伏笔', yi: '宜复习 宜写长评', ji: '忌催更过头' },
  { level: '中吉', text: '宜安利,把咖啡厅分享给一只新猪咪吧', yi: '宜分享 宜安利', ji: '忌吃独食' },
  { level: '中吉', text: '今天的灵感像猫一样,追它它就跑,坐下来它自己来', yi: '宜摸鱼 宜等待', ji: '忌硬憋' },
  { level: '中吉', text: '去明信片墙逛逛,会有让你笑出声的一张', yi: '宜逛墙 宜点赞', ji: '忌路过不看' },
  { level: '中吉', text: '今天和群友斗图会赢,表情包库存自动+10', yi: '宜斗图 宜存图', ji: '忌认输' },
  { level: '中吉', text: '一句"猫猫酱辛苦了",今天的功德就圆满了', yi: '宜问候 宜暖心', ji: '忌只催不爱' },
  { level: '中吉', text: '今天的评论区会有神回复,记得去围观', yi: '宜围观 宜跟帖', ji: '忌沉默' },
  { level: '中吉', text: '睡前读一章甜的,今晚的梦也是甜的', yi: '宜夜读 宜早睡', ji: '忌虐文' },
  { level: '中吉', text: '今天立的 flag 不会倒,比如"看一章就睡"', yi: '宜立flag 宜自律', ji: '忌真香' },
  { level: '中吉', text: '投喂一只群友,善意会绕一圈回到你身上', yi: '宜投喂 宜互助', ji: '忌自闭' },
  { level: '中吉', text: '今天的签名档该换了,用猫猫酱的金句吧', yi: '宜换签名 宜引用', ji: '忌土味' },
  { level: '中吉', text: '适合整理收藏夹,会发现遗忘的宝藏章节', yi: '宜整理 宜考古', ji: '忌吃灰' },
  // ===== 小吉(12) =====
  { level: '小吉', text: '平平淡淡才是真,去留言板说句晚安也好', yi: '宜签到 宜点赞', ji: '忌熬夜' },
  { level: '小吉', text: '今天的猪咪也是可爱的猪咪,这就够了', yi: '宜摸鱼 宜贴贴', ji: '忌 emo' },
  { level: '小吉', text: '适合二刷旧糖,旧糖新嗑,越嗑越甜', yi: '宜二刷 宜回味', ji: '忌挑食' },
  { level: '小吉', text: '今天的一切小确丧,都会被首页的猫治愈', yi: '宜撸猫 宜深呼吸', ji: '忌生闷气' },
  { level: '小吉', text: '点一杯小说慢慢喝,别催,好东西都值得等', yi: '宜慢读 宜品茶', ji: '忌囫囵吞枣' },
  { level: '小吉', text: '给猫猫酱的朋友圈点个赞,好运+1', yi: '宜点赞 宜冒泡', ji: '忌窥屏' },
  { level: '小吉', text: '今天适合把"哈哈哈"换成走心的一句话', yi: '宜走心 宜长评', ji: '忌划水' },
  { level: '小吉', text: '慢慢来,猪猪也会有春天,猪咪也是', yi: '宜躺平 宜晒太阳', ji: '忌焦虑' },
  { level: '小吉', text: '今天的糖分摄入刚刚好,不多不少一章就好', yi: '宜适量 宜细品', ji: '忌贪杯' },
  { level: '小吉', text: '悄悄努力攒鱼干,升级那天惊艳全群', yi: '宜攒鱼干 宜低调', ji: '忌炫耀' },
  { level: '小吉', text: '今天的云朵像棉花糖,适合拍照发群里', yi: '宜抬头 宜拍照', ji: '忌低头' },
  { level: '小吉', text: '普通的一天,也是猫猫酱笔下温柔的那种', yi: '宜感受 宜记录', ji: '忌烦躁' },
  // ===== 末吉(6) =====
  { level: '末吉', text: '小心!今天点开小说会一口气看到凌晨三点', yi: '宜早睡 宜喝水', ji: '忌开新坑' },
  { level: '末吉', text: '今天的催更可能会被已读不回,但爱不会', yi: '宜佛系 宜自愈', ji: '忌夺命连环催' },
  { level: '末吉', text: '今天的刀子比糖多,看文前请备好纸巾', yi: '宜备纸 宜抱团', ji: '忌深夜看虐' },
  { level: '末吉', text: '摇签的手气一般,但摇尾巴的运气不错', yi: '宜撸猫转运', ji: '忌再摇一次' },
  { level: '末吉', text: '小心群里的大佬,今天的斗图你可能会输', yi: '宜猥琐发育', ji: '忌正面硬刚' },
  { level: '末吉', text: '今天的更新可能在明天,但猫猫酱的爱一直在', yi: '宜等待 宜信任', ji: '忌催命' },
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

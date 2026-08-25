<script setup>
import { ref, computed, onMounted } from 'vue'
import PawPrint from './PawPrint.vue'
import SectionTitle from './SectionTitle.vue'
import { site } from '../data.js'

const emit = defineEmits(['read', 'comments'])

// ===== 更新订阅 =====
const nick = ref(localStorage.getItem('catcafe_nick') || '')
const mySubs = ref(new Set())

onMounted(async () => {
  if (!nick.value) return
  try {
    const res = await fetch(`/api/subscriptions?nick=${encodeURIComponent(nick.value)}`)
    mySubs.value = new Set(await res.json())
  } catch { /* 静默 */ }
})

async function toggleSub(n) {
  if (!nick.value) {
    alert('先去「猪咪聚集地」里输入昵称查询一次头衔,就能订阅更新提醒啦~')
    return
  }
  try {
    const res = await fetch('/api/subscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ file: n.file, nick: nick.value }),
    })
    const data = await res.json()
    if (res.ok) {
      if (data.subscribed) mySubs.value.add(n.file)
      else mySubs.value.delete(n.file)
      mySubs.value = new Set(mySubs.value)
    }
  } catch { /* 静默 */ }
}

const cats = ['全部', '连载中', '已完结', '番外']
const activeCat = ref('全部')
const filtered = computed(() =>
  activeCat.value === '全部'
    ? site.novels
    : site.novels.filter((n) => (n.cat || '已完结') === activeCat.value)
)
</script>

<template>
  <section id="novels" class="menu-section">
    <div class="container">
      <SectionTitle title="今日特供 · 小说菜单" sub="点单即下载 · 本区作品含轻度成人向内容,未成年猪咪请自觉绕行喵~" />

      <div class="cat-tabs" v-reveal>
        <button
          v-for="t in cats"
          :key="t"
          :class="{ active: activeCat === t }"
          @click="activeCat = t"
        >
          {{ t }}
        </button>
      </div>

      <div class="menu-board" v-reveal>
        <div v-for="n in filtered" :key="n.file" class="menu-item">
          <div class="item-info">
            <div class="item-head">
              <h3 class="font-cute">{{ n.title }}</h3>
              <span class="cup">{{ n.cup }}</span>
              <span v-if="n.cat" class="cat-tag">{{ n.cat }}</span>
              <span v-if="n.chapters && n.chapters.length" class="cup">共{{ n.chapters.length }}章</span>
            </div>
            <p class="desc">{{ n.desc }}</p>
          </div>
          <div class="item-actions">
            <button
              class="sub-btn"
              :class="{ on: mySubs.has(n.file) }"
              :title="mySubs.has(n.file) ? '已订阅,更新会提醒你' : '订阅更新提醒'"
              @click="toggleSub(n)"
            >
              {{ mySubs.has(n.file) ? '🔔 已订阅' : '🔕 订阅' }}
            </button>
            <button class="btn btn-ghost read" @click="emit('read', n)">在线阅读</button>
            <button class="btn btn-ghost read" @click="emit('comments', n)">💬 评论</button>
            <a class="btn btn-primary download" :href="`/downloads/${n.file}`" download>
              <PawPrint :size="18" color="#fff" />
              下载
            </a>
          </div>
        </div>
        <p v-if="!filtered.length" class="empty-tip">这个分类还空着,等猫猫酱上新喵~</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.menu-section {
  background: var(--bg-deep);
}

.menu-board {
  max-width: 860px;
  margin: 0 auto;
  background: var(--card);
  border: 3px solid var(--pink-soft);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 12px 32px;
  position: relative;
}

/* 菜单板顶部挂绳孔装饰 */
.menu-board::before {
  content: '';
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  width: 60px;
  height: 14px;
  border-radius: 8px 8px 0 0;
  background: var(--pink-soft);
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 24px 0;
  border-bottom: 2px dashed var(--pink-pale);
}

.menu-item:last-child {
  border-bottom: none;
}

.item-info {
  flex: 1;
}

.item-head {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

h3 {
  font-size: 1.15rem;
  color: var(--ink);
}

.cup {
  font-size: 0.8rem;
  color: var(--pink-deep);
  background: var(--pink-pale);
  border-radius: 999px;
  padding: 2px 12px;
  white-space: nowrap;
}

.cat-tag {
  font-size: 0.8rem;
  color: #fff;
  background: var(--pink);
  border-radius: 999px;
  padding: 2px 12px;
  white-space: nowrap;
}

.cat-tabs {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.cat-tabs button {
  border: 2px solid var(--pink-soft);
  background: var(--surface-2);
  color: var(--ink);
  border-radius: 999px;
  padding: 6px 20px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}

.cat-tabs button.active {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
}

.empty-tip {
  text-align: center;
  color: var(--muted);
  padding: 24px 0;
}

.desc {
  color: var(--muted);
  font-size: 0.92rem;
  margin-top: 6px;
}

.item-actions {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.read {
  padding: 12px 20px;
}

.sub-btn {
  border: 2px solid var(--pink-soft);
  background: var(--surface-2);
  color: var(--muted);
  border-radius: 999px;
  padding: 10px 16px;
  font-size: 0.88rem;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.sub-btn.on {
  background: var(--pink-pale);
  border-color: var(--pink);
  color: var(--pink-deep);
}

@media (max-width: 720px) {
  .menu-board {
    padding: 8px 20px;
  }
  .menu-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 14px;
  }
  .item-actions {
    align-self: stretch;
  }
  .item-actions .btn {
    flex: 1;
    justify-content: center;
  }
}
</style>

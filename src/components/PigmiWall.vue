<script setup>
import { ref, computed, onMounted } from 'vue'
import SectionTitle from './SectionTitle.vue'
import PigmiFace from './PigmiFace.vue'
import PawPrint from './PawPrint.vue'

const posts = ref([])
const nick = ref('')
const note = ref('')
const file = ref(null)
const fileInput = ref(null)
const sending = ref(false)
const msg = ref('')
const msgOk = ref(false)

// 本机已点过赞的明信片(localStorage 记录,防重复点)
const liked = ref(new Set(JSON.parse(localStorage.getItem('catcafe_likes') || '[]')))

// 点赞多的排前面,同数按时间新→旧(文件名前缀是时间戳)
const sorted = computed(() =>
  [...posts.value].sort((a, b) => (b.likes || 0) - (a.likes || 0) || b.img.localeCompare(a.img))
)

onMounted(load)

async function load() {
  try {
    const res = await fetch('/api/wall')
    posts.value = await res.json()
  } catch {
    /* 静默降级 */
  }
}

async function like(p) {
  if (liked.value.has(p.img)) return showMsg('这张你已经拍过爪啦 🐾', false)
  try {
    const res = await fetch('/api/wall/like', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ img: p.img }),
    })
    const data = await res.json()
    if (res.ok) {
      p.likes = data.likes
      liked.value.add(p.img)
      localStorage.setItem('catcafe_likes', JSON.stringify([...liked.value]))
    } else {
      showMsg(data.error || '点赞失败', false)
    }
  } catch {
    showMsg('网络打了个盹,稍后再试', false)
  }
}

async function send() {
  if (!file.value || sending.value) return
  sending.value = true
  try {
    const fd = new FormData()
    fd.append('file', file.value)
    fd.append('nick', nick.value)
    fd.append('note', note.value)
    const res = await fetch('/api/wall', { method: 'POST', body: fd })
    const data = await res.json()
    if (res.ok) {
      posts.value = [data, ...posts.value]
      nick.value = note.value = ''
      file.value = null
      if (fileInput.value) fileInput.value.value = ''
      showMsg('明信片寄出啦,已贴到墙上 💌', true)
    } else {
      showMsg(data.error || '寄件失败,稍后再试', false)
    }
  } catch {
    showMsg('网络打了个盹,稍后再试', false)
  } finally {
    sending.value = false
  }
}

function showMsg(text, ok) {
  msg.value = text
  msgOk.value = ok
  setTimeout(() => (msg.value = ''), 3000)
}
</script>

<template>
  <section id="wall">
    <div class="container">
      <SectionTitle title="猪咪明信片墙" sub="寄一张你的明信片,让猫猫酱看到可爱的你" />

      <!-- 寄件表单 -->
      <div class="send-card hd-card" v-reveal>
        <div class="send-form">
          <input v-model="nick" maxlength="20" placeholder="你的昵称(不填就是匿名猪咪)" />
          <input v-model="note" maxlength="60" placeholder="一句话…(60字以内,写给墙上大家看的)" />
          <div class="send-row">
            <input ref="fileInput" type="file" accept="image/*" @change="file = $event.target.files[0]" />
            <button class="btn btn-primary" :disabled="!file || sending" @click="send">
              {{ sending ? '寄件中…' : '💌 寄出明信片' }}
            </button>
          </div>
          <p v-if="msg" class="send-msg" :class="{ ok: msgOk }">{{ msg }}</p>
        </div>
      </div>

      <!-- 明信片网格:按点赞数排序 -->
      <div v-if="posts.length" class="wall-grid" v-reveal>
        <figure v-for="(p, i) in sorted" :key="p.img" class="postcard" :class="`tilt-${i % 3}`">
          <div class="photo">
            <img :src="`/wall/${p.img}`" :alt="p.nick + '的明信片'" loading="lazy" />
            <span v-if="i === 0 && p.likes > 0" class="crown font-cute">👑 人气王</span>
          </div>
          <figcaption>
            <p v-if="p.note" class="note">{{ p.note }}</p>
            <div class="who">
              <b>{{ p.nick }}</b>
              <button
                class="like-btn"
                :class="{ liked: liked.has(p.img) }"
                :aria-label="`给${p.nick}的明信片点赞`"
                @click="like(p)"
              >
                <PawPrint :size="16" :color="liked.has(p.img) ? '#fff' : '#e85d7f'" />
                {{ p.likes || 0 }}
              </button>
            </div>
          </figcaption>
        </figure>
      </div>

      <div v-else class="empty-wall" v-reveal>
        <PigmiFace :size="90" />
        <p>墙上还空空的,来贴第一张明信片喵~</p>
      </div>
    </div>
  </section>
</template>

<style scoped>
#wall {
  background: var(--bg-deep);
}

.send-card {
  max-width: 640px;
  margin: 0 auto 40px;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  box-shadow: var(--shadow);
  padding: 24px 28px;
}

.send-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.send-form input[type='text'],
.send-form input:not([type]) {
  border: 2px solid var(--pink-pale);
  border-radius: 12px;
  padding: 9px 14px;
  font-size: 0.95rem;
  font-family: inherit;
  color: var(--ink);
  outline: none;
}

.send-form input:focus {
  border-color: var(--pink);
}

.send-row {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.send-row input[type='file'] {
  flex: 1;
  font-size: 0.85rem;
  color: var(--muted);
  min-width: 200px;
}

.send-msg {
  font-size: 0.85rem;
  color: var(--pink-deep);
}

.send-msg.ok {
  color: #4caf7d;
}

.wall-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 28px;
  max-width: 960px;
  margin: 0 auto;
}

.postcard {
  background: #fff;
  border-radius: 8px;
  padding: 12px 12px 14px;
  box-shadow: 0 8px 20px rgba(233, 93, 127, 0.15);
  transition: transform 0.25s ease;
}

.postcard:hover {
  transform: rotate(0deg) scale(1.03) !important;
}

.tilt-0 { transform: rotate(-1.6deg); }
.tilt-1 { transform: rotate(1.2deg); }
.tilt-2 { transform: rotate(-0.8deg); }

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

.note {
  font-size: 0.88rem;
  color: var(--ink);
  margin: 8px 0 4px;
  word-break: break-word;
}

.who {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.who b {
  color: var(--pink-deep);
  font-size: 0.85rem;
}

.like-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 2px solid var(--pink-soft);
  background: #fff;
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 3px 12px;
  font-size: 0.82rem;
  cursor: pointer;
  transition: transform 0.15s, background 0.15s;
}

.like-btn:hover {
  transform: scale(1.08);
}

.like-btn.liked {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
  cursor: default;
}

.photo {
  position: relative;
}

.crown {
  position: absolute;
  top: 8px;
  left: 8px;
  background: rgba(255, 255, 255, 0.92);
  color: var(--pink-deep);
  font-size: 0.75rem;
  border-radius: 999px;
  padding: 2px 10px;
  box-shadow: 0 2px 8px rgba(233, 93, 127, 0.25);
}

.empty-wall {
  text-align: center;
  color: var(--muted);
  padding: 20px 0 10px;
}

@media (max-width: 720px) {
  .wall-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
  .send-card {
    padding: 18px;
  }
}
</style>

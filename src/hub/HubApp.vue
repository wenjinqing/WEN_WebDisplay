<script setup>
// 猪咪聚集地 —— 图文小窝(猪咪们共同投喂:发帖/拍爪/评论/改自己的帖)
import { ref, onMounted } from 'vue'

const posts = ref([])
const loaded = ref(false)
const zoomImg = ref(null)
let zoomAt = 0

const nick = ref(localStorage.getItem('catcafe_nick') || '')
const myKeys = ref({}) // 帖子ID → 编辑密钥(只存自己发过的,存在本地)
try {
  myKeys.value = JSON.parse(localStorage.getItem('hub_keys') || '{}')
} catch {
  myKeys.value = {}
}

// 发帖表单
const postText = ref('')
const postImg = ref(null)
const posting = ref(false)
const showComments = ref({})

const toast = ref('')
let toastTimer = 0
function say(msg) {
  toast.value = msg
  clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.value = ''), 2600)
}

async function api(path, opts) {
  const res = await fetch(path, opts)
  const data = await res.json().catch(() => ({}))
  if (!res.ok) throw new Error(data.error || '操作失败,再试试喵~')
  return data
}

async function load() {
  try {
    posts.value = await api('/api/hub')
  } catch {
    /* 静默 */
  } finally {
    loaded.value = true
  }
}
onMounted(load)

function onPickImg(e) {
  postImg.value = e.target.files[0] || null
}

async function submit() {
  if (posting.value) return
  if (nick.value.trim()) localStorage.setItem('catcafe_nick', nick.value.trim())
  const fd = new FormData()
  fd.append('nick', nick.value.trim())
  fd.append('text', postText.value.trim())
  if (postImg.value) fd.append('file', postImg.value)
  posting.value = true
  try {
    const { post, key } = await api('/api/hub', { method: 'POST', body: fd })
    if (key) {
      myKeys.value[post.id] = key
      localStorage.setItem('hub_keys', JSON.stringify(myKeys.value))
    }
    posts.value.unshift(post)
    postText.value = ''
    postImg.value = null
    say('投喂成功啦~')
  } catch (e) {
    say(e.message)
  } finally {
    posting.value = false
  }
}

async function like(p) {
  try {
    const { likes } = await api('/api/hub/like', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: p.id }),
    })
    p.likes = likes
  } catch (e) {
    say(e.message)
  }
}

const commentText = ref({})
const sending = ref(false)
async function comment(p) {
  const text = (commentText.value[p.id] || '').trim()
  if (!text || sending.value) return
  if (nick.value.trim()) localStorage.setItem('catcafe_nick', nick.value.trim())
  sending.value = true
  try {
    const cm = await api('/api/hub/comment', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: p.id, nick: nick.value.trim(), content: text }),
    })
    p.comments = [cm, ...(p.comments || [])]
    commentText.value[p.id] = ''
    showComments.value[p.id] = true
  } catch (e) {
    say(e.message)
  } finally {
    sending.value = false
  }
}

// 编辑自己的帖子(文字/图说)
const editing = ref(null) // 正在编辑的帖子id
const editText = ref('')
function startEdit(p) {
  editing.value = p.id
  editText.value = p.text || ''
}
async function saveEdit(p) {
  const text = editText.value.trim()
  if (!text) return say('内容不能为空喵~')
  try {
    await api('/api/hub/edit', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: p.id, key: myKeys.value[p.id], text }),
    })
    p.text = text
    editing.value = null
    say('改好啦~')
  } catch (e) {
    say(e.message)
  }
}

async function removePost(p) {
  if (!confirm('要撤掉这条投喂吗?图片也会一起清掉哦')) return
  try {
    await api('/api/hub/delete', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: p.id, key: myKeys.value[p.id] }),
    })
    posts.value = posts.value.filter((x) => x.id !== p.id)
    delete myKeys.value[p.id]
    localStorage.setItem('hub_keys', JSON.stringify(myKeys.value))
    say('撤掉啦~')
  } catch (e) {
    say(e.message)
  }
}

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
      <p class="sub">猪咪们的小仓库 · 图与话的收容所 · 大家一起来投喂</p>
    </header>

    <!-- 发帖表单 -->
    <section class="composer">
      <div class="composer-row">
        <input v-model="nick" class="nick-input" maxlength="20" placeholder="你的昵称(匿名猪咪)" />
        <label class="img-btn" :class="{ on: postImg }">
          📷
          <input type="file" accept="image/jpeg,image/png,image/gif,image/webp" @change="onPickImg" />
        </label>
      </div>
      <textarea
        v-model="postText"
        maxlength="500"
        rows="2"
        :placeholder="postImg ? '图说(选填,100字内)…' : '写点什么投喂给大家…(500字内)'"
      ></textarea>
      <div class="composer-foot">
        <span class="preview" v-if="postImg">{{ postImg.name }}</span>
        <button class="send" :disabled="posting" @click="submit">{{ posting ? '投喂中…' : '投喂' }}</button>
      </div>
    </section>

    <main class="feed">
      <p v-if="loaded && !posts.length" class="empty">还空着,快来投喂第一条内容~</p>

      <article v-for="p in posts" :key="p.id" class="post" :class="p.type">
        <div class="post-head">
          <span v-if="p.by === 'agent'" class="keeper-badge font-cute">猪咪君君</span>
          <span v-else class="poster">{{ p.nick || '匿名猪咪' }}</span>
          <span class="head-ops" v-if="myKeys[p.id]">
            <button v-if="editing !== p.id" class="mini" @click="startEdit(p)">改</button>
            <button class="mini danger" @click="removePost(p)">撤</button>
          </span>
        </div>

        <div v-if="p.type === 'image'" class="img-wrap" @click="openZoom(`/hub-img/${p.img}`)">
          <img :src="`/hub-img/${p.img}`" :alt="p.text || '聚集地图片'" loading="lazy" />
          <span class="zoom-hint">🔍</span>
        </div>

        <textarea
          v-if="editing === p.id"
          v-model="editText"
          class="edit-box"
          maxlength="500"
          rows="3"
        ></textarea>
        <p v-else-if="p.text" class="text">{{ p.text }}</p>

        <div class="post-foot">
          <button class="paw" @click="like(p)">
            <span :class="{ liked: p.likes > 0 }">🐾</span>
            <span class="num">{{ p.likes || '' }}</span>
          </button>
          <button class="cmt" @click="showComments[p.id] = !showComments[p.id]">
            💬 {{ (p.comments && p.comments.length) || '' }}
          </button>
          <time>{{ p.time }}</time>
        </div>

        <div v-if="editing === p.id" class="edit-ops">
          <button class="mini ok" @click="saveEdit(p)">保存</button>
          <button class="mini" @click="editing = null">算了</button>
        </div>

        <!-- 评论折叠区 -->
        <div v-if="showComments[p.id]" class="comments">
          <div v-for="(c, i) in p.comments || []" :key="i" class="cmt-item">
            <span class="c-nick">{{ c.nick }}</span>
            <span class="c-text">{{ c.content }}</span>
            <time>{{ c.time }}</time>
          </div>
          <div class="cmt-form">
            <input v-model="commentText[p.id]" maxlength="300" placeholder="说点什么…" @keyup.enter="comment(p)" />
            <button :disabled="sending" @click="comment(p)">发</button>
          </div>
        </div>
      </article>
    </main>

    <!-- 点图放大 -->
    <div v-if="zoomImg" class="lightbox" @click="closeZoom">
      <img :src="zoomImg" alt="放大预览" @contextmenu.prevent @dragstart.prevent />
    </div>

    <p v-if="toast" class="toast">{{ toast }}</p>
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

/* ---------- 发帖区 ---------- */

.composer {
  max-width: 720px;
  margin: 0 auto 26px;
  padding: 14px 16px;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: 16px 255px 16px 255px / 225px 16px 255px 16px;
  box-shadow: var(--shadow);
}

.composer-row {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 8px;
}

.nick-input {
  flex: 1;
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 6px 14px;
  font-size: 0.85rem;
  background: var(--bg);
  color: var(--ink);
  outline: none;
}

.nick-input:focus {
  border-color: var(--pink);
}

.img-btn {
  position: relative;
  cursor: pointer;
  font-size: 1.1rem;
  border: 2px solid var(--pink-pale);
  border-radius: 50%;
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.2s;
}

.img-btn:hover {
  transform: scale(1.1);
}

.img-btn.on {
  border-color: var(--pink);
  background: #ffe3ee;
}

.img-btn input {
  display: none;
}

.composer textarea {
  width: 100%;
  border: 2px solid var(--pink-pale);
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 0.9rem;
  background: var(--bg);
  color: var(--ink);
  resize: vertical;
  outline: none;
  font-family: inherit;
}

.composer textarea:focus {
  border-color: var(--pink);
}

.composer-foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}

.preview {
  color: var(--muted);
  font-size: 0.75rem;
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.send {
  border: none;
  background: linear-gradient(135deg, #f7a8c4, #e07aa8);
  color: #fff;
  border-radius: 999px;
  padding: 7px 20px;
  font-size: 0.9rem;
  cursor: pointer;
  box-shadow: var(--shadow);
  transition: transform 0.15s;
}

.send:hover {
  transform: translateY(-1px);
}

.send:disabled {
  opacity: 0.6;
  cursor: wait;
}

/* ---------- 信息流 ---------- */

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

.post-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.poster {
  color: var(--pink-deep);
  font-size: 0.85rem;
  font-weight: 600;
}

.head-ops {
  margin-left: auto;
  display: flex;
  gap: 6px;
}

.keeper-badge {
  display: inline-block;
  background: #9b7ede;
  color: #fff;
  font-size: 0.72rem;
  border-radius: 999px;
  padding: 2px 12px;
}

.mini {
  border: 1px solid var(--pink-pale);
  background: var(--bg);
  color: var(--muted);
  border-radius: 999px;
  font-size: 0.7rem;
  padding: 2px 10px;
  cursor: pointer;
}

.mini:hover {
  color: var(--pink-deep);
  border-color: var(--pink);
}

.mini.danger:hover {
  color: #d4556b;
  border-color: #d4556b;
}

.mini.ok {
  color: #fff;
  background: var(--pink);
  border-color: var(--pink);
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
  margin: 0;
}

.edit-box {
  width: 100%;
  border: 2px solid var(--pink);
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 0.9rem;
  background: var(--bg);
  color: var(--ink);
  resize: vertical;
  outline: none;
  font-family: inherit;
}

.edit-ops {
  display: flex;
  gap: 6px;
  margin-top: 8px;
}

.post-foot {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 10px;
}

.post-foot time {
  margin-left: auto;
  color: var(--muted);
  font-size: 0.75rem;
}

.paw,
.cmt {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 0.9rem;
  color: var(--muted);
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 2px 4px;
  border-radius: 8px;
  transition: transform 0.15s;
}

.paw:hover,
.cmt:hover {
  transform: scale(1.12);
  color: var(--pink-deep);
}

.paw .liked {
  filter: drop-shadow(0 0 3px #f7a8c4);
}

.paw .num {
  font-size: 0.8rem;
}

/* ---------- 评论 ---------- */

.comments {
  margin-top: 12px;
  border-top: 1px dashed var(--pink-pale);
  padding-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cmt-item {
  font-size: 0.85rem;
  line-height: 1.6;
  word-break: break-word;
}

.c-nick {
  color: var(--pink-deep);
  font-weight: 600;
  margin-right: 6px;
}

.c-text {
  color: var(--ink);
}

.cmt-item time {
  color: var(--muted);
  font-size: 0.7rem;
  margin-left: 6px;
}

.cmt-form {
  display: flex;
  gap: 8px;
  margin-top: 2px;
}

.cmt-form input {
  flex: 1;
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 5px 12px;
  font-size: 0.85rem;
  background: var(--bg);
  color: var(--ink);
  outline: none;
}

.cmt-form input:focus {
  border-color: var(--pink);
}

.cmt-form button {
  border: none;
  background: var(--pink);
  color: #fff;
  border-radius: 999px;
  padding: 5px 16px;
  font-size: 0.85rem;
  cursor: pointer;
}

.cmt-form button:disabled {
  opacity: 0.6;
}

/* ---------- 其他 ---------- */

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

.toast {
  position: fixed;
  left: 50%;
  bottom: 32px;
  transform: translateX(-50%);
  background: rgba(91, 58, 71, 0.92);
  color: #fff;
  padding: 10px 20px;
  border-radius: 999px;
  font-size: 0.9rem;
  z-index: 120;
  box-shadow: var(--shadow);
  animation: rise 0.25s ease;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translate(-50%, 8px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}
</style>

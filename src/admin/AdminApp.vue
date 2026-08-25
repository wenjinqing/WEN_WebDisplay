<script setup>
import { ref, reactive, onMounted } from 'vue'

const token = ref(localStorage.getItem('catcafe_token') || '')
const password = ref('')
const loginError = ref('')
const loggingIn = ref(false)

const c = reactive({})
const loaded = ref(false)
const saving = ref(false)
const toast = ref('')
const tab = ref('basic')

const tabs = [
  { key: 'basic', label: '基本资料' },
  { key: 'notices', label: '公告管理' },
  { key: 'novels', label: '小说上架' },
  { key: 'gallery', label: '插画管理' },
  { key: 'wall', label: '明信片管理' },
  { key: 'msgs', label: '留言管理' },
  { key: 'account', label: '修改密码' },
]

const wallPosts = ref([])
const msgs = ref([])
const replyText = reactive({})

async function loadMsgs() {
  const res = await fetch('/api/messages')
  msgs.value = await res.json()
}

async function delMsg(m) {
  const res = await fetch('/api/admin/messages/delete', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ time: m.time, nick: m.nick }),
  })
  if (res.status === 401) return logout()
  if (res.ok) {
    msgs.value = msgs.value.filter((x) => !(x.time === m.time && x.nick === m.nick))
    showToast('留言已删除')
  }
}

async function replyMsg(m) {
  const text = (replyText[m.time] || '').trim()
  if (!text) return
  const res = await fetch('/api/admin/messages/reply', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ time: m.time, nick: m.nick, reply: text }),
  })
  if (res.status === 401) return logout()
  if (res.ok) {
    m.reply = text
    replyText[m.time] = ''
    showToast('回复成功')
  }
}

async function loadWall() {
  const res = await fetch('/api/wall')
  wallPosts.value = await res.json()
}

async function delWall(img) {
  const res = await fetch('/api/admin/wall/delete', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ img }),
  })
  if (res.status === 401) return logout()
  if (res.ok) {
    wallPosts.value = wallPosts.value.filter((p) => p.img !== img)
    showToast('已撤下')
  }
}

// 新条目表单
const newNotice = reactive({ date: new Date().toISOString().slice(0, 10), tag: '公告', text: '' })
const newNovel = reactive({ title: '', desc: '', file: '', cup: '中杯 · 微糖', cat: '连载中' })
const newChapter = reactive({ title: '', file: null })
const expandedNovel = ref(-1)

async function addChapter(n) {
  if (!newChapter.title.trim() || !newChapter.file) return
  try {
    const fname = await uploadFile(newChapter.file, 'novel')
    n.chapters = [...(n.chapters || []), { title: newChapter.title, file: fname }]
    Object.assign(newChapter, { title: '', file: null })
    showToast('章节已上传,记得点保存')
  } catch (e) {
    showToast(e.message)
  }
}
const newArt = reactive({ title: '', note: '', file: null })

const oldPass = ref('')
const newPass = ref('')

onMounted(() => {
  if (token.value) loadContent()
})

async function login() {
  if (!password.value || loggingIn.value) return
  loggingIn.value = true
  loginError.value = ''
  try {
    const res = await fetch('/api/admin/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password: password.value }),
    })
    const data = await res.json()
    if (res.ok) {
      token.value = data.token
      localStorage.setItem('catcafe_token', data.token)
      loadContent()
    } else {
      loginError.value = data.error || '登录失败'
    }
  } catch {
    loginError.value = '网络错误,稍后再试'
  } finally {
    loggingIn.value = false
  }
}

function logout() {
  token.value = ''
  localStorage.removeItem('catcafe_token')
}

let contentSnapshot = '' // 服务器内容快照,用于检测外部改动

async function loadContent() {
  const res = await fetch('/api/admin/content-full', { headers: authHeaders() })
  const text = await res.text()
  contentSnapshot = text
  Object.assign(c, JSON.parse(text))
  loadWall()
  loadMsgs()
  loaded.value = true
}

function authHeaders() {
  return { Authorization: `Bearer ${token.value}` }
}

async function save() {
  saving.value = true
  try {
    // 防覆盖:保存前先对比服务器最新内容,若已被别处改动则刷新并提示
    const fresh = await (await fetch('/api/admin/content-full', { headers: authHeaders() })).text()
    if (contentSnapshot && fresh !== contentSnapshot) {
      contentSnapshot = fresh
      Object.assign(c, JSON.parse(fresh))
      showToast('内容刚在别处更新过,已为你刷新,请确认后再点保存 ⚠️')
      return
    }
    const res = await fetch('/api/admin/content', {
      method: 'PUT',
      headers: { ...authHeaders(), 'Content-Type': 'application/json' },
      body: JSON.stringify(c),
    })
    if (res.status === 401) return logout()
    if (res.ok) {
      contentSnapshot = await (await fetch('/api/admin/content-full', { headers: authHeaders() })).text() // 同步新快照
      showToast('已保存!前台刷新即可看到')
    } else {
      showToast('保存失败')
    }
  } catch {
    showToast('网络错误,保存失败')
  } finally {
    saving.value = false
  }
}

async function uploadFile(file, type) {
  const fd = new FormData()
  fd.append('type', type)
  fd.append('file', file)
  const res = await fetch('/api/admin/upload', { method: 'POST', headers: authHeaders(), body: fd })
  if (res.status === 401) {
    logout()
    showToast('登录已过期,请重新登录')
    throw new Error('登录已过期')
  }
  const data = await res.json()
  if (!res.ok) throw new Error(data.error || '上传失败')
  return data.file
}

function addNotice() {
  if (!newNotice.text.trim()) return
  c.notices = [{ ...newNotice }, ...(c.notices || [])]
  newNotice.text = ''
}

async function addNovel() {
  if (!newNovel.title.trim() || !newNovel.file) return
  try {
    const fname = await uploadFile(newNovel.file, 'novel')
    c.novels = [...(c.novels || []), { title: newNovel.title, desc: newNovel.desc, cup: newNovel.cup, cat: newNovel.cat, file: fname }]
    Object.assign(newNovel, { title: '', desc: '', file: '', cup: '中杯 · 微糖', cat: '连载中' })
    showToast('小说文件已上传,记得点保存')
  } catch (e) {
    showToast(e.message)
  }
}

async function addArt() {
  if (!newArt.title.trim() || !newArt.file) return
  try {
    const fname = await uploadFile(newArt.file, 'gallery')
    // 清掉占位卡,加入真实图
    c.gallery = (c.gallery || []).filter((g) => g.img)
    c.gallery.push({ img: fname, title: newArt.title, note: newArt.note })
    Object.assign(newArt, { title: '', note: '', file: null })
    showToast('插画已上传,记得点保存')
  } catch (e) {
    showToast(e.message)
  }
}

async function changePass() {
  if (newPass.value.length < 6) return showToast('新密码至少 6 位')
  const res = await fetch('/api/admin/password', {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    body: JSON.stringify({ old: oldPass.value, new: newPass.value }),
  })
  const data = await res.json()
  if (res.ok) {
    oldPass.value = newPass.value = ''
    showToast('密码已更新')
  } else {
    showToast(data.error || '修改失败')
  }
}

// 导出全部数据备份(zip)
async function exportData() {
  try {
    const res = await fetch('/api/admin/export', { headers: authHeaders() })
    if (res.status === 401) return logout()
    const blob = await res.blob()
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = `catcafe-backup-${new Date().toISOString().slice(0, 10)}.zip`
    a.click()
    URL.revokeObjectURL(a.href)
    showToast('备份包已开始下载')
  } catch {
    showToast('导出失败,稍后再试')
  }
}

function showToast(text) {
  toast.value = text
  setTimeout(() => (toast.value = ''), 3000)
}
</script>

<template>
  <div class="admin">
    <!-- 登录页 -->
    <div v-if="!token" class="login-card">
      <h1 class="font-cute">店主通道</h1>
      <p class="sub">爱丽丝猫猫酱专属,猪咪请回前台~</p>
      <input
        v-model="password"
        type="password"
        placeholder="输入店主密码"
        @keyup.enter="login"
      />
      <p v-if="loginError" class="err">{{ loginError }}</p>
      <button class="btn btn-primary" :disabled="loggingIn" @click="login">
        {{ loggingIn ? '验证中…' : '开门营业' }}
      </button>
      <a href="/" class="back">← 回前台</a>
    </div>

    <!-- 后台 -->
    <div v-else-if="loaded" class="panel">
      <header class="panel-head">
        <h1 class="font-cute">店主后台</h1>
        <div>
          <a href="/" target="_blank" class="back">看前台 ↗</a>
          <button class="btn btn-ghost small" @click="logout">退出</button>
        </div>
      </header>

      <nav class="tabs">
        <button
          v-for="t in tabs"
          :key="t.key"
          :class="{ active: tab === t.key }"
          @click="tab = t.key"
        >
          {{ t.label }}
        </button>
      </nav>

      <!-- 基本资料 -->
      <section v-if="tab === 'basic'" class="tab-body">
        <label>站点标题 <input v-model="c.title" /></label>
        <label>作者名 <input v-model="c.author" /></label>
        <label>P站主页 <input v-model="c.authorPixiv" /></label>
        <label>首页标语 <input v-model="c.slogan" /></label>
        <label>店主状态(显示在首页)
          <select v-model="c.authorStatus">
            <option>赶稿中</option>
            <option>摸鱼中</option>
            <option>冬眠中</option>
            <option>爆更中</option>
          </select>
        </label>
        <label>粉丝群号 <input v-model="c.fanClub.qq" /></label>
        <label>
          作者简介(每行一段)
          <textarea :value="(c.aboutAuthor || []).join('\n')" rows="4"
            @input="c.aboutAuthor = $event.target.value.split('\n')" />
        </label>
      </section>

      <!-- 公告管理 -->
      <section v-else-if="tab === 'notices'" class="tab-body">
        <div class="add-row">
          <input v-model="newNotice.date" type="date" />
          <select v-model="newNotice.tag">
            <option>公告</option><option>新坑</option><option>更新</option><option>活动</option>
          </select>
          <input v-model="newNotice.text" placeholder="公告内容…" class="grow" />
          <button class="btn btn-primary small" @click="addNotice">添加</button>
        </div>
        <div v-for="(n, i) in c.notices" :key="i" class="row">
          <span class="pill">{{ n.tag }}</span>
          <time>{{ n.date }}</time>
          <span class="grow">{{ n.text }}</span>
          <button class="del" @click="c.notices.splice(i, 1)">删除</button>
        </div>
      </section>

      <!-- 小说上架 -->
      <section v-else-if="tab === 'novels'" class="tab-body">
        <div class="add-form">
          <input v-model="newNovel.title" placeholder="作品标题" />
          <input v-model="newNovel.desc" placeholder="一句话简介" />
          <input v-model="newNovel.cup" placeholder="杯型标签(如:中杯·微糖)" />
          <select v-model="newNovel.cat">
            <option>连载中</option><option>已完结</option><option>番外</option>
          </select>
          <input type="file" accept=".txt,.md" @change="newNovel.file = $event.target.files[0]" />
          <button class="btn btn-primary small" :disabled="!newNovel.file" @click="addNovel">
            上传并加入菜单
          </button>
        </div>
        <div v-for="(n, i) in c.novels" :key="i" class="novel-block">
          <div class="row">
            <span class="grow"><b>{{ n.title }}</b> · {{ n.file }}</span>
            <button class="del" @click="expandedNovel = expandedNovel === i ? -1 : i">
              {{ expandedNovel === i ? '收起' : `章节(${(n.chapters || []).length})` }}
            </button>
            <button class="del" @click="c.novels.splice(i, 1)">下架</button>
          </div>
          <div v-if="expandedNovel === i" class="chapter-box">
            <div v-for="(ch, ci) in n.chapters || []" :key="ci" class="row">
              <span class="grow">{{ ci + 1 }}. {{ ch.title }} · {{ ch.file }}</span>
              <button class="del" @click="n.chapters.splice(ci, 1)">删除</button>
            </div>
            <div class="reply-row">
              <input v-model="newChapter.title" placeholder="章节名(如:第一章 相遇)" />
              <input type="file" accept=".txt,.md" @change="newChapter.file = $event.target.files[0]" />
              <button class="btn btn-primary small" :disabled="!newChapter.file" @click="addChapter(n)">
                加章节
              </button>
            </div>
            <p class="hint-text">有章节的作品,阅读器会显示章节切换;没有章节则整本阅读</p>
          </div>
        </div>
      </section>

      <!-- 插画管理 -->
      <section v-else-if="tab === 'gallery'" class="tab-body">
        <div class="add-form">
          <input v-model="newArt.title" placeholder="作品名" />
          <input v-model="newArt.note" placeholder="备注(可选)" />
          <input type="file" accept="image/*" @change="newArt.file = $event.target.files[0]" />
          <button class="btn btn-primary small" :disabled="!newArt.file" @click="addArt">
            上传并挂上墙
          </button>
        </div>
        <div class="art-grid">
          <figure v-for="(g, i) in c.gallery" :key="i">
            <img v-if="g.img" :src="`/gallery/${g.img}`" :alt="g.title" />
            <div v-else class="ph">占位</div>
            <figcaption>{{ g.title }}</figcaption>
            <button class="del" @click="c.gallery.splice(i, 1)">取下</button>
          </figure>
        </div>
      </section>

      <!-- 明信片管理 -->
      <section v-else-if="tab === 'wall'" class="tab-body">
        <p class="hint-text">群友寄来的明信片,不合适的可以撤下(会同时删除图片文件)</p>
        <div v-if="!wallPosts.length" class="hint-text">还没有明信片~</div>
        <div class="art-grid">
          <figure v-for="p in wallPosts" :key="p.img">
            <img :src="`/wall/${p.img}`" :alt="p.nick" />
            <figcaption>{{ p.nick }} · {{ p.note }}</figcaption>
            <button class="del" @click="delWall(p.img)">撤下</button>
          </figure>
        </div>
      </section>

      <!-- 留言管理 -->
      <section v-else-if="tab === 'msgs'" class="tab-body">
        <p class="hint-text">猪咪们的留言,可以回复或删除</p>
        <div v-if="!msgs.length" class="hint-text">还没有留言~</div>
        <div v-for="m in msgs" :key="m.time + m.nick" class="msg-row">
          <div class="msg-main">
            <div class="msg-meta"><b>{{ m.nick }}</b><time>{{ m.time }}</time></div>
            <p>{{ m.content }}</p>
            <p v-if="m.reply" class="has-reply">已回复:{{ m.reply }}</p>
            <div class="reply-row">
              <input v-model="replyText[m.time]" :placeholder="m.reply ? '修改回复…' : '回复这条留言…'" />
              <button class="btn btn-primary small" @click="replyMsg(m)">回复</button>
              <button class="del" @click="delMsg(m)">删除</button>
            </div>
          </div>
        </div>
      </section>

      <!-- 修改密码 -->
      <section v-else class="tab-body">
        <label>原密码 <input v-model="oldPass" type="password" /></label>
        <label>新密码(至少6位) <input v-model="newPass" type="password" /></label>
        <button class="btn btn-primary" @click="changePass">更新密码</button>
        <div class="export-box">
          <p class="hint-text">数据备份:留言、明信片、积分、公告等全部内容打包下载,建议每月存一份</p>
          <button class="btn btn-ghost" @click="exportData">导出数据备份</button>
        </div>
      </section>

      <footer v-if="tab !== 'account' && tab !== 'wall' && tab !== 'msgs'" class="panel-foot">
        <button class="btn btn-primary" :disabled="saving" @click="save">
          {{ saving ? '保存中…' : '保存全部修改' }}
        </button>
      </footer>

      <transition name="fade">
        <p v-if="toast" class="toast">{{ toast }}</p>
      </transition>
    </div>
  </div>
</template>

<style scoped>
.admin {
  min-height: 100vh;
  background: var(--bg);
  padding: 40px 16px;
  display: flex;
  justify-content: center;
}

.login-card {
  align-self: center;
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 40px;
  width: min(400px, 100%);
  display: flex;
  flex-direction: column;
  gap: 14px;
  text-align: center;
}

.login-card h1 {
  color: var(--ink);
  font-size: 1.6rem;
}

.sub {
  color: var(--muted);
  font-size: 0.9rem;
}

.err {
  color: var(--pink-deep);
  font-size: 0.85rem;
}

.back {
  color: var(--muted);
  font-size: 0.85rem;
  text-decoration: none;
}

.panel {
  width: min(860px, 100%);
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.panel-head h1 {
  font-size: 1.5rem;
  color: var(--ink);
}

.panel-head > div {
  display: flex;
  gap: 12px;
  align-items: center;
}

.tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.tabs button {
  border: 2px solid var(--pink-pale);
  background: #fff;
  color: var(--ink);
  border-radius: 999px;
  padding: 8px 18px;
  cursor: pointer;
  font-size: 0.9rem;
}

.tabs button.active {
  background: var(--pink);
  border-color: var(--pink);
  color: #fff;
}

.tab-body {
  background: var(--card);
  border: 2px solid var(--pink-pale);
  border-radius: var(--radius);
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

label {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 0.9rem;
  color: var(--muted);
}

input, textarea, select {
  border: 2px solid var(--pink-pale);
  border-radius: 12px;
  padding: 9px 14px;
  font-size: 0.95rem;
  font-family: inherit;
  color: var(--ink);
  outline: none;
}

input:focus, textarea:focus {
  border-color: var(--pink);
}

.add-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.grow {
  flex: 1;
  min-width: 140px;
}

.add-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-bottom: 14px;
  border-bottom: 2px dashed var(--pink-pale);
}

.row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.9rem;
  padding: 8px 0;
  border-bottom: 1px dashed var(--pink-pale);
}

.pill {
  background: var(--pink-pale);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 2px 10px;
  font-size: 0.75rem;
  flex-shrink: 0;
}

time {
  color: var(--muted);
  font-size: 0.8rem;
  flex-shrink: 0;
}

.del {
  border: none;
  background: none;
  color: var(--pink-deep);
  cursor: pointer;
  font-size: 0.85rem;
  flex-shrink: 0;
}

.del:hover {
  text-decoration: underline;
}

.art-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 14px;
}

.art-grid figure {
  text-align: center;
}

.art-grid img, .art-grid .ph {
  width: 100%;
  aspect-ratio: 4/3;
  object-fit: cover;
  border-radius: 10px;
  background: var(--pink-pale);
}

.art-grid .ph {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--muted);
  font-size: 0.8rem;
}

.art-grid figcaption {
  font-size: 0.85rem;
  margin: 6px 0 2px;
}

.panel-foot {
  margin-top: 20px;
  text-align: center;
}

.btn.small {
  padding: 8px 18px;
  font-size: 0.9rem;
}

.hint-text {
  color: var(--muted);
  font-size: 0.88rem;
}

.msg-row {
  border: 1px solid var(--pink-pale);
  border-radius: 12px;
  padding: 12px 16px;
}

.msg-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.msg-meta b {
  color: var(--pink-deep);
  font-size: 0.9rem;
}

.msg-main > p {
  font-size: 0.92rem;
  margin-bottom: 8px;
  word-break: break-word;
}

.has-reply {
  color: var(--muted);
  font-size: 0.85rem;
}

.reply-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.reply-row input {
  flex: 1;
}

.export-box {
  border-top: 2px dashed var(--pink-pale);
  padding-top: 16px;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
}

.novel-block {
  border-bottom: 1px dashed var(--pink-pale);
  padding-bottom: 8px;
}

.chapter-box {
  background: var(--bg);
  border-radius: 12px;
  padding: 12px 16px;
  margin: 4px 0 10px;
}

.toast {
  position: fixed;
  bottom: 28px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--ink);
  color: #fff;
  padding: 10px 24px;
  border-radius: 999px;
  font-size: 0.9rem;
  box-shadow: var(--shadow);
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>

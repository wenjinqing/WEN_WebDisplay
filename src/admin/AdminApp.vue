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
  { key: 'account', label: '修改密码' },
]

// 新条目表单
const newNotice = reactive({ date: new Date().toISOString().slice(0, 10), tag: '公告', text: '' })
const newNovel = reactive({ title: '', desc: '', file: '', cup: '中杯 · 微糖' })
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

async function loadContent() {
  const res = await fetch('/api/content')
  Object.assign(c, await res.json())
  loaded.value = true
}

function authHeaders() {
  return { Authorization: `Bearer ${token.value}` }
}

async function save() {
  saving.value = true
  try {
    const res = await fetch('/api/admin/content', {
      method: 'PUT',
      headers: { ...authHeaders(), 'Content-Type': 'application/json' },
      body: JSON.stringify(c),
    })
    if (res.status === 401) return logout()
    showToast(res.ok ? '已保存!前台刷新即可看到 ✅' : '保存失败')
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
    c.novels = [...(c.novels || []), { title: newNovel.title, desc: newNovel.desc, cup: newNovel.cup, file: fname }]
    Object.assign(newNovel, { title: '', desc: '', file: '', cup: '中杯 · 微糖' })
    showToast('小说文件已上传,记得点保存 ✅')
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
    showToast('插画已上传,记得点保存 ✅')
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
    showToast('密码已更新 ✅')
  } else {
    showToast(data.error || '修改失败')
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
      <h1 class="font-cute">🐾 店主通道</h1>
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
        <h1 class="font-cute">🐾 店主后台</h1>
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
          <input type="file" accept=".txt,.md" @change="newNovel.file = $event.target.files[0]" />
          <button class="btn btn-primary small" :disabled="!newNovel.file" @click="addNovel">
            上传并加入菜单
          </button>
        </div>
        <div v-for="(n, i) in c.novels" :key="i" class="row">
          <span class="grow"><b>{{ n.title }}</b> · {{ n.file }}</span>
          <button class="del" @click="c.novels.splice(i, 1)">下架</button>
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

      <!-- 修改密码 -->
      <section v-else class="tab-body">
        <label>原密码 <input v-model="oldPass" type="password" /></label>
        <label>新密码(至少6位) <input v-model="newPass" type="password" /></label>
        <button class="btn btn-primary" @click="changePass">更新密码</button>
      </section>

      <footer v-if="tab !== 'account'" class="panel-foot">
        <button class="btn btn-primary" :disabled="saving" @click="save">
          {{ saving ? '保存中…' : '💾 保存全部修改' }}
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

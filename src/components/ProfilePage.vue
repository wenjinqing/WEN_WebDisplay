<script setup>
// APP「我的」页:绑定昵称(免登录) + 检查更新 / 应用内升级
import { ref, onMounted } from 'vue'
import { App } from '@capacitor/app'
import { Filesystem, Directory } from '@capacitor/filesystem'
import { FileOpener } from '@capawesome-team/capacitor-file-opener'
import { Browser } from '@capacitor/browser'
import PawPrint from './PawPrint.vue'
import { profile, setNick } from '../app.js'

const nickInput = ref(profile.nick)
const nickSaved = ref(false)

function saveNick() {
  const n = nickInput.value.trim().slice(0, 16)
  if (!n) return
  setNick(n)
  nickSaved.value = true
  setTimeout(() => (nickSaved.value = false), 2000)
}

// ---------- 检查更新 ----------
const curVersion = ref('')
const curBuild = ref(0)
const checking = ref(false)
const updateInfo = ref(null) // { versionName, versionCode, notes, url }
const noUpdate = ref(false)
const downloading = ref(false)
const progress = ref(0)
const updateError = ref('')

onMounted(async () => {
  try {
    const info = await App.getInfo()
    curVersion.value = info.version
    curBuild.value = Number(info.build) || 0
  } catch {
    curVersion.value = '未知'
  }
})

async function checkUpdate() {
  checking.value = true
  noUpdate.value = false
  updateInfo.value = null
  updateError.value = ''
  try {
    const res = await fetch('https://alicefans.cn/downloads/app/version.json?t=' + Date.now())
    if (!res.ok) throw new Error('bad status')
    const v = await res.json()
    if (v.versionCode > curBuild.value) {
      updateInfo.value = v
    } else {
      noUpdate.value = true
    }
  } catch {
    updateError.value = '检查失败,看看网络是不是被猫踩了'
  } finally {
    checking.value = false
  }
}

// 流式下载 APK(带进度)→ 写入缓存目录 → 调起系统安装器
async function doUpdate() {
  const v = updateInfo.value
  if (!v) return
  downloading.value = true
  progress.value = 0
  updateError.value = ''
  try {
    const res = await fetch(v.url)
    if (!res.ok) throw new Error('download failed')
    const total = Number(res.headers.get('content-length')) || 0
    const reader = res.body.getReader()
    const chunks = []
    let received = 0
    for (;;) {
      const { done, value } = await reader.read()
      if (done) break
      chunks.push(value)
      received += value.length
      if (total) progress.value = Math.round((received / total) * 100)
    }
    // 合并并转 base64(分块转,避免 40MB 一次性爆栈)
    let binary = ''
    const CHUNK = 0x8000
    const all = new Uint8Array(received)
    let offset = 0
    for (const c of chunks) {
      all.set(c, offset)
      offset += c.length
    }
    let b64 = ''
    for (let i = 0; i < all.length; i += CHUNK) {
      b64 += String.fromCharCode.apply(null, all.subarray(i, i + CHUNK))
    }
    b64 = btoa(b64)

    const fileName = 'catcafe-update-' + v.versionCode + '.apk'
    await Filesystem.writeFile({
      path: fileName,
      data: b64,
      directory: Directory.Cache,
      recursive: true,
    })
    const uri = (await Filesystem.getUri({ path: fileName, directory: Directory.Cache })).uri
    await FileOpener.openFile({
      path: uri,
      mimeType: 'application/vnd.android.package-archive',
    })
  } catch (err) {
    // 应用内安装失败就退回系统浏览器下载,并把原因显示出来方便排查
    updateError.value = '应用内安装没拉起(' + ((err && err.message) || err) + '),已帮你打开浏览器下载'
    try {
      await Browser.open({ url: v.url })
    } catch {
      window.open(v.url, '_blank')
    }
  } finally {
    downloading.value = false
  }
}
</script>

<template>
  <section class="profile container">
    <div class="card">
      <div class="who">
        <PawPrint :size="44" color="#f9718f" />
        <div>
          <div class="nick font-cute">{{ profile.nick || '未绑定昵称' }}</div>
          <div class="sub">{{ profile.nick ? '这只猪咪有名字啦' : '绑个昵称,留言/游戏全店通用' }}</div>
        </div>
      </div>
      <div class="nick-row">
        <input v-model="nickInput" maxlength="16" placeholder="输入你的昵称" />
        <button class="btn" @click="saveNick">{{ nickSaved ? '已绑定 ✓' : '绑定' }}</button>
      </div>
      <p class="hint">不用注册不用登录,昵称只存在你自己手机里。</p>
    </div>

    <div class="card">
      <div class="ver-row">
        <span>当前版本</span>
        <span class="ver">v{{ curVersion }} ({{ curBuild }})</span>
      </div>
      <button class="btn wide" :disabled="checking || downloading" @click="checkUpdate">
        {{ checking ? '正在检查…' : '检查更新' }}
      </button>

      <div v-if="noUpdate" class="tip ok">已经是最新版啦,猫猫说很满意</div>

      <div v-if="updateInfo" class="update-box">
        <div class="new-ver">发现新版本 v{{ updateInfo.versionName }}</div>
        <div v-if="updateInfo.notes" class="notes">{{ updateInfo.notes }}</div>
        <button class="btn wide pink" :disabled="downloading" @click="doUpdate">
          {{ downloading ? `下载中 ${progress}%` : '立即更新' }}
        </button>
        <div v-if="downloading" class="bar"><i :style="{ width: progress + '%' }" /></div>
        <p class="hint">下载完成后会弹出系统安装界面,点「安装」即覆盖升级,数据不丢。</p>
      </div>

      <div v-if="updateError" class="tip err">{{ updateError }}</div>
    </div>
  </section>
</template>

<style scoped>
.profile {
  padding-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card {
  background: var(--card);
  border: 2px dashed var(--pink-soft);
  border-radius: 18px;
  padding: 20px;
}

.who {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.nick {
  font-size: 1.3rem;
  color: var(--ink);
}

.sub {
  font-size: 0.85rem;
  color: var(--ink-soft, #9b8a90);
  margin-top: 2px;
}

.nick-row {
  display: flex;
  gap: 10px;
}

.nick-row input {
  flex: 1;
  border: 2px solid var(--pink-soft);
  border-radius: 999px;
  padding: 10px 16px;
  font-size: 0.95rem;
  font-family: inherit;
  background: transparent;
  color: var(--ink);
  outline: none;
}

.nick-row input:focus {
  border-color: var(--pink-deep);
}

.btn {
  border: none;
  border-radius: 999px;
  padding: 10px 20px;
  background: var(--pink-pale);
  color: var(--pink-deep);
  font-size: 0.95rem;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  white-space: nowrap;
}

.btn.wide {
  width: 100%;
  margin-top: 12px;
  padding: 12px;
}

.btn.pink {
  background: var(--pink-deep);
  color: #fff;
}

.btn:disabled {
  opacity: 0.6;
}

.hint {
  font-size: 0.8rem;
  color: var(--ink-soft, #9b8a90);
  margin: 10px 0 0;
  line-height: 1.6;
}

.ver-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.95rem;
  color: var(--ink);
}

.ver {
  color: var(--ink-soft, #9b8a90);
}

.tip {
  margin-top: 12px;
  text-align: center;
  font-size: 0.9rem;
}

.tip.ok {
  color: #4caf7d;
}

.tip.err {
  color: #e0678a;
}

.update-box {
  margin-top: 14px;
  border-top: 1px dashed var(--pink-soft);
  padding-top: 14px;
}

.new-ver {
  font-weight: 700;
  color: var(--pink-deep);
}

.notes {
  font-size: 0.85rem;
  color: var(--ink);
  margin-top: 6px;
  line-height: 1.6;
  white-space: pre-line;
}

.bar {
  margin-top: 10px;
  height: 6px;
  border-radius: 3px;
  background: var(--pink-pale);
  overflow: hidden;
}

.bar i {
  display: block;
  height: 100%;
  background: var(--pink-deep);
  transition: width 0.2s;
}
</style>

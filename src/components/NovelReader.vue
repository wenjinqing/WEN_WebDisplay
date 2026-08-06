<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps({
  novel: { type: Object, required: true },
})
const emit = defineEmits(['close'])

// 有 chapters 就用章节列表,否则把整本当成单章
const chapters = computed(() =>
  props.novel.chapters && props.novel.chapters.length
    ? props.novel.chapters
    : [{ title: props.novel.title, file: props.novel.file }]
)

const chIdx = ref(0)
const loading = ref(true)
const error = ref('')
const pages = ref([])
const page = ref(0)
const fontSize = ref(17)
const contentEl = ref(null)

const total = computed(() => pages.value.length)
const currentParas = computed(() => pages.value[page.value] || [])
const progressKey = computed(() => `read_${chapters.value[chIdx.value].file}`)

function paginate(text) {
  const paras = text.split(/\n+/).map((p) => p.trim()).filter(Boolean)
  const result = []
  let buf = []
  let len = 0
  for (const p of paras) {
    if (len + p.length > 550 && buf.length) {
      result.push(buf)
      buf = []
      len = 0
    }
    buf.push(p)
    len += p.length
  }
  if (buf.length) result.push(buf)
  return result
}

async function loadChapter() {
  loading.value = true
  error.value = ''
  try {
    const res = await fetch(`/downloads/${chapters.value[chIdx.value].file}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    pages.value = paginate(await res.text())
    // 恢复该章节的阅读进度
    const saved = Number(localStorage.getItem(progressKey.value))
    page.value = saved > 0 && saved < pages.value.length ? saved : 0
  } catch (e) {
    error.value = '小说加载失败了,刷新一下再试试喵~'
    pages.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  document.body.style.overflow = 'hidden'
  window.addEventListener('keydown', onKey)
  loadChapter()
})

onUnmounted(() => {
  document.body.style.overflow = ''
  window.removeEventListener('keydown', onKey)
})

watch(chIdx, loadChapter)

watch(page, (p) => {
  contentEl.value?.scrollTo({ top: 0 })
  localStorage.setItem(progressKey.value, String(p))
})

function prev() {
  if (page.value > 0) page.value--
}
function next() {
  if (page.value < total.value - 1) {
    page.value++
  } else if (chIdx.value < chapters.value.length - 1) {
    chIdx.value++ // 本章读完自动进下一章
  }
}
function onKey(e) {
  if (e.key === 'Escape') emit('close')
  if (e.key === 'ArrowLeft') prev()
  if (e.key === 'ArrowRight') next()
}
</script>

<template>
  <div class="reader-mask" @click.self="emit('close')">
    <div class="reader" role="dialog" :aria-label="novel.title">
      <header class="reader-head">
        <h3 class="font-cute">{{ novel.title }}</h3>
        <select v-if="chapters.length > 1" v-model="chIdx" class="ch-select" aria-label="选择章节">
          <option v-for="(ch, i) in chapters" :key="ch.file" :value="i">{{ ch.title }}</option>
        </select>
        <button class="close" aria-label="关闭" @click="emit('close')">✕</button>
      </header>

      <div ref="contentEl" class="reader-body" :style="{ fontSize: fontSize + 'px' }">
        <p v-if="loading" class="state">正在端上桌……</p>
        <p v-else-if="error" class="state">{{ error }}</p>
        <template v-else>
          <p v-for="(para, i) in currentParas" :key="i" class="para">{{ para }}</p>
          <p v-if="page === total - 1 && chIdx < chapters.length - 1" class="state next-hint">
            — 本章完,点「下一页」进入下一章 —
          </p>
          <p v-else-if="page === total - 1" class="state next-hint">— 全书完,感谢阅读 —</p>
        </template>
      </div>

      <footer class="reader-foot">
        <button class="ctl" :disabled="page === 0" @click="prev">‹ 上一页</button>
        <span class="pageinfo">{{ total ? `${page + 1} / ${total}` : '' }}</span>
        <button class="ctl" :disabled="page >= total - 1 && chIdx >= chapters.length - 1" @click="next">
          下一页 ›
        </button>
        <span class="sep" />
        <button class="ctl small" :disabled="fontSize <= 14" @click="fontSize -= 1">A-</button>
        <button class="ctl small" :disabled="fontSize >= 22" @click="fontSize += 1">A+</button>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.reader-mask {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(91, 58, 71, 0.45);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.reader {
  width: min(720px, 100%);
  height: min(86vh, 900px);
  background: #fff9f4;
  border-radius: var(--radius);
  border: 2px solid var(--pink-soft);
  box-shadow: 0 24px 64px rgba(91, 58, 71, 0.3);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.reader-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px;
  border-bottom: 2px dashed var(--pink-pale);
  background: #fff;
}

.reader-head h3 {
  font-size: 1.05rem;
  color: var(--ink);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 1;
  min-width: 0;
}

.ch-select {
  border: 2px solid var(--pink-pale);
  border-radius: 999px;
  padding: 5px 12px;
  font-size: 0.85rem;
  color: var(--pink-deep);
  background: #fff;
  outline: none;
  max-width: 180px;
  flex-shrink: 0;
}

.close {
  flex-shrink: 0;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: none;
  background: var(--pink-pale);
  color: var(--pink-deep);
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}

.close:hover {
  background: var(--pink-soft);
}

.reader-body {
  flex: 1;
  overflow-y: auto;
  padding: 28px 32px;
  line-height: 2;
  color: #4a323d;
}

.para {
  margin-bottom: 1em;
  text-indent: 2em;
}

.state {
  text-align: center;
  color: var(--muted);
  margin-top: 40px;
  text-indent: 0;
}

.next-hint {
  margin-top: 24px;
  font-size: 0.85em;
}

.reader-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  border-top: 2px dashed var(--pink-pale);
  background: #fff;
}

.ctl {
  border: none;
  background: var(--pink-pale);
  color: var(--pink-deep);
  border-radius: 999px;
  padding: 8px 16px;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background 0.2s;
}

.ctl:hover:not(:disabled) {
  background: var(--pink-soft);
}

.ctl:disabled {
  opacity: 0.4;
  cursor: default;
}

.ctl.small {
  padding: 6px 12px;
}

.pageinfo {
  color: var(--muted);
  font-size: 0.85rem;
  min-width: 56px;
  text-align: center;
}

.sep {
  flex: 1;
}

@media (max-width: 720px) {
  .reader {
    height: 100%;
    border-radius: 0;
    border: none;
  }
  .reader-mask {
    padding: 0;
  }
  .reader-body {
    padding: 20px;
  }
  .reader-foot {
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
    padding: 10px 12px;
  }
  .ctl {
    padding: 8px 12px;
    font-size: 0.85rem;
  }
  .sep {
    display: none;
  }
}
</style>

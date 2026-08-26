// ============================================================
// APP 模式共享状态:是否运行在 Capacitor 原生壳里 + 个人资料
// ============================================================
import { reactive } from 'vue'
import { Capacitor } from '@capacitor/core'

export const isApp = Capacitor.isNativePlatform()

// 个人资料:免登录,昵称存本地,与留言板(catcafe_nick)通用
export const profile = reactive({
  nick: localStorage.getItem('catcafe_nick') || '',
})

export function setNick(name) {
  profile.nick = name
  if (name) localStorage.setItem('catcafe_nick', name)
}

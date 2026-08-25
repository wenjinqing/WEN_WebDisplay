import type { CapacitorConfig } from '@capacitor/cli'

const config: CapacitorConfig = {
  appId: 'cn.alicefans.catcafe',
  appName: '涩猫咖啡厅',
  webDir: 'dist',
  android: {
    // WebView 内允许跳回线上站点(图片/下载文件都在服务器上)
    allowMixedContent: false,
  },
}

export default config

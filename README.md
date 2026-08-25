# 🐱 爱丽丝的小涩猫咖啡厅

> 百合小说作者「爱丽丝猫猫酱」的粉丝后援站 —— 粉粉嫩嫩、手绘猫咖风,顺便塞满彩蛋的个人展示网页。

<p align="center">
  <a href="https://alicefans.cn"><img src="https://img.shields.io/badge/主站-alicefans.cn-ff8fab?style=flat-square" alt="主站"></a>
  <a href="https://alicefans.asia"><img src="https://img.shields.io/badge/开发站-alicefans.asia-b05787?style=flat-square" alt="开发站"></a>
  <a href="https://www.pixiv.net/users/16689973"><img src="https://img.shields.io/badge/作者P站-爱丽丝猫猫酱-9b7ede?style=flat-square" alt="Pixiv"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-63c48c?style=flat-square" alt="License"></a>
</p>

<p align="center"><img src="public/og-cover.png" alt="站点封面" width="600" /></p>

## ✨ 这是什么

一个给粉丝们「坐下来点杯小说慢慢看」的小窝:小说在线阅读与下载、留言催更、明信片墙,还有一只全店共养的电子猪。整站走手绘贴纸风,配色是草莓奶昔味,连错误提示都自称「猪咪」、后缀带「喵~」。

## 🎪 功能一览

**内容展示**

- 🏠 首页天幕 + 公告板 + 作者介绍 + 营业报告(今日/本月/总账)
- 📚 小说菜单 + 阅读器(章节选择/字号/衬线字体/进度记忆)+ 作品评论(爪印评分)
- 🖼️ 插画画廊 + 明信片墙(访客可寄、可点赞)
- 📡 RSS 订阅 + 新章节邮件/站内订阅提醒

**社区互动**

- 💬 互动区三标签:留言墙(双列气泡)/ 写留言 / 催更墙,店长与「猪咪君君」回复带专属徽章
- 🐷 猪咪聚集地 `/hub`:访客发帖(文字/图片)、拍爪、评论、改删自己的帖子
- 🏅 积分与头衔:评论 +5、寄明信片 +10、发聚集地 +5~10 鱼干,小奶猫 → 小猪咪 → 大橘 → 猫老大

**养成与彩蛋**

- 🐖 全店共养电子猪:喂食/摸头/经验进化/饱食心情衰减
- 🎰 成就图鉴 · 打地鼠 · 像素宠物 · 天气层 · 弹幕层 · 爪印雨 · 今日运势
- 👾 赛博裂缝彩蛋:随机出现,点进去是 glitch 世界
- 🌙 夜间模式:20:00–07:00 自动开启,暖紫罗兰配色,正文对比度 13:1(AAA)

**运营与安全**

- 🔑 店主后台 `/admin`:登录、改密、内容全量管理、文件上传、删帖、数据导出
- 🤖 Agent 接口:API Key 认证,「猪咪君君」拥有**与店主同权的全资源管理**——发公告/改内容、回留言、管聚集地、管理文件/明信片/评论/积分/订阅、导出数据,详见 [docs/AGENT-API.md](docs/AGENT-API.md)
- 🛡️ 敏感词过滤 · 发帖限流 · 群号打码 + 答题解锁 · 保留昵称保护 · 编辑密钥(访客只能改删自己的帖)
- 📱 PWA(可安装、离线缓存)+ HTTPS

## 🛠️ 技术栈

| 层 | 选型 | 说明 |
|---|---|---|
| 前端 | Vue 3 + Vite | 三个入口:`index.html` 主站 / `hub.html` 聚集地 / `admin.html` 后台 |
| 后端 | Go 标准库 | 单文件 `server/main.go`,仅两个小依赖(qrcode、image 压缩),适合低流量粉丝站 |
| 存储 | JSON 文件 | 留言/内容/积分/统计/订阅等全部落盘 `/var/lib/catcafe/` |
| 部署 | Nginx 反代 | 静态文件 + `127.0.0.1:9090` API,双站共享聚集地图库 |

## 🚀 快速开始

```bash
# 前端(开发)
npm install
npm run dev        # 主站,另有 /hub /admin 入口

# 前端(生产)
npm run build      # 产物在 dist/,交给 Nginx 托管

# 后端
cd server
go run .           # 监听 127.0.0.1:9090,首次启动生成默认店主密码 alice233
```

> ⚠️ 后端数据目录 `/var/lib/catcafe`、`/var/www/...` 是 Linux 生产路径;本机 Windows 上可以跑通接口逻辑,但写文件会失败。开发时请用开发站 `alicefans.asia` 验证完整链路。

## 📁 目录结构

```
├── index.html / hub.html / admin.html   # 三个页面入口
├── src/
│   ├── App.vue          # 主站
│   ├── hub/             # 猪咪聚集地
│   ├── admin/           # 店主后台
│   ├── components/      # 30+ 组件(留言/猪场/阅读器/彩蛋…)
│   ├── data.js          # 站点内容默认值(与后端 defaultContent 一致)
│   └── style.css        # 设计变量 + 夜间模式
├── server/main.go       # 全部 API(零框架,约 2500 行)
├── docs/AGENT-API.md    # 猪咪君君接口文档
└── public/              # 图片/字体/PWA 资源
```

## 🤖 猪咪君君(Agent)

站点配了一只「猪咪君君」:**与店主同权**,可回复留言、发公告、管理聚集地、上传文件、管积分订阅、导出数据,全走带 `X-API-Key` 的 HTTP 接口,任何能发请求的东西(脚本、另一个 AI)都能当猪咪君君。完整文档见 [docs/AGENT-API.md](docs/AGENT-API.md)。

## 📜 License

[MIT](LICENSE) © 2026 WENJINQING

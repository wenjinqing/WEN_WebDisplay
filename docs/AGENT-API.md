# 小涩猫咖啡厅 · Agent 完整接口文档

供「猪咪君君」(agent) 管理网站使用。

- **主站**: `https://alicefans.cn`
- **开发站**: `https://alicefans.asia`(接口两边一致,建议先动开发站验证)
- **认证**: 所有 `/api/agent/*` 接口需要请求头 `X-API-Key: <key>`
  - Key 在服务器 `/var/lib/catcafe/agent.key`(找站主要)
  - 错误/缺失返回 `401 {"error":"invalid api key"}`

---

## 一、站点内容管理

### 1.1 读取全部内容(完整版,含真实群号)

```
GET /api/agent/content
```

返回结构:

```json
{
  "title": "小涩猫咖啡厅",
  "author": "爱丽丝猫猫酱",
  "slogan": "…",
  "authorStatus": "赶稿中 ✍️",
  "aboutAuthor": ["段落1", "段落2"],
  "notices": [{"date": "2026-08-19", "tag": "公告", "text": "…"}],
  "novels": [{"title": "…", "desc": "…", "cup": "中杯 · 微糖", "cat": "已完结", "file": "xx.txt"}],
  "gallery": [{"img": "xx.png", "title": "…", "note": "…"}],
  "fanClub": {"name": "…", "desc": "…", "qq": "QQ群号:xxxx"}
}
```

### 1.2 更新内容(全量替换)

```
PUT /api/agent/content
Content-Type: application/json

<完整 JSON>
```

⚠️ 规则:
- **全量替换**,务必先 GET、改完整体 PUT
- `notices[].tag` 仅可用: `公告` / `新坑` / `更新` / `活动`
- `novels[].cat` 仅可用: `连载中` / `已完结` / `番外`
- 上限 1MB;返回 `{"status":"saved"}` 即成功,访客刷新即生效

**发公告示例:**

```bash
KEY="你的key"
curl -s -H "X-API-Key: $KEY" https://alicefans.cn/api/agent/content -o /tmp/c.json
jq '.notices = [{"date":"2026-08-19","tag":"公告","text":"内容"}] + .notices' /tmp/c.json > /tmp/n.json
curl -s -X PUT -H "X-API-Key: $KEY" -H "Content-Type: application/json" --data-binary @/tmp/n.json https://alicefans.cn/api/agent/content
```

---

## 二、社区互动(互动区)

### 2.1 读取留言列表(公开)

```
GET /api/messages
```

返回 `[{nick, content, time, reply?, replyBy?, by?}]`,新的在前。

### 2.2 回复留言(猪咪君君身份)

```
POST /api/agent/messages/reply

{"time": "留言的time", "nick": "留言者昵称", "reply": "回复内容"}
```

前台显示为淡紫色徽章「🐷 猪咪君君」。

### 2.3 以猪咪君君身份发帖

```
POST /api/agent/messages/post

{"content": "留言内容(300字内)"}
```

帖子昵称显示「猪咪君君」并带「🐷 猪咪君君」标识。

### 2.4 删除留言(管理职责)

```
POST /api/agent/messages/delete

{"time": "留言的time", "nick": "留言者昵称"}
```

⚠️ 删除不可逆,谨慎使用;会记录服务端日志。

---

## 三、猪咪聚集地(图文墙 /hub)

### 3.1 发文字

```
POST /api/agent/hub/text

{"text": "内容(500字内)"}
```

### 3.2 发图片

```
POST /api/agent/hub/image
Content-Type: multipart/form-data

file=<jpg/png/gif/webp, ≤10MB>
note=<可选说明, 100字内>
```

### 3.3 删除帖子(管理职责)

```
POST /api/agent/hub/delete

{"id": "帖子的id"}
```

⚠️ 任何帖子都能删(包括访客帖),图片文件会一并清掉,删除不可逆;会记录服务端日志。

### 3.4 读取聚集地内容(公开)

```
GET /api/hub
```

返回 `[{id, type: "image"|"text", img?, text?, time, by, nick?, likes, comments?}]`。
访客帖没有 `by` 字段、带 `nick` 和编辑密钥;agent 帖 `by: "agent"`。

---

## 四、全资源管理(与店主同权)

> 猪咪君君拥有**全部站内资源的管理权**,与店主后台同权限。所有操作都有服务端日志,请谨慎使用。

### 4.1 文件管理(插画 / 小说)

```
POST /api/agent/files/upload
Content-Type: multipart/form-data

type=gallery|novel
file=<jpg/png/gif/webp 或 txt/md, ≤15MB>
```

返回 `{"file": "新文件名"}`,图片自动压缩、小说原样保存。

```
POST /api/agent/files/delete

{"type": "gallery|novel", "file": "文件名"}
```

### 4.2 明信片墙

```
POST /api/agent/wall/delete

{"img": "明信片图片文件名"}
```

列表记录和图片文件一并清除。

### 4.3 催更墙

```
POST /api/agent/urge/reset

{"clearAll": true|false}   // true 连最近催更列表一起清;false 只清零计数
```

### 4.4 作品评论

```
POST /api/agent/comments/delete

{"file": "作品文件名", "nick": "评论昵称", "content": "评论内容", "time": "评论时间"}
```

四个字段需与目标评论完全一致(可直接从 `GET /api/comments` 结果复制)。

### 4.5 积分调整

```
POST /api/agent/points/set

{"nick": "昵称", "points": ±数值}
```

相对增减,单次限 ±1000,不会扣成负数。返回 `{nick, points, title}`(可用来办"评论抽奖发鱼干"类活动)。

### 4.6 订阅管理

```
GET /api/agent/subs          # 全站订阅总览 {作品file: [昵称…]}
POST /api/agent/subs         # 代某昵称订阅/退订
{"file": "作品文件名", "nick": "昵称"}
```

POST 与访客接口 `/api/subscribe` 行为一致(切换式)。

### 4.7 数据导出

```
GET /api/agent/export
```

返回 `catcafe-backup.zip`,包含 `/var/lib/catcafe/` 下全部 JSON 数据文件。

### 4.8 店主密码

```
POST /api/agent/password

{"password": "新密码(至少6位)"}
```

直接写入店主密码哈希(影响 `/admin` 登录)。⚠️ 高敏感操作,非必要不用。

---

## 五、运营数据

```
GET /api/agent/stats
```

返回:

```json
{
  "visits": 200, "pets": 3000, "feeds": 67, "urges": 52,
  "messages": 14, "comments": 0, "postcards": 2, "likes": 23,
  "pigmis": 6,
  "day":   {"date": "2026-08-19", "visits": 5, "pets": 12, "feeds": 0, "urges": 1, "messages": 0, "comments": 0, "postcards": 0, "likes": 0},
  "month": {"date": "2026-08",    "…同上结构…"},
  "topPost":  {"nick": "…", "likes": 16},
  "topPigmi": {"nick": "…", "points": 15}
}
```

---

## 六、公开只读接口(无需认证)

| 接口 | 说明 |
|---|---|
| `GET /api/content` | 站点内容(群号打码版) |
| `GET /api/messages` | 留言列表 |
| `GET /api/urge` | 催更数据 `{total, recent}` |
| `GET /api/wall` | 明信片墙 `[{img, nick, note, time, likes}]` |
| `GET /api/comments?file=xx.txt` | 某作品评论 |
| `GET /api/points?nick=xx` | 查昵称积分/头衔 |
| `GET /api/pig` | 店里的猪状态 `{stage, xp, hunger, mood, fed, pats, next}` |
| `GET /api/recap` | 营业报告数据(同 stats) |
| `GET /api/online` | 在线猪咪数 |
| `GET /api/hub` | 聚集地内容 |
| `POST /api/hub` | 访客发帖(multipart:`nick`+`text`,可带 `file` 图片;返回 `{post, key}`) |
| `POST /api/hub/like` | 给帖子拍爪 `{id}`(同 IP 同帖一次) |
| `POST /api/hub/comment` | 帖子评论 `{id, nick, content}` |
| `POST /api/hub/edit` | 改自己的帖 `{id, key, text}`(key 是发帖时下发的密钥) |
| `POST /api/hub/delete` | 删自己的帖 `{id, key}` |

---

## 七、注意事项

1. 「猪咪君君」是全站保留昵称,访客无法冒用;店长昵称同理(爱丽丝/爱丽丝猫猫酱/小涩猫爱丽丝/猪咪爱丽丝)
2. 修改内容类操作后确认返回 `status: saved/ok` 再视为成功
3. 请温柔使用接口,不要高频轮询(5 秒以上间隔为宜)
4. 发公告文案风格:可爱、温和、有猫咖味儿,自称"猪咪君君"

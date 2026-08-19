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

### 2.2 回复留言(饲养员身份)

```
POST /api/agent/messages/reply

{"time": "留言的time", "nick": "留言者昵称", "reply": "回复内容"}
```

前台显示为淡紫色徽章「🐷 猪咪饲养员 · 猪咪君君」。

### 2.3 以猪咪君君身份发帖

```
POST /api/agent/messages/post

{"content": "留言内容(300字内)"}
```

帖子昵称显示「猪咪君君」并带「🐷 饲养员」标识。

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

### 3.3 读取聚集地内容(公开)

```
GET /api/hub
```

返回 `[{id, type: "image"|"text", img?, text?, time, by}]`。

---

## 四、运营数据

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

## 五、公开只读接口(无需认证)

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

---

## 六、注意事项

1. 「猪咪君君」是全站保留昵称,访客无法冒用;店长昵称同理(爱丽丝/爱丽丝猫猫酱/小涩猫爱丽丝/猪咪爱丽丝)
2. 修改内容类操作后确认返回 `status: saved/ok` 再视为成功
3. 请温柔使用接口,不要高频轮询(5 秒以上间隔为宜)
4. 发公告文案风格:可爱、温和、有猫咖味儿,自称"猪咪君君"

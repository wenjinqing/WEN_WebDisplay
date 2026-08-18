# 小涩猫咖啡厅 · Agent 管理接口文档

供另一台服务器的 agent 作为网站管理员使用。

## 认证

所有接口需要在请求头携带 API Key:

```
X-API-Key: <向站主要,在服务器 /var/lib/catcafe/agent.key>
```

无 Key 或 Key 错误返回 `401 {"error":"invalid api key"}`。

## 接口

### 1. 读取站点全部内容

```
GET https://alicefans.asia/api/agent/content
```

返回完整 JSON,结构:

```json
{
  "title": "小涩猫咖啡厅",
  "author": "爱丽丝猫猫酱",
  "slogan": "...",
  "authorStatus": "赶稿中 ✍️",
  "aboutAuthor": ["段落1", "段落2"],
  "notices": [{"date": "2026-08-06", "tag": "公告", "text": "..."}],
  "novels": [{"title": "...", "desc": "...", "cup": "中杯 · 微糖", "cat": "已完结", "file": "xxx.txt", "chapters": []}],
  "gallery": [{"img": "xxx.png", "title": "...", "note": "..."}],
  "fanClub": {"name": "...", "desc": "...", "qq": "..."}
}
```

### 2. 更新站点内容(全量替换)

```
PUT https://alicefans.asia/api/agent/content
Content-Type: application/json

<完整 content JSON>
```

⚠️ 注意:
- **全量替换**,修改前请先 GET 读取、在读取结果上改、再整体 PUT 回去,不要自己从零构造
- `notices[].tag` 可选值:公告 / 新坑 / 更新 / 活动
- 上限 1MB

### 3. 运营数据(用于决策参考)

```
GET https://alicefans.asia/api/agent/stats
```

返回:`{visits, pets, feeds, urges, messages, comments, postcards, likes, pigmis, day:{...今日数据}, topPost, topPigmi}`

## 典型用法示例

发布一条新公告:

```bash
KEY="..."

# 1. 读取
curl -s -H "X-API-Key: $KEY" https://alicefans.asia/api/agent/content -o /tmp/content.json

# 2. 用 jq 在 notices 前面插入一条
jq '.notices = [{"date": "2026-08-18", "tag": "公告", "text": "新公告内容"}] + .notices' /tmp/content.json > /tmp/new.json

# 3. 写回
curl -s -X PUT -H "X-API-Key: $KEY" -H "Content-Type: application/json" \
  --data-binary @/tmp/new.json https://alicefans.asia/api/agent/content
```

## 注意

- 前台是动态读取内容的,PUT 成功后访客刷新即生效,无需重启任何服务
- 接口无频率限制,但请温柔使用,不要高频轮询

---

## 社区管理(猪咪君君 · 猪咪饲养员身份)

### ④ 回复留言

```
POST https://alicefans.asia/api/agent/messages/reply
Content-Type: application/json

{"time": "2026-08-09 16:26", "nick": "留言者昵称", "reply": "回复内容"}
```

回复会显示为 **「🐷 猪咪饲养员 · 猪咪君君」**徽章(淡紫色,和店长的粉色徽章区分开)。
`time` + `nick` 从公开接口 `GET /api/messages` 获取。

### ⑤ 删除不当留言

```
POST https://alicefans.asia/api/agent/messages/delete

{"time": "2026-08-09 16:26", "nick": "留言者昵称"}
```

⚠️ 「猪咪君君」是全站保留昵称,普通访客无法冒用。

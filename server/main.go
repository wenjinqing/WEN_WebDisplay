// 爱丽丝的小涩猫咖啡厅 —— API 服务
// 功能:留言板 / 催更墙 / 站点内容管理 / 店主后台登录与文件上传
// 零依赖,仅标准库;数据存 JSON 文件,适合低流量粉丝站
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	dataFile    = "/var/lib/catcafe/data.json"    // 留言/催更数据
	contentFile = "/var/lib/catcafe/content.json" // 站点内容
	passFile    = "/var/lib/catcafe/admin.hash"   // 店主密码哈希
	wallFile    = "/var/lib/catcafe/wall.json"    // 明信片墙数据
	galleryDir  = "/var/www/catcafe/gallery"      // 插画目录
	novelDir    = "/var/www/downloads"            // 小说目录
	wallDir     = "/var/www/catcafe/wall"         // 明信片图片目录

	maxMessages = 200
	maxUrges    = 50
	maxNickLen  = 20
	maxTextLen  = 300
	minInterval = 5 * time.Second
	tokenTTL    = 7 * 24 * time.Hour
	listenAddr  = "127.0.0.1:9090"

	defaultAdminPass = "alice233" // 首次启动的默认密码,登录后请立即修改
)

// 站点内容默认值(与前端 src/data.js 一致,content.json 不存在时使用)
const defaultContent = `{
  "title": "爱丽丝的小涩猫咖啡厅",
  "author": "爱丽丝猫猫酱",
  "authorPixiv": "https://www.pixiv.net/users/16689973",
  "slogan": "百合小说作家长期营业中,猪咪们里边请~",
  "aboutAuthor": [
    "爱丽丝猫猫酱,一位专注于百合题材的轻小说作者。",
    "笔下是女孩子之间酸酸甜甜、偶尔让人脸红心跳的故事。文字软软的,后劲足足的。",
    "长期驻扎 P 站更新,欢迎来咖啡厅坐坐,点一杯小说慢慢看。"
  ],
  "notices": [
    {"date": "2026-08-06", "tag": "公告", "text": "咖啡厅正式开业啦!本站是猫猫酱的粉丝后援页,小说会持续上架,欢迎收藏~"}
  ],
  "novels": [],
  "gallery": [],
  "fanClub": {"name": "猪咪饲养基地", "desc": "店里常客都是可爱的猪咪。进来一起催更、吸猫、聊猫猫酱的新坑!", "qq": "QQ群号:1054390069"}
}`

type Message struct {
	Nick    string `json:"nick"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Reply   string `json:"reply,omitempty"` // 店主回复
}

type Urge struct {
	Nick string `json:"nick"`
	Time string `json:"time"`
}

type Store struct {
	Messages  []Message `json:"messages"`
	UrgeTotal int       `json:"urgeTotal"`
	Urges     []Urge    `json:"urges"`
}

var (
	mu       sync.Mutex
	store    Store
	rateMu   sync.Mutex
	lastPost = map[string]time.Time{}
	lastWall = map[string]time.Time{} // 明信片墙单独限流(60秒/次)
	likedBy  = map[string]bool{}      // 点赞记录:"ip|img"(重启后清零,可接受)
	tokenMu  sync.Mutex
	tokens   = map[string]time.Time{} // token → 过期时间
	wallMu   sync.Mutex
	wall     []WallPost
)

type WallPost struct {
	Img   string `json:"img"`
	Nick  string `json:"nick"`
	Note  string `json:"note"`
	Time  string `json:"time"`
	Likes int    `json:"likes"`
}

func loadWall() {
	data, err := os.ReadFile(wallFile)
	if err == nil {
		json.Unmarshal(data, &wall)
	}
}

func saveWall() {
	data, _ := json.MarshalIndent(wall, "", "  ")
	tmp := wallFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err == nil {
		os.Rename(tmp, wallFile)
	}
}

// ---------- 工具 ----------

func loadStore() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		store = Store{}
		return
	}
	if err := json.Unmarshal(data, &store); err != nil {
		log.Printf("数据文件损坏,从零开始: %v", err)
		store = Store{}
	}
}

func saveStore() {
	data, _ := json.MarshalIndent(store, "", "  ")
	tmp := dataFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		log.Printf("保存失败: %v", err)
		return
	}
	os.Rename(tmp, dataFile)
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func rateLimited(r *http.Request) bool {
	ip := clientIP(r)
	rateMu.Lock()
	defer rateMu.Unlock()
	if t, ok := lastPost[ip]; ok && time.Since(t) < minInterval {
		return true
	}
	lastPost[ip] = time.Now()
	return false
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func clean(s string, max int) string {
	s = strings.TrimSpace(s)
	r := []rune(s)
	if len(r) > max {
		r = r[:max]
	}
	return string(r)
}

// ---------- 店主认证 ----------

func hashPass(p string) string {
	sum := sha256.Sum256([]byte(p))
	return hex.EncodeToString(sum[:])
}

func checkPass(p string) bool {
	data, err := os.ReadFile(passFile)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == hashPass(p)
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	tok := hex.EncodeToString(b)
	tokenMu.Lock()
	tokens[tok] = time.Now().Add(tokenTTL)
	tokenMu.Unlock()
	return tok
}

func authed(r *http.Request) bool {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return false
	}
	tok := strings.TrimPrefix(h, "Bearer ")
	tokenMu.Lock()
	defer tokenMu.Unlock()
	exp, ok := tokens[tok]
	if !ok || time.Now().After(exp) {
		delete(tokens, tok)
		return false
	}
	return true
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if !authed(r) {
		writeJSON(w, 401, map[string]string{"error": "请先登录"})
		return false
	}
	return true
}

// ---------- 留言板 & 催更墙 ----------

func handleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, 200, store.Messages)
	case http.MethodPost:
		var m Message
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeJSON(w, 400, map[string]string{"error": "请求格式不对"})
			return
		}
		m.Nick = clean(m.Nick, maxNickLen)
		m.Content = clean(m.Content, maxTextLen)
		if m.Content == "" {
			writeJSON(w, 400, map[string]string{"error": "留言内容不能为空喵"})
			return
		}
		if m.Nick == "" {
			m.Nick = "匿名猪咪"
		}
		if rateLimited(r) {
			writeJSON(w, 429, map[string]string{"error": "喝口水休息一下,发得太快啦"})
			return
		}
		m.Time = time.Now().Format("2006-01-02 15:04")
		mu.Lock()
		store.Messages = append([]Message{m}, store.Messages...)
		if len(store.Messages) > maxMessages {
			store.Messages = store.Messages[:maxMessages]
		}
		saveStore()
		mu.Unlock()
		addPoints(m.Nick, 5) // 留言 +5 鱼干
		writeJSON(w, 200, m)
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

func handleUrge(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, 200, map[string]any{"total": store.UrgeTotal, "recent": store.Urges})
	case http.MethodPost:
		var u Urge
		json.NewDecoder(r.Body).Decode(&u)
		u.Nick = clean(u.Nick, maxNickLen)
		if u.Nick == "" {
			u.Nick = "匿名猪咪"
		}
		if rateLimited(r) {
			writeJSON(w, 429, map[string]string{"error": "刚刚催过啦,猫猫酱看到啦"})
			return
		}
		u.Time = time.Now().Format("2006-01-02 15:04")
		mu.Lock()
		store.UrgeTotal++
		store.Urges = append([]Urge{u}, store.Urges...)
		if len(store.Urges) > maxUrges {
			store.Urges = store.Urges[:maxUrges]
		}
		total := store.UrgeTotal
		saveStore()
		mu.Unlock()
		addPoints(u.Nick, 2) // 催更 +2 鱼干
		writeJSON(w, 200, map[string]any{"total": total, "urge": u})
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

// ---------- 图片压缩 ----------

// 超过 1200px 宽的图片自动等比缩小并重新压缩;GIF 原样保存以保留动图
func saveImageCompressed(data []byte, ext, dst string) error {
	if ext == ".gif" {
		return os.WriteFile(dst, data, 0644)
	}
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("图片解码失败: %w", err)
	}
	b := src.Bounds()
	const maxW = 1200
	if b.Dx() > maxW {
		h := b.Dy() * maxW / b.Dx()
		scaled := image.NewRGBA(image.Rect(0, 0, maxW, h))
		draw.CatmullRom.Scale(scaled, scaled.Bounds(), src, b, draw.Over, nil)
		src = scaled
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	if ext == ".png" { // PNG 可能带透明通道,保留格式
		return png.Encode(f, src)
	}
	return jpeg.Encode(f, src, &jpeg.Options{Quality: 82})
}

// ---------- 站点内容 ----------

// GET /api/content — 公开,前台启动时拉取
func handleContent(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(contentFile)
	if err != nil {
		data = []byte(defaultContent)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// PUT /api/admin/content — 登录后整体替换站点内容
func handlePutContent(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPut {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	body := http.MaxBytesReader(w, r.Body, 1<<20) // 1MB 上限
	data, err := io.ReadAll(body)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "内容太大了"})
		return
	}
	var check map[string]any
	if err := json.Unmarshal(data, &check); err != nil {
		writeJSON(w, 400, map[string]string{"error": "内容不是合法 JSON"})
		return
	}
	pretty, _ := json.MarshalIndent(check, "", "  ")
	tmp := contentFile + ".tmp"
	if err := os.WriteFile(tmp, pretty, 0644); err != nil {
		writeJSON(w, 500, map[string]string{"error": "保存失败"})
		return
	}
	os.Rename(tmp, contentFile)
	writeJSON(w, 200, map[string]string{"status": "saved"})
}

// ---------- 店主登录 ----------

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	if rateLimited(r) {
		writeJSON(w, 429, map[string]string{"error": "试得太频繁了,歇一会儿"})
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !checkPass(req.Password) {
		writeJSON(w, 401, map[string]string{"error": "密码不对喵"})
		return
	}
	writeJSON(w, 200, map[string]string{"token": newToken()})
}

func handleChangePass(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !checkPass(req.Old) {
		writeJSON(w, 401, map[string]string{"error": "原密码不对"})
		return
	}
	if len(req.New) < 6 {
		writeJSON(w, 400, map[string]string{"error": "新密码至少 6 位"})
		return
	}
	if err := os.WriteFile(passFile, []byte(hashPass(req.New)), 0600); err != nil {
		writeJSON(w, 500, map[string]string{"error": "保存失败"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

// ---------- 文件上传/删除 ----------

var allowExt = map[string]map[string]bool{
	"gallery": {".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true},
	"novel":   {".txt": true, ".md": true},
}

var typeDir = map[string]string{
	"gallery": galleryDir,
	"novel":   novelDir,
}

func safeName(name string) string {
	name = filepath.Base(name) // 去掉路径,防穿越
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '.', r == '-', r == '_', r >= 0x4e00 && r <= 0x9fff: // 允许中文文件名
			return r
		}
		return '-'
	}, name)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	typ := r.FormValue("type")
	dir, ok := typeDir[typ]
	if !ok {
		writeJSON(w, 400, map[string]string{"error": "type 只能是 gallery 或 novel"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 15<<20) // 15MB 上限
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		writeJSON(w, 400, map[string]string{"error": "文件太大或格式不对(上限15MB)"})
		return
	}
	f, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "没有收到文件"})
		return
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowExt[typ][ext] {
		writeJSON(w, 400, map[string]string{"error": "不支持的文件类型: " + ext})
		return
	}
	name := fmt.Sprintf("%d-%s", time.Now().Unix(), safeName(header.Filename))
	data, err := io.ReadAll(f)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "文件读取失败"})
		return
	}
	if typ == "gallery" { // 只有图片走压缩,小说文本原样保存
		err = saveImageCompressed(data, ext, filepath.Join(dir, name))
	} else {
		err = os.WriteFile(filepath.Join(dir, name), data, 0644)
	}
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "文件处理失败,换一个试试"})
		return
	}
	log.Printf("店主上传了文件: %s/%s", typ, name)
	writeJSON(w, 200, map[string]string{"file": name})
}

func handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Type string `json:"type"`
		File string `json:"file"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "请求格式不对"})
		return
	}
	dir, ok := typeDir[req.Type]
	if !ok || strings.Contains(req.File, "/") || strings.Contains(req.File, "..") {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	os.Remove(filepath.Join(dir, req.File))
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// ---------- 明信片墙 ----------

// GET /api/wall — 拉取明信片列表
// POST /api/wall — 群友寄明信片(multipart: file + nick + note)
func handleWall(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		wallMu.Lock()
		defer wallMu.Unlock()
		if wall == nil {
			wall = []WallPost{}
		}
		writeJSON(w, 200, wall)
	case http.MethodPost:
		ip := clientIP(r)
		rateMu.Lock()
		if t, ok := lastWall[ip]; ok && time.Since(t) < 60*time.Second {
			rateMu.Unlock()
			writeJSON(w, 429, map[string]string{"error": "明信片一分钟只能寄一张哦"})
			return
		}
		rateMu.Unlock()

		r.Body = http.MaxBytesReader(w, r.Body, 6<<20) // 约5MB上限
		if err := r.ParseMultipartForm(6 << 20); err != nil {
			writeJSON(w, 400, map[string]string{"error": "图片太大了(上限5MB)"})
			return
		}
		f, header, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "别忘了选一张图片喵"})
			return
		}
		defer f.Close()

		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowExt["gallery"][ext] {
			writeJSON(w, 400, map[string]string{"error": "只支持 jpg/png/gif/webp 图片"})
			return
		}
		name := fmt.Sprintf("%d-%s", time.Now().Unix(), safeName(header.Filename))
		data, err := io.ReadAll(f)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "文件读取失败"})
			return
		}
		if err := saveImageCompressed(data, ext, filepath.Join(wallDir, name)); err != nil {
			writeJSON(w, 400, map[string]string{"error": "图片处理失败,换一张试试"})
			return
		}

		nick := clean(r.FormValue("nick"), maxNickLen)
		if nick == "" {
			nick = "匿名猪咪"
		}
		post := WallPost{
			Img:  name,
			Nick: nick,
			Note: clean(r.FormValue("note"), 60),
			Time: time.Now().Format("2006-01-02 15:04"),
		}
		wallMu.Lock()
		wall = append([]WallPost{post}, wall...)
		if len(wall) > 30 {
			for _, old := range wall[30:] { // 删掉最老的图文件
				os.Remove(filepath.Join(wallDir, old.Img))
			}
			wall = wall[:30]
		}
		saveWall()
		wallMu.Unlock()

		rateMu.Lock()
		lastWall[ip] = time.Now()
		rateMu.Unlock()

		log.Printf("新明信片: %s 寄出了 %s", nick, name)
		addPoints(nick, 10) // 寄明信片 +10 鱼干
		writeJSON(w, 200, post)
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

// POST /api/wall/like {img} — 给明信片点赞(同一 IP 对同一张只能点一次)
func handleWallLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Img string `json:"img"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		strings.Contains(req.Img, "/") || strings.Contains(req.Img, "..") {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	key := clientIP(r) + "|" + req.Img
	rateMu.Lock()
	if likedBy[key] {
		rateMu.Unlock()
		writeJSON(w, 429, map[string]string{"error": "这张你已经拍过爪啦"})
		return
	}
	rateMu.Unlock()

	wallMu.Lock()
	found := false
	likes := 0
	for i := range wall {
		if wall[i].Img == req.Img {
			wall[i].Likes++
			likes = wall[i].Likes
			found = true
			break
		}
	}
	if found {
		saveWall()
	}
	wallMu.Unlock()
	if !found {
		writeJSON(w, 404, map[string]string{"error": "明信片不存在"})
		return
	}

	rateMu.Lock()
	likedBy[key] = true
	rateMu.Unlock()
	writeJSON(w, 200, map[string]int{"likes": likes})
}

// POST /api/admin/wall/delete {img} — 店主撤下明信片
func handleWallDelete(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Img string `json:"img"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		strings.Contains(req.Img, "/") || strings.Contains(req.Img, "..") {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	wallMu.Lock()
	for i, p := range wall {
		if p.Img == req.Img {
			wall = append(wall[:i], wall[i+1:]...)
			break
		}
	}
	saveWall()
	wallMu.Unlock()
	os.Remove(filepath.Join(wallDir, req.Img))
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// ---------- 留言管理(店主) ----------

// POST /api/admin/messages/delete {time, nick} — 删除留言
func handleMessageDelete(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	var req struct {
		Time string `json:"time"`
		Nick string `json:"nick"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	mu.Lock()
	for i, m := range store.Messages {
		if m.Time == req.Time && m.Nick == req.Nick {
			store.Messages = append(store.Messages[:i], store.Messages[i+1:]...)
			break
		}
	}
	saveStore()
	mu.Unlock()
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// POST /api/admin/messages/reply {time, nick, reply} — 回复留言
func handleMessageReply(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	var req struct {
		Time  string `json:"time"`
		Nick  string `json:"nick"`
		Reply string `json:"reply"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	mu.Lock()
	found := false
	for i := range store.Messages {
		if store.Messages[i].Time == req.Time && store.Messages[i].Nick == req.Nick {
			store.Messages[i].Reply = clean(req.Reply, maxTextLen)
			found = true
			break
		}
	}
	if found {
		saveStore()
	}
	mu.Unlock()
	if !found {
		writeJSON(w, 404, map[string]string{"error": "留言不存在"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

// ---------- 访问统计 ----------

type Stats struct {
	Total      int    `json:"total"`
	Today      string `json:"today"`
	TodayCount int    `json:"todayCount"`
	Pets       int    `json:"pets"`  // 猫猫被撸次数
	Feeds      int    `json:"feeds"` // 全店投喂次数
}

var (
	statsMu sync.Mutex
	stats   Stats
)

const (
	statsFile    = "/var/lib/catcafe/stats.json"
	commentsFile = "/var/lib/catcafe/comments.json" // 作品评论
	pointsFile   = "/var/lib/catcafe/points.json"   // 猪咪积分
)

func loadStats() {
	data, err := os.ReadFile(statsFile)
	if err == nil {
		json.Unmarshal(data, &stats)
	}
}

// GET /api/hit — 页面打开时调用,累计访问量
func handleHit(w http.ResponseWriter, r *http.Request) {
	statsMu.Lock()
	today := time.Now().Format("2006-01-02")
	if stats.Today != today {
		stats.Today = today
		stats.TodayCount = 0
	}
	stats.Total++
	stats.TodayCount++
	data, _ := json.Marshal(stats)
	os.WriteFile(statsFile, data, 0644)
	out := stats
	statsMu.Unlock()
	writeJSON(w, 200, out)
}

// ---------- 作品评论 ----------

type Comment struct {
	Nick    string `json:"nick"`
	Content string `json:"content"`
	Score   int    `json:"score"` // 爪印评分 1-5
	Time    string `json:"time"`
}

var (
	commentsMu sync.Mutex
	comments   = map[string][]Comment{} // key: 小说文件名
)

func loadComments() {
	data, err := os.ReadFile(commentsFile)
	if err == nil {
		json.Unmarshal(data, &comments)
	}
}

func saveComments() {
	data, _ := json.MarshalIndent(comments, "", "  ")
	tmp := commentsFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err == nil {
		os.Rename(tmp, commentsFile)
	}
}

// GET /api/comments?file=xxx — 拉取某作品的评论
// POST /api/comments {file, nick, content, score} — 发评论
func handleComments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f := r.URL.Query().Get("file")
		commentsMu.Lock()
		list := comments[f]
		commentsMu.Unlock()
		if list == nil {
			list = []Comment{}
		}
		writeJSON(w, 200, list)
	case http.MethodPost:
		var cm Comment
		var req struct {
			File string `json:"file"`
			Comment
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, 400, map[string]string{"error": "请求格式不对"})
			return
		}
		if strings.Contains(req.File, "/") || strings.Contains(req.File, "..") || req.File == "" {
			writeJSON(w, 400, map[string]string{"error": "参数不对"})
			return
		}
		cm.Nick = clean(req.Nick, maxNickLen)
		cm.Content = clean(req.Content, maxTextLen)
		cm.Score = req.Score
		if cm.Content == "" {
			writeJSON(w, 400, map[string]string{"error": "评论内容不能为空喵"})
			return
		}
		if cm.Score < 1 || cm.Score > 5 {
			cm.Score = 5
		}
		if cm.Nick == "" {
			cm.Nick = "匿名猪咪"
		}
		if rateLimited(r) {
			writeJSON(w, 429, map[string]string{"error": "发得太快啦,歇一会儿"})
			return
		}
		cm.Time = time.Now().Format("2006-01-02 15:04")
		commentsMu.Lock()
		comments[req.File] = append([]Comment{cm}, comments[req.File]...)
		if len(comments[req.File]) > 100 {
			comments[req.File] = comments[req.File][:100]
		}
		saveComments()
		commentsMu.Unlock()
		addPoints(cm.Nick, 5) // 评论 +5 鱼干
		writeJSON(w, 200, cm)
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

// ---------- 猪咪积分/头衔 ----------

var (
	pointsMu sync.Mutex
	points   = map[string]int{} // key: 昵称
)

func loadPoints() {
	data, err := os.ReadFile(pointsFile)
	if err == nil {
		json.Unmarshal(data, &points)
	}
}

func savePoints() {
	data, _ := json.MarshalIndent(points, "", "  ")
	tmp := pointsFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err == nil {
		os.Rename(tmp, pointsFile)
	}
}

// 匿名猪咪不累计积分(没法认领)
func addPoints(nick string, n int) {
	if nick == "" || nick == "匿名猪咪" {
		return
	}
	pointsMu.Lock()
	points[nick] += n
	savePoints()
	pointsMu.Unlock()
}

// 头衔:小奶猫 → 小猪咪 → 大橘 → 猫老大
func titleFor(p int) (string, int) {
	switch {
	case p < 50:
		return "小奶猫", 50
	case p < 150:
		return "小猪咪", 150
	case p < 300:
		return "大橘", 300
	default:
		return "猫老大", 0
	}
}

// GET /api/points?nick=xxx — 查询积分和头衔
func handlePoints(w http.ResponseWriter, r *http.Request) {
	nick := clean(r.URL.Query().Get("nick"), maxNickLen)
	if nick == "" {
		writeJSON(w, 400, map[string]string{"error": "填一下昵称喵"})
		return
	}
	pointsMu.Lock()
	p := points[nick]
	pointsMu.Unlock()
	title, nextAt := titleFor(p)
	writeJSON(w, 200, map[string]any{
		"nick": nick, "points": p, "title": title, "nextAt": nextAt,
	})
}

// GET /api/pet — 撸猫总数
// POST /api/pet — 撸一下
func handlePet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		statsMu.Lock()
		n := stats.Pets
		statsMu.Unlock()
		writeJSON(w, 200, map[string]int{"pets": n})
	case http.MethodPost:
		// 娱乐计数器,允许快速点,开心就好
		statsMu.Lock()
		stats.Pets++
		data, _ := json.Marshal(stats)
		os.WriteFile(statsFile, data, 0644)
		n := stats.Pets
		statsMu.Unlock()
		writeJSON(w, 200, map[string]int{"pets": n})
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

// POST /api/feed — 投喂一次,返回全店总投喂数
func handleFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	statsMu.Lock()
	stats.Feeds++
	data, _ := json.Marshal(stats)
	os.WriteFile(statsFile, data, 0644)
	n := stats.Feeds
	statsMu.Unlock()
	writeJSON(w, 200, map[string]int{"feeds": n})
}

// POST /api/points/add {nick, n} — 娱乐加分(捡鱼干等),n 上限 10
func handlePointsAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Nick string `json:"nick"`
		N    int    `json:"n"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	nick := clean(req.Nick, maxNickLen)
	if nick == "" || nick == "匿名猪咪" {
		writeJSON(w, 400, map[string]string{"error": "需要昵称才能攒鱼干"})
		return
	}
	if req.N < 1 || req.N > 10 {
		req.N = 5
	}
	if rateLimited(r) {
		writeJSON(w, 429, map[string]string{"error": "手速太快啦"})
		return
	}
	pointsMu.Lock()
	points[nick] += req.N
	total := points[nick]
	savePoints()
	pointsMu.Unlock()
	title, _ := titleFor(total)
	writeJSON(w, 200, map[string]any{"points": total, "title": title})
}

// ---------- 在线猪咪 ----------

var (
	onlineMu sync.Mutex
	visitors = map[string]time.Time{}
)

// GET /api/online — 心跳+返回 90 秒内的活跃猪咪数
func handleOnline(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)
	onlineMu.Lock()
	visitors[ip] = time.Now()
	n := 0
	for v, t := range visitors {
		if time.Since(t) < 90*time.Second {
			n++
		} else {
			delete(visitors, v) // 顺手清理
		}
	}
	onlineMu.Unlock()
	writeJSON(w, 200, map[string]int{"online": n})
}

// ---------- 启动 ----------

func main() {
	for _, d := range []string{"/var/lib/catcafe", galleryDir, novelDir, wallDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatal(err)
		}
	}
	loadStore()
	loadWall()
	loadStats()
	loadComments()
	loadPoints()

	// 首次启动:写入默认店主密码哈希
	if _, err := os.Stat(passFile); os.IsNotExist(err) {
		os.WriteFile(passFile, []byte(hashPass(defaultAdminPass)), 0600)
		log.Printf("!!! 已生成默认店主密码 %q,请登录后台后立即修改", defaultAdminPass)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/urge", handleUrge)
	mux.HandleFunc("/api/wall", handleWall)
	mux.HandleFunc("/api/wall/like", handleWallLike)
	mux.HandleFunc("/api/hit", handleHit)
	mux.HandleFunc("/api/pet", handlePet)
	mux.HandleFunc("/api/feed", handleFeed)
	mux.HandleFunc("/api/online", handleOnline)
	mux.HandleFunc("/api/comments", handleComments)
	mux.HandleFunc("/api/points", handlePoints)
	mux.HandleFunc("/api/points/add", handlePointsAdd)
	mux.HandleFunc("/api/content", handleContent)
	mux.HandleFunc("/api/admin/login", handleLogin)
	mux.HandleFunc("/api/admin/content", handlePutContent)
	mux.HandleFunc("/api/admin/password", handleChangePass)
	mux.HandleFunc("/api/admin/upload", handleUpload)
	mux.HandleFunc("/api/admin/delete-file", handleDeleteFile)
	mux.HandleFunc("/api/admin/wall/delete", handleWallDelete)
	mux.HandleFunc("/api/admin/messages/delete", handleMessageDelete)
	mux.HandleFunc("/api/admin/messages/reply", handleMessageReply)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	log.Printf("猫咖 API 营业中 → %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

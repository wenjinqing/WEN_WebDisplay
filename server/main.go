// 爱丽丝的小涩猫咖啡厅 —— API 服务
// 功能:留言板 / 催更墙 / 站点内容管理 / 店主后台登录与文件上传
// 零依赖,仅标准库;数据存 JSON 文件,适合低流量粉丝站
package main

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/skip2/go-qrcode"
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
	tokenTTL    = 30 * 24 * time.Hour
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
	Reply   string `json:"reply,omitempty"`   // 回复内容
	ReplyBy string `json:"replyBy,omitempty"` // 回复者:空=店长,猪咪君君=饲养员agent
	By      string `json:"by,omitempty"`      // 发帖者:空=访客,agent=猪咪君君
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

const tokensFile = "/var/lib/catcafe/tokens.json"
const hubFile = "/var/lib/catcafe/hub.json" // 猪咪聚集地图文
const hubDir = "/var/www/shared/hub"        // 聚集地图片目录(双站共享)
const agentKeyFile = "/var/lib/catcafe/agent.key"

var agentKey string // 机器管理员 API Key

func loadAgentKey() {
	data, err := os.ReadFile(agentKeyFile)
	if err == nil {
		agentKey = strings.TrimSpace(string(data))
	}
}

// X-API-Key 校验(常量时间比较防时序侧漏)
func agentOK(r *http.Request) bool {
	if agentKey == "" {
		return false
	}
	got := r.Header.Get("X-API-Key")
	return subtle.ConstantTimeCompare([]byte(got), []byte(agentKey)) == 1
}

// 登录态持久化:重启服务不掉登录
func loadTokens() {
	data, err := os.ReadFile(tokensFile)
	if err != nil {
		return
	}
	raw := map[string]int64{}
	if json.Unmarshal(data, &raw) != nil {
		return
	}
	tokenMu.Lock()
	defer tokenMu.Unlock()
	for tok, exp := range raw {
		if t := time.Unix(exp, 0); time.Now().Before(t) {
			tokens[tok] = t
		}
	}
}

func saveTokens() {
	raw := map[string]int64{}
	for tok, exp := range tokens {
		raw[tok] = exp.Unix()
	}
	data, _ := json.Marshal(raw)
	os.WriteFile(tokensFile, data, 0600)
}

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

// 店长专属昵称,猪咪们不能使用
var reservedNicks = map[string]bool{
	"爱丽丝":     true,
	"爱丽丝猫猫酱": true,
	"小涩猫爱丽丝": true,
	"猪咪爱丽丝":  true,
	"猪咪君君":   true, // 饲养员 agent 专用
}

func isReserved(nick string) bool {
	return reservedNicks[strings.TrimSpace(nick)]
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
	saveTokens()
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
		if isReserved(m.Nick) {
			writeJSON(w, 400, map[string]string{"error": "这个名字属于店长,换一个喵~"})
			return
		}
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
		bumpDaily("messages")
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
		if isReserved(u.Nick) {
			writeJSON(w, 400, map[string]string{"error": "这个名字属于店长,换一个喵~"})
			return
		}
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
		bumpDaily("urges")
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

const qqMask = "答题验证后可见"

// GET /api/content — 公开,前台启动时拉取(群号打码)
func handleContent(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(contentFile)
	if err != nil {
		data = []byte(defaultContent)
	}
	// 隐藏群号:fanClub.qq 替换为提示语
	var c map[string]any
	if json.Unmarshal(data, &c) == nil {
		if fc, ok := c["fanClub"].(map[string]any); ok {
			if _, has := fc["qq"]; has {
				fc["qq"] = qqMask
				data, _ = json.Marshal(c)
			}
		}
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
		if isReserved(nick) {
			writeJSON(w, 400, map[string]string{"error": "这个名字属于店长,换一个喵~"})
			return
		}
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
		bumpDaily("postcards")
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
	bumpDaily("likes")

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

type Daily struct {
	Date      string `json:"date"`
	Visits    int    `json:"visits"`
	Pets      int    `json:"pets"`
	Feeds     int    `json:"feeds"`
	Urges     int    `json:"urges"`
	Messages  int    `json:"messages"`
	Comments  int    `json:"comments"`
	Postcards int    `json:"postcards"`
	Likes     int    `json:"likes"`
}

type Stats struct {
	Total      int    `json:"total"`
	Today      string `json:"today"`
	TodayCount int    `json:"todayCount"`
	Pets       int    `json:"pets"`  // 猫猫被撸次数
	Feeds      int    `json:"feeds"` // 全店投喂次数
	Day        Daily  `json:"day"`   // 今日数据
	Month      Daily  `json:"month"` // 本月数据
}

// 每日计数:跨天自动清零,跨月自动清月账
func bumpDaily(field string) {
	statsMu.Lock()
	defer statsMu.Unlock()
	today := time.Now().Format("2006-01-02")
	if stats.Day.Date != today {
		stats.Day = Daily{Date: today}
	}
	thisMonth := time.Now().Format("2006-01")
	if stats.Month.Date != thisMonth {
		stats.Month = Daily{Date: thisMonth}
	}
	bumpOne(&stats.Day, field)
	bumpOne(&stats.Month, field)
	data, _ := json.Marshal(stats)
	os.WriteFile(statsFile, data, 0644)
}

func bumpOne(d *Daily, field string) {
	switch field {
	case "visits":
		d.Visits++
	case "pets":
		d.Pets++
	case "feeds":
		d.Feeds++
	case "urges":
		d.Urges++
	case "messages":
		d.Messages++
	case "comments":
		d.Comments++
	case "postcards":
		d.Postcards++
	case "likes":
		d.Likes++
	}
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
	out := stats
	statsMu.Unlock()
	bumpDaily("visits")
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
		if isReserved(cm.Nick) {
			writeJSON(w, 400, map[string]string{"error": "这个名字属于店长,换一个喵~"})
			return
		}
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
		bumpDaily("comments")
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
		n := stats.Pets
		statsMu.Unlock()
		bumpDaily("pets")
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
	n := stats.Feeds
	statsMu.Unlock()
	bumpDaily("feeds")
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
	if isReserved(nick) {
		writeJSON(w, 400, map[string]string{"error": "这个名字属于店长哦"})
		return
	}
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

// ---------- 年度账本 & 数据导出 ----------

// GET /api/recap — 年度回忆数据汇总
func handleRecap(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	msgCount := len(store.Messages)
	urgeTotal := store.UrgeTotal
	mu.Unlock()

	wallMu.Lock()
	postCount := len(wall)
	likeTotal := 0
	var topPost *WallPost
	for i := range wall {
		likeTotal += wall[i].Likes
		if topPost == nil || wall[i].Likes > topPost.Likes {
			topPost = &wall[i]
		}
	}
	wallMu.Unlock()

	statsMu.Lock()
	visits := stats.Total
	pets := stats.Pets
	feeds := stats.Feeds
	statsMu.Unlock()

	pointsMu.Lock()
	pigmiCount := len(points)
	topNick, topPoints := "", 0
	for nick, p := range points {
		if p > topPoints {
			topNick, topPoints = nick, p
		}
	}
	pointsMu.Unlock()

	commentsMu.Lock()
	commentCount := 0
	for _, list := range comments {
		commentCount += len(list)
	}
	commentsMu.Unlock()

	statsMu.Lock()
	day := stats.Day
	month := stats.Month
	statsMu.Unlock()
	out := map[string]any{
		"messages":  msgCount,
		"urges":     urgeTotal,
		"comments":  commentCount,
		"postcards": postCount,
		"likes":     likeTotal,
		"visits":    visits,
		"pets":      pets,
		"feeds":     feeds,
		"pigmis":    pigmiCount,
		"day":       day,
		"month":     month,
	}
	if topPost != nil && topPost.Likes > 0 {
		out["topPost"] = map[string]any{"nick": topPost.Nick, "likes": topPost.Likes}
	}
	if topNick != "" {
		out["topPigmi"] = map[string]any{"nick": topNick, "points": topPoints}
	}
	writeJSON(w, 200, out)
}

// GET /api/admin/export — 导出全部数据(zip)
func handleExport(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	files, _ := filepath.Glob("/var/lib/catcafe/*.json")
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		fw, _ := zw.Create(filepath.Base(f))
		fw.Write(data)
	}
	zw.Close()
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=catcafe-backup.zip")
	w.Write(buf.Bytes())
}

// GET /api/qr — 本站二维码
func handleQR(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	png, err := qrcode.Encode(scheme+"://"+r.Host+"/", qrcode.Medium, 220)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "生成失败"})
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(png)
}

// ---------- 机器管理员(agent)接口 ----------

// GET /api/agent/content — 读取全站内容
// PUT /api/agent/content — 整体更新内容(公告/小说/插画/资料等)
func handleAgentContent(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	switch r.Method {
	case http.MethodGet:
		// agent 返回完整内容(含真实群号)
		data, err := os.ReadFile(contentFile)
		if err != nil {
			data = []byte(defaultContent)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	case http.MethodPut:
		body := http.MaxBytesReader(w, r.Body, 1<<20)
		data, err := io.ReadAll(body)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "content too large"})
			return
		}
		var check map[string]any
		if err := json.Unmarshal(data, &check); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid json"})
			return
		}
		// 防止打码群号覆盖真实值(agent 读到的是打码版)
		if fc, ok := check["fanClub"].(map[string]any); ok && fc["qq"] == qqMask {
			var old map[string]any
			if oldData, err := os.ReadFile(contentFile); err == nil && json.Unmarshal(oldData, &old) == nil {
				if ofc, ok := old["fanClub"].(map[string]any); ok {
					fc["qq"] = ofc["qq"]
				}
			}
		}
		pretty, _ := json.MarshalIndent(check, "", "  ")
		tmp := contentFile + ".tmp"
		if err := os.WriteFile(tmp, pretty, 0644); err != nil {
			writeJSON(w, 500, map[string]string{"error": "save failed"})
			return
		}
		os.Rename(tmp, contentFile)
		log.Printf("agent 更新了站点内容")
		writeJSON(w, 200, map[string]string{"status": "saved"})
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
}

// GET /api/agent/stats — 运营数据(留言数/催更/访问量等,供 agent 决策)
func handleAgentStats(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	handleRecap(w, r)
}

// GET /api/admin/content-full — 店主后台读取完整内容(含真实群号)
func handleAdminContentFull(w http.ResponseWriter, r *http.Request) {
	if !requireAuth(w, r) {
		return
	}
	data, err := os.ReadFile(contentFile)
	if err != nil {
		data = []byte(defaultContent)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// POST /api/agent/messages/reply {time, nick, reply} — 以"猪咪君君(饲养员)"身份回复留言
func handleAgentReply(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Time  string `json:"time"`
		Nick  string `json:"nick"`
		Reply string `json:"reply"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "请求格式不对"})
		return
	}
	mu.Lock()
	found := false
	for i := range store.Messages {
		if store.Messages[i].Time == req.Time && store.Messages[i].Nick == req.Nick {
			store.Messages[i].Reply = clean(req.Reply, maxTextLen)
			store.Messages[i].ReplyBy = "猪咪君君"
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

// POST /api/agent/messages/delete {time, nick} — 删除不当留言(管理职责)
func handleAgentMsgDelete(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Time string `json:"time"`
		Nick string `json:"nick"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "请求格式不对"})
		return
	}
	mu.Lock()
	found := false
	for i, m := range store.Messages {
		if m.Time == req.Time && m.Nick == req.Nick {
			store.Messages = append(store.Messages[:i], store.Messages[i+1:]...)
			found = true
			break
		}
	}
	if found {
		saveStore()
		log.Printf("agent 删除了留言: %s %s", req.Nick, req.Time)
	}
	mu.Unlock()
	if !found {
		writeJSON(w, 404, map[string]string{"error": "留言不存在"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// ---------- 全店养猪 ----------

const pigFile = "/var/lib/catcafe/pig.json"

type PigState struct {
	XP      int     `json:"xp"`
	Fed     int     `json:"fed"`
	Pats    int     `json:"pats"`
	Hunger  float64 `json:"hunger"` // 0-100 饱食度
	Mood    float64 `json:"mood"`   // 0-100 心情
	Updated int64   `json:"updated"`
}

var (
	pigMu sync.Mutex
	pig   PigState
)

var pigStages = []struct {
	At   int
	Name string
}{
	{0, "猪崽"},
	{50, "小猪咪"},
	{150, "圆润猪咪"},
	{300, "猪王"},
}

func pigStage() (string, int, int) {
	stage := pigStages[0].Name
	next := -1
	for i, st := range pigStages {
		if pig.XP >= st.At {
			stage = st.Name
			if i+1 < len(pigStages) {
				next = pigStages[i+1].At
			}
		}
	}
	if pig.XP >= pigStages[len(pigStages)-1].At {
		next = 0 // 已满级
	}
	return stage, next, pig.XP
}

func loadPig() {
	data, err := os.ReadFile(pigFile)
	if err != nil {
		pig = PigState{Hunger: 80, Mood: 80, Updated: time.Now().Unix()}
		return
	}
	json.Unmarshal(data, &pig)
}

func savePig() {
	data, _ := json.MarshalIndent(pig, "", "  ")
	tmp := pigFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err == nil {
		os.Rename(tmp, pigFile)
	}
}

// 随时间掉饱食/心情(读取时惰性结算)
func pigDecay() {
	now := time.Now().Unix()
	if pig.Updated == 0 {
		pig.Updated = now
	}
	hours := float64(now-pig.Updated) / 3600
	if hours > 0 {
		pig.Hunger = math.Max(0, pig.Hunger-hours*4)
		pig.Mood = math.Max(0, pig.Mood-hours*3)
		pig.Updated = now
	}
}

var lastFeed = map[string]time.Time{} // 喂食冷却(5分钟/IP)
var lastPat = map[string]time.Time{}  // 摸头冷却(1分钟/IP)

func cooldownOK(m map[string]time.Time, r *http.Request, d time.Duration) bool {
	ip := clientIP(r)
	rateMu.Lock()
	defer rateMu.Unlock()
	if t, ok := m[ip]; ok && time.Since(t) < d {
		return false
	}
	m[ip] = time.Now()
	return true
}

// GET /api/pig — 猪状态
// POST /api/pig/feed — 喂食(+5经验,5分钟/次)
// POST /api/pig/pet — 摸头(+2经验,1分钟/次)
func handlePig(w http.ResponseWriter, r *http.Request) {
	pigMu.Lock()
	pigDecay()
	pigMu.Unlock()

	if r.Method == http.MethodGet {
		pigMu.Lock()
		stage, next, xp := pigStage()
		out := map[string]any{
			"stage": stage, "xp": xp, "next": next,
			"hunger": int(pig.Hunger), "mood": int(pig.Mood),
			"fed": pig.Fed, "pats": pig.Pats,
		}
		pigMu.Unlock()
		writeJSON(w, 200, out)
		return
	}
	writeJSON(w, 405, map[string]string{"error": "method not allowed"})
}

func handlePigFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	if !cooldownOK(lastFeed, r, 5*time.Minute) {
		writeJSON(w, 429, map[string]string{"error": "刚喂过啦,它还在嚼,5分钟后再来"})
		return
	}
	pigMu.Lock()
	pigDecay()
	oldStage, _, _ := pigStage()
	pig.XP += 5
	pig.Fed++
	pig.Hunger = math.Min(100, pig.Hunger+25)
	pig.Mood = math.Min(100, pig.Mood+10)
	newStage, _, _ := pigStage()
	savePig()
	out := map[string]any{"stage": newStage, "xp": pig.XP, "hunger": int(pig.Hunger), "mood": int(pig.Mood), "fed": pig.Fed, "pats": pig.Pats}
	if newStage != oldStage {
		out["evolved"] = newStage // 进化了!
	}
	pigMu.Unlock()
	writeJSON(w, 200, out)
}

func handlePigPet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	if !cooldownOK(lastPat, r, time.Minute) {
		writeJSON(w, 429, map[string]string{"error": "摸得太勤啦,1分钟后再来"})
		return
	}
	pigMu.Lock()
	pigDecay()
	pig.XP += 2
	pig.Pats++
	pig.Mood = math.Min(100, pig.Mood+15)
	stage, _, _ := pigStage()
	savePig()
	out := map[string]any{"stage": stage, "xp": pig.XP, "hunger": int(pig.Hunger), "mood": int(pig.Mood), "fed": pig.Fed, "pats": pig.Pats}
	pigMu.Unlock()
	writeJSON(w, 200, out)
}

// ---------- 入群答题验证 ----------

// 作品 → 主角(答案表,答对任意一个主角即算对)
var gateNovels = []struct {
	Title string
	Heros []string
}{
	{"我才不会做宠物呢", []string{"璃沢", "莫琳娜"}},
	{"白天是黑长直萝莉优等生,晚上被富家大小姐当做小宠物", []string{"洛小雪"}},
	{"诱受萝莉的百合修仙传奇", []string{"苏泱"}},
	{"战败魔法少女的涩涩艺术品博物馆", []string{"奈叶"}},
}

// GET /api/gate/questions — 返回作品名列表(不含答案)
func handleGateQuestions(w http.ResponseWriter, r *http.Request) {
	titles := make([]string, len(gateNovels))
	for i, n := range gateNovels {
		titles[i] = n.Title
	}
	writeJSON(w, 200, titles)
}

// POST /api/gate/verify {answer: "..."} — 答对任意作品名或主角名即放行
func handleGateVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
		return
	}
	var req struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	if rateLimited(r) {
		writeJSON(w, 429, map[string]string{"error": "试得太快啦,歇一会儿再试"})
		return
	}
	// 归一化:去书名号、空格、大小写
	norm := strings.ToLower(strings.NewReplacer("《", "", "》", "", " ", "", "　", "").Replace(strings.TrimSpace(req.Answer)))
	pass := false
	for _, n := range gateNovels {
		if norm == strings.ToLower(n.Title) {
			pass = true
			break
		}
		for _, h := range n.Heros {
			if norm == strings.ToLower(h) {
				pass = true
				break
			}
		}
	}
	if !pass {
		writeJSON(w, 200, map[string]any{"pass": false})
		return
	}
	var c map[string]any
	qq := ""
	if data, err := os.ReadFile(contentFile); err == nil && json.Unmarshal(data, &c) == nil {
		if fc, ok := c["fanClub"].(map[string]any); ok {
			qq, _ = fc["qq"].(string)
		}
	}
	writeJSON(w, 200, map[string]any{"pass": true, "qq": qq})
}

// ---------- 猪咪聚集地(图文墙) ----------

type HubPost struct {
	ID   string `json:"id"`
	Type string `json:"type"` // image | text
	Img  string `json:"img,omitempty"`
	Text string `json:"text,omitempty"`
	Time string `json:"time"`
	By   string `json:"by"` // agent=猪咪君君(饲养员)
}

var (
	hubMu sync.Mutex
	hub   []HubPost
)

func loadHub() {
	data, err := os.ReadFile(hubFile)
	if err == nil {
		json.Unmarshal(data, &hub)
	}
}

func saveHub() {
	data, _ := json.MarshalIndent(hub, "", "  ")
	tmp := hubFile + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err == nil {
		os.Rename(tmp, hubFile)
	}
}

// GET /api/hub — 公开:聚集地内容列表
func handleHub(w http.ResponseWriter, r *http.Request) {
	hubMu.Lock()
	defer hubMu.Unlock()
	if hub == nil {
		hub = []HubPost{}
	}
	writeJSON(w, 200, hub)
}

// POST /api/agent/hub/text {text} — agent 发文字
func handleAgentHubText(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	text := clean(req.Text, 500)
	if text == "" {
		writeJSON(w, 400, map[string]string{"error": "内容不能为空"})
		return
	}
	post := HubPost{
		ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Type: "text",
		Text: text,
		Time: time.Now().Format("2006-01-02 15:04"),
		By:   "agent",
	}
	hubMu.Lock()
	hub = append([]HubPost{post}, hub...)
	if len(hub) > 200 {
		hub = hub[:200]
	}
	saveHub()
	hubMu.Unlock()
	log.Printf("agent 在聚集地发了文字")
	writeJSON(w, 200, post)
}

// POST /api/agent/hub/image — agent 发图片(multipart file + note)
func handleAgentHubImage(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSON(w, 400, map[string]string{"error": "文件太大(上限10MB)"})
		return
	}
	f, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "没有收到文件"})
		return
	}
	defer f.Close()
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowExt["gallery"][ext] {
		writeJSON(w, 400, map[string]string{"error": "只支持 jpg/png/gif/webp"})
		return
	}
	os.MkdirAll(hubDir, 0755)
	name := fmt.Sprintf("%d-%s", time.Now().Unix(), safeName(header.Filename))
	data, _ := io.ReadAll(f)
	if err := saveImageCompressed(data, ext, filepath.Join(hubDir, name)); err != nil {
		writeJSON(w, 400, map[string]string{"error": "图片处理失败"})
		return
	}
	post := HubPost{
		ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Type: "image",
		Img:  name,
		Text: clean(r.FormValue("note"), 100),
		Time: time.Now().Format("2006-01-02 15:04"),
		By:   "agent",
	}
	hubMu.Lock()
	hub = append([]HubPost{post}, hub...)
	if len(hub) > 200 {
		hub = hub[:200]
	}
	saveHub()
	hubMu.Unlock()
	log.Printf("agent 在聚集地发了图片: %s", name)
	writeJSON(w, 200, post)
}

// POST /api/agent/messages/post {content} — 以猪咪君君身份在留言板发帖
func handleAgentMsgPost(w http.ResponseWriter, r *http.Request) {
	if !agentOK(r) {
		writeJSON(w, 401, map[string]string{"error": "invalid api key"})
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, map[string]string{"error": "参数不对"})
		return
	}
	content := clean(req.Content, maxTextLen)
	if content == "" {
		writeJSON(w, 400, map[string]string{"error": "内容不能为空"})
		return
	}
	m := Message{
		Nick:    "猪咪君君",
		Content: content,
		Time:    time.Now().Format("2006-01-02 15:04"),
		By:      "agent",
	}
	mu.Lock()
	store.Messages = append([]Message{m}, store.Messages...)
	if len(store.Messages) > maxMessages {
		store.Messages = store.Messages[:maxMessages]
	}
	saveStore()
	mu.Unlock()
	writeJSON(w, 200, m)
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
	loadTokens()
	loadAgentKey()
	loadPig()
	loadHub()

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
	mux.HandleFunc("/api/recap", handleRecap)
	mux.HandleFunc("/api/agent/content", handleAgentContent)
	mux.HandleFunc("/api/agent/stats", handleAgentStats)
	mux.HandleFunc("/api/admin/content-full", handleAdminContentFull)
	mux.HandleFunc("/api/hub", handleHub)
	mux.HandleFunc("/api/agent/hub/text", handleAgentHubText)
	mux.HandleFunc("/api/agent/hub/image", handleAgentHubImage)
	mux.HandleFunc("/api/agent/messages/post", handleAgentMsgPost)
	mux.HandleFunc("/api/agent/messages/reply", handleAgentReply)
	mux.HandleFunc("/api/agent/messages/delete", handleAgentMsgDelete)
	mux.HandleFunc("/api/gate/questions", handleGateQuestions)
	mux.HandleFunc("/api/gate/verify", handleGateVerify)
	mux.HandleFunc("/api/pig", handlePig)
	mux.HandleFunc("/api/pig/feed", handlePigFeed)
	mux.HandleFunc("/api/pig/pet", handlePigPet)
	mux.HandleFunc("/api/qr", handleQR)
	mux.HandleFunc("/api/admin/export", handleExport)
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

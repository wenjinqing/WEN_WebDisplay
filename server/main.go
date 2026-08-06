// 爱丽丝的小涩猫咖啡厅 —— API 服务
// 功能:留言板 / 催更墙 / 站点内容管理 / 店主后台登录与文件上传
// 零依赖,仅标准库;数据存 JSON 文件,适合低流量粉丝站
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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
	tokenMu  sync.Mutex
	tokens   = map[string]time.Time{} // token → 过期时间
	wallMu   sync.Mutex
	wall     []WallPost
)

type WallPost struct {
	Img  string `json:"img"`
	Nick string `json:"nick"`
	Note string `json:"note"`
	Time string `json:"time"`
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
		writeJSON(w, 200, map[string]any{"total": total, "urge": u})
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
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
	dst, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "服务器保存失败"})
		return
	}
	defer dst.Close()
	io.Copy(dst, f)
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
		dst, err := os.Create(filepath.Join(wallDir, name))
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": "服务器保存失败"})
			return
		}
		io.Copy(dst, f)
		dst.Close()

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
		writeJSON(w, 200, post)
	default:
		writeJSON(w, 405, map[string]string{"error": "method not allowed"})
	}
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

// ---------- 启动 ----------

func main() {
	for _, d := range []string{"/var/lib/catcafe", galleryDir, novelDir, wallDir} {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatal(err)
		}
	}
	loadStore()
	loadWall()

	// 首次启动:写入默认店主密码哈希
	if _, err := os.Stat(passFile); os.IsNotExist(err) {
		os.WriteFile(passFile, []byte(hashPass(defaultAdminPass)), 0600)
		log.Printf("!!! 已生成默认店主密码 %q,请登录后台后立即修改", defaultAdminPass)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/urge", handleUrge)
	mux.HandleFunc("/api/wall", handleWall)
	mux.HandleFunc("/api/content", handleContent)
	mux.HandleFunc("/api/admin/login", handleLogin)
	mux.HandleFunc("/api/admin/content", handlePutContent)
	mux.HandleFunc("/api/admin/password", handleChangePass)
	mux.HandleFunc("/api/admin/upload", handleUpload)
	mux.HandleFunc("/api/admin/delete-file", handleDeleteFile)
	mux.HandleFunc("/api/admin/wall/delete", handleWallDelete)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	log.Printf("猫咖 API 营业中 → %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

// 爱丽丝的小涩猫咖啡厅 —— 留言板 & 催更墙 API
// 零依赖,仅标准库;数据存 JSON 文件,适合低流量粉丝站
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	dataFile     = "/var/lib/catcafe/data.json"
	maxMessages  = 200            // 留言最多保留条数
	maxUrges     = 50             // 催更记录最多保留条数
	maxNickLen   = 20             // 昵称最大字符数
	maxTextLen   = 300            // 内容最大字符数
	minInterval  = 5 * time.Second // 同一 IP 两次提交的最小间隔
	listenAddr   = "127.0.0.1:9090"
)

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
	mu    sync.Mutex
	store Store
	// 简单的 IP 提交频率限制
	rateMu   sync.Mutex
	lastPost = map[string]time.Time{}
)

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
	// 经过 nginx 代理,取真实 IP
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

// GET /api/messages — 拉取留言(新的在前)
// POST /api/messages {nick, content} — 发表留言
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

// GET /api/urge — 返回 {total, recent}
// POST /api/urge {nick} — 催更一次
func handleUrge(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, 200, map[string]any{"total": store.UrgeTotal, "recent": store.Urges})
	case http.MethodPost:
		var u Urge
		json.NewDecoder(r.Body).Decode(&u) // body 可空
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

func main() {
	if err := os.MkdirAll("/var/lib/catcafe", 0755); err != nil {
		log.Fatal(err)
	}
	loadStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/messages", handleMessages)
	mux.HandleFunc("/api/urge", handleUrge)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	log.Printf("猫咖 API 营业中 → %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

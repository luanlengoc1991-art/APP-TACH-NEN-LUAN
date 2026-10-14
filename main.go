// LUAN GROUP Layer Studio: Windows portable desktop launcher.
// Opens a local-only HTTP server in Edge/Chrome app mode. No credentials, uploads, or cloud APIs.
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

//go:embed web/*
var webFiles embed.FS
var lastSeen atomic.Int64

func main() {
	noLaunch := flag.Bool("no-browser", false, "start local-only server but do not launch browser (tests)")
	flag.Parse()
	sub, err := fs.Sub(webFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	// Stable origin is necessary: browser IndexedDB is scoped by host AND port.
	// Fallback to an ephemeral port only if another process occupies our port.
	listener, err := net.Listen("tcp", "127.0.0.1:49753")
	if err != nil {
		log.Println("Cổng 49753 đang bận; dữ liệu phiên lưu có thể khác phiên trước:", err)
		listener, err = net.Listen("tcp", "127.0.0.1:0")
	}
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d/", port)
	lastSeen.Store(time.Now().Unix())
	handler := http.NewServeMux()
	handler.HandleFunc("/open-chatgpt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		origin, err := neturl.Parse(r.Header.Get("Origin"))
		// Reject any cross-site invocation that tries to launch a browser on the machine.
		if err != nil || origin.Scheme != "http" || origin.Host != r.Host || r.Header.Get("Sec-Fetch-Site") != "same-origin" {
			http.Error(w, "Forbidden", 403)
			return
		}
		if err := openAppWindow("https://chatgpt.com/"); err != nil {
			http.Error(w, "Cannot open browser: "+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	handler.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Luan-Application", "layer-studio-1.0")
		_, _ = w.Write([]byte("LUAN LAYER STUDIO v1"))
	})
	handler.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		lastSeen.Store(time.Now().Unix())
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(204)
	})
	fileServer := http.FileServer(http.FS(sub))
	handler.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' blob: data:; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; script-src 'self'; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		// Only local static assets; a foreign website cannot read the user's image data.
		lastSeen.Store(time.Now().Unix())
		fileServer.ServeHTTP(w, r)
	}))
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Println("server:", err)
		}
	}()
	fmt.Println("LUAN GROUP AI Layer Studio running locally at:", url)
	if !*noLaunch {
		if err := openAppWindow(url); err != nil {
			log.Println("Không thể tự mở cửa sổ:", err, " — Hãy mở:", url)
		}
	}
	if *noLaunch {
		select {}
	}
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if time.Now().Unix()-lastSeen.Load() > 180 {
			log.Println("Không còn cửa sổ ứng dụng — dừng máy chủ cục bộ.")
			break
		}
	}
}

func openAppWindow(url string) error {
	if runtime.GOOS == "windows" {
		paths := []string{
			"msedge.exe", "chrome.exe",
			`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
			filepath.Join(os.Getenv("LOCALAPPDATA"), `Microsoft\Edge\Application\msedge.exe`),
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			filepath.Join(os.Getenv("LOCALAPPDATA"), `Google\Chrome\Application\chrome.exe`),
		}
		for _, p := range paths {
			if strings.TrimSpace(p) == "" {
				continue
			}
			if _, err := exec.LookPath(p); err != nil {
				continue
			}
			cmd := exec.Command(p, "--app="+url, "--window-size=1510,940", "--no-first-run")
			if err := cmd.Start(); err == nil {
				return nil
			}
		}
		// Fallback: Windows default browser (still only local traffic).
		return exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url).Start()
	}
	if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Start()
	}
	for _, b := range []string{"chromium", "google-chrome", "xdg-open"} {
		if _, err := exec.LookPath(b); err == nil {
			if b == "xdg-open" {
				return exec.Command(b, url).Start()
			}
			return exec.Command(b, "--app="+url).Start()
		}
	}
	return fmt.Errorf("không tìm thấy trình duyệt")
}

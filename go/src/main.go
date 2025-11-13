package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// ロガーの初期化（Go 1.21+ の slog パッケージ）
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	slog.Info("Go サンプルアプリケーション開始")

	// HTTP サーバーのセットアップ
	mux := http.NewServeMux()

	// ハンドラの登録
	mux.HandleFunc("GET /", handleRoot)
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /hello/{name}", handleHello)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバー起動をゴルーチンで実行
	go func() {
		slog.Info("HTTP サーバー起動", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("サーバーエラー", "error", err)
		}
	}()

	// グレースフルシャットダウンの処理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	slog.Info("シャットダウンシグナル受信")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("シャットダウンエラー", "error", err)
	}

	slog.Info("アプリケーション終了")
}

// handleRoot はルートパスのハンドラ
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World! 🚀\n"))
	slog.DebugContext(r.Context(), "ルートアクセス", "method", r.Method, "path", r.URL.Path)
}

// handleHealth はヘルスチェックのハンドラ
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// handleHello は名前付きパラメータを受け取るハンドラ
func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, " + name + "!\n"))
	slog.InfoContext(r.Context(), "Hello ハンドラ実行", "name", name)
}

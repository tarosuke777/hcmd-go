package server

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 各カテゴリごとに「賢いハンドラー」を生成
	// これにより、同じURLパスで GET/POST 両方対応できる
	mux.Handle("/books/", NewSmartHandler("books", true, 10))
	mux.Handle("/images/", NewSmartHandler("images", true, 10))
	mux.Handle("/videos/", NewSmartHandler("videos", true, 100))
	mux.Handle("/backups/", NewSmartHandler("backups", false, 500))
	return mux
}

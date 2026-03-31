package server

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 各カテゴリごとに「賢いハンドラー」を生成
	// これにより、同じURLパスで GET/POST 両方対応できる
	mux.Handle("/books/", NewSmartHandler("books"))
	mux.Handle("/images/", NewSmartHandler("images"))
	mux.Handle("/videos/", NewSmartHandler("videos"))
	return mux
}

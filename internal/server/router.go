package server

import (
	"home/internal/handler" // プロジェクト名に合わせて変更してください
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 各カテゴリごとに「賢いハンドラー」を生成
	// これにより、同じURLパスで GET/POST 両方対応できる
	mux.Handle("/books/", handler.NewSmartHandler("books"))
	mux.Handle("/images/", handler.NewSmartHandler("images"))
	mux.Handle("/videos/", handler.NewSmartHandler("videos"))
	return mux
}

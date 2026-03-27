package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SmartHandler struct {
	SaveDir string
	Prefix  string
}

func NewSmartHandler(dirName string) *SmartHandler {
	// 実行ファイル基準の絶対パスを作成
	exePath, _ := os.Executable()
	absPath := filepath.Join(filepath.Dir(exePath), dirName)
	_ = os.MkdirAll(absPath, os.ModePerm)

	prefix := fmt.Sprintf("/%s/", dirName)

	return &SmartHandler{SaveDir: absPath, Prefix: prefix}
}

func (h *SmartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GETリクエスト：標準のファイルサーバーとして振る舞う
		// http.StripPrefixでURLの /images/ などを削ってディレクトリを参照
		fs := http.FileServer(http.Dir(h.SaveDir))

		strippedPath := r.URL.Path[len(h.Prefix)-1:]
		log.Printf("[DEBUG] GET Request: Prefix=%s, OriginalPath=%s, Looking for file: %s in %s",
			h.Prefix, r.URL.Path, strippedPath, h.SaveDir)

		http.StripPrefix(h.Prefix, fs).ServeHTTP(w, r)

	case http.MethodPost:
		// POSTリクエスト：アップロード処理
		h.handleUpload(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ServeHTTP を実装することで http.Handler インターフェースを満たします
func (h *SmartHandler) handleUpload(w http.ResponseWriter, r *http.Request) {

	// 5MB制限
	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	file, header, err := r.FormFile("myImage")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 1. 現在時刻を取得
	now := time.Now()

	// 2. 時刻を文字列にフォーマット (例: 20240327-153005-123)
	// Go独自のレイアウト指定方法（2006, 01, 02...）を使います
	timestamp := now.Format("20060102-150405-000")

	// 3. 元の拡張子を取得
	ext := filepath.Ext(header.Filename)

	// 4. 新しいファイル名を組み立てる
	// 同時実行が心配な場合は、最後に少しランダムな文字を足すか、元の名前を繋げます
	newFileName := fmt.Sprintf("%s%s", timestamp, ext)

	// 保存パスの作成
	dstPath := filepath.Join(h.SaveDir, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to copy", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded: %s", header.Filename)
}

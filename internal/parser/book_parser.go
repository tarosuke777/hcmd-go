package parser

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BookInfo は本全体の情報
type BookInfo struct {
	Title      string
	Author     *string
	TotalPages int
	Pages      []PageInfo
}

// PageInfo は1ページ分の情報
type PageInfo struct {
	PageNumber int
	FilePath   string
}

// ParseBookDirectories はカレントディレクトリの各フォルダを「1冊の本」として解析します
func CurrentFolderToBook() (*BookInfo, error) {
	// 1. 現在の絶対パスを取得
	absPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 2. パスの最後（フォルダ名）をタイトルにする
	title := filepath.Base(absPath)

	// 3. フォルダ内のファイルを取得
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var pages []PageInfo
	pageCount := 1

	// ソート（01.jpg, 02.JPG を正しく並べるため）
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			pages = append(pages, PageInfo{
				PageNumber: pageCount,
				FilePath:   entry.Name(), // フォルダ内で実行するのでファイル名だけでOK
			})
			pageCount++
		}
	}

	return &BookInfo{
		Title:      title,
		Author:     nil,
		TotalPages: len(pages),
		Pages:      pages,
	}, nil
}
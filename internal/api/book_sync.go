package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home/internal/parser"
	"net/http"
	"time"
)

// 本の1ページ分のデータ構造
type BookPageRequest struct {
    PageNumber int    `json:"page_number"`
    FilePath   string `json:"file_path"`
}

// 本の一括登録用データ構造
type BookStoreRequest struct {
    Title  string            `json:"title"`
    Author *string            `json:"author"`
    Pages  []BookPageRequest `json:"pages"`
}

const apiBookStoreURL = "http://192.168.10.10/hv/api/books/store"

// SyncBooksToAPI 本のデータを一括で送信する
func SyncBooksToAPI(info *parser.BookInfo) error {
// 1. 送信用の構造体 (Request型) に詰め替える
    var pageRequests []BookPageRequest
    for _, p := range info.Pages {
        pageRequests = append(pageRequests, BookPageRequest{
            PageNumber: p.PageNumber,
            FilePath:   p.FilePath,
        })
    }

    requestData := BookStoreRequest{
        Title:  info.Title,
        Author: info.Author,
        Pages:  pageRequests,
    }

    jsonData, err := json.Marshal(requestData)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", apiBookStoreURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json") // Laravel側にAPIだと伝えるために重要！

    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("book sync failed with status: %d", resp.StatusCode)
    }

    fmt.Printf("Book '%s' 同期成功\n", requestData.Title)
    return nil
}
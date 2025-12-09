package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// APIに送信するデータ構造体
type VideoRequest struct {
	Title     string `json:"title"`
	Name      string `json:"name"`
	FileName  string `json:"file_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

const (
	targetDir = "./"
	apiURL    = "http://192.168.10.10/api/videos/store"
)

var re = regexp.MustCompile(`^(.+?)\s(\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2})\.(mp4|mov|avi|webm)$`)

// フォルダを走査し、各ファイル情報をAPIにPOST送信する
func SyncVideosToAPI() {
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		matches := re.FindStringSubmatch(fileName)

		if len(matches) == 0 {
			return nil
		}

		title := strings.TrimSpace(matches[1])
		dbDateTime := formatToSQLDateTime(matches[2])

		// JSON用の構造体を作成
		payload := VideoRequest{
			Title:     title,
			Name:      "", // 必要に応じて設定
			FileName:  fileName,
			CreatedAt: dbDateTime,
			UpdatedAt: dbDateTime,
		}

		// API呼び出しの実行
		err = sendToAPI(payload)
		if err != nil {
			log.Printf("API送信エラー (%s): %v", fileName, err)
		} else {
			fmt.Printf("送信成功: %s\n", fileName)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("フォルダ走査エラー: %v", err)
	}
	fmt.Printf("--- ログ: 全ての処理が完了しました。---\n")
}

// HTTP POSTでJSONを送信するヘルパー関数
func sendToAPI(data VideoRequest) error {
	// 構造体をJSONに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")
	// 認証が必要な場合はここに追加
	// req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

	// リクエストの送信
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// ステータスコードのチェック
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("APIがエラーを返しました: %s", resp.Status)
	}

	return nil
}

// (既存のヘルパー関数群)
func formatToSQLDateTime(rawDate string) string {
	parts := strings.Split(rawDate, " ")
	if len(parts) != 2 {
		return rawDate
	}
	timePart := strings.ReplaceAll(parts[1], "-", ":")
	return parts[0] + " " + timePart
}
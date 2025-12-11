package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home/internal/parser"
	"log"
	"net/http"
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

const apiURL = "http://192.168.10.10/hv/api/videos/store"

// SyncVideosToAPI はフォルダを走査し、各ファイル情報をAPIにPOST送信します。
func SyncVideosToAPI() {
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	// WalkAndParse 関数に処理ロジックを渡す
	err := parser.WalkAndParse(parser.TargetDir, func(info parser.VideoInfo) error {
		// JSON用の構造体を作成
		payload := VideoRequest{
			Title:     info.Title,
			Name:      "", // 必要に応じて設定
			FileName:  info.FileName,
			CreatedAt: info.DBDateTime,
			UpdatedAt: info.DBDateTime,
		}

		// API呼び出しの実行
		err := sendToAPI(payload)
		if err != nil {
			log.Printf("API送信エラー (%s): %v", info.FileName, err)
		} else {
			fmt.Printf("送信成功: %s\n", info.FileName)
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
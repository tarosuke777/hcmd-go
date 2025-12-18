package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"home/internal/parser"
	"io"
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

type MaxTimestampResponse struct {
    MaxCreatedAt string `json:"max_created_at"`
}

const (
    apiStoreURL = "http://192.168.10.10/hv/api/videos/store"
    apiMaxTimeURL = "http://192.168.10.10/hv/api/videos/max-timestamp"
    timeLayout  = "2006-01-02 15:04:05" // Laravelのデフォルトフォーマットに合わせる
)

// SyncVideosToAPI はフォルダを走査し、各ファイル情報をAPIにPOST送信します。
func SyncVideosToAPI() {
	fmt.Printf("--- ログ: 最新タイムスタンプの取得を開始します ---\n")

    // 1. 最新の時刻をAPIから取得
    maxTime, fetchApiErr := fetchMaxTimestamp()
    if fetchApiErr != nil {
        log.Printf("最新時刻の取得に失敗したため、全件送信を試みます: %v", fetchApiErr)
        // 失敗時に中断するか、古い時刻(1970年など)にするか選べます
        maxTime = time.Unix(0, 0)
    }

    fmt.Printf("最新の同期済み時刻: %v\n", maxTime.Format(timeLayout))
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	// WalkAndParse 関数に処理ロジックを渡す
	err := parser.WalkAndParse(parser.TargetDir, func(info parser.VideoInfo) error {

		// info.DBDateTime を time.Time に変換
        fileTime, parseErr := time.Parse(timeLayout, info.DBDateTime)
        if parseErr != nil {
            log.Printf("時刻パースエラー (%s): %v", info.FileName, parseErr)
            return nil 
        }

		// 3. 比較: ファイル時刻 > DB最新時刻 の場合のみ送信
        if fileTime.After(maxTime) {

			// JSON用の構造体を作成
			payload := VideoRequest{
				Title:     info.Title,
				Name:      "", // 必要に応じて設定
				FileName:  info.FileName,
				CreatedAt: info.DBDateTime,
				UpdatedAt: info.DBDateTime,
			}

			// API呼び出しの実行
			apiErr := sendToAPI(payload)
			if apiErr != nil {
				log.Printf("API送信エラー (%s): %v", info.FileName, apiErr)
			} else {
				fmt.Printf("送信成功: %s\n", info.FileName)
			}
		} else {
			fmt.Printf("スキップ: 既に同期済み (%s)\n", info.FileName)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("フォルダ走査エラー: %v", err)
	}
	fmt.Printf("--- ログ: 全ての処理が完了しました。---\n")
}

// Laravelから最新の時刻を取得する関数
func fetchMaxTimestamp() (time.Time, error) {
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(apiMaxTimeURL)
    if err != nil {
        return time.Time{}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return time.Time{}, fmt.Errorf("status error: %d", resp.StatusCode)
    }

    var result MaxTimestampResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return time.Time{}, err
    }

    // 文字列を time.Time に変換
    if result.MaxCreatedAt == "" {
        return time.Unix(0, 0), nil // データが空なら最小値を返す
    }

    return time.Parse(timeLayout, result.MaxCreatedAt)
}

// HTTP POSTでJSONを送信するヘルパー関数
func sendToAPI(data VideoRequest) error {
	// 構造体をJSONに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", apiStoreURL, bytes.NewBuffer(jsonData))
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

		bodyBytes, readErr := io.ReadAll(resp.Body)
        if readErr != nil {
            // ボディの読み取り自体に失敗した場合
            return fmt.Errorf("APIがエラーを返しました: %s, ボディの読み取りに失敗: %v", resp.Status, readErr)
        }

		bodyString := string(bodyBytes)

		return fmt.Errorf("APIがエラーを返しました: %s\nエラー詳細JSON:\n%s", resp.Status, bodyString)
	}

	return nil
}
package video

import (
	"fmt"
	"home/internal/hvapi"
	"log"
	"time"
)

// APIに送信するデータ構造体
type VideoRequest struct {
	Title     string        `json:"title"`
	Name      string        `json:"name"`
	FileName  string        `json:"file_name"`
	CreatedAt hvapi.APITime `json:"created_at"`
	UpdatedAt hvapi.APITime `json:"updated_at"`
}

type MaxTimestampResponse struct {
	MaxCreatedAt string `json:"max_created_at"`
}

const (
	apiStoreURL   = "http://192.168.10.10/hv/api/videos/store"
	apiMaxTimeURL = "http://192.168.10.10/hv/api/videos/max-timestamp"
)

// SyncVideosToAPI はフォルダを走査し、各ファイル情報をAPIにPOST送信します。
func SyncVideosToAPI() {
	fmt.Printf("--- ログ: 最新タイムスタンプの取得を開始します ---\n")

	// 1. 最新の時刻をAPIから取得
	maxTime, err := hvapi.FetchMaxTimestamp(apiMaxTimeURL)
	if err != nil {
		log.Printf("最新時刻の取得に失敗したため、全件送信を試みます: %v", err)
		// 失敗時に中断するか、古い時刻(1970年など)にするか選べます
		maxTime = time.Unix(0, 0)
	}

	fmt.Printf("最新の同期済み時刻: %v\n", hvapi.APITime(maxTime))
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	// WalkAndParse 関数に処理ロジックを渡す
	WalkAndParse("./", func(info VideoInfo) error {

		// 3. 比較: ファイル時刻 > DB最新時刻 の場合のみ送信
		if info.FileTime.After(maxTime) {

			hvApiTime := hvapi.APITime(info.FileTime)
			// JSON用の構造体を作成
			payload := VideoRequest{
				Title:     info.Title,
				Name:      "", // 必要に応じて設定
				FileName:  info.FileName,
				CreatedAt: hvApiTime,
				UpdatedAt: hvApiTime,
			}

			// API呼び出しの実行
			apiErr := hvapi.SendToAPI(apiStoreURL, payload)
			if apiErr != nil {
				log.Printf("API送信エラー (%s): %v", info.FileName, apiErr)
			} else {
				fmt.Printf("送信成功: %s\n", info.FileName)
			}
		}
		return nil
	})

	fmt.Printf("--- ログ: 全ての処理が完了しました。---\n")
}

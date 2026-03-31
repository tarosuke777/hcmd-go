package image

import (
	"fmt"
	"home/internal/hvapi"
	"log"
	"time"
)

// APIに送信するデータ構造体
type ImageRequest struct {
	FileName  string `json:"file_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MaxTimestampResponse struct {
	MaxCreatedAt string `json:"max_created_at"`
}

const (
	apiStoreURL   = "http://192.168.10.10/hv/api/images/store"
	apiMaxTimeURL = "http://192.168.10.10/hv/api/images/max-timestamp"
	timeLayout    = "2006-01-02 15:04:05" // Laravelのデフォルトフォーマットに合わせる
)

// SyncImagesToAPI はフォルダを走査し、各ファイル情報をAPIにPOST送信します。
func SyncImagesToAPI() {
	fmt.Printf("--- ログ: 最新タイムスタンプの取得を開始します ---\n")

	// 1. 最新の時刻をAPIから取得
	maxTime, err := hvapi.FetchMaxTimestamp(apiMaxTimeURL)
	if err != nil {
		log.Printf("最新時刻の取得に失敗したため、全件送信を試みます: %v", err)
		// 失敗時に中断するか、古い時刻(1970年など)にするか選べます
		maxTime = time.Unix(0, 0)
	}

	fmt.Printf("最新の同期済み時刻: %v\n", maxTime.Format(timeLayout))
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	// WalkAndParse 関数に処理ロジックを渡す
	WalkAndParse("./", func(info ImageInfo) error {
		fileTime, _ := time.Parse(hvapi.TimeLayout, info.DBDateTime)

		if fileTime.After(maxTime) {
			payload := ImageRequest{
				FileName:  info.FileName,
				CreatedAt: info.DBDateTime,
				UpdatedAt: info.DBDateTime,
			}

			// ジェネリクスにより ImageRequest 型として送信
			if err := hvapi.SendToAPI(apiStoreURL, payload); err != nil {
				log.Printf("API送信エラー (%s): %v", info.FileName, err)
			} else {
				fmt.Printf("送信成功: %s\n", info.FileName)
			}
		}
		return nil
	})
	fmt.Printf("--- ログ: 全ての処理が完了しました。---\n")
}

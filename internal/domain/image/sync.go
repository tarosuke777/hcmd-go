package image

import (
	"fmt"
	"home/internal/hvapi"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	uploadURL     = "http://192.168.10.11:8080/images/" // 画像アップロード用のURL (handler.goのエンドポイント)
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

func UploadAndSyncImages() {

	fmt.Printf("--- ログ: アップロードとAPI登録を開始します ---\n")

	files, err := os.ReadDir("./")
	if err != nil {
		log.Fatalf("ディレクトリの読み込みに失敗: %v", err)
	}

	fmt.Printf("--- ログ: 一括アップロードとAPI登録を開始します ---\n")

	for _, file := range files {
		// ディレクトリはスキップ
		if file.IsDir() {
			continue
		}

		// 拡張子チェック (画像・動画と思われるもの)
		ext := strings.ToLower(filepath.Ext(file.Name()))
		switch ext {
		case ".jpg", ".jpeg", ".png", ".heic":
			// 処理を継続
		default:
			continue
		}

		fileName := file.Name()

		// 2. Webサーバーへアップロードして新しいファイル名を取得
		newFileName, err := hvapi.UploadFile(uploadURL, fileName)
		if err != nil {
			log.Printf("アップロード失敗 (%s): %v", fileName, err)
			continue
		}

		matches := re.FindStringSubmatch(newFileName)
		if len(matches) < 3 {
			log.Printf("サーバーからのレスポンス形式が不正です (%s)", newFileName)
			continue
		}

		// matches[1]: 日付(YYYYMMDD), matches[2]: 時間(HHmmss)
		datePart, timePart := matches[1], matches[2]
		dbTime := fmt.Sprintf("%s-%s-%s %s:%s:%s",
			datePart[0:4], datePart[4:6], datePart[6:8],
			timePart[0:2], timePart[2:4], timePart[4:6])

		// 5. HV API(10.10)へ送信
		payload := ImageRequest{
			FileName:  newFileName,
			CreatedAt: dbTime,
			UpdatedAt: dbTime,
		}

		// 5. HV APIへDB登録
		if err := hvapi.SendToAPI(apiStoreURL, payload); err != nil {
			log.Printf("API登録エラー (%s): %v", newFileName, err)
		} else {
			fmt.Printf("成功: %s -> %s (DB日時: %s)\n", fileName, newFileName, dbTime)

			// --- 追加: ファイルの削除処理 ---
			if err := os.Remove(fileName); err != nil {
				log.Printf("ファイル削除失敗 (%s): %v", fileName, err)
			} else {
				fmt.Printf("削除完了: %s\n", fileName)
			}
		}
	}
	fmt.Printf("--- ログ: 全ての処理が完了しました ---\n")
}

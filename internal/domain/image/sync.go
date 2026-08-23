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
	FileName  string        `json:"file_name"`
	CreatedAt hvapi.APITime `json:"created_at"`
	UpdatedAt hvapi.APITime `json:"updated_at"`
}

type MaxTimestampResponse struct {
	MaxCreatedAt string `json:"max_created_at"`
}

const (
	uploadURL     = "http://192.168.10.11:8080/images/" // 画像アップロード用のURL (handler.goのエンドポイント)
	apiStoreURL   = "https://hv.home.arpa/hv/api/images/store"
	apiMaxTimeURL = "https://hv.home.arpa/hv/api/images/max-timestamp"
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

	fmt.Printf("最新の同期済み時刻: %v\n", hvapi.APITime(maxTime))
	fmt.Printf("--- ログ: APIへのデータ送信を開始します ---\n")

	// WalkAndParse 関数に処理ロジックを渡す
	WalkAndParse("./", func(info ImageInfo) error {

		if info.FileTime.After(maxTime) {
			hvApiTime := hvapi.APITime(info.FileTime)
			// JSON用の構造体を作成
			payload := ImageRequest{
				FileName:  info.FileName,
				CreatedAt: hvApiTime,
				UpdatedAt: hvApiTime,
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

		// matches[1]: "20260328", matches[2]: "130702"
		// これを繋げて "20260328-130702" という形でパースする
		rawDateTime := matches[1] + "-" + matches[2]

		// 文字列の並び順通りのレイアウトを指定
		// 20060102 (年月日) - 150405 (時分秒)
		fileTime, err := time.Parse("20060102-150405", rawDateTime)
		if err != nil {
			log.Printf("日時変換エラー (%s): %v", matches[0], err)
			fileTime = time.Time{} // エラー時はゼロ値
		}

		hvApiTime := hvapi.APITime(fileTime)

		// 5. HV API(10.10)へ送信
		payload := ImageRequest{
			FileName:  newFileName,
			CreatedAt: hvApiTime,
			UpdatedAt: hvApiTime,
		}

		// 5. HV APIへDB登録
		if err := hvapi.SendToAPI(apiStoreURL, payload); err != nil {
			log.Printf("API登録エラー (%s): %v", newFileName, err)
		} else {
			fmt.Printf("成功: %s -> %s (DB日時: %s)\n", fileName, newFileName, hvApiTime)

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

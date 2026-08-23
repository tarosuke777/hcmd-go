package backup

import (
	"fmt"
	"home/internal/hvapi"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	uploadURL = "http://192.168.10.11:8080/backups/" // バックアップデータのアップロード用のURL (handler.goのエンドポイント)
	oldDir    = "./old"                              // 移動先のディレクトリ
)

func UploadAndSyncBackups() {

	fmt.Printf("--- ログ: アップロードを開始します ---\n")

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

		// 拡張子チェック (SQLファイルと圧縮ファイル)
		ext := strings.ToLower(filepath.Ext(file.Name()))
		switch ext {
		case ".sql", ".gz":
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

		log.Printf("アップロードファイル：%s", newFileName)

		// 3. アップロード成功後、ファイルを old ディレクトリへ移動
		dstPath := filepath.Join(oldDir, fileName)
		if err := os.Rename(fileName, dstPath); err != nil {
			log.Printf("ファイルの移動に失敗 (%s -> %s): %v", fileName, dstPath, err)
			continue
		}

		log.Printf("ファイルを移動しました: %s", dstPath)

	}
	fmt.Printf("--- ログ: 全ての処理が完了しました ---\n")
}

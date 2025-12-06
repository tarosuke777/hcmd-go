package sqlexport // main.goと同じパッケージ名

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// --- SQL生成処理用 定数と変数 ---
const targetDir = "./"
const outputFile = "insert_videos.sql"
const insertTemplate = `INSERT INTO videos (title, name, file_name, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s');`
var re = regexp.MustCompile(`^(.+?)\s(\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2})\.(mp4|mov|avi|webm)$`)

// フォルダを走査し、SQL INSERT文を生成してファイルに出力する関数
// この関数を main.go から呼び出します。
func GenerateInsertSQLs() {
	// ... (元の generateInsertSQLs の中身をそのまま記述) ...
    
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("出力ファイル %s のオープンに失敗しました: %v", outputFile, err)
	}
	defer f.Close()

	fmt.Printf("--- ログ: INSERT文をファイルと標準出力に出力します ---\n")
	fmt.Printf("出力ファイル: %s\n", outputFile)

	// ... 省略 ... フォルダ走査とSQL生成のロジック ...
    
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		os.Mkdir(targetDir, 0755)
		createDummyFile(targetDir, "別の動画タイトル１ 2025-09-13 18-06-25.mp4")
		createDummyFile(targetDir, "別の動画タイトル２ 2025-10-01 10-30-00.mp4")
	}

	err = filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
        // ... ファイル走査ロジック ...
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		matches := re.FindStringSubmatch(fileName)
		
		if len(matches) == 0 {
			fmt.Printf("-- SKIP: ファイル名形式が一致しません: %s\n", fileName)
			return nil
		}

		title := strings.TrimSpace(matches[1]) 
		name := "" 
		file_name := fileName
		
		now := time.Now().Format("2006-01-02 15:04:05")

		insertSQL := fmt.Sprintf(
			insertTemplate,
			escapeSQL(title),
			escapeSQL(name),
			escapeSQL(file_name),
			now,
			now,
		)

		fmt.Println(insertSQL)
		
		_, err = f.WriteString(insertSQL + "\n")
		if err != nil {
			log.Printf("ファイル書き込みエラー: %v", err)
		}
		
		return nil
	})

	if err != nil {
		log.Fatalf("フォルダの走査中にエラーが発生しました: %v", err)
	}
	fmt.Printf("--- ログ: SQL生成処理が完了しました。---\n")
}

// シングルクォートをエスケープするヘルパー関数
func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// ダミーファイルを作成するヘルパー関数
func createDummyFile(dir, fileName string) {
	filePath := filepath.Join(dir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			log.Printf("ダミーファイル作成エラー: %v", err)
			return
		}
		f.Close()
	}
}
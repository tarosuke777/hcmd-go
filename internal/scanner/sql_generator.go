package scanner

import (
	"fmt"
	"home/internal/parser"
	"log"
	"os"
)

// --- SQL生成処理用 定数と変数 ---
const OutputFile = "insert_videos.sql" // ここに移動
const insertTemplate = `INSERT INTO videos (title, name, file_name, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s');`

// GenerateInsertSQLs はフォルダを走査し、SQL INSERT文を生成してファイルに出力する関数です。
func GenerateInsertSQLs() {
	f, err := os.OpenFile(OutputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("出力ファイル %s のオープンに失敗しました: %v", OutputFile, err)
	}
	defer f.Close()

	fmt.Printf("--- ログ: INSERT文をファイルと標準出力に出力します ---\n")
	fmt.Printf("出力ファイル: %s\n", OutputFile)

	// WalkAndParse 関数に処理ロジックを渡す
	err = parser.WalkAndParse(parser.TargetDir, func(info parser.VideoInfo) error {
		insertSQL := fmt.Sprintf(
			insertTemplate,
			parser.EscapeSQL(info.Title),
			"", // name (空)
			parser.EscapeSQL(info.FileName),
			info.DBDateTime, // created_at
			info.DBDateTime, // updated_at
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
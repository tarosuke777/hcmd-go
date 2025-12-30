package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// VideoInfo は解析された動画ファイルの情報を含む構造体です。
type VideoInfo struct {
	Title      string
	RawDate    string // 元の日時文字列 (例: 2025-09-13 18-06-25)
	FileName   string
	DBDateTime string // SQL/API向けの日時文字列 (例: 2025-09-13 18:06:25)
}

// --- 定数 ---
const TargetDir = "./"

// ファイル名解析用の正規表現
// 例: "別の動画タイトル１ 2025-09-13 18-06-25.mp4"
var re = regexp.MustCompile(`^(.+?)\s(\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2})\.(mp4|mov|avi|webm)$`)

// --- 公開関数 ---

// WalkAndParse は指定されたフォルダを走査し、正規表現に一致するファイル情報を処理関数に渡します。
func WalkAndParse(targetDir string, processor func(info VideoInfo) error) error {
	return filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
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

		// matches[1]: タイトル
		// matches[2]: 日時 (例: 2025-09-13 18-06-25)
		title := strings.TrimSpace(matches[1])
		rawDate := matches[2]
		dbDateTime := formatToSQLDateTime(rawDate)

		videoInfo := VideoInfo{
			Title:      title,
			RawDate:    rawDate,
			FileName:   fileName,
			DBDateTime: dbDateTime,
		}

		// 処理関数を実行
		return processor(videoInfo)
	})
}
// EscapeSQL はシングルクォートをエスケープするヘルパー関数です。
func EscapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
// --- 非公開ヘルパー関数 ---

// formatToSQLDateTime はファイル名の日時形式(HH-mm-ss)をSQL形式(HH:mm:ss)に変換します。
func formatToSQLDateTime(rawDate string) string {
	// スペースで日付と時間を分ける
	parts := strings.Split(rawDate, " ")
	if len(parts) != 2 {
		return rawDate
	}
	// 時間部分のハイフンをコロンに置換
	timePart := strings.ReplaceAll(parts[1], "-", ":")
	return parts[0] + " " + timePart
}
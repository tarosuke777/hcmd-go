package video

import (
	"home/internal/hvapi"
	"regexp"
	"strings"
)

// VideoInfo は解析された動画ファイルの情報を含む構造体です。
type VideoInfo struct {
	Title      string
	FileName   string
	DBDateTime string // SQL/API向けの日時文字列 (例: 2025-09-13 18:06:25)
}

// ファイル名解析用の正規表現
// 例: "別の動画タイトル１ 2025-09-13 18-06-25.mp4"
var re = regexp.MustCompile(`^(.+?)\s(\d{4}-\d{2}-\d{2}\s\d{2}-\d{2}-\d{2})\.(mp4|mov|avi|webm)$`)

// --- 公開関数 ---

// WalkAndParse は指定されたフォルダを走査し、正規表現に一致するファイル情報を処理関数に渡します。
func WalkAndParse(targetDir string, processor func(info VideoInfo) error) error {
	return hvapi.GenericWalkAndParse(targetDir, re, func(matches []string) VideoInfo {
		// matches[1]: タイトル, matches[2]: 日時
		rawDate := matches[2]
		// 18-06-25 -> 18:06:25 に変換
		parts := strings.Split(rawDate, " ")
		dbTime := parts[0] + " " + strings.ReplaceAll(parts[1], "-", ":")

		return VideoInfo{
			Title:      strings.TrimSpace(matches[1]),
			FileName:   matches[0],
			DBDateTime: dbTime,
		}
	}, processor)
}

package image

import (
	"home/internal/hvapi"
	"log"
	"regexp"
	"time"
)

// VideoInfo は解析された動画ファイルの情報を含む構造体です。
type ImageInfo struct {
	FileName string
	FileTime time.Time // SQL/API向けの日時文字列 (例: 2025-09-13 18:06:25)
}

// ファイル名解析用の正規表現
// 対象例: "20260328-130702-249.jpg"
// グループ1: 日付(8桁), グループ2: 時間(6桁), グループ3: ミリ秒等(3桁), グループ4: 拡張子
var re = regexp.MustCompile(`(?i)^(\d{8})-(\d{6})-(\d{3})\.(jpg|jpeg|png|heic|mp4|mov|avi|webm)$`)

// --- 公開関数 ---

// WalkAndParse は指定されたフォルダを走査し、正規表現に一致するファイル情報を処理関数に渡します。
func WalkAndParse(targetDir string, processor func(info ImageInfo) error) error {
	return hvapi.GenericWalkAndParse(targetDir, re, func(matches []string) ImageInfo {

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

		return ImageInfo{
			FileName: matches[0], // ファイル名全体
			FileTime: fileTime,   // パースした日時
		}
	}, processor)
}

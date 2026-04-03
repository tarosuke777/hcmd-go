package image

import (
	"fmt"
	"home/internal/hvapi"
	"regexp"
)

// VideoInfo は解析された動画ファイルの情報を含む構造体です。
type ImageInfo struct {
	FileName   string
	DBDateTime string // SQL/API向けの日時文字列 (例: 2025-09-13 18:06:25)
}

// ファイル名解析用の正規表現
// 対象例: "20260328-130702-249.jpg"
// グループ1: 日付(8桁), グループ2: 時間(6桁), グループ3: ミリ秒等(3桁), グループ4: 拡張子
var re = regexp.MustCompile(`(?i)^(\d{8})-(\d{6})-(\d{3})\.(jpg|jpeg|png|heic|mp4|mov|avi|webm)$`)

// --- 公開関数 ---

// WalkAndParse は指定されたフォルダを走査し、正規表現に一致するファイル情報を処理関数に渡します。
func WalkAndParse(targetDir string, processor func(info ImageInfo) error) error {
	return hvapi.GenericWalkAndParse(targetDir, re, func(matches []string) ImageInfo {
		// matches[1]: 日付, matches[2]: 時間
		datePart, timePart := matches[1], matches[2]

		dbTime := fmt.Sprintf("%s-%s-%s %s:%s:%s",
			datePart[0:4], datePart[4:6], datePart[6:8],
			timePart[0:2], timePart[2:4], timePart[4:6])

		return ImageInfo{
			FileName:   matches[0], // ファイル名全体
			DBDateTime: dbTime,
		}
	}, processor)
}

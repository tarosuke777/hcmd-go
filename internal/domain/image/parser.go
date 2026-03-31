package image

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

// VideoInfo は解析された動画ファイルの情報を含む構造体です。
type ImageInfo struct {
	FileName   string
	DBDateTime string // SQL/API向けの日時文字列 (例: 2025-09-13 18:06:25)
}

// --- 定数 ---
const TargetDir = "./"

// ファイル名解析用の正規表現
// 対象例: "20260328-130702-249.jpg"
// グループ1: 日付(8桁), グループ2: 時間(6桁), グループ3: ミリ秒等(3桁), グループ4: 拡張子
var re = regexp.MustCompile(`^(\d{8})-(\d{6})-(\d{3})\.(jpg|jpeg|png|mp4|mov)$`)

// --- 公開関数 ---

// WalkAndParse は指定されたフォルダを走査し、正規表現に一致するファイル情報を処理関数に渡します。
func WalkAndParse(targetDir string, processor func(info ImageInfo) error) error {
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

		// matches[1]: 20260328 (日付)
		// matches[2]: 130702 (時間)
		// matches[3]: 249 (ミリ秒など)
		dbDateTime := formatToSQLDateTime(matches[1], matches[2])

		imageInfo := ImageInfo{
			FileName:   fileName,
			DBDateTime: dbDateTime,
		}

		// 処理関数を実行
		return processor(imageInfo)
	})
}

// --- 非公開ヘルパー関数 ---

// formatToSQLDateTime は連続した数字の文字列を SQL 形式 (YYYY-MM-DD HH:mm:ss) に変換します。
func formatToSQLDateTime(datePart, timePart string) string {
	if len(datePart) != 8 || len(timePart) != 6 {
		return datePart + " " + timePart
	}

	// YYYYMMDD -> YYYY-MM-DD
	formattedDate := fmt.Sprintf("%s-%s-%s", datePart[0:4], datePart[4:6], datePart[6:8])
	// HHmmss -> HH:mm:ss
	formattedTime := fmt.Sprintf("%s:%s:%s", timePart[0:2], timePart[2:4], timePart[4:6])

	return formattedDate + " " + formattedTime
}

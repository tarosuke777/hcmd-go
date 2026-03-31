package hvapi

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

// FileInfo は、videoやimageで共通して必要な最小限の情報です
type BaseFileInfo struct {
	FileName   string
	DBDateTime string
}

// GenericWalkAndParse は、正規表現と、マッチ結果を解析する関数(parser)を受け取って走査します
func GenericWalkAndParse[T any](targetDir string, re *regexp.Regexp, parser func(matches []string) T, processor func(info T) error) error {
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
			return nil // マッチしない場合はスキップ
		}

		// マッチした結果を、各パッケージ独自の構造体(T)に変換
		data := parser(matches)

		return processor(data)
	})
}

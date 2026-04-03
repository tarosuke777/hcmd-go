package hvapi

import (
	"fmt"
	"time"
)

// APITime は API とのやり取り専用の時刻型です
type APITime time.Time

// MarshalJSON で API が求める形式に自動変換されるようにします
func (t APITime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeLayout))
	return []byte(stamp), nil
}

// String メソッドを追加して、デバッグやログ出力で見やすくします
func (t APITime) String() string {
	return time.Time(t).Format(timeLayout)
}

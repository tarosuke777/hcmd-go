package hvapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

type MaxTimestampResponse struct {
	MaxCreatedAt string `json:"max_created_at"`
}

// FetchMaxTimestamp は指定されたURLから最新の時刻を取得します
func FetchMaxTimestamp(url string) (time.Time, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("status error: %d", resp.StatusCode)
	}

	var result MaxTimestampResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return time.Time{}, err
	}

	if result.MaxCreatedAt == "" {
		return time.Unix(0, 0), nil
	}

	return time.Parse(TimeLayout, result.MaxCreatedAt)
}

// SendToAPI は任意の構造体データをJSONとして指定されたURLにPOSTします (ジェネリクス使用)
func SendToAPI[T any](url string, data T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s, Details: %s", resp.Status, string(bodyBytes))
	}

	return nil
}

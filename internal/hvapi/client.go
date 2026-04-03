package hvapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const timeLayout = "2006-01-02 15:04:05"

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

	return time.Parse(timeLayout, result.MaxCreatedAt)
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

func UploadFile(url string, filePath string) (string, error) {
	// 1. ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 2. multipart/form-data の作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// "myImage" フィールドを作成 (handler.go の r.FormFile("myImage") に対応)
	part, err := writer.CreateFormFile("myImage", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// ファイルの中身を part にコピー
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// フォームの書き込みを完了させて、バウンダリを閉じます
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 3. POST送信
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	// multipart の Content-Type (バウンダリを含む) を設定
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// 4. レスポンスボディ (変換後のファイル名) を文字列として読み取って返す
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(respBody), nil
}

package main

import (
	"fmt"
	"home/internal/network"
	"home/internal/scanner"
	"os"
	"os/exec"
)

func main() {
	// コマンドライン引数をチェック
	// 期待する引数は "hms" です。
	if len(os.Args) < 2 {
		fmt.Println("Usage: home <service>")
		fmt.Println("Example: home hms")
		return
	}

	service := os.Args[1] // 最初の引数 ("hms") を取得
	command := ""
	if len(os.Args) >= 3 {
		command = os.Args[2]
	}
	var url string

	// 引数に基づいて開くURLを決定
	switch service {
	case "hms":
		url = "http://192.168.10.10/hms"
	case "hv":
		if command == "insert" {
			fmt.Printf("--- INFO: 'home hv insert' コマンドが検出されました。SQL生成処理を開始します。 ---\n")
			// 外部ファイル（sqlexport.go）で定義された関数を呼び出す
			scanner.GenerateInsertSQLs() 
			return
		}
		if command == "magic" {
			macAddress := os.Getenv("HV_MAC_ADDRESS")
			fmt.Printf("--- INFO: %s へのマジックパケット送信を開始します ---\n", macAddress)
			err := network.SendMagicPacket(macAddress)
			if err != nil {
				fmt.Printf("Error sending magic packet: %v\n", err)
			}
			return
		}
		url = "http://192.168.10.11/videos/v2/"
	default:
		fmt.Printf("Unknown service: %s\n", service)
		return
	}

	// Windowsの 'start' コマンドを使用してブラウザでURLを開く
	// cmd.exeの 'start' コマンドは、指定されたファイルをデフォルトの関連付けられたプログラムで開きます。
	// この場合、URLなのでデフォルトのブラウザで開かれます。
	// `/c` はコマンド実行後ウィンドウを閉じるためのオプションですが、ここでは 'start' が非同期にブラウザを起動するため不要です。
	// 正しくは、`start` コマンド自体を `cmd.exe` を使って実行します。

	// cmd /c start "" "URL" の形式で実行します。
	// 最初の "" はウィンドウタイトルとして必要です。
	cmd := exec.Command("cmd", "/c", "start", "", url)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error launching browser: %v\n", err)
		return
	}

	fmt.Printf("Launched browser for service '%s' at %s\n", service, url)
}
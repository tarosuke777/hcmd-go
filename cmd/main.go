package main

import (
	"fmt"
	"home/internal/api"
	"home/internal/network"
	"home/internal/parser"
	"home/internal/scanner"
	"log"
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
		if command == "sql" {
			fmt.Printf("--- INFO: 'home hv sql' コマンドが検出されました。SQL生成処理を開始します。 ---\n")
			scanner.GenerateInsertSQLs() 
			return
		}
		if command == "api" {
			fmt.Printf("--- INFO: 'home hv api' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			api.SyncVideosToAPI() 
			return
		}

		if command == "book" {
			fmt.Printf("--- INFO: 'home hv book' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			book, err := parser.CurrentFolderToBook()
			if err != nil {
				log.Fatal(err)
			}

		    // あとがAPI 用の構造体に詰め替えて送るだけ
			err = api.SyncBooksToAPI(book)
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
		url = "http://192.168.10.10/hv/videos/v2/"
	case "hb":
		url = "http://192.168.10.10/hb/"
	case "hc":
		url = "http://192.168.10.10/hc/"
	case "jenkins":
		url = "http://192.168.10.10/jenkins/"
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
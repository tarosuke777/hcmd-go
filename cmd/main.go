package main

import (
	"fmt"
	"home/internal/domain/backup"
	"home/internal/domain/book"
	"home/internal/domain/image"
	"home/internal/domain/video"
	"home/internal/jenkins"
	"home/internal/network"
	"home/internal/server"
	"log"
	"net/http"
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
		url = "https://192.168.10.10/hms"
	case "hv":
		fmt.Printf("--- Debug: HVモードに入りました (command: %s) ---\n", command)
		if command == "video" {
			fmt.Printf("--- INFO: 'home hv video' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			video.SyncVideosToAPI()
			return
		}

		if command == "image" {
			fmt.Printf("--- INFO: 'home hv image' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			image.SyncImagesToAPI()
			return
		}

		if command == "book" {
			fmt.Printf("--- INFO: 'home hv book' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			b, err := book.CurrentFolderToBook()
			if err != nil {
				log.Fatal(err)
			}

			err = book.SyncBooksToAPI(b)

			if err != nil {
				log.Fatal(err)
			}
			return
		}
		url = "https://192.168.10.10/hv/videos/"
	case "hb":
		url = "https://192.168.10.10/hb/"
	case "hc":
		url = "http://192.168.10.10/hc/"
	case "jenkins":

		if command == "build" {
			if len(os.Args) < 4 {
				fmt.Println("Usage: home jenkins build <job-name>")
				return
			}
			jobName := os.Args[3]

			jenkinsUrl := os.Getenv("JENKINS_URL")

			fmt.Printf("--- INFO: Jenkins ジョブ '%s' の実行を開始します ---\n", jobName)

			// クライアントの初期化
			client, err := jenkins.NewClientFromEnv(jenkinsUrl)
			if err != nil {
				log.Fatalf("Jenkins client init error: %v", err)
			}

			// ビルド実行
			if err := client.TriggerBuild(jobName); err != nil {
				log.Fatalf("Build error: %v", err)
			}

			fmt.Printf("Successfully triggered build for '%s'\n", jobName)
			return
		}

		url = "https://192.168.10.10/jenkins/"

	case "magic":
		macAddress := os.Getenv("HV_MAC_ADDRESS")
		fmt.Printf("--- INFO: %s へのマジックパケット送信を開始します ---\n", macAddress)
		err := network.SendMagicPacket(macAddress)
		if err != nil {
			fmt.Printf("Error sending magic packet: %v\n", err)
		}
		return
	case "upload":
		if command == "image" {
			fmt.Printf("--- INFO: 'home upload image' コマンドが検出されました。api呼び出し処理を開始します。 ---\n")
			image.UploadAndSyncImages()
			return
		}

		if command == "backup" {
			fmt.Printf("--- INFO: 'home upload backup' コマンドが検出されました。upload処理を開始します。 ---\n")
			backup.UploadAndSyncBackups()
			return
		}

		return
	case "web":
		port := ":8080"
		// ルーターを初期化（内部で books, images, videos を設定済み）
		router := server.NewRouter()

		fmt.Printf("--- Web Server Started ---\n")
		fmt.Printf("Endpoints: /books/, /images/, /videos/\n")

		if err := http.ListenAndServe(port, router); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
		return
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

Set ws = CreateObject("Wscript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

' 1. プロジェクトのルートディレクトリを絶対パスで指定
' (末尾に \ があってもなくても大丈夫です)
baseDir = "D:\content\Captures"

' 2. 実行ファイルのフルパスを組み立てる
exePath = fso.BuildPath(baseDir, "home.exe")

' 3. 作業ディレクトリ（カレントフォルダ）を baseDir に固定する
' これにより、Goプログラム内の相対パスがすべてここを起点にします
ws.CurrentDirectory = baseDir

' 4. 実行！ (0: 非表示, False: 終了を待たない)
' 引数 "web" を付けて起動します
ws.Run chr(34) & exePath & chr(34) & " web", 0, False
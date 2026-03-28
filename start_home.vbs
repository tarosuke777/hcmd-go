Set ws = CreateObject("Wscript.Shell")
' 実行ファイルの場所を指定（パスは自分の環境に合わせて書き換えてください）
' 第2引数の 0 が「画面を一切出さない」という魔法の数字です
ws.run "cmd /c ""D:\content\Captures\home.exe"" web", 0
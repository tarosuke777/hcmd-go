Set ws = CreateObject("Wscript.Shell")

' taskkill コマンドを非表示(0)で実行
' /F : 強制終了
' /IM : イメージ名（プロセス名）指定
ws.Run "taskkill /F /IM home.exe", 0, True

' 完了通知が欲しい場合は、下の行の ' を消してください
' MsgBox "home.exe を停止しました。"
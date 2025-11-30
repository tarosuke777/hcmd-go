# Doc

https://go.dev/doc/

# Build Command

go build -o home.exe cmd/main.go

# Add Windows Path(Power Shell)
$path = (Get-ItemProperty -Path 'HKCU:\Environment' -Name Path).Path
$newPath = $path + ";" + (Get-Location).Path
Set-ItemProperty -Path 'HKCU:\Environment' -Name Path -Value $newPath
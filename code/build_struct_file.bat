echo "start generate code!"
mkdir "..\config-data"
start go build
start ./code.exe -savePath="./../config-data" -readPath="./../csv"
set pathStr=%~dp0
start gofmt -s -w  "%pathStr%..\config-data\."
pause

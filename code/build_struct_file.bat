echo "start generate code!"
mkdir "..\config-data"
go build
code.exe -savePath="./../config-data" -readPath="./../csv"
set pathStr=%~dp0
gofmt -s -w  "%pathStr%..\config-data\."
pause

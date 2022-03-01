echo "start generate code!"
go_build_genStruct.exe -savePath="./../config-data" -readPath="./../csv"
set pathStr=%~dp0
gofmt -s -w  "%pathStr%..\config-data\."
pause

echo "start generate code!"
mkdir ".\config-data"
code.exe -configDataPath="./config-data" -csvPath="./csv"
set pathStr=%~dp0
gofmt -s -w  "%pathStr%config-data\."
pause

@echo off

echo "start generate code!"

start genStruct.exe -savePath="./configData" -readPath="./csv"

timeout /nobreak /t 3
exit

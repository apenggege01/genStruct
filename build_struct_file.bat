@echo off

echo "start generate code!"

start go_build_genStruct.exe -savePath="./configData" -readPath="./csv"

timeout /nobreak /t 3
exit

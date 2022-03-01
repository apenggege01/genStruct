#!/bin/bash
rm -rf ../config-data
mkdir ../config-data
go build
chmod +x ./code
./code -savePath="./../config-data" -readPath="./../csv"
gofmt -s -w  "../config-data/"


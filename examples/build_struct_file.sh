#!/bin/bash
rm -rf ./config-data
mkdir ./config-data
chmod +x ./code
./code -savePath="./config-data" -readPath="./csv"
gofmt -s -w  "./config-data/"


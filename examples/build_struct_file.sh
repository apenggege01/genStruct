#!/bin/bash
rm -rf ./config-data
mkdir ./config-data
chmod +x ./code
./code -configDataPath="./config-data" -csvPath="./csv"
gofmt -s -w  "./config-data/"


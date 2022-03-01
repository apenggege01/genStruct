#!/bin/bash

./go_build_genStruct -savePath="./../config-data" -readPath="./../csv"
set pathStr=%~dp0
gofmt -s -w  "../config-data/"


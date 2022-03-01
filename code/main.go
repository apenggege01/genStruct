package main

import (
	"flag"
	"fmt"
	"genStruct/code/tool"
)

var (
	savePath = flag.String("savePath", "", "Path to save the makefile")
	readPath = flag.String("readPath", "", "The path of reading Excel")
)

func doWhile(step string) {

	for {
		fmt.Println("savePath, readPath or allType is nil step:" + step)
	}
}

func main() {
	flag.Parse()
	if *savePath == "" || *readPath == "" {
		fmt.Println("savePath, readPath or allType is nil")
		return
	}
	fmt.Printf("savePath is :%s readPath is :%s", *savePath, *readPath)

	gt := tool.Generate{}
	err := gt.GenerateStruct(*readPath, *savePath)

	if err != nil {
		fmt.Printf("something err:%v\n", err)
		return
	}
}

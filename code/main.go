package main

import (
	"flag"
	"fmt"
	"github.com/apenggege01/genStruct/code/tool"
)

var (
	configDataPath = flag.String("configDataPath", "", "Path to save the makefile")
	csvPath = flag.String("csvPath", "", "The path of reading Excel")
)

func main() {
	flag.Parse()
	if *configDataPath == "" || *csvPath == "" {
		fmt.Println("configDataPath, csvPath or allType is nil")
		return
	}
	fmt.Printf("configDataPath is :%s csvPath is :%s", *configDataPath, *csvPath)

	gt := tool.Generate{}
	err := gt.GenerateStruct(*csvPath, *configDataPath)

	if err != nil {
		fmt.Printf("something err:%v\n", err)
		return
	}
}

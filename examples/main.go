package main

import (
	"fmt"
	"github.com/apenggege01/genStruct/examples/config-data"
)

func main() {

	// 引入生成的 configData 包
	var configManager configData.ConfigManager
	// 指定 csv 的路径
	errList := configManager.InitAllConfig("./csv/")
	if errList!= nil && len(errList) > 0 {
		fmt.Errorf("InitAllConfig err :\"%v ", errList)
	}

	//  获取某个 csv 配置的使用示例，例如 test1 表，查询 id 为 "1" 的字段
	test1CfgData, err := configManager.Test1Cfg.GetTest1CfgByID("1")
	if err != nil {
		fmt.Errorf("get test1Cfg data err :\"%v", err)
		return
	}
	fmt.Errorf("test1Cfg data is :\"%v", test1CfgData.RecycleReward)

	test2CfgData, err := configManager.Test2Cfg.GetTest2CfgByID(float32(1))
	if err != nil {
		fmt.Errorf("get test2Cfg data err :\"%v", err)
		return
	}
	fmt.Errorf("test2Cfg data is :\"%v", test2CfgData.RecycleReward)

	Test3Cfg, err := configManager.Test3Cfg.GetTest3CfgByID(int32(1))
	if err != nil {
		fmt.Errorf("get test3Cfg data err :\"%v", err)
		return
	}
	fmt.Errorf("test3Cfg data is :\"%v", Test3Cfg.RecycleReward)


	// 热更新接口调用示例
	fileNameList := make([]string, 0)
	fileNameList = append(fileNameList, "transfer")
	reloadErroList := configManager.ReloadConfig(fileNameList)
	if reloadErroList != nil && len(reloadErroList) > 0 {
		fmt.Errorf("ReloadConfig reloadErroList is  :\"%v", reloadErroList)
	}
}

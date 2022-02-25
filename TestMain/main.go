package main

import (
	"fmt"
	"genStruct/configData"
)

func main() {

	// 引入生成的configData包
	var configManager configData.ConfigManager
	// 指定csv的路径
	configManager.InitAllConfig("../csv/")

	//  获取某个csv配置的使用示例，例如transger表，查询id为1的字段
	transfer, er := configManager.TransferCfg.GetTransferCfgByID(1)
	if er != "" {
		return
	}
	fmt.Errorf("data is :\"%vtransfer", transfer.RecycleReward)

	// 热更新接口调用示例
	fileNameList := make([]string, 0)
	fileNameList = append(fileNameList, "transfer")
	erroFileList := configManager.ReloadConfig(fileNameList)
	if len(erroFileList) > 0 {
		fmt.Errorf("ReloadConfig err  :\"%v", erroFileList[0])
	}
}

package configData

	TransferCfg TransferDic
}
func (this *ConfigManager)InitAllConfig(csvRoot string){
	SetCSVPath(csvRoot)
	this.TransferCfg.InitTransferDic()
}

package configData

import (
	"errors"
	parse "github.com/apenggege01/genStruct"
)

type ConfigManager struct {
	      
	 Test1Cfg  Test1Dic
	 Test2Cfg  Test2Dic 
	 TransferCfg  TransferDic 
	
}

func (this *ConfigManager)InitAllConfig(csvRoot string) []error{
	errList := make([]error, 0)
	parse.SetCSVPath(csvRoot)
	
	var err error
	      
	err = this.Test1Cfg.InitTest1Dic()
	if err != nil{
		errList = append(errList, err)
	}
	err = this.Test2Cfg.InitTest2Dic()
	if err != nil{
		errList = append(errList, err)
	}
	err = this.TransferCfg.InitTransferDic()
	if err != nil{
		errList = append(errList, err)
	}
	
	
	return errList
}

func (this *ConfigManager)ReloadConfig(fileNameList []string) []error {
	erroList := make([]error, 0)
	if fileNameList == nil || len(fileNameList) <= 0{
		return append(erroList, errors.New("ReloadConfig fileNameList is empty"))
	}
	for _,fileName := range fileNameList{
		switch fileName {	
		      
		 case "test1":
			err := this.Test1Cfg.InitTest1Dic()
			if err != nil{
				erroList = append(erroList, err)
			}
		 case "test2":
			err := this.Test2Cfg.InitTest2Dic()
			if err != nil{
				erroList = append(erroList, err)
			}
		 case "transfer":
			err := this.TransferCfg.InitTransferDic()
			if err != nil{
				erroList = append(erroList, err)
			}
		
		default:
			erroList = append(erroList, errors.New(fileName + " file is not exist"))
		}
	}
	return erroList
}
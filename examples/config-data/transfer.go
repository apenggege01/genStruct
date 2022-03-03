
package configData

import "errors"
import "fmt"

type Transfer struct {
	                    
		Id int32 `csv:"id"` // 	显示顺序1 
		BadgeId int32 `csv:"badgeId"` // 	徽章编号 
		RuneType int64 `csv:"runeType"` // 	可镶嵌符文类型类型 
		RecycleReward []int32 `csv:"recycleReward"` // 	满星回收奖励 
		SkillId string `csv:"skillId"` // 	描述 
		LvInfo []string `csv:"lvInfo"` // 	各个等级说明 
		LightMax []float32 `csv:"lightMax"` // 	徽章升阶上限 
		AddHp float32 `csv:"addHp"` // 	自己给自己加血修正（百分比） 
		BeAddHp int `csv:"beAddHp"` // 	别人给自己加血修正（百分比） 
		ToAddHp int `csv:"toAddHp"` // 	我给别人加血修正（百分比） 
	
}

type TransferDic struct {
	TableName string
	MapRowsData map[int32] Transfer
}

func (this *TransferDic)InitTransferDic()(error){
	this.TableName = "transfer.csv"
	mapRowsData := make(map[int32]Transfer, 0)
	csv, err := NewWithOpts(this.TableName, Transfer{}, Comma(','), LazyQuotes(true), SkipLine(4))
	if err != nil {
		this.MapRowsData = make(map[int32]Transfer, 0)
		return err
	}

	for i := 0; i < csv.LineLen; i++ {
		temp := Transfer{}
		err = csv.Parse(temp)
		if err != nil {
			return err
		}
		if _,ok:= mapRowsData[temp.Id]; ok{
			return errors.New(this.TableName + fmt.Sprintf("%+v" , temp.Id))
		}
		mapRowsData[temp.Id] = temp
	}
	this.MapRowsData  = mapRowsData
	return nil
}
func (this *TransferDic) GetTransferCfgByID(id int32) (Transfer, error) {
	cfg, isok := this.MapRowsData[id]
	if !isok {
		return cfg, errors.New("rows not exist")
	}
	return cfg, nil
}
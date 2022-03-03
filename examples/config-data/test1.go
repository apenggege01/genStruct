
package configData

import "errors"
import "fmt"

type Test1 struct {
	                    
		Id string `csv:"id"` // 	显示顺序1 
		BadgeId int `csv:"badgeId"` // 	徽章编号 
		RuneType int `csv:"runeType"` // 	可镶嵌符文类型类型 
		RecycleReward []int `csv:"recycleReward"` // 	满星回收 
		SkillId int `csv:"skillId"` // 	普攻，对应skill 
		RuneMax int `csv:"runeMax"` // 	符文等级上限 
		LightMax int `csv:"lightMax"` // 	徽章升阶上限 
		AddHp int `csv:"addHp"` // 	自己给自己加血修正（百分比） 
		BeAddHp int `csv:"beAddHp"` // 	别人给自己加血修正（百分比） 
		ToAddHp int `csv:"toAddHp"` // 	我给别人加血修正（百分比） 
	
}

type Test1Dic struct {
	TableName string
	MapRowsData map[string] Test1
}

func (this *Test1Dic)InitTest1Dic()(error){
	this.TableName = "test1.csv"
	mapRowsData := make(map[string]Test1, 0)
	csv, err := NewWithOpts(this.TableName, Test1{}, Comma(','), LazyQuotes(true), SkipLine(4))
	if err != nil {
		this.MapRowsData = make(map[string]Test1, 0)
		return err
	}

	for i := 0; i < csv.LineLen; i++ {
		temp := Test1{}
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
func (this *Test1Dic) GetTest1CfgByID(id string) (Test1, error) {
	cfg, isok := this.MapRowsData[id]
	if !isok {
		return cfg, errors.New("rows not exist")
	}
	return cfg, nil
}
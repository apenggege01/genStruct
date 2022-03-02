package tool

import (
	"encoding/csv"
	"fmt"
	configData "genStruct/code/template-file"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const (
	TemplteFileName       = "template-file"
	ConfigManagerFileName = "config_manager.go"
)

/**
 * 将csv中的前四列转化为struct
 * 第一行字段类型		如 int
 * 第二行字段名称		如 配置表字段说明
 * 第三行字段名		如 id
 * 第四行配置表导出类型（s, c, all 	s 表示服务端使用 c 表示客户端使用 all 表示都使用）
 */
const (
	TablelineKey     = 0 // 配置表字段名字
	TablelineComment = 1 // 配置表字段说明
	TablelineType    = 2 // 配置表字段类型
	TablelineExpTyp  = 3 // 配置表导出类型（s, c, all 	s 表示服务端使用 c 表示客户端使用 all 表示都使用）
	lineNumber       = 3 // 每个工作表需要读取的行数
)

var genConfigManagerTmplate = `
package configData

import "errors"

type ConfigManager struct {
	{{range .CapitalStructNameDic}}  {{end}}
	{{range $lowStructName, $CapitalStructName := .CapitalStructNameDic}} {{$CapitalStructName}}Cfg  {{$CapitalStructName}}Dic 
	{{end}}
}

func (this *ConfigManager)InitAllConfig(csvRoot string) []error{
	errList := make([]error, 0)
	SetCSVPath(csvRoot)
	
	var err error
	{{range .CapitalStructNameDic}}  {{end}}
	{{range $lowStructName, $CapitalStructName := .CapitalStructNameDic}}err = this.{{$CapitalStructName}}Cfg.Init{{$CapitalStructName}}Dic()
	if err != nil{
		errList = append(errList, err)
	}
	{{end}}
	
	return errList
}

func (this *ConfigManager)ReloadConfig(fileNameList []string) []error {
	erroList := make([]error, 0)
	if fileNameList == nil || len(fileNameList) <= 0{
		return append(erroList, errors.New("ReloadConfig fileNameList is empty"))
	}
	for _,fileName := range fileNameList{
		switch fileName {	
		{{range .CapitalStructNameDic}}  {{end}}
		{{range $lowStructName, $CapitalStructName := .CapitalStructNameDic}} case "{{$lowStructName}}":
			err := this.{{$CapitalStructName}}Cfg.Init{{$CapitalStructName}}Dic()
			if err != nil{
				erroList = append(erroList, err)
			}
		{{end}}
		default:
			erroList = append(erroList, errors.New(fileName + " file is not exist"))
		}
	}
	return erroList
}`

var genStructTmplate = `
package configData

import "errors"
import "fmt"

type {{.CapitalStructName}} struct {
	{{range .FileInfoList}}  {{end}}
	{{range $idx, $value := .FileInfoList}}{{$value}} 
	{{end}}
}

type {{.CapitalStructName}}Dic struct {
	TableName string
	MapRowsData map[{{.PrimeKeyType}}] {{.CapitalStructName}}
}

func (this *{{.CapitalStructName}}Dic)Init{{.CapitalStructName}}Dic()(error){
	this.TableName = "{{.LowStructName}}.csv"
	mapRowsData := make(map[{{.PrimeKeyType}}]{{.CapitalStructName}}, 0)
	csv, err := NewWithOpts(this.TableName, {{.CapitalStructName}}{}, Comma(','), LazyQuotes(true), SkipLine(4))
	if err != nil {
		this.MapRowsData = make(map[{{.PrimeKeyType}}]{{.CapitalStructName}}, 0)
		return err
	}

	for i := 0; i < csv.LineLen; i++ {
		temp := {{.CapitalStructName}}{}
		err = csv.Parse(temp)
		if err != nil {
			return err
		}
		if _,ok:= mapRowsData[temp.Id]; ok{
			return errors.New(this.TableName + fmt.Sprintf({{.SimpleStr}} , temp.Id))
		}
		mapRowsData[temp.Id] = temp
	}
	this.MapRowsData  = mapRowsData
	return nil
}
func (this *{{.CapitalStructName}}Dic) Get{{.CapitalStructName}}CfgByID(id {{.PrimeKeyType}}) ({{.CapitalStructName}}, error) {
	cfg, isok := this.MapRowsData[id]
	if !isok {
		return cfg, errors.New("rows not exist")
	}
	return cfg, nil
}`

func (this *Generate)GenerateStructByTemplte(primeKeyType, structName string, fileInfoList []string) error{
	data := map[string]interface{}{
		"CapitalStructName": firstRuneToUpper(structName),
		"LowStructName": structName,
		"PrimeKeyType": primeKeyType,
		"FileInfoList": fileInfoList,
		"SimpleStr": "\"%+v\"",
	}

	tmpl, _ := template.New("test").Parse(genStructTmplate)
	fmt.Println("")
	fw, err := os.OpenFile(filepath.Join(this.savePath,structName+".go"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
	}
	defer fw.Close()
	tmpl.Execute(fw, data)

	return nil
}

type Generate struct {
	savePath      string   // 生成文件的保存路径
	FileNameSlice []string // 生成的公共文件的内容，所有结构自动注册到一个地方，服务引用这个文件就可以自动加载所有配置文件
}

func GetFileNameByFullName(fileName string) string {
	return strings.Split(fileName, ".")[0]
}

// 读取 csv 文件
func (this *Generate) GenerateStruct(readPath, savePath string) error {
	if savePath == "" {
		return fmt.Errorf("ReadExcel|savePath is nil")
	}
	this.savePath = savePath
	this.FileNameSlice = make([]string, 0, 100)

	var enc mahonia.Decoder
	enc = mahonia.NewDecoder("utf8")
	files, err := ioutil.ReadDir(readPath)
	if err != nil {
		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
	}
	for _, fileInfo := range files {
		if path.Ext(fileInfo.Name()) != ".csv" || hasChineseOrDefault(fileInfo.Name()) {
			continue
		}
		filePath := readPath + "/" + fileInfo.Name()

		file, openErr := os.Open(filePath)
		if openErr != nil {
			return fmt.Errorf("\nerror %v  open file %s", openErr, filePath)
		}
		defer file.Close()

		r := csv.NewReader(file)
		r.Comma = ','
		r.LazyQuotes = true
		//针对大文件，一行一行的读取文件
		sheetData := make([][]string, 0)
		line := 0
		for ; line < lineNumber; line++ {
			row, err := r.Read()
			if err != nil && err != io.EOF {
				fmt.Errorf("can not read, err is %+v", err)
			}
			if err == io.EOF {
				break
			}
			sheetData = append(sheetData, row)
		}

		// 判断表格中内容的行数是否小于需要读取的行数
		if line != lineNumber {
			return fmt.Errorf("ReadCSV %s |sheet.MaxRow:%d < lineNumber:%d", file.Name(), line, lineNumber)
		}

		if len(sheetData[TablelineType]) != len(sheetData[TablelineComment]) ||
			len(sheetData[TablelineType]) != len(sheetData[TablelineKey]) {
			return fmt.Errorf("ReadCSV %s sheetTitle len unequ", file.Name())
		}
		CellDatas := make([][]string, 0)
		for line = 0; line < len(sheetData[TablelineType]); line++ {
			cellData := make([]string, lineNumber, lineNumber)
			cellData[TablelineKey] = sheetData[TablelineKey][line]
			cellData[TablelineComment] = enc.ConvertString(sheetData[TablelineComment][line])
			cellData[TablelineType] = sheetData[TablelineType][line]
			//cellData[TablelineExpTyp] = sheetData[TablelineExpTyp][line]

			CellDatas = append(CellDatas, cellData)
		}

		err := this.SplicingData(configData.GetTypeName(sheetData[TablelineType][0]), CellDatas,
			GetFileNameByFullName(fileInfo.Name()))
		if err != nil {
			return fmt.Errorf("fileName:\"%v\" is err:%v", fileInfo.Name(), err)
		}
		this.FileNameSlice = append(this.FileNameSlice, GetFileNameByFullName(fileInfo.Name()))
	}

	//生成公共加载文件
	this.GenCommFile()
	//拷贝模板
	this.CopytemplateFile()
	return nil
}

func (this *Generate) GenCommFile() error {
	capitalStructNameDic := make(map[string]string, 0)
	for _, structName := range this.FileNameSlice {
		capitalStructNameDic[structName] = firstRuneToUpper(structName)
	}
	data := map[string]interface{}{
		"CapitalStructNameDic": capitalStructNameDic,
	}

	tmpl, _ := template.New("").Parse(genConfigManagerTmplate)
	fmt.Println("")
	fw, err := os.OpenFile(filepath.Join(this.savePath,ConfigManagerFileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
	}
	defer fw.Close()
	tmpl.Execute(fw, data)
	return nil
}

// 读取excel
//func (this *Generate) ReadExcel(readPath, savePath string) error {
//	if savePath == "" {
//		return fmt.Errorf("ReadExcel|savePath is nil")
//	}
//	this.savePath = savePath
//	files, err := ioutil.ReadDir(readPath)
//	if err != nil {
//		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
//	}
//	for _, file := range files {
//		if path.Ext(file.Name()) != ".csv" || hasChineseOrDefault(file.Name()) {
//			continue
//		}
//		filePath := readPath + "\\" + file.Name()
//		wb, err := xlsx.OpenFile(filePath)
//		if err != nil {
//			return fmt.Errorf("ReadExcel|xlsx.OpenFile is err :%v", err)
//		}
//		// 遍历工作表
//		for _, sheet := range wb.Sheets {
//			if hasChineseOrDefault(sheet.Name) {
//				continue
//			}
//			sheetData := make([][]string, 0)
//			// 判断表格中内容的行数是否小于需要读取的行数
//			if sheet.MaxRow < lineNumber {
//				return fmt.Errorf("ReadExcel|sheet.MaxRow:%d < lineNumber:%d", sheet.MaxRow, lineNumber)
//			}
//			// 遍历列
//			for i := 0; i < sheet.MaxCol; i++ {
//				// 判断某一列的第一行是否为空
//				if sheet.Cell(0, i).Value == "" {
//					continue
//				}
//				cellData := make([]string, 0)
//				// 遍历行
//				for j := 0; j < lineNumber; j++ {
//					cellData = append(cellData, sheet.Cell(j, i).Value)
//				}
//				sheetData = append(sheetData, cellData)
//			}
//			err := this.SplicingData(sheetData, sheet.Name)
//			if err != nil {
//				return fmt.Errorf("fileName:\"%v\" is err:%v", file.Name(), err)
//			}
//		}
//	}
//	if this.data == "" {
//		return fmt.Errorf("ReadExcel|this.data is nil")
//	}
//	//err = this.WriteNewFile(this.data)
//	//if err != nil {
//	//	return err
//	//}
//	return nil
//}

// 拼装 struct
func (this *Generate) SplicingData(primeKeyType string, data [][]string, structName string) error {
	fileInfoList := make([]string, 0)
	for _, value := range data {
		if len(value) != lineNumber {
			return fmt.Errorf("SplicingData|sheetName:%v col's len:%d is err", value, len(value))
		}
		err := this.CheckType(value[TablelineType], structName)
		if err != nil {
			return err
		}
		fileInfo := fmt.Sprintf("\t%s %s `csv:\"%s\"` // \t%s", firstRuneToUpper(value[TablelineKey]),
			configData.GetTypeName(value[TablelineType]), value[TablelineKey], value[TablelineComment])
		fileInfoList = append(fileInfoList, fileInfo)
	}
	return this.GenerateStructByTemplte(primeKeyType, structName, fileInfoList)
}

// 复制模板文件到生成路径
func (this *Generate) CopytemplateFile() error {
	files, err := ioutil.ReadDir(filepath.Join(TemplteFileName))
	if err != nil {
		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
	}
	for _, fileInfo := range files {
		fileName := fileInfo.Name()
		filePath := "./" + TemplteFileName + "/" + fileName
		source, openErr := os.Open(filePath)
		if openErr != nil {
			return fmt.Errorf("%s open file error %v", filePath, openErr)
		}
		defer source.Close()

		destination, err := os.OpenFile(this.savePath+"\\"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
	}
	return err
}

// 检测解析出来的字段类型是否符合要求
func (this *Generate) CheckType(dataType, structName string) error {
	name := configData.CheckTypeName(dataType)
	if name == "" {
		return fmt.Errorf("CheckType|struct:\"%v\" dataType:\"%v\" is not in provide dataType", structName, dataType)
	}
	return nil
}

// 字符串首字母转换成大写
func firstRuneToUpper(str string) string {
	data := []byte(str)
	for k, v := range data {
		if k == 0 {
			first := []byte(strings.ToUpper(string(v)))
			newData := data[1:]
			data = append(first, newData...)
			break
		}
	}
	return string(data[:])
}

// 判断是否存在汉字或者是否为默认的工作表
func hasChineseOrDefault(r string) bool {
	if strings.Index(r, "Sheet") != -1 {
		return true
	}
	for _, v := range []rune(r) {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

package tool

import (
	"encoding/csv"
	"fmt"
	configData "genStruct/templateFile"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

/**
 * 将csv中的前四列转化为struct
 * 第一行字段类型		如 int
 * 第二行字段名称		如 配置表字段说明
 * 第三行字段名		如 id
 * 第四行配置表导出类型（s, c, all 	s 表示服务端使用 c 表示客户端使用 all 表示都使用）
 */
const(
	TablelineType    = 0    // 配置表字段类型
	TablelineComment = 1    // 配置表字段说明
	TablelineKey     = 2    // 配置表字段名字
	TablelineExpTyp  = 3    // 配置表导出类型（s, c, all 	s 表示服务端使用 c 表示客户端使用 all 表示都使用）
	FormateChangeRow = "\n" // 换行
	lineNumber           = 4 // 每个工作表需要读取的行数
)

// 结构模板
const(
	structBegin          = "type %s struct {\n"                   // 结构体开始
	structValue          = "\t%s %s	`csv:\"%s\" client:\"%s\"`"   // 结构体的内容
	structValueForServer = "\t%s %s	`csv:\"%s\"`"                 // 服务端使用的结构体内容
	structRemarks        = "\t// %s"                              // 结构体备注
	structValueEnd       = "\n"                                   // 结构体内容结束
	structEnd            = "}\n"                                  // 结构体结束
	header 				 = "package configData\n\r" 			  // 文件头
)

// 管理对象结构
const (
	DicBegin          = "type %sDic struct {\n"          // 结构体开始
	DicValue          = "\tTableName string\n"         // 表名字
	DicValueForServer = "\tMapRowsData map[int] *%s\n" // 解析出来的字典
	DicValueEnd       = "\n"                             // 结构体内容结束
	DicEnd            = "}\n"                            // 结构体结束
)

// 初始化方法
const (
	InitStructDicMgr = "var %s"
	InitBegin        = "func (this *%sDic)Init%sDic(){\n"           // 初始化函数开始
	InitFilePath     = "\tthis.TableName = \"%s.csv\"\n"          // 表名
	TempMapRows      = "\tmapRowsData := make(map[int]*%s)\n" // 提取字典数据

	readCsv     = "\tcsv, err := NewWithOpts(this.TableName, %s{}, Comma(','), LazyQuotes(true), SkipLine(4))\n"
	readCsvRet  = "\tif err != nil {\n"
	readCsvRet1 = "\t\tthis.MapRowsData = make(map[int]*%s)\n" // 提取字典数据
	readCsvRet2 = "\t\treturn\n"
	readCsvRet3 = "\t}\n"

	readCsvRet4  = "\tfor i := 0; i < csv.LineLen; i++ {\n"
	readCsvRet5  = "\t\ttemp := &%s{}\n"
	readCsvRet6  = "\t\terr = csv.Parse(temp)\n"
	readCsvRet7  = "\t\tif err != nil {\n"
	readCsvRet8  = "\t\t\tbreak\n"
	readCsvRet9  = "\t\t}\n"
	readCsvRet10 = "\t\tmapRowsData[temp.Id] = temp\n"
	readCsvRet11 = "\t}\n"
	readCsvRet12 = "\tthis.MapRowsData  = mapRowsData\n"

	//这块是解析函数
	InitEnd = "}\n" //初始化完成
)

// 初始化方法
const (
	GetBegin     = "func (this *%sDic) Get%sCfgByID(id int) (*%s, string) {\n" // 查询函数开始
	GetDoGet     = "    cfg, isok := this.MapRowsData[id]\n"                   // 提取字典数据
	GetIsOK      = "    if !isok {\n"                                          // 判断是否提取成功
	GetRetErrVal = "		return nil, \"ErrRowNotExist\"\n"                       // 返回错误
	GetEnd       = "	}\n"                                                      //提取完成
	GetRetVal    = "	return cfg, \"\"\n"                                       // 返回错误
	GetFuncEnd   = "}\n"
)


// 通用管理文件
//type ConfigManager struct {
//
//	TransferCfg TransferDic
//
//}
//
//func (this *ConfigManager)InitAllConfig(){
//	this.TransferCfg.InitTransferDic()
//
//}
//
//func (this *ConfigManager)ReloadConfig(fileNameList []string) []string {
//	erroFileList := make([]string, 0)
//	if fileNameList == nil || len(fileNameList) <= 0{
//		return erroFileList
//	}
//	for _,fileName := range fileNameList{
//		switch fileName {
//			case %s :
//				this.%sCfg.Init%sDic()
//			default:
//				erroFileList = append(erroFileList, fileName)
//  	}
//	}
//  return erroFileList
//}


// 公共管理代码生成
const (
	ConfigManager1     = "type ConfigManager struct {\n" // 公共管理代码生成
	ConfigManager2     = "\t%sCfg %sDic\n" // 管理结构定义
	ConfigManager3     = "}\n" // 公共管理代码生成

	Init1     = "func (this *ConfigManager)InitAllConfig(csvRoot string){\n" // 公共管理代码生成
	Init2     = "\tSetCSVPath(csvRoot)\n" // 公共管理代码生成
	Init3     = "\tthis.%sCfg.Init%sDic()\n" // 公共管理代码生成
	Init4     = "}\n" // 公共管理代码生成
)

const(
	ReloadConst1 = "func (this *ConfigManager)ReloadConfig(fileNameList []string) []string {\n"
	ReloadConst2 = "\terroFileList := make([]string, 0)\n"
	ReloadConst3 = "\tif fileNameList == nil || len(fileNameList) <= 0{\n"
	ReloadConst4 = 	"\t\treturn erroFileList\n"
	ReloadConst5 = 	"\t}\n"
	ReloadConst6 = 	"\tfor _,fileName := range fileNameList{\n"
	ReloadConst7 = 	"\t\tswitch fileName {\n"
	ReloadConst8 = 	"\t\t\tcase \"%s\" :\n"
	ReloadConst9 = 	"\t\t\t\tthis.%sCfg.Init%sDic()\n"
	ReloadConst10 = "\t\t\tdefault:\n"
	ReloadConst11 = "\t\t\t\terroFileList = append(erroFileList, fileName)\n"
	ReloadConst12 = "\t\t}\n"
	ReloadConst13 = "\t}\n"
	ReloadConst14 = "\treturn erroFileList\n"
	ReloadConst15 = "}\n"
)

type Generate struct {
	savePath string // 生成文件的保存路径
	FileNameSlice     []string // 生成的公共文件的内容，所有结构自动注册到一个地方，服务引用这个文件就可以自动加载所有配置文件
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

		filePath := readPath + "\\" + fileInfo.Name()
		file, openErr := os.Open(filePath)
		if openErr != nil {
			return fmt.Errorf("%s open file error %v", readPath+"\\"+file.Name(), openErr)
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

		// 检查长度说明
		if len(sheetData[TablelineType]) != len(sheetData[TablelineComment]) || len(sheetData[TablelineKey]) != len(sheetData[TablelineExpTyp]) || len(sheetData[TablelineType]) != len(sheetData[TablelineKey]) {
			return fmt.Errorf("ReadCSV %s sheetTitle len unequ", file.Name())
		}

		CellDatas := make([][]string, 0)
		for line = 0; line < len(sheetData[TablelineType]); line++ {
			cellData := make([]string, 0)
			cellData = append(cellData, sheetData[TablelineType][line])
			cellData = append(cellData, enc.ConvertString(sheetData[TablelineComment][line]))
			cellData = append(cellData, sheetData[TablelineKey][line])
			cellData = append(cellData, sheetData[TablelineExpTyp][line])

			CellDatas = append(CellDatas, cellData)
		}

		err := this.SplicingData(CellDatas, GetFileNameByFullName(fileInfo.Name()))
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

// 生成公共文件
func (this *Generate) GenCommFile() error {

	structData := ""
	//生成表解析函数和获取函数
	structData += fmt.Sprintf(ConfigManager1)
	for _, structName := range this.FileNameSlice {
		structData += fmt.Sprintf(ConfigManager2, firstRuneToUpper(structName), firstRuneToUpper(structName))
	}
	structData += fmt.Sprintf(ConfigManager3)

	structData += fmt.Sprintf(FormateChangeRow)
	structData += fmt.Sprintf(Init1)
	structData += fmt.Sprintf(Init2)

	for _, structName := range this.FileNameSlice {
		structData += fmt.Sprintf(Init3, firstRuneToUpper(structName), firstRuneToUpper(structName))
	}
	structData += fmt.Sprintf(Init4)


	// 热更新调用
	structData += fmt.Sprintf(FormateChangeRow)
	structData += fmt.Sprintf(ReloadConst1)
	structData += fmt.Sprintf(ReloadConst2)
	structData += fmt.Sprintf(ReloadConst3)
	structData += fmt.Sprintf(ReloadConst4)
	structData += fmt.Sprintf(ReloadConst5)
	structData += fmt.Sprintf(ReloadConst6)
	structData += fmt.Sprintf(ReloadConst7)
	for _, structName := range this.FileNameSlice {
		structData += fmt.Sprintf(ReloadConst8, structName)
		structData += fmt.Sprintf(ReloadConst9, firstRuneToUpper(structName), firstRuneToUpper(structName))
	}
	structData += fmt.Sprintf(ReloadConst10)
	structData += fmt.Sprintf(ReloadConst11)
	structData += fmt.Sprintf(ReloadConst12)
	structData += fmt.Sprintf(ReloadConst13)
	structData += fmt.Sprintf(ReloadConst14)
	structData += fmt.Sprintf(ReloadConst15)

	err := this.WriteNewFile(structData, "ConfigManager.go")
	if err != nil {
		return err
	}
	//this.data += structData
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
func (this *Generate) SplicingData(data [][]string, structName string) error {
	structData := fmt.Sprintf(structBegin, firstRuneToUpper(structName))
	for _, value := range data {
		if len(value) != lineNumber {
			return fmt.Errorf("SplicingData|sheetName:%v col's len:%d is err", value, len(value))
		}
		err := this.CheckType(value[TablelineType], structName)
		if err != nil {
			return err
		}
		switch value[TablelineExpTyp] {
		case "all":
			structData += fmt.Sprintf(structValue, firstRuneToUpper(value[TablelineKey]), configData.GetTypeName(value[TablelineType]), value[TablelineKey], value[TablelineKey])
			if value[TablelineComment] != "" {
				structData += fmt.Sprintf(structRemarks, value[TablelineComment])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "s":
			structData += fmt.Sprintf(structValueForServer, firstRuneToUpper(value[TablelineKey]), configData.GetTypeName(value[TablelineType]), value[TablelineKey])
			if value[TablelineComment] != "" {
				structData += fmt.Sprintf(structRemarks, value[TablelineComment])
			}
			structData += fmt.Sprintf(structValueEnd)
		case "c":
			continue
		default:
			return fmt.Errorf("SplicingData|value[TablelineExpTyp]:\"%v\" is not in s,c,all", value[TablelineExpTyp])
		}
	}
	structData += structEnd

	//生成表解析函数和获取函数
	structData += fmt.Sprintf(FormateChangeRow)
	structData += fmt.Sprintf(DicBegin, firstRuneToUpper(structName))
	structData += fmt.Sprintf(DicValue)
	structData += fmt.Sprintf(DicValueForServer, firstRuneToUpper(structName))
	structData += fmt.Sprintf(DicValueEnd)
	structData += fmt.Sprintf(DicEnd)

	//初始化解析函数
	structData += fmt.Sprintf(FormateChangeRow)
	structData += fmt.Sprintf(InitBegin, firstRuneToUpper(structName), firstRuneToUpper(structName))
	structData += fmt.Sprintf(InitFilePath, structName)
	structData += fmt.Sprintf(TempMapRows, firstRuneToUpper(structName))

	structData += fmt.Sprintf(readCsv, firstRuneToUpper(structName))
	structData += fmt.Sprintf(readCsvRet)
	structData += fmt.Sprintf(readCsvRet1, firstRuneToUpper(structName))
	structData += fmt.Sprintf(readCsvRet2)
	structData += fmt.Sprintf(readCsvRet3)
	structData += fmt.Sprintf(readCsvRet4)
	structData += fmt.Sprintf(readCsvRet5, firstRuneToUpper(structName))
	structData += fmt.Sprintf(readCsvRet6)
	structData += fmt.Sprintf(readCsvRet7)
	structData += fmt.Sprintf(readCsvRet8)
	structData += fmt.Sprintf(readCsvRet9)
	structData += fmt.Sprintf(readCsvRet10)
	structData += fmt.Sprintf(readCsvRet11)
	structData += fmt.Sprintf(readCsvRet12)

	structData += fmt.Sprintf(InitEnd)

	// 提取字典数据方法
	structData += fmt.Sprintf(FormateChangeRow)
	structData += fmt.Sprintf(GetBegin, firstRuneToUpper(structName), firstRuneToUpper(structName), firstRuneToUpper(structName))
	structData += fmt.Sprintf(GetDoGet)
	structData += fmt.Sprintf(GetIsOK)
	structData += fmt.Sprintf(GetRetErrVal)
	structData += fmt.Sprintf(GetEnd)
	structData += fmt.Sprintf(GetRetVal)
	structData += fmt.Sprintf(GetFuncEnd)

	err := this.WriteNewFile(structData, structName+".go")
	if err != nil {
		return err
	}
	//this.data += structData
	return nil
}

// 拼装好的struct写入新的文件
func (this *Generate) WriteNewFile(data string, fileName string) error {
	str := strings.Split(this.savePath, "\\")
	if len(str) == 0 {
		return fmt.Errorf("WriteNewFile|len(str) is 0")
	}
	//header = fmt.Sprintf(header, str[len(str)-1])
	data = header + data
	fw, err := os.OpenFile(this.savePath+"\\"+fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("WriteNewFile|OpenFile is err:%v", err)
	}
	defer fw.Close()
	_, err = fw.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("WriteNewFile|Write is err:%v", err)
	}
	return nil
}

// 复制模板文件到生成路径
func (this *Generate) CopytemplateFile() error {

	files, err := ioutil.ReadDir("./templateFile")
	if err != nil {
		return fmt.Errorf("ReadExcel|ReadDir is err:%v", err)
	}
	for _, fileInfo := range files {
		//if path.Ext(fileInfo.Name()) != ".csv" || hasChineseOrDefault(fileInfo.Name()) {
		//	continue
		//}

		fileName := fileInfo.Name()
		filePath :=  "./templateFile/" + fileName
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

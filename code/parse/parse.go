package configData

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var invalidValue = reflect.Value{}
var FileRoot = "../csv/"

// Default converters for basic types.
var converters = map[reflect.Kind]func(string) reflect.Value{
	reflect.Bool:    convertBool,
	reflect.Float32: convertFloat32,
	reflect.Float64: convertFloat64,
	reflect.Int:     convertInt,
	reflect.Int8:    convertInt8,
	reflect.Int16:   convertInt16,
	reflect.Int32:   convertInt32,
	reflect.Int64:   convertInt64,
	reflect.String:  convertString,
	reflect.Uint:    convertUint,
	reflect.Uint8:   convertUint8,
	reflect.Uint16:  convertUint16,
	reflect.Uint32:  convertUint32,
	reflect.Uint64:  convertUint64,
}
var (
	boolSliceType    = reflect.TypeOf([]bool{})
	float32SliceType = reflect.TypeOf([]float32{})
	float64SliceType = reflect.TypeOf([]float64{})
	intSliceType     = reflect.TypeOf([]int{})
	int8SliceType    = reflect.TypeOf([]int8{})
	int16SliceType   = reflect.TypeOf([]int16{})
	int32SliceType   = reflect.TypeOf([]int32{})
	int64SliceType   = reflect.TypeOf([]int64{})
	stringSliceType  = reflect.TypeOf([]string{})
	uintSliceType    = reflect.TypeOf([]uint{})
	uint8SliceType   = reflect.TypeOf([]uint8{})
	uint16SliceType  = reflect.TypeOf([]uint16{})
	uint32SliceType  = reflect.TypeOf([]uint32{})
	uint64SliceType  = reflect.TypeOf([]uint64{})
)

var sliceConvertes = map[reflect.Type]func([]string) reflect.Value{
	boolSliceType:    convertBools,
	float32SliceType: convertFloat32s,
	float64SliceType: convertFloat64s,
	intSliceType:     convertInts,
	int8SliceType:    convertInt8s,
	int16SliceType:   convertInt16s,
	int32SliceType:   convertInt32s,
	int64SliceType:   convertInt64s,
	stringSliceType:  convertStrings,
	uint8SliceType:   convertUint8s,
	uint16SliceType:  convertUint16s,
	uint32SliceType:  convertUint32s,
	uint64SliceType:  convertUint64s,
}

//名字映射类型
var TypeToName = map[string]string{
	"Bool":    "bool",
	"Float32": "float32",
	"Float64": "float64",
	"Int":     "int",
	"Int8":    "int8",
	"Int16":   "int16",
	"Int32":   "int32",
	"Int64":   "int64",
	"String":  "string",
	"Uint":    "uint",
	"Uint8":   "uint8",
	"Uint16":  "uint16",
	"Uint32":  "uint32",
	"Uint64":  "uint64",

	"BoolArray":     "[]bool",
	"Float32SArray": "[]float32",
	"Float64SArray": "[]float64",
	"IntArray":      "[]int",
	"Int8Array":     "[]int8",
	"Int16Array":    "[]int16",
	"Int32Array":    "[]int32",
	"Int64Array":    "[]int64",
	"StringArray":   "[]string",
	"IintArray":     "[]uint",
	"Uint8Array":    "[]uint8",
	"Uint16Array":   "[]uint16",
	"Uint32Array":   "[]uint32",
	"Uint64Array":   "[]uint64",
}

//根据类型名字获取实际类型
func GetTypeName(name string) string {
	tname, ok := TypeToName[name]
	if ok {
		return tname
	}
	return name
}

//根据类型名字获取实际类型
func CheckTypeName(name string) string {
	tname, ok := TypeToName[name]
	if ok {
		return tname
	}
	return ""
}

func convertBool(value string) reflect.Value {
	if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertBools(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertBool(values[i])
	}
	ret := reflect.MakeSlice(boolSliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertFloat32(value string) reflect.Value {
	if v, err := strconv.ParseFloat(value, 32); err == nil {
		return reflect.ValueOf(float32(v))
	}
	return invalidValue
}

func convertFloat32s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertFloat32(values[i])
	}
	ret := reflect.MakeSlice(float32SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertFloat64(value string) reflect.Value {
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertFloat64s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertFloat64(values[i])
	}
	ret := reflect.MakeSlice(float64SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertInt(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 0); err == nil {
		return reflect.ValueOf(int(v))
	}
	return invalidValue
}

func convertInts(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertInt(values[i])
	}
	ret := reflect.MakeSlice(intSliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertInt8(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 8); err == nil {
		return reflect.ValueOf(int8(v))
	}
	return invalidValue
}

func convertInt8s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertInt8(values[i])
	}
	ret := reflect.MakeSlice(int8SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertInt16(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 16); err == nil {
		return reflect.ValueOf(int16(v))
	}
	return invalidValue
}

func convertInt16s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertInt16(values[i])
	}
	ret := reflect.MakeSlice(int16SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertInt32(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 32); err == nil {
		return reflect.ValueOf(int32(v))
	}
	return invalidValue
}

func convertInt32s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertInt32(values[i])
	}
	ret := reflect.MakeSlice(int32SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertInt64(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertInt64s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertInt64(values[i])
	}
	ret := reflect.MakeSlice(int64SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertString(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func convertStrings(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertString(values[i])
	}
	ret := reflect.MakeSlice(stringSliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertUint(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 0); err == nil {
		return reflect.ValueOf(uint(v))
	}
	return invalidValue
}

func convertUInts(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertUint(values[i])
	}
	ret := reflect.MakeSlice(uintSliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertUint8(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 8); err == nil {
		return reflect.ValueOf(uint8(v))
	}
	return invalidValue
}
func convertUint8s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertUint8(values[i])
	}
	ret := reflect.MakeSlice(uint8SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertUint16(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 16); err == nil {
		return reflect.ValueOf(uint16(v))
	}
	return invalidValue
}

func convertUint16s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertUint16(values[i])
	}
	ret := reflect.MakeSlice(uint16SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertUint32(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 32); err == nil {
		return reflect.ValueOf(uint32(v))
	}
	return invalidValue
}

func convertUint32s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertUint32(values[i])
	}
	ret := reflect.MakeSlice(uint32SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

func convertUint64(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertUint64s(values []string) reflect.Value {
	items := make([]reflect.Value, len(values))
	for i := range values {
		items[i] = convertUint64(values[i])
	}
	ret := reflect.MakeSlice(uint64SliceType, 0, 0)
	return reflect.Append(ret, items...)
}

type FieldDefine struct {
	reflect.StructField
	FieldIndex int
}

type Csv4g struct {
	name       string
	fields     []*FieldDefine
	lines      [][]string
	lineNo     int
	LineLen    int
	lineOffset int
}

type Option struct {
	Comma      rune
	LazyQuotes bool
	SkipLine   int
}

func Comma(r rune) func(*Option)      { return func(opt *Option) { opt.Comma = r } }
func LazyQuotes(b bool) func(*Option) { return func(opt *Option) { opt.LazyQuotes = b } }
func SkipLine(l int) func(*Option)    { return func(opt *Option) { opt.SkipLine = l } }

func SetCSVPath(csvRoot string) {
	FileRoot = csvRoot
}

func NewWithOpts(fileName string, o interface{}, options ...func(*Option)) (*Csv4g, error) {
	file, openErr := os.Open(FileRoot + fileName)
	if openErr != nil {
		return nil, fmt.Errorf("%s open file error %v", fileName, openErr)
	}
	defer file.Close()

	defaultOpt := &Option{Comma: ',', LazyQuotes: false, SkipLine: 4}
	for _, opt := range options {
		opt(defaultOpt)
	}

	r := csv.NewReader(file)
	r.Comma = defaultOpt.Comma
	r.LazyQuotes = defaultOpt.LazyQuotes
	var err error
	var fields []string

	offset := defaultOpt.SkipLine
	for i := 0; i < offset; i++ {
		row, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("%s skip line error %v", file.Name(), err)
		}

		//第三行是字段名字
		if i == 2 {
			fields = row // first line is field's description
		}
	}

	tType := reflect.TypeOf(o)
	if tType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%v must be a struct, cannot be an interface or pointer", tType.Elem().Name())
	}
	ret := &Csv4g{
		name:       file.Name(),
		fields:     make([]*FieldDefine, 0),
		lineNo:     0,
		lineOffset: offset + 1}

Out:
	for i := 0; i < tType.NumField(); i++ {
		f := tType.Field(i)
		tagStr := f.Tag.Get("csv")
		fieldName := f.Name
		canSkip := false
		if tagStr != "" {
			tags := strings.Split(tagStr, ",")
			for _, tag := range tags {
				switch tag {
				//这个以后可以做扩展，现在用不到
				case "-":
					continue Out
				case "omitempty":
					canSkip = true
				default:
					fieldName = tag
				}
			}
		}
		fd := &FieldDefine{f, 0}
		index := -1
		for j := range fields {
			if fields[j] == fieldName {
				index = j
				break
			}
		}
		if index == -1 {
			if !canSkip {
				return nil, fmt.Errorf("%s cannot find field %s", file.Name(), f.Name)
			}
			continue
		}
		fd.FieldIndex = index
		ret.fields = append(ret.fields, fd)
	}

	var lines [][]string
	lines, err = r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s read error %v", file.Name(), err)
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("%s has no data!", file.Name())
	}
	ret.lines = lines
	ret.LineLen = len(lines)
	return ret, nil
}

func (this *Csv4g) Parse(obj interface{}) (err error) {
	if this.lineNo >= len(this.lines) {
		return io.EOF
	}
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("%s error on parse line %d [%v]", this.name, this.lineNo+this.lineOffset+1, x)
			return
		}
		this.lineNo++
	}()
	values := this.lines[this.lineNo]
	elem := reflect.ValueOf(obj).Elem()
	for index := range this.fields {
		f := elem.FieldByIndex(this.fields[index].Index)
		value := values[this.fields[index].FieldIndex]
		if conv, ok := converters[f.Kind()]; ok {
			v := conv(value)
			f.Set(v)
		} else {
			if f.Kind() == reflect.Slice {
				if sliceConv, ok := sliceConvertes[f.Type()]; ok {
					f.Set(sliceConv(strings.Split(value, "|")))
				} else {
					err = fmt.Errorf("%s:[%d] unsupported field set %v -> %v :[%d].",
						this.name, this.lineNo+this.lineOffset, this.fields[index], value)
				}
			} else {
				err = fmt.Errorf("%s:[%d] unsupported field set %v -> %v :[%d].",
					this.name, this.lineNo+this.lineOffset, this.fields[index], value)
			}
		}
	}
	return
}

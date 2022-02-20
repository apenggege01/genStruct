## 指定的csv转化成go文件中的struct ##
该程序使用bat批处理文件运行，使用之前需要先修改批处理文件的路径

	::生成配置结构文件、配置管理器的路径
    set savePath=F:\code\src\genStruct\tool
	
	::目标excel文件路径
    set readPath=F:\code\src\genStruct

工程目录说明
configData：生成的代码库，里面的代码不允许用户更改
csv：测试用的配置csv路径
templateFile：通用解析库主要包括解析方法
TestMain：使用测试项目，初次使用可以参考里面内容
tool：代码模板生成工具
main.go:代码模板生成工具入口函数
build_struct_file.bat:生成结构脚本


build_struct_file.bat:
参数说明：详见bat文件内部说明


configData目录文件说明
Parse.go:templateFile目录下拷贝过来的公共模板库
configManger.go：配置文件管理生成器

具体使用示例，请查看TestMain示例用法


后续扩展：
1.目前只做了简单的首个字段做主键（必须是id类型是int），后续可以扩展符合组件代码生成，方便有相关需求的配置扩展
2.日志函数可以设定，保证和项目输出日志耦合
3.热更新方法扩展，可以根据需求添加热更新接口，实现reloadall配置或者单个配置热更新
4.客户端代码，或者别的语言代码生成，可以根据自己语言规则，自行扩充，思路相同
5.如果只想做部分表的自动加载，可以生成注册函数，调用者自动注册，然后调用给初始化函数加载所有配置
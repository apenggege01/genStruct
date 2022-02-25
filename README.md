## 指定的csv转化成go文件中的struct ##
## 对应文件要转码utf8 ##
该程序使用bat批处理文件运行，使用之前需要先修改批处理文件的路径

	::生成配置结构文件、配置管理器的路径（也可以用相对路径）
    set savePath="C:\Users\Administrator\Desktop\genStruct\configData"
	
	::配置文件（excel）路径（也可以用相对路径）
    set readPath="C:\Users\Administrator\Desktop\genStruct\csv"

一、工程目录以及文件说明
configData：生成的代码库，里面的代码不允许用户更改；
			configManager.go 这个是所有配置文件的管理类
							1.自动注册所有配置文件结构
							2.InitAllConfig 解析所有文件到内存
							3.ReloadConfig重新加载文件 参数：需要重新加载的配置文件名字列表
			parse.go 这个是模板文件，用来把配置中的字符串转成自定义的数据结构，templateFile文件夹下拷贝过来的，里面目前有一些基础数据的解析，复杂数据可根据项目自定义解析方法
			其他文件为配置生成的对应go文件
			这些文件都包含两个结构、三个接口：1.xxx结构就是对应的每一行配置的一个struct；
											  2.xxxDic结构就是以主键为key xxx结构指针为值的字典，初始化后囊括了整个xxx配置文件的数据（一个配置文件就是一个字典）；
											  3.InitxxxDic接口，用来解析xxx配置文件到xxxDic中；热更新某个文件后，直接调用此接口就可以重新解析配置文件数据；
											  4.GetxxxCfgByID接口根据主键id取出某行数据对应的xxx结构
											  
csv：测试用的配置csv路径，这个路径可以根据项目自行定义
templateFile：通用解析库主要包括解析方法，里面存放 parse.go 文件
			parse.go 文件用来自定解析规则，开发者可以根据常用的结构自行在里面定义解析规则：
			详情请看代码
			
TestMain：使用测试项目，初次使用可以参考里面内容
tool：代码模板生成工具 codeTemplate
main.go:代码模板生成工具入口函数
build_struct_file.bat:windows下生成结构脚本


build_struct_file.bat:
参数说明：详见bat文件内部说明

二、示例
具体使用示例，请查看TestMain示例用法


后续扩展：
1.目前只做了简单的首个字段做主键（必须是id类型是int），后续可以扩展符合组件代码生成，方便有相关需求的配置扩展
2.日志函数可以设定，保证和项目输出日志耦合
3.客户端代码，或者别的语言代码生成，可以根据自己语言规则，自行扩充，思路相同
4.如果只想做部分表的自动加载，可以生成注册函数，调用者自动注册，然后调用给初始化函数加载所有配置

思考几个问题
1.reload有没有必要加锁（现在是指针替换我觉得没必要加）
2.数据get的时候，是否应该拷贝出来，防止被调用者修改；带来的问题是每次调用都存在数据拷贝问题，影响效率
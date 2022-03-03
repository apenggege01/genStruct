## 指定的csv转化成go文件struct ##
## 对应CSV文件要转码utf8!!! ##
# **工程目录以及文件说明:**

## code: csv转go结构代码生成工具库

​	1.parse：里面存放 parse.go 文件

​		parse.go 文件用来自定解析规则，开发者可以根据常用的结构自行在里面定义解析规则：
​		详情请看代码		

​	2.tool:里面存放 gen_struct_tool.go 文件

​		 gen_struct_tool.go csv结构生成工具代码，里面包含了模板生成代码

### examples：使用测试项目，初次使用可以参考里面内容

##### 	config-data：生成的代码库，里面的代码不允许用户更改（实际路径根据自己项目目录结构自由设定）；

​		1.config_manager.go 这个是所有配置文件的管理类

​						a.自动注册所有配置文件结构

​						b.InitAllConfig 解析所有文件到内存

​						c.ReloadConfig重新加载文件 参数：需要重新加载的配置文件名字列表

​		2.其他文件为配置生成的对应go文件，这些文件都包含两个结构、三个接口：

​						a.xxx结构就是对应的每一行配置的一个struct；
​						b.xxxDic结构就是以主键为key xxx结构指针为值的字典，初始化后囊括了整个xxx配置文件的数据（一个配置文件一个字典);
​						c.InitxxxDic接口，用来解析xxx配置文件到xxxDic中；热更新某个文件后，直接调用此接口就可以重新解析配置文件数据；
​						d.GetxxxCfgByID接口根据主键id取出某行数据对应的xxx结构(主键id类型可以自由定义);

##### 	CSV：配置文件存放目录（实际路径根据自己项目目录结构自由设定）

##### 	build_struct_file.sh: linux下生成结构脚本，它会执行 code 编译后的二进制文件生成 config-data代码

##### 	build_struct_file.bat: windows下生成结构脚本，它会执行 code 编译后的二进制文件生成 config-data代码

​		脚本需要配置两个路径

生成配置结构文件、配置管理器的路径（也可以用相对路径）
​		set configDataPath="xxx\config-data"

​		配置文件（excel）路径（也可以用相对路径）
​		set csvPath="xxx\csv"



------

## 使用示例：

```go
// 引入生成的 configData 包
	var configManager configData.ConfigManager
	// 指定 csv 的路径
	errList := configManager.InitAllConfig("./csv/")
	if errList!= nil && len(errList) > 0 {
		fmt.Errorf("InitAllConfig err :\"%v ", errList)
	}

	//  获取某个 csv 配置的使用示例，例如 transger 表，查询 id 为1的字段
	Test3Cfg, err := configManager.Test3Cfg.GetTest3CfgByID(int32(1))
	if err != nil {
		fmt.Errorf("get transfer data err :\"%v", err)
		return
	}
	fmt.Errorf("transfer data is :\"%v", Test3Cfg.RecycleReward)

	test1CfgData, err := configManager.Test1Cfg.GetTest1CfgByID("1")
	if err != nil {
		fmt.Errorf("get test1Cfg data err :\"%v", err)
		return
	}
	fmt.Errorf("test1Cfg data is :\"%v", test1CfgData.RecycleReward)

	test2CfgData, err := configManager.Test2Cfg.GetTest2CfgByID(float32(1))
	if err != nil {
		fmt.Errorf("get test1Cfg data err :\"%v", err)
		return
	}
	fmt.Errorf("test1Cfg data is :\"%v", test2CfgData.RecycleReward)

	// 热更新接口调用示例
	fileNameList := make([]string, 0)
	fileNameList = append(fileNameList, "transfer")
	reloadErroList := configManager.ReloadConfig(fileNameList)
	if reloadErroList != nil && len(reloadErroList) > 0 {
		fmt.Errorf("ReloadConfig reloadErroList is  :\"%v", reloadErroList)
	}
```


后续扩展：
1.目前只做了简单的首个字段做主键（支持多类型例如string int32, float32等等），后续可以扩展多主键代码生成，方便有相关需求的配置扩展
2.客户端代码，或者别的语言代码生成，可以根据自己语言规则，自行扩充(详细请看 /core/tool/gen_struct_tool.go 模板解析成代码部分逻辑)
3.如果只想做部分表的自动加载，可以生成注册函数，调用者自动注册，然后调用给初始化函数加载所有配置（已经支持所有表的热更新，更新顺序按照传入的表名字顺序）

思考几个问题
1.reload 有没有必要加锁（现在是指针替换我觉得没必要加）
2.数据 get 的时候，是否应该拷贝出来，防止被调用者修改；带来的问题是每次调用都存在数据拷贝问题，影响效率（目前是返回数据，不返回指针，安全第一的原则）
## 简介
此项目是一个开源的golang的日志功能项目，go的版本是1.17

## 核心技术及组件
logrus v1.8.1

file-rotatelogs v2.4.0+incompatible

## 获取
go get github.com/yingying0708/GoLogger

## 使用
log := gologhelper.GetLogHelper("项目名称", "日志路径").SetConsolePrint(true).SetLogLevel("info").SetWhen("m").SetBackupCount(2)

log.Info("info")

## 方法介绍
SetConsolePrint(true)  控制台打印，默认是false


SetLogLevel("info")    只生成info的日志文件，默认是error，日志级别从小到大排列是：track，debug，info，warn，error。例：设置info ，会打印info、warn、error的类型的日志


SetWhen("m")           设置按什么进行分割，默认是D，D：天，H：小时，m：分钟


SetBackupCount(2)      保留多少个文件，默认是15

## 流水日志

字段规范

字段名称  | 字段类型  | 是否必填  | 字段中文含义  | 备注  
  ----  |  ------ | ------ |  ------- | --------   
transaction_id  | String  | 是 | 流水 id |   
dialog_type | String | 是  | 日志类型(in/out)  | in代表内部流水out 代表外部流水内部流水：当前记录日志的系统访问其他服务的调用行为记录，属于内部流水外部流水：当前记录日志的系统作为服务提供方，被其他系统访问时，记录的处理记录日志，属于外部流水如果该字段为空，默认认为当前日志属于外部流水日志
address | String  | 是 | 请求地址  | 请求访问的地址
fcode | String  | 是 | 调用方系统编码 | 发起请求的调用方系统编码, 参考：新架构平台编码规范 fcode
tcode | String  | 是 | 接收请求方系统编码 | 被调用方编码, 参考：新架构平台编码规范tcode
method_code | String  | 是 | 方法编码  | 收到请求的接口方法编码, 参考：新架构平台编码规范 method
http_method | String  | 是 | 请求方式  | GET/POST/PUT/DELETE
request_time  | String  | 是 | 请求时间  | 
request_headers | String  | 是 | 请求头 |
request_payload |  String  | 是  | 请求体参数  |
response_headers  | String  | 否 | 响应头 |
response_payload  | String |  是 |  响应体参数 |
response_time | String  | 是 | 响应时间  |
response_code |String | 否 | 业务级响应码（错误码） | 请参考【新架构平台编码规范】中【四、错误编码规范】
response_remark |String |  否 |  业务级响应描述文字 |
http_status_code  |String | 否 | HTTP 状态码  | HTTP 状态码（常见的状态码比如：200-服务器成功返回网页，404–请求的网页不存在，503–服务不可用等）
total_time  | String |  是 |  请求处理总耗时 |  总耗时时间
key_type  | String  | 否 | 参数类型  |
key_param | String |  否 |  参数值 |  核心参数类型对应的值
traceId | String  | 否 |   |
key_name  | String  | 否  |    |
log_type  | String  | 否 |    |

## 流水日志调用

以内部流水日志为例：

```
type InfoLogFile struct {
	Dialog_type      string `json:"dialog_type"`      //in or out
	Request_time     string `json:"request_time"`     //请求时间
	Response_time    string `json:"response_time"`    //响应时间
	Address          string `json:"address"`          //请求地址
	Http_method      string `json:"http_method"`      //请求方式
	Request_payload  string `json:"request_payload"`  //请求体参数
	Response_payload string `json:"response_payload"` //响应体参数
	Request_headers  string `json:"request_headers"`  //请求头
	Response_headers string `json:"response_headers"` //响应头
	Response_code    string `json:"response_code"`    //业务级响应码（错误码）
	Response_remark  string `json:"response_remark"`  //业务级响应描述文字
	Http_status_code string `json:"http_status_code"` //HTTP 状态码
	Total_time       string `json:"total_time"`       //请求处理总耗时
	Method_code      string `json:"method_code"`      //方法编码
	Transaction_id   string `json:"transaction_id"`   //流水 id
	Key_type         string `json:"key_type"`         //参数类型
	Key_param        string `json:"key_param"`        //参数值
	Fcode            string `json:"fcode"`            //调用方系统编码
	Tcode            string `json:"tcode"`            //接收请求方系统编码
	TraceId          string `json:"traceId"`
	Key_name         string `json:"key_name"`
	Log_type         string `json:"log_type"`
}
```

```
infoLog := InfoLogFile{
			Dialog_type:      "in",
			Request_time:     time.Now().Local().String(),
			Response_time:    time.Now().Local().String(),
			Address:          "localhost",
			Http_method:      "GET",
			Request_payload:  "{bodydata:qqq}",
			Response_payload: "res-pay",
			Request_headers:  "{header:http1.1}",
			Response_headers: "{responseheader:http1.2}",
			Response_code:    "200",
			Response_remark:  "remark",
			Http_status_code: "200",
			Total_time:       "1",
			Method_code:      "get",
			Transaction_id:   "1111",
			Key_type:         "k-t01",
			Key_param:        "k-p:{p1:qq,p2:cc}",
			Fcode:            "fcode001",
			Tcode:            "tcode002",
			TraceId:          "010",
			Key_name:         "test",
			Log_type:         "infolog",
		}
		msg := ""
		jsonBytes, err := json.Marshal(infoLog)
		if err == nil {
			msg = string(jsonBytes)
		}
		extra := `extra={"www":"1","key_name":"其他想打印到日志里的信息"}`
		log.Internal_log(msg, extra)
```
  
## 测试

使用gin框架测试

```
package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	gologger "github.com/yingying0708/GoLogger"
)

var wg sync.WaitGroup
var log *gologger.GoLogHelper

func main() {
	log = gologger.GetGoLogHelper("gin_test_log", "C:\\goginlogs\\").SetConsolePrint(true).SetLogLevel("info").SetWhen("m").SetBackupCount(2)

	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	wg.Add(7)
	r.GET("/healthcheck", Healthcheck)

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")

}

func Healthcheck(c *gin.Context) {
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			time.Sleep(time.Second)
			log.Trace(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 200; i < 400; i++ {
			time.Sleep(time.Second)
			log.Debug(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 400; i < 600; i++ {
			time.Sleep(time.Second)
			log.Info(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 600; i < 800; i++ {
			time.Sleep(time.Second)
			log.Warn(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 800; i < 1000; i++ {
			time.Sleep(time.Second)
			log.Error(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1000; i < 1200; i++ {
			time.Sleep(time.Second)
			extra := `extra={"www":"0009","key_name":"其他想打印到日志里的信息"}`
			log.Error(i, extra)
		}
	}()
	go func() {
		defer wg.Done()
		//内部流水日志
		infoLog := InfoLogFile{
			Dialog_type:      "in",
			Request_time:     time.Now().Local().String(),
			Response_time:    time.Now().Local().String(),
			Address:          "localhost",
			Http_method:      "GET",
			Request_payload:  "{bodydata:qqq}",
			Response_payload: "res-pay",
			Request_headers:  "{header:http1.1}",
			Response_headers: "{responseheader:http1.2}",
			Response_code:    "200",
			Response_remark:  "remark",
			Http_status_code: "200",
			Total_time:       "1",
			Method_code:      "get",
			Transaction_id:   "1111",
			Key_type:         "k-t01",
			Key_param:        "k-p:{p1:qq,p2:cc}",
			Fcode:            "fcode001",
			Tcode:            "tcode002",
			TraceId:          "010",
			Key_name:         "test",
			Log_type:         "infolog",
		}
		msg := ""
		jsonBytes, err := json.Marshal(infoLog)
		if err == nil {
			msg = string(jsonBytes)
		}
		extra := `extra={"www":"1","key_name":"其他想打印到日志里的信息"}`
		log.Internal_log(msg, extra)
	}()
	wg.Wait()
}

type InfoLogFile struct {
	Dialog_type      string `json:"dialog_type"`      //in or out
	Request_time     string `json:"request_time"`     //请求时间
	Response_time    string `json:"response_time"`    //响应时间
	Address          string `json:"address"`          //请求地址
	Http_method      string `json:"http_method"`      //请求方式
	Request_payload  string `json:"request_payload"`  //请求体参数
	Response_payload string `json:"response_payload"` //响应体参数
	Request_headers  string `json:"request_headers"`  //请求头
	Response_headers string `json:"response_headers"` //响应头
	Response_code    string `json:"response_code"`    //业务级响应码（错误码）
	Response_remark  string `json:"response_remark"`  //业务级响应描述文字
	Http_status_code string `json:"http_status_code"` //HTTP 状态码
	Total_time       string `json:"total_time"`       //请求处理总耗时
	Method_code      string `json:"method_code"`      //方法编码
	Transaction_id   string `json:"transaction_id"`   //流水 id
	Key_type         string `json:"key_type"`         //参数类型
	Key_param        string `json:"key_param"`        //参数值
	Fcode            string `json:"fcode"`            //调用方系统编码
	Tcode            string `json:"tcode"`            //接收请求方系统编码
	TraceId          string `json:"traceId"`
	Key_name         string `json:"key_name"`
	Log_type         string `json:"log_type"`
}

```

### 输出结果
#### gin_test_log_p1_error.log.202112201110.log

内容：

{"app_name":"gin_test_log_code","funcName":"Healthcheck","hostname":"DESKTOP-0JRA26N","key_name":"其他想打印到日志里的信息","level":"ERROR","lineno":76,"log_time":"2021-12-20 11:10:44","logger":"C:/GitTest/DEMO/main.go","module":"main","msg":1001,"www":"0009"}


{"logger":"C:/GitTest/DEMO/main.go","lineno":68,"app_name":"gin_test_log_code","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 11:10:44","hostname":"DESKTOP-0JRA26N","level":"ERROR","msg":801}

#### gin_test_log_p1_info.log.202112201110.log

内容（包含内部流水日志和普通的info日志）：

{"address":"localhost","app_name":"gin_test_log_info","dialog_type":"in","fcode":"fcode001","funcName":"Healthcheck","hostname":"DESKTOP-0JRA26N","http_method":"GET","http_status_code":"200","key_name":"其他想打印到日志里的信息","key_param":"k-p:{p1:qq,p2:cc}","key_type":"k-t01","level":"INFO","lineno":112,"log_time":"2021-12-20 11:10:42","log_type":"infolog","logger":"C:/GitTest/DEMO/main.go","method_code":"get","module":"main","msg":"","request_headers":"{header:http1.1}","request_payload":"{bodydata:qqq}","request_time":"2021-12-20 11:10:42.1730631 +0800 CST","response_code":"200","response_headers":"{responseheader:http1.2}","response_payload":"res-pay","response_remark":"remark","response_time":"2021-12-20 11:10:42.2710584 +0800 CST","tcode":"tcode002","total_time":"1","traceId":"010","transaction_id":"1111","www":"1"}


{"logger":"C:/GitTest/DEMO/main.go","lineno":54,"app_name":"gin_test_log_info","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 11:10:43","hostname":"DESKTOP-0JRA26N","level":"INFO","msg":400}


#### gin_test_log_p1_warn.log.202112201110.log

内容：

{"logger":"C:/GitTest/DEMO/main.go","lineno":61,"app_name":"gin_test_log_code","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 11:10:43","hostname":"DESKTOP-0JRA26N","level":"WARN","msg":600}

根据测试相应的日志写入相应的文件中

控制台日志内容：

{"app_name":"gin_test_log_code","fields.level":"WARN","funcName":"com/yingying0708/GoLogger","hostname":"DESKTOP-0JRA26N","level":"info","lineno":150,"log_time":"2021-12-20 11:10:54","logger":"C:/Users/59231/go/pkg/mod/github.com/yingying0708/!go!logger@v0.0.1/GoLogger.go","module":"github","msg":"611","time":"2021-12-20 11:11:10"}

{"app_name":"gin_test_log_code","fields.level":"ERROR","funcName":"com/yingying0708/GoLogger","hostname":"DESKTOP-0JRA26N","key_name":"其他想打印到日志里
的信息","level":"info","lineno":158,"log_time":"2021-12-20 11:10:55","logger":"C:/Users/59231/go/pkg/mod/github.com/yingying0708/!go!logger@v0.0.1/GoLogger.go","module":"github","msg":"1012","time":"2021-12-20 11:11:10","www":"0009"}

{"app_name":"gin_test_log_code","fields.level":"ERROR","funcName":"com/yingying0708/GoLogger","hostname":"DESKTOP-0JRA26N","level":"info","lineno":158,"log_time":"2021-12-20 11:10:55","logger":"C:/Users/59231/go/pkg/mod/github.com/yingying0708/!go!logger@v0.0.1/GoLogger.go","module":"github","msg":"812","time":"2021-12-20 11:11:10"}

{"app_name":"gin_test_log_code","fields.level":"WARN","funcName":"com/yingying0708/GoLogger","hostname":"DESKTOP-0JRA26N","level":"info","lineno":150,"log_time":"2021-12-20 11:10:55","logger":"C:/Users/59231/go/pkg/mod/github.com/yingying0708/!go!logger@v0.0.1/GoLogger.go","module":"github","msg":"612","time":"2021-12-20 11:11:10"}

{"app_name":"gin_test_log_info","fields.level":"INFO","funcName":"com/yingying0708/GoLogger","hostname":"DESKTOP-0JRA26N","level":"info","lineno":142,"log_time":"2021-12-20 11:10:55","logger":"C:/Users/59231/go/pkg/mod/github.com/yingying0708/!go!logger@v0.0.1/GoLogger.go","module":"github","msg":"412","time":"2021-12-20 11:11:10"}


## 简介
此项目是一个开源的golang的日志功能项目，go的版本是1.17，此项目建议在go版本是1.14及以上版本中使用

## 获取
go get github.com/yingying0708/GoLogger

## 使用
```
log := gologhelper.GetLogHelper("项目名称", "日志路径").SetConsolePrint(true).SetLogLeve("info").SetWhen("m").SetBackupCount(2)

log.Info("info")

log.Trace("trace")

log.Warn("warn")

log.Error("error")

extra := Extra{
	K_N: "k_error",
	K_V: 0,
}

log.Error("msg", extra)

//内部流水日志
infoLogIn := gologgerModel.InfoLogFile{} //实例化流水日志结构体
extraIn :=  Extra{
	K_N: "k_info",
	K_V: 1,
}
log.Internal_log(infoLogIn, extraIn)

//外部流水日志
infoLogEx := gologgerModel.InfoLogFile{} //实例化流水日志结构体
extraEx :=  Extra{
	K_N: "k_info",
	K_V: 2,
}
log.External_log(infoLogEx, extraEx)

```
输出结果：
```
{"logger":"C:/GitTest/DEMO/main.go","lineno":54,"app_name":"gin_test_log_info","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 14:40:59","hostname":"DESKTOP-0JRA26N","level":"INFO","msg":"info"}

{"address":"localhost","app_name":"gin_test_log_info","dialog_type":"in","fcode":"fcode001","funcName":"Healthcheck","hostname":"DESKTOP-0JRA26N","http_method":"GET","http_status_code":"200","k_N":"k_info","key_param":"k-p:{p1:qq,p2:cc}","key_type":"k-t01","level":"INFO","lineno":113,"log_time":"2021-12-20 15:22:42","log_type":"infolog","logger":"C:/GitTest/DEMO/main.go","method_code":"get","module":"main","msg":"","request_headers":"{header:http1.1}","request_payload":"{bodydata:qqq}","request_time":"2021-12-20 15:22:42.4934171 +0800 CST","response_code":"200","response_headers":"{responseheader:http1.2}","response_payload":"res-pay","response_remark":"remark","response_time":"2021-12-20 15:22:42.5909523 +0800 CST","tcode":"tcode002","total_time":"1","traceId":"010","transaction_id":"1111","K_V":"1"}

{"address":"localhost","app_name":"gin_test_log_info","dialog_type":"out","fcode":"fcode001","funcName":"Healthcheck","hostname":"DESKTOP-0JRA26N","http_method":"GET","http_status_code":"200","k_N":"k_info","key_param":"k-p:{p1:qq,p2:cc}","key_type":"k-t01","level":"INFO","lineno":183,"log_time":"2021-12-20 15:22:42","log_type":"infolog","logger":"C:/GitTest/DEMO/main.go","method_code":"get","module":"main","msg":"","request_headers":"{header:http1.1}","request_payload":"{bodydata:qqq}","request_time":"2021-12-20 15:22:42.4934171 +0800 CST","response_code":"200","response_headers":"{responseheader:http1.2}","response_payload":"res-pay","response_remark":"remark","response_time":"2021-12-20 15:22:42.5909523 +0800 CST","tcode":"tcode002","total_time":"1","traceId":"010","transaction_id":"1111","K_V":"2"}

{"logger":"C:/GitTest/DEMO/main.go","lineno":40,"app_name":"gin_test_log_trace","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 14:40:17","hostname":"DESKTOP-0JRA26N","level":"TRACE","msg":"trace"}

{"logger":"C:/GitTest/DEMO/main.go","lineno":61,"app_name":"gin_test_log_code","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 14:40:59","hostname":"DESKTOP-0JRA26N","level":"WARN","msg":"warn"}

{"logger":"C:/GitTest/DEMO/main.go","lineno":68,"app_name":"gin_test_log_code","module":"main","funcName":"Healthcheck","log_time":"2021-12-20 14:41:00","hostname":"DESKTOP-0JRA26N","level":"ERROR","msg":"error"}

{"app_name":"gin_test_log_code","funcName":"Healthcheck","hostname":"DESKTOP-0JRA26N","K_N":"k_error","level":"ERROR","lineno":76,"log_time":"2021-12-20 14:41:00","logger":"C:/GitTest/DEMO/main.go","module":"main","msg":"msg","K_V":0}

```

## 方法介绍
SetConsolePrint(true)  控制台打印，默认是false


SetLogLevel("info")    只生成info的日志文件，默认是error，日志级别从小到大排列是：track，debug，info，warn，error。例：设置info ，会打印info、warn、error的类型的日志


SetWhen("m")           设置按什么进行分割，默认是D，D：天，H：小时，m：分钟


SetBackupCount(2)      保留多少个文件，默认是15
  
## 线程安全

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
	gologgerModel "github.com/yingying0708/GoLogger/LogModel"
)

var log *gologger.GoLogHelper

func main() {
	log = gologger.GetGoLogHelper("gin_test_log", "C:\\goginlogs\\").SetConsolePrint(true).SetLogLevel("trace").SetWhen("m").SetBackupCount(2)

	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	r.GET("/healthcheck", Healthcheck)

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")

}

func Healthcheck(c *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(17)
	go func() {
		defer wg.Done()
		log.Trace("trace")
	}()
	go func() {
		defer wg.Done()
		log.Warn("warn")
	}()
	go func() {
		defer wg.Done()
		log.Info("info")
	}()
	go func() {
		defer wg.Done()
		log.Debug("debug")
	}()
	go func() {
		defer wg.Done()
		log.Error("error")
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			time.Sleep(time.Second)
			log.Trace(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1400; i < 1600; i++ {
			extra := Extra01{
				K_N: "k_trace",
				K_V: i,
			}
			log.Trace(extra)
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
		extra := `input={"t1:","t1"}`
		log.Info(extra)
	}()
	go func() {
		defer wg.Done()
		extra := `{"t1:","t1"}`
		log.Info(extra)
	}()
	go func() {
		defer wg.Done()
		for i := 1600; i < 1800; i++ {
			extra := Extra01{
				K_N: "k_info",
				K_V: i,
			}
			log.Info(i, extra)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 2600; i < 2800; i++ {
			extra := Extra01{
				K_N: "k_info",
				K_V: i,
			}
			log.Info(extra)
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
			extra01 := Extra01{
				K_N: "k_error",
				K_V: i,
			}
			log.Error(i, extra01)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1200; i < 1400; i++ {
			//内部流水日志
			infoLog := gologgerModel.InfoLogFile{
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
			extra := Extra{
				Key_Name0000:   "key00001",
				Text_Name0001:  "text00001",
				Value_Test0001: i,
			}
			extra01 := Extra01{
				K_N: "k_n",
				K_V: i,
			}
			extra02 := Extra01{
				K_N: "k_n02",
				K_V: i,
			}
			//后面的结构体会覆盖前面的结构体的数据
			log.Internal_log(infoLog, extra, extra01, extra02)
		}
	}()
	wg.Wait()
}

type Extra struct {
	Key_Name0000   string
	Text_Name0001  string
	Value_Test0001 int
}

type Extra01 struct {
	K_N string
	K_V int
}

```
需验证各类型的日志存在对应的文件中且对应的日志文件中仅有对应类型的日志信息，使用上面的代码进行测试，启动程序，在浏览器输入http://localhost:8000/healthcheck ，得到的日志文件上传到TestResultFiles文件夹中，经查看生成的日志，info类型的日志存在info的文件中，error类型的日志存在error的文件中，各类型的日志存在对应的文件中，且对应生成的日志文件中并未有其他类型的日志信息存在，符合预期，故推断此程序是线程安全的。
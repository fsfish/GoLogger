package loghelper

import (
	model "github.com/fsfish/GoLogger/LogModel"
	"os"
	"runtime"
	"strings"
	"time"
)

func GetPrintLogFile(appname, level string, msg interface{}) interface{} {
	pc, file, line, _ := runtime.Caller(3)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		return `{"msg":"获取hostname失败"}`
	}
	module := strings.Split(f.Name(), ".")[0]
	funcName := strings.Split(f.Name(), ".")[1]
	log_time := time.Now().Format("2006-01-02 15:04:05.000")

	entity := model.LogFile{
		Logger:   file,
		FuncName: funcName,
		LineNo:   line,
		App_Name: appname,
		Module:   module,
		Log_Time: log_time,
		HOSTNAME: hostname,
		Level:    level,
		Msg:      msg,
	}
	return entity
}

package loghelper

import (
	"os"
	"runtime"
	"strings"
	"time"
)

func GetPrintLogFileCustom(appname, level string, msg interface{}, fields map[string]interface{}) interface{} {
	pc, file, line, _ := runtime.Caller(3)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		return `{"msg":"获取hostname失败"}`
	}
	module := strings.Split(f.Name(), ".")[0]
	funcName := strings.Split(f.Name(), ".")[1]
	log_time := time.Now().Format("2006-01-02 15:04:05.000")

	fields["logger"] = file
	fields["lineno"] = line
	fields["app_name"] = appname
	fields["module"] = module
	fields["funcName"] = funcName
	fields["log_time"] = log_time
	fields["hostname"] = hostname
	fields["level"] = level
	fields["msg"] = msg
	return fields
}

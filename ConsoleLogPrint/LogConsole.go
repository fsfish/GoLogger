package loghelper

import (
	"os"
	"runtime"
	"strings"
	"time"
)

func GetPrintLogConsole(appname, level string) map[string]interface{} {
	fields := make(map[string]interface{})
	pc, file, line, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		fields["error_msg"] = "获取hostname失败"
		return fields

	}
	module := strings.Split(f.Name(), ".")[0]
	funcName := strings.Split(f.Name(), ".")[1]
	log_time := time.Now().Format("2006-01-02 15:04:05")
	fields["logger"] = file
	fields["lineno"] = line
	fields["app_name"] = appname
	fields["module"] = module
	fields["funcName"] = funcName
	fields["log_time"] = log_time
	fields["hostname"] = hostname
	fields["level"] = level

	return fields
}

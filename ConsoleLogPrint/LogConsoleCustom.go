package loghelper

import (
	common "github.com/yingying0708/GoLogger/LogCommon"
	"os"
	"runtime"
	"strings"
	"time"
)

func GetPrintLogConsoleCustom(appname, level string, extra map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	pc, file, line, _ := runtime.Caller(3)
	f := runtime.FuncForPC(pc)
	hostname, err := os.Hostname()
	if err != nil {
		fields["error_msg"] = "获取hostname失败"
		return fields
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
	fields = common.MergeMap(extra, fields)
	return fields
}

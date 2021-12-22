package loghelper

// 日志级别
const (
	Log_Trace = iota
	Log_Debug
	Log_Info
	Log_Warn
	Log_Error
)

// 日志级别字段
const LevelTrace string = "TRACE"
const LevelDebug string = "DEBUG"
const LevelInfo string = "INFO"
const LevelWarn string = "WARN"
const LevelError string = "ERROR"

// 根据输入的日志级别，返回匹配的自定义常数
func GetLogLevel(loglevel string) int {
	res := Log_Error
	if loglevel == "info" {
		res = Log_Info
	}
	if loglevel == "warn" {
		res = Log_Warn
	}
	if loglevel == "debug" {
		res = Log_Debug
	}
	if loglevel == "trace" {
		res = Log_Trace
	}
	if loglevel == "error" {
		res = Log_Error
	}
	return res
}

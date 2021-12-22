package loghelper

import (
	"strings"
)

// 日志帮助类
type LogHelper struct {
	AppName      string
	ConsolePrint bool
	LogPath      string
	BackupCount  int
	When         string
	LogLevel     string
}

func GetLogHelper(app_name, log_path, leverstr string) *LogHelper {
	return &LogHelper{
		AppName:      app_name,
		LogPath:      log_path,
		ConsolePrint: false,
		BackupCount:  15,
		When:         "D",
		LogLevel:     leverstr,
	}
}

//设置when(D:天，H：小时，M：分钟，默认是D)
func (log *LogHelper) SetWhen(when string) *LogHelper {
	if when != "" {
		log.When = strings.ToUpper(when)
	}
	return log
}

//设置日志级别（error,debug,info,trace,warn 默认是error）
func (log *LogHelper) SetLogLevel(level string) *LogHelper {
	levelstr := "error"
	if level != "" {
		levelstr = strings.ToLower(level)
	}
	log.LogLevel = levelstr
	return log
}

//设置是否控制台打印默认是false
func (log *LogHelper) SetConsolePrint(isPrint bool) *LogHelper {
	log.ConsolePrint = isPrint
	return log
}

//设置多少个文件后进行回滚操作默认是15
func (log *LogHelper) SetBackupCount(backupCount int) *LogHelper {
	if backupCount > 0 {
		log.BackupCount = backupCount
	}
	return log
}

//根据规则生成日志的appname与项目的appname不一样
func GetAppName(name, level string) string {
	appname := name + "_code"
	levelstr := strings.ToLower(level)
	if levelstr == "trace" {
		appname = name + "_trace"
	}
	if levelstr == "info" {
		appname = name + "_info"
	}
	return appname
}


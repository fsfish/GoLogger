package loghelper

import (
	"encoding/json"
	"time"

	consolehelper "github.com/yingying0708/GoLogger/ConsoleLogPrint"
	filehelper "github.com/yingying0708/GoLogger/FileLogPrint"
	common "github.com/yingying0708/GoLogger/LogCommon"
	service "github.com/yingying0708/GoLogger/LogService"

	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var logs = logrus.New()
var logsConsole = logrus.New()

//初始化
func init() {
	//设置日志格式为json
	logs.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
	})
	logsConsole.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:03:04",
	})
}

type GoLogHelper struct {
	errorLog         *service.LogHelper
	traceLog         *service.LogHelper
	warnLog          *service.LogHelper
	infoLog          *service.LogHelper
	debugLog         *service.LogHelper
	LogLevel         string
	ErrorWriter      *rotatelogs.RotateLogs
	RotateLogsWriter *rotatelogs.RotateLogs
}

func GetGoLogHelper(app_name, log_path string) *GoLogHelper {
	return &GoLogHelper{
		errorLog: service.GetLogHelper(app_name, log_path, strings.ToLower(common.LevelError)),
		traceLog: service.GetLogHelper(app_name, log_path, strings.ToLower(common.LevelTrace)),
		warnLog:  service.GetLogHelper(app_name, log_path, strings.ToLower(common.LevelWarn)),
		infoLog:  service.GetLogHelper(app_name, log_path, strings.ToLower(common.LevelInfo)),
		debugLog: service.GetLogHelper(app_name, log_path, strings.ToLower(common.LevelDebug)),
	}
}

//设置writer
func getWriter(log *service.LogHelper) *rotatelogs.RotateLogs {
	logPath := log.LogPath + log.AppName + "_p1_" + log.LogLevel + ".log"
	writer, _ := rotatelogs.New(
		logPath+".%Y%m%d.log",
		rotatelogs.WithRotationCount(uint(log.BackupCount)),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if log.When == "H" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H.log",
			rotatelogs.WithRotationCount(uint(log.BackupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
		)
	}
	if log.When == "M" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H%M.log",
			rotatelogs.WithRotationCount(uint(log.BackupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		)
	}
	return writer
}

//设置日志级别（error,debug,info,trace,warn 默认是error）
func (log *GoLogHelper) SetLogLevel(level string) *GoLogHelper {
	levelstr := "error"
	if level != "" {
		levelstr = strings.ToLower(level)
	}
	log.LogLevel = levelstr
	return log
}

//设置是否控制台打印默认是false
func (log *GoLogHelper) SetConsolePrint(isPrint bool) *GoLogHelper {
	log.errorLog.SetConsolePrint(isPrint)
	log.debugLog.SetConsolePrint(isPrint)
	log.infoLog.SetConsolePrint(isPrint)
	log.traceLog.SetConsolePrint(isPrint)
	log.warnLog.SetConsolePrint(isPrint)
	return log
}

//设置多少个文件后进行回滚操作默认是15
func (log *GoLogHelper) SetBackupCount(backupCount int) *GoLogHelper {
	if backupCount > 0 {
		log.errorLog.SetBackupCount(backupCount)
		log.debugLog.SetBackupCount(backupCount)
		log.infoLog.SetBackupCount(backupCount)
		log.traceLog.SetBackupCount(backupCount)
		log.warnLog.SetBackupCount(backupCount)
	}
	return log
}

//设置when(D:天，H：小时，M：分钟，默认是D)
func (log *GoLogHelper) SetWhen(when string) *GoLogHelper {
	if when != "" {
		log.errorLog.SetWhen(when)
		log.debugLog.SetWhen(when)
		log.infoLog.SetWhen(when)
		log.traceLog.SetWhen(when)
		log.warnLog.SetWhen(when)
	}
	return log
}

//Trace
func (log *GoLogHelper) Trace(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Trace {
		msg, extra := getParams(param...)
		log.printTraceLog(log.traceLog, msg, extra)
	}
}

//Debug
func (log *GoLogHelper) Debug(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Debug {
		msg, extra := getParams(param...)
		log.printDebugLog(log.debugLog, msg, extra)
	}
}

//Info
func (log *GoLogHelper) Info(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Info {
		msg, extra := getParams(param...)
		log.printInfoLog(log.infoLog, msg, extra)
	}
}

//Warn
func (log *GoLogHelper) Warn(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Warn {
		msg, extra := getParams(param...)
		log.printWarnLog(log.warnLog, msg, extra)
	}
}

//Error
func (log *GoLogHelper) Error(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Error {
		msg, extra := getParams(param...)
		log.printErrorLog(log.errorLog, msg, extra)
	}
}

//内部流水日志
func (log *GoLogHelper) Internal_log(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra == nil {
		extra = ConvertIToM(msg)
	}
	if msg != nil && extra != nil {
		extra = common.MergeMap(ConvertIToM(msg), extra)
	}
	if extra != nil {
		if validateRequireItem(extra) {
			log.printInfoLog(log.infoLog, "", extra)
		} else {
			log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
		}
	} else {
		log.printInfoLog(log.infoLog, msg, nil)
	}
}

//外部流水日志
func (log *GoLogHelper) External_log(param ...interface{}) {
	msg, extra := getParams(param...)
	if msg != nil && extra == nil {
		extra = ConvertIToM(msg)
	}
	if msg != nil && extra != nil {
		extra = common.MergeMap(ConvertIToM(msg), extra)
	}
	if extra != nil {
		if validateRequireItem(extra) {
			log.printInfoLog(log.infoLog, "", extra)
		} else {
			log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
		}
	} else {
		log.printInfoLog(log.infoLog, msg, nil)
	}
}

//打印
func (log *GoLogHelper) printErrorLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
	var entity interface{}
	levelstr := common.LevelError
	appname := service.GetAppName(logModel.AppName, levelstr)
	//控制台打印
	if logModel.ConsolePrint {
		var fields map[string]interface{}
		if extra != nil {
			fields = consolehelper.GetPrintLogConsoleCustom(appname, levelstr, extra)
		} else {
			fields = consolehelper.GetPrintLogConsole(appname, levelstr)
		}
		logsConsole.WithFields(fields).Println(msg)
	}
	//打印到文件
	if extra != nil {
		entity = filehelper.GetPrintLogFileCustom(appname, levelstr, msg, extra)
	} else {
		entity = filehelper.GetPrintLogFile(appname, levelstr, msg)
	}
	writer := getWriter(logModel)
	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	logs.SetOutput(writer)
}
func (log *GoLogHelper) printWarnLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
	var entity interface{}
	levelstr := common.LevelWarn
	appname := service.GetAppName(logModel.AppName, levelstr)
	//控制台打印
	if logModel.ConsolePrint {
		var fields map[string]interface{}
		if extra != nil {
			fields = consolehelper.GetPrintLogConsoleCustom(appname, levelstr, extra)
		} else {
			fields = consolehelper.GetPrintLogConsole(appname, levelstr)
		}
		logsConsole.WithFields(fields).Println(msg)
	}
	//打印到文件
	if extra != nil {
		entity = filehelper.GetPrintLogFileCustom(appname, levelstr, msg, extra)
	} else {
		entity = filehelper.GetPrintLogFile(appname, levelstr, msg)
	}
	writer := getWriter(logModel)
	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	logs.SetOutput(writer)
}
func (log *GoLogHelper) printTraceLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
	var entity interface{}
	levelstr := common.LevelTrace
	appname := service.GetAppName(logModel.AppName, levelstr)
	//控制台打印
	if logModel.ConsolePrint {
		var fields map[string]interface{}
		if extra != nil {
			fields = consolehelper.GetPrintLogConsoleCustom(appname, levelstr, extra)
		} else {
			fields = consolehelper.GetPrintLogConsole(appname, levelstr)
		}
		logsConsole.WithFields(fields).Println(msg)
	}
	//打印到文件
	if extra != nil {
		entity = filehelper.GetPrintLogFileCustom(appname, levelstr, msg, extra)
	} else {
		entity = filehelper.GetPrintLogFile(appname, levelstr, msg)
	}
	writer := getWriter(logModel)
	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	logs.SetOutput(writer)
}
func (log *GoLogHelper) printDebugLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
	var entity interface{}
	levelstr := common.LevelDebug
	appname := service.GetAppName(logModel.AppName, levelstr)
	//控制台打印
	if logModel.ConsolePrint {
		var fields map[string]interface{}
		if extra != nil {
			fields = consolehelper.GetPrintLogConsoleCustom(appname, levelstr, extra)
		} else {
			fields = consolehelper.GetPrintLogConsole(appname, levelstr)
		}
		logsConsole.WithFields(fields).Println(msg)
	}
	//打印到文件
	if extra != nil {
		entity = filehelper.GetPrintLogFileCustom(appname, levelstr, msg, extra)
	} else {
		entity = filehelper.GetPrintLogFile(appname, levelstr, msg)
	}
	writer := getWriter(logModel)
	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	logs.SetOutput(writer)
}
func (log *GoLogHelper) printInfoLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
	var entity interface{}
	levelstr := common.LevelInfo
	appname := service.GetAppName(logModel.AppName, levelstr)
	//控制台打印
	if logModel.ConsolePrint {
		var fields map[string]interface{}
		if extra != nil {
			fields = consolehelper.GetPrintLogConsoleCustom(appname, levelstr, extra)
		} else {
			fields = consolehelper.GetPrintLogConsole(appname, levelstr)
		}
		logsConsole.WithFields(fields).Println(msg)
	}
	//打印到文件
	if extra != nil {
		entity = filehelper.GetPrintLogFileCustom(appname, levelstr, msg, extra)
	} else {
		entity = filehelper.GetPrintLogFile(appname, levelstr, msg)
	}
	writer := getWriter(logModel)
	if res, err := json.Marshal(&entity); err == nil {
		writer.Write(res)
		writer.Write([]byte("\n"))
	}
	logs.SetOutput(writer)
}

//拆分参数(第一个参数默认是msg，后面的参数都是附加参数)
func getParams(param ...interface{}) (interface{}, map[string]interface{}) {
	if len(param) > 0 {
		msg := getParamMsg(param[0])
		//只有一个参数
		if len(param) == 1 {
			return msg, nil
		}
		//循环后面的参数
		var extraList []map[string]interface{}
		for i := 1; i < len(param); i++ {
			item := ConvertIToM(param[i])
			if item != nil {
				extraList = append(extraList, item)
			}
		}
		//多个集合整合
		var extra map[string]interface{}
		if len(extraList) > 0 {
			extra = common.MergeMap(extraList...)
		}
		return msg, extra
	}
	return nil, nil
}

//流水日志必填项的校验
func validateRequireItem(fields map[string]interface{}) (isFlag bool) {
	if len(fields) <= 0 {
		isFlag = false
		return
	}
	if _, ok := fields["transaction_id"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["address"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["fcode"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["tcode"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["method_code"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["http_method"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["request_time"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["request_headers"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["request_payload"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["response_payload"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["response_time"]; !ok {
		isFlag = false
		return
	}
	if _, ok := fields["total_time"]; !ok {
		isFlag = false
		return
	}
	isFlag = true
	return
}

//获取参数Msg
func getParamMsg(msg interface{}) interface{} {
	msgResult := ConvertIToM(msg)
	if msgResult == nil {
		return msg
	}
	return msgResult
}

//将interface{}转为map[string]interface{}
func ConvertIToM(msg interface{}) map[string]interface{} {
	var mapResult map[string]interface{}
	err := mapstructure.Decode(msg, &mapResult)
	if err != nil {
		return nil
	}
	return mapResult
}

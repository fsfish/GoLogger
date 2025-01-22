package loghelper

import (
	"encoding/json"
	"strconv"
	"time"

	consolehelper "github.com/fsfish/GoLogger/ConsoleLogPrint"
	filehelper "github.com/fsfish/GoLogger/FileLogPrint"
	common "github.com/fsfish/GoLogger/LogCommon"
	model "github.com/fsfish/GoLogger/LogModel"
	service "github.com/fsfish/GoLogger/LogService"

	"strings"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var logs = logrus.New()
var logsConsole = logrus.New()

// 初始化
func init() {
	//设置日志格式为json
	logs.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logsConsole.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
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

// 设置writer，返回分割日志实例
func getWriter(log *service.LogHelper) *rotatelogs.RotateLogs {
	logPath := log.LogPath + log.AppName + "_p1_" + log.LogLevel + ".log"
	//按天：D生成
	writer, _ := rotatelogs.New(
		logPath+".%Y%m%d.log",
		rotatelogs.WithRotationCount(uint(log.BackupCount)),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	//按小时：H生成
	if log.When == "H" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H.log",
			rotatelogs.WithRotationCount(uint(log.BackupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Minute),
		)
	}
	//按分钟：M生成
	if log.When == "M" {
		writer, _ = rotatelogs.New(
			logPath+".%Y%m%d%H%M.log",
			rotatelogs.WithRotationCount(uint(log.BackupCount)),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		)
	}
	return writer
}

// 设置日志级别（error,debug,info,trace,warn 默认是error）
func (log *GoLogHelper) SetLogLevel(level string) *GoLogHelper {
	levelstr := "error"
	if level != "" {
		levelstr = strings.ToLower(level)
	}
	log.LogLevel = levelstr
	return log
}

// 设置是否控制台打印默认是false
func (log *GoLogHelper) SetConsolePrint(isPrint bool) *GoLogHelper {
	log.errorLog.SetConsolePrint(isPrint)
	log.debugLog.SetConsolePrint(isPrint)
	log.infoLog.SetConsolePrint(isPrint)
	log.traceLog.SetConsolePrint(isPrint)
	log.warnLog.SetConsolePrint(isPrint)
	return log
}

// 设置多少个文件后进行回滚操作默认是15
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

// 设置when(D:天，H：小时，M：分钟，默认是D)
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

// Trace
func (log *GoLogHelper) Trace(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Trace {
		msg, extra := getParams(param...)
		log.printTraceLog(log.traceLog, msg, extra)
	}
}

// Debug
func (log *GoLogHelper) Debug(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Debug {
		msg, extra := getParams(param...)
		log.printDebugLog(log.debugLog, msg, extra)
	}
}

// Info
func (log *GoLogHelper) Info(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Info {
		msg, extra := getParams(param...)
		log.printInfoLog(log.infoLog, msg, extra)
	}
}

// Warn
func (log *GoLogHelper) Warn(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Warn {
		msg, extra := getParams(param...)
		log.printWarnLog(log.warnLog, msg, extra)
	}
}

// Error
func (log *GoLogHelper) Error(param ...interface{}) {
	if common.GetLogLevel(log.LogLevel) <= common.Log_Error {
		msg, extra := getParams(param...)
		log.printErrorLog(log.errorLog, msg, extra)
	}
}

// 内部流水日志msg参数应为结构体
func (log *GoLogHelper) Internal_log(param ...interface{}) {
	msg, extra := getParams(param...)
	//msg为流水日志结构体
	msgS := ConvertMToS(msg)
	if msgS == nil {
		log.printInfoLog(log.infoLog, "流水日志中第一个参数必须为流水日志的结构体，请重新填写", nil)
	} else {
		//验证必填项
		if isValidate(msgS) {
			if msg != nil && extra == nil {
				extra = ConvertIToM(msg)
			}
			if msg != nil && extra != nil {
				extra = common.MergeMap(ConvertIToM(msg), extra)
			}
			if extra != nil {
				if validateRequireItem(extra) {
					log.printInfoLog(log.infoLog, extra["msg"], extra)
				} else {
					log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
				}
			} else {
				log.printInfoLog(log.infoLog, msg, nil)
			}
		} else {
			log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
		}
	}
}

// 外部流水日志
func (log *GoLogHelper) External_log(param ...interface{}) {
	msg, extra := getParams(param...)
	//msg为流水日志结构体
	msgS := ConvertMToS(msg)
	if msgS == nil {
		log.printInfoLog(log.infoLog, "流水日志中第一个参数必须为流水日志的结构体，请重新填写", nil)
	} else {
		//验证必填项
		if isValidate(msgS) {
			if msg != nil && extra == nil {
				extra = ConvertIToM(msg)
			}
			if msg != nil && extra != nil {
				extra = common.MergeMap(ConvertIToM(msg), extra)
			}
			if extra != nil {
				if validateRequireItem(extra) {
					log.printInfoLog(log.infoLog, extra["msg"], extra)
				} else {
					log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
				}
			} else {
				log.printInfoLog(log.infoLog, msg, nil)
			}
		} else {
			log.printInfoLog(log.infoLog, "流水日志中必填项没有传入完全，请核查必填项", nil)
		}
	}
}

// 流水日志msg参数应为结构体
func (log *GoLogHelper) ServiceLog(param ...interface{}) {
	msg, _ := getParams(param...)
	log.serviceInfoLog(log.infoLog, msg, nil)
}

// 打印
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

func (log *GoLogHelper) serviceInfoLog(logModel *service.LogHelper, msg interface{}, extra map[string]interface{}) {
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
	if str, ok := msg.(string); ok {
	        writer := getWriter(logModel)
		writer.Write(string(msg))
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

// 拆分参数(第一个参数默认是msg，后面的参数都是附加参数)
func getParams(param ...interface{}) (interface{}, map[string]interface{}) {
	if len(param) > 0 {
		// msg参数
		msg := getParamMsg(param[0])
		//只有一个参数
		if len(param) == 1 {
			return msg, nil
		}
		//循环后面的参数(msg后的额外参数)
		var extraList []map[string]interface{}
		//不可转换的参数排序使用
		orderNumber := 1
		for i := 1; i < len(param); i++ {
			item := ConvertIToM(param[i])
			if item != nil {
				extraList = append(extraList, item)
			} else {
				//不可以转换的参数处理
				extraListUnStructItem := make(map[string]interface{}, 1)
				key := "extraParam" + strconv.Itoa(orderNumber)
				extraListUnStructItem[key] = param[i]
				//将不可转换的参数存到集合中
				extraList = append(extraList, extraListUnStructItem)
				//序号加1处理
				orderNumber = orderNumber + 1
			}
		}
		//多个集合整合(将额外参数整合到一起)
		var extra map[string]interface{}
		if len(extraList) > 0 {
			extra = common.MergeMap(extraList...)
		}
		//返回msg参数和额外参数
		return msg, extra
	}
	return nil, nil
}

// 流水日志必填项的校验
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
	if _, ok := fields["msg"]; !ok {
		isFlag = false
		return
	}
	isFlag = true
	return
}

// 获取参数Msg
func getParamMsg(msg interface{}) interface{} {
	//接收msg参数，调用ConvertIToM函数，如果返回的是nil，说明不可以转换，使用传入的msg
	msgResult := ConvertIToM(msg)
	if msgResult == nil {
		return msg
	}
	return msgResult
}

// 将interface{}转为map[string]interface{}
func ConvertIToM(msg interface{}) map[string]interface{} {
	//如果不可以转换返回nil，不处理异常，转换可以成功，返回转换后的内容
	var mapResult map[string]interface{}
	err := mapstructure.Decode(msg, &mapResult)
	if err != nil {
		return nil
	}
	return mapResult
}

// 将map[string]interface{}转为结构体
func ConvertMToS(field interface{}) *model.InfoLogFile {
	var mapResult *model.InfoLogFile
	err := mapstructure.Decode(field, &mapResult)
	if err != nil {
		return nil
	}
	return mapResult
}

// 校验必填内容是否都填
func isValidate(entity *model.InfoLogFile) (isFlag bool) {
	if entity == nil {
		return isFlag
	}
	if len(entity.Transaction_id) < 1 {
		isFlag = false
		return
	}
	if len(entity.Address) < 1 {
		isFlag = false
		return
	}
	if len(entity.Fcode) < 1 {
		isFlag = false
		return
	}
	if len(entity.Tcode) < 1 {
		isFlag = false
		return
	}
	if len(entity.Method_code) < 1 {
		isFlag = false
		return
	}
	if len(entity.Http_method) < 1 {
		isFlag = false
		return
	}
	if len(entity.Request_time) < 1 {
		isFlag = false
		return
	}
	if len(entity.Request_headers) < 1 {
		isFlag = false
		return
	}
	if len(entity.Request_payload) < 1 {
		isFlag = false
		return
	}
	if len(entity.Response_payload) < 1 {
		isFlag = false
		return
	}
	if len(entity.Response_time) < 1 {
		isFlag = false
		return
	}
	if len(entity.Total_time) < 1 {
		isFlag = false
		return
	}
	if len(entity.Msg) < 1 {
		isFlag = false
		return
	}
	isFlag = true
	return isFlag
}

package loghelper

//打印的数据结构体
type LogFile struct {
	Logger   string      `json:"logger"`
	LineNo   int         `json:"lineno"`
	App_Name string      `json:"app_name"`
	Module   string      `json:"module"`
	FuncName string      `json:"funcName"`
	Log_Time string      `json:"log_time"`
	HOSTNAME string      `json:"hostname"`
	Level    string      `json:"level"`
	Msg      interface{} `json:"msg"`
}

package loghelper


type InfoLogFile struct {
	Dialog_type      string `json:"dialog_type"`      //in or out
	Request_time     string `json:"request_time"`     //请求时间
	Response_time    string `json:"response_time"`    //响应时间
	Address          string `json:"address"`          //请求地址
	Http_method      string `json:"http_method"`      //请求方式
	Request_payload  string `json:"request_payload"`  //请求体参数
	Response_payload string `json:"response_payload"` //响应体参数
	Request_headers  string `json:"request_headers"`  //请求头
	Response_headers string `json:"response_headers"` //响应头
	Response_code    string `json:"response_code"`    //业务级响应码（错误码）
	Response_remark  string `json:"response_remark"`  //业务级响应描述文字
	Http_status_code string `json:"http_status_code"` //HTTP 状态码
	Total_time       string `json:"total_time"`       //请求处理总耗时
	Method_code      string `json:"method_code"`      //方法编码
	Transaction_id   string `json:"transaction_id"`   //流水 id
	Key_type         string `json:"key_type"`         //参数类型
	Key_param        string `json:"key_param"`        //参数值
	Fcode            string `json:"fcode"`            //调用方系统编码
	Tcode            string `json:"tcode"`            //接收请求方系统编码
	TraceId          string `json:"traceId"`
	Key_name         string `json:"key_name"`
	Log_type         string `json:"log_type"`
}

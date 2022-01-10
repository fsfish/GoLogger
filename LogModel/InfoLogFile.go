package loghelper

type InfoLogFile struct {
	Dialog_type      string `mapstructure:"dialog_type"`      //in or out
	Request_time     string `mapstructure:"request_time"`     //请求时间
	Response_time    string `mapstructure:"response_time"`    //响应时间
	Address          string `mapstructure:"address"`          //请求地址
	Http_method      string `mapstructure:"http_method"`      //请求方式
	Request_payload  string `mapstructure:"request_payload"`  //请求体参数
	Response_payload string `mapstructure:"response_payload"` //响应体参数
	Request_headers  string `mapstructure:"request_headers"`  //请求头
	Response_headers string `mapstructure:"response_headers"` //响应头
	Response_code    string `mapstructure:"response_code"`    //业务级响应码（错误码）
	Response_remark  string `mapstructure:"response_remark"`  //业务级响应描述文字
	Http_status_code string `mapstructure:"http_status_code"` //HTTP 状态码
	Total_time       string `mapstructure:"total_time"`       //请求处理总耗时
	Method_code      string `mapstructure:"method_code"`      //方法编码
	Transaction_id   string `mapstructure:"transaction_id"`   //流水 id
	Key_type         string `mapstructure:"key_type"`         //参数类型
	Key_param        string `mapstructure:"key_param"`        //参数值
	Fcode            string `mapstructure:"fcode"`            //调用方系统编码
	Tcode            string `mapstructure:"tcode"`            //接收请求方系统编码
	TraceId          string `mapstructure:"traceId"`
	Key_name         string `mapstructure:"key_name"`
	Log_type         string `mapstructure:"log_type"`
	Msg              string `mapstructure:"msg"`
}

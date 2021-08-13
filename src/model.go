package account_service

type AcStatus struct {
	Status    string	`json:"status"`		//空调状态 1代表运行 0代表停止 2代表有错误 3代表通讯故障
	Address   string	`json:"address"`	//空调地址
	ErrorCode string	`json:"errorCode"`	//错误代码
	OnOff     int32		`json:"onOff"`		//是否开关 1=开 0=关
	RunMode   int32		`json:"runMode"`	//运行模式 1=制冷 2=制热 3=送风 4=除湿
	FanSpeed  int32		`json:"fanSpeed"`	//风速 1=高风 2=中风 4=低风
	TempSet   int32		`json:"tempSet"` 	//设定温度
}

type AcInfo struct {
	Account  string			`json:"account"`	// 登陆账号
	AcStatus []*AcStatus	`json:"acStatus"`	// 飞弈返回设备状态信息
}

type AcInfoResponse struct {
	Code 	int32			`json:"code"`
	Msg 	string			`json:"msg"`
	Data    []*AcInfo		`json:"data"`
}

type AcSetParams struct {
	Account 	string	`json:"account"`	// 空调所属账号
	Action  	string	`json:"action"`		// 控制就传Set 锁定传Lock
	OnOff   	int32	`json:"onOff"`		// 开关 1=开 0=关
	Temp    	int32	`json:"temp"`		// 设置温度
	WorkMode 	int32	`json:"workMode"`	// 模式
	Speed    	int32	`json:"speed"`		// 风速
	SelectedAc 	string	`json:"selectedAc"`	// 空调对应的地址，空调 address 如果是多个用#分隔
}

type AcSetResponse struct {
	Code 	int32			`json:"code"`
	Msg  	string			`json:"msg"`
	Data 	interface{}		`json:"data"`
}

type ElecSumParams struct {
	Account 	string		`json:"account"`
	Address		string		`json:"address"`
	FromDate 	string		`json:"fromDate"`
	ToDate		string		`json:"toDate"`
}

type ElecSum struct {
	UsageSum	float64		`json:"usageSum"`	//该时间内空调内机分摊的总电量，单位度
	Address 	string		`json:"address"`	//地址
	FromDate	float64		`json:"fromDate"`	//空调内机分摊的电费（电量乘以每度的价格）
	DivideDate  string		`json:"divideDate"`
}

type ElecSumResponse struct {
	Code 	int32			`json:"code"`
	Msg  	string			`json:"msg"`
	Data 	*ElecSum		`json:"data"`
}

type Scheme struct {
	RequestUrl string		`json:"requestUrl"`
	Account    string		`json:"account"`
	AppKey     string		`json:"appKey"`
	AppSecret  string		`json:"appSecret"`
}
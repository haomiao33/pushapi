package applepush

// 文档中心-应用服务-消息推送 服务器API地址以及参数
// https://dev.mi.com/distribute/doc/details?pId=1559
const (
	// APNs API 地址
	apnsURL = "https://api.sandbox.push.apple.com"
	// Bundle Identifier
	bundleIdentifier = "com.xiaogongqiu.app"
)

type SendRes struct {
	ApnsId string `json:"apnsId,omitempty"` // APNs消息ID
}

type SendReq struct {
	DeviceToken string `json:"deviceToken,omitempty"` // 接收者token
	Data        Data   `json:"data,omitempty"`        // 推送数据
}

type Data struct {
	Aps     Aps `json:"aps,omitempty"`
	Payload int `json:"payload,omitempty"` // 自定义数据
}

type Aps struct {
	Alert Alert  `json:"alert,omitempty"` // 消息内容
	Badge int    `json:"badge,omitempty"` // 角标
	Sound string `json:"sound,omitempty"` // 声音
}

type Alert struct {
	Title string `json:"title,omitempty"` // 标题
	Body  string `json:"body,omitempty"`  // 内容
}

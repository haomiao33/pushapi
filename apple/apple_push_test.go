package applepush

import (
	"testing"
)

func TestClient_Send(t *testing.T) {
	// APNs 认证密钥路径
	authKeyPath := "xx.p8"
	// Key ID 和 Team ID
	keyID := "xx"
	teamID := "xx"
	// 设备 Token
	deviceToken := "xx"
	c := NewClient(authKeyPath, keyID, teamID)
	_, err := c.Send(&SendReq{
		DeviceToken: deviceToken,
		Data: Data{
			Aps: Aps{
				Alert: Alert{
					Title: "你好啊",
					Body:  "吃饭了没有？",
				},
				Badge: 1,
				Sound: "default",
			},
			Payload: 123456,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

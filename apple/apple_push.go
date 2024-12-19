package applepush

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/http2"
)

type Client struct {
	privateKey *ecdsa.PrivateKey
	jwtToken   string
	keyId      string
	teamId     string
}

func NewClient(privateKeyPath string, keyId string, teamId string) *Client {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}
	return &Client{
		privateKey: privateKey,
		keyId:      keyId,
		teamId:     teamId,
	}
}

func (c *Client) Send(param *SendReq) (*SendRes, error) {
	if c.jwtToken == "" {
		token, err := c.generateJWT(c.privateKey)
		if err != nil {
			log.Printf("Error generating JWT: %v", err)
			return nil, err
		}
		c.jwtToken = token
	}

	// 将推送数据转为 JSON
	pushDataJSON, err := json.Marshal(param.Data)
	if err != nil {
		return nil, err
	}

	// 创建 http2 请求
	tr := &http2.Transport{}
	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/3/device/%s", apnsURL, param.DeviceToken),
		bytes.NewReader(pushDataJSON))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("apns-topic", bundleIdentifier)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.jwtToken))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 打印响应状态码
	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	// 读取并解析响应体
	var responseBody []byte
	_, err = resp.Body.Read(responseBody)
	if err != nil {
		return nil, err
	}

	// 打印响应体
	fmt.Println("Response Body:", string(responseBody))
	return &SendRes{
		ApnsId: resp.Header.Get("apns-id"),
	}, nil
}

// 读取认证密钥并解析
func loadPrivateKey(path string) (*ecdsa.PrivateKey, error) {
	keyData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	// 使用 ParsePKCS8PrivateKey 来解析 PKCS#8 格式的私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// 检查私钥类型并转换为 *ecdsa.PrivateKey
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type ecdsa.PrivateKey")
	}

	return ecdsaPrivateKey, nil
}

// 生成 JWT 令牌
func (c *Client) generateJWT(privateKey *ecdsa.PrivateKey) (string, error) {
	// 设置 JWT 负载
	payload := jwt.MapClaims{
		"iss": c.teamId, // 团队 ID
		"iat": time.Now().Unix(),
	}

	// 使用 ES256 算法生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodES256, payload)
	token.Header["kid"] = c.keyId

	// 使用私钥签名 JWT
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

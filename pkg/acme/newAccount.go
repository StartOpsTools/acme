package acme

// 3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NewAcctRequest struct {
	Protected string `json:"protected"`
	Payload string `json:"payload"`
	Signature string `json:"signature"`
}

// protected base64
type NewAcctRequestProtected struct {
	Nonce string `json:"nonce"` // NewNonce 请求中 Replay-Nonce 值
	Url string `json:"url"` // 本接口地址
	Alg string `json:"alg"` // 签名算法  RS256
	Jwk newAcctRequestProtectedJwk `json:"jwk"` // 账户公钥数据
}

type newAcctRequestProtectedJwk struct {
	E string `json:"e"`			// e 和 n、kty 可以确定一份公钥  AQAB
	Kty string `json:"kty"`		// e 和 n、kty 可以确定一份公钥  RSA
	N string `json:"n"`			// e 和 n、kty 可以确定一份公钥
}

// payload base64
type NewAcctRequestPayload struct {
	Contact []string `json:"contact"`
	TermsOfServiceAgreed bool `json:"termsOfServiceAgreed"`
}


func (newAcctRequest NewAcctRequest) NewAccount(newAccountUrl string){
	newAcctRequestByte, err := json.Marshal(newAcctRequest)
	requestParameter := bytes.NewBuffer(newAcctRequestByte)
	//fmt.Println("requestParameter: ", requestParameter)
	resp, err := http.Post(newAccountUrl, "application/jose+json", requestParameter)
	if err != nil {
		return
	}
	
	respBody := resp.Body
	defer resp.Body.Close()
	
	respBodyByte, err := ioutil.ReadAll(respBody)
	if err != nil {
		return
	}
	
	fmt.Println("respBodyByte: ", string(respBodyByte))
	
}


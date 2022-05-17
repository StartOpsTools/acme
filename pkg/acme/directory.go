package acme

// 1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DirectoryResponse struct {
	KeyChange   string                `json:"keyChange"` // 更换 KeyPair 的接口
	Meta        directoryResponseMeta `json:"meta"`
	NewAccount  string                `json:"newAccount"` // 上传 account key 时候调用
	NewNonce    string                `json:"newNonce"`   // 业务流程（如申请、撤销）时候需要预请求获取一次性令牌的一个接口
	NewOrder    string                `json:"newOrder"`   // 下单的接口
	RenewalInfo string                `json:"renewalInfo"`
	RevokeCert  string                `json:"revokeCert"` // 撤销证书的接口
	//RandomUrl  string                `json:"randomUrl"`
}

type directoryResponseMeta struct {
	CaaIdentities           []string `json:"caaIdentities"` // 可选，字符串数组
	TermsOfService          string   `json:"termsOfService"`
	Website                 string   `json:"website"`
	ExternalAccountRequired bool     `json:"externalAccountRequired"`
}

func (directoryResponse DirectoryResponse) String() string {
	bytes, _ := json.Marshal(directoryResponse)
	return string(bytes)
}

// 获取 directory 信息
func Directory(directoryUrl string) (DirectoryResponse, error) {
	var directoryResponse DirectoryResponse

	resp, err := http.Get(directoryUrl)
	if err != nil {
		return directoryResponse, err
	}

	respBody := resp.Body

	defer resp.Body.Close()

	respBodyByte, err := ioutil.ReadAll(respBody)
	if err != nil {
		return directoryResponse, err
	}
	//fmt.Println("respBodyByte: ", string(respBodyByte))
	err = json.Unmarshal(respBodyByte, &directoryResponse)
	if err != nil {
		//fmt.Println("json unmarshal err: ", err)
		return directoryResponse, err
	}

	return directoryResponse, nil
}

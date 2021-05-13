package acme

// 1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DirectoryResponse struct {
	KeyChange  string                `json:"keyChange"`
	Meta       directoryResponseMeta `json:"meta"`
	NewAccount string                `json:"newAccount"`
	NewNonce   string                `json:"newNonce"`
	NewOrder   string                `json:"newOrder"`
	RevokeCert string                `json:"revokeCert"`
	RandomUrl  string                `json:"randomUrl"`
}

type directoryResponseMeta struct {
	CaaIdentities  []string `json:"caaIdentities"`
	TermsOfService string   `json:"termsOfService"`
	Website        string   `json:"website"`
}

// 获取 directory 信息
func Directory(directoryUrl string) (DirectoryResponse, error) {
	var directoryResponse DirectoryResponse
	var respMap map[string]interface{}
	respMap = make(map[string]interface{})
	
	directoryResponseKey := []string{"meta", "newAccount", "newNonce", "newOrder", "revokeCert", "keyChange"}
	
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
	err = json.Unmarshal(respBodyByte, &respMap)
	json.Unmarshal(respBodyByte, &directoryResponse)
	
	if err != nil {
		fmt.Println("json unmarshal err: ", err)
		return directoryResponse, err
	}
	
	for _, v := range directoryResponseKey {
		delete(respMap, v)
	}
	
	for _, v := range respMap {
		directoryResponse.RandomUrl = v.(string)
	}
	return directoryResponse, nil
}

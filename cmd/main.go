package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/StartOpsTools/acme/pkg/acme"
)

var (
	//let Encrypt 正式环境 API
	letEncryptDirectoryProdUrl = "https://acme-v02.api.letsencrypt.org/directory"
	//let Encrypt 测试环境 API
	letEncryptDirectoryTestUrl = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

func main() {
	// directory - get
	directoryResponse, err := acme.Directory(letEncryptDirectoryTestUrl)
	if err != nil {
		fmt.Println("request directory error, err: ", err.Error())
		return
	}
	// newNonce - head
	newNonceUrl := directoryResponse.NewNonce
	replayNonce, err := acme.NewNonce(newNonceUrl)
	if err != nil {
		fmt.Println("request nonceUrl error, replayNonce failed, err: ", err.Error())
		return
	}
	// newAccount - post
	
	newAccountUrl := directoryResponse.NewAccount
	
	var newAcctRequestProtected acme.NewAcctRequestProtected
	newAcctRequestProtected.Nonce = replayNonce
	newAcctRequestProtected.Url = newAccountUrl
	newAcctRequestProtected.Alg = "RS256"
	newAcctRequestProtected.Jwk.E = "AQAB"
	newAcctRequestProtected.Jwk.Kty = "RSA"
	newAcctRequestProtected.Jwk.N = "uxogbZMm3hXq5c3tgIkqtrKyk69ZWKOlfmbHmSErd4NHNaSywfKey2l3-tAPywL_iI6IHEW4ViYT7ha41h6q6fGxwuRc2v9AGegV24RMwZXlZwihLBn3h507WXQYIA5auSJS6NWkF0ERsx2_095MJCrsH5HvqbCyocrqwh97mHy2zBN4IOPn-BEgb2xf8_AslgzQcVZ5UjzZxQ-p5lyzwOpfDPXY8-bsqFfk9-Mu1kzDtCTas6J35em2GhaHVN85i_tA2LNfhbz1bNTopt2wZCFRv9eMhlxBw_RhmCAkpmcVnL-a6pCnCWDSL8aWw2mGn5suXE607DWpyXHvmWJx5Q"
	newAcctRequestProtectedByte, err := json.Marshal(newAcctRequestProtected)
	if err != nil {
		return
	}
	newAcctRequestProtectedString := base64.StdEncoding.EncodeToString(newAcctRequestProtectedByte)
	
	var newAcctRequestPayload acme.NewAcctRequestPayload
	newAcctRequestPayloadContact := []string{"mailto: qx@mail.com"}
	newAcctRequestPayload.Contact = newAcctRequestPayloadContact
	newAcctRequestPayload.TermsOfServiceAgreed = true
	newAcctRequestPayloadByte, err := json.Marshal(newAcctRequestPayload)
	if err != nil {
		return
	}
	newAcctRequestPayloadString := base64.StdEncoding.EncodeToString(newAcctRequestPayloadByte)
	
	var newAcctRequest acme.NewAcctRequest
	newAcctRequest.Signature = "bzW7CyT2ln1Ms17vfc8kgDgy9yNlwFWFshEUVrFFBUJVWVROvlBjG6CFDXldCGkPvTJCJ1mRvXoMsk3uQ5FhTNaMpBsv7QVL9N03LuLZax1RE_2eyjj5ATMGq8wzXqgAjzC7ueJfcHYuFbuR33NsDXMWZmpEqKi4cf8h-GBDWnuKKUUN0-N-2aQo_JEbRrk5dMbZ1ACJ2Evy46y7-0jCmrHLgnGK7nPK6fmp__6WfNOK_PzjFj3t65NXqio68w2IqXy4CONqx0miqLdgyZfk-HGCN9c4ACNGEhj_mJ9pUEiABaQJagcOqkIADEESg6tIaw9pj0HjqxxlpHTzgrGwxQ"
	newAcctRequest.Protected = newAcctRequestProtectedString
	newAcctRequest.Payload = newAcctRequestPayloadString
	
	newAcctRequest.NewAccount(newAccountUrl)
	
}

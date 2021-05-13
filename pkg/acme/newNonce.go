package acme

// 2

import (
	"errors"
	"net/http"
)

func NewNonce(newNonceUrl string) (string, error){
	resp, err := http.Head(newNonceUrl)
	if err != nil {
		return "", err
	}
	
	contentLength := resp.ContentLength
	if contentLength > 0 {
		return "", errors.New("contentLength greater than 0")
	}
	
	replayNonce := resp.Header.Get("Replay-Nonce")
	return replayNonce, nil
}

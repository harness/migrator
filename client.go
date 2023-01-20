package main

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func Post(reqUrl string, auth string, body interface{}) ([]byte, error) {
	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)
	log.WithFields(log.Fields{
		"body": string(postBody),
	}).Debug("The request body")
	req, err := http.NewRequest("POST", reqUrl, requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"body": string(respBody),
	}).Debug("The response body")
	if resp.StatusCode != 200 {
		return nil, errors.New("received non 200 response code. The response code was " + strconv.Itoa(resp.StatusCode))
	}
	return respBody, nil
}

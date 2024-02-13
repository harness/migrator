package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Post(reqUrl string, auth string, body interface{}) (respBodyObj ResponseBody, err error) {
	postBody, _ := json.Marshal(body)
	requestBody := bytes.NewBuffer(postBody)
	log.WithFields(log.Fields{
		"url":  reqUrl,
		"body": string(postBody),
	}).Trace("The request details")
	req, err := http.NewRequest("POST", reqUrl, requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(AuthHeaderKey(auth), auth)
	return handleResp(req)
}

func Get(reqUrl string, auth string) (respBodyObj ResponseBody, err error) {
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return
	}
	log.WithFields(log.Fields{
		"url": reqUrl,
	}).Trace("The request details")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(AuthHeaderKey(auth), auth)
	return handleResp(req)
}

func Delete(reqUrl string, auth string, body interface{}) (respBodyObj ResponseBody, err error) {
	var requestBody *bytes.Buffer
	if body != nil {
		postBody, _ := json.Marshal(body)
		requestBody = bytes.NewBuffer(postBody)
		log.WithFields(log.Fields{
			"url":  reqUrl,
			"body": string(postBody),
		}).Trace("The request details")
	} else {
		log.WithFields(log.Fields{
			"url": reqUrl,
		}).Trace("The request details")
	}
	var req *http.Request
	if requestBody != nil {
		req, err = http.NewRequest("DELETE", reqUrl, requestBody)
	} else {
		req, err = http.NewRequest("DELETE", reqUrl, nil)
	}
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(AuthHeaderKey(auth), auth)
	if requestBody == nil {
		req.ContentLength = 0
	}
	return handleResp(req)
}

func handleResp(req *http.Request) (respBodyObj ResponseBody, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.WithFields(log.Fields{
		"body": string(respBody),
	}).Trace("The response body")
	err = json.Unmarshal(respBody, &respBodyObj)
	if err != nil {
		log.Fatalln("There was error while parsing the response from server. Exiting...", err)
	}
	if resp.StatusCode != 200 {
		if len(respBodyObj.Message) > 0 {
			log.Error(respBodyObj.Message)
		} else if len(respBodyObj.Messages) > 0 && len(respBodyObj.Messages[0].Message) > 0 {
			log.Error(respBodyObj.Messages[0].Message)
		}
		return respBodyObj, errors.New("received non 200 response code. The response code was " + strconv.Itoa(resp.StatusCode))
	}

	return respBodyObj, nil
}

func AuthHeaderKey(auth string) string {
	if strings.HasPrefix(auth, "Bearer ") {
		return "Authorization"
	}
	return "x-api-key"
}

func GetWithAuth(host string, query string, authMethod string, base64Auth string, certPath string, keyPath string) (body []byte, err error) {
	baseURL := "https://" + host + "/api/v1/" + query

	var client *http.Client

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Configure client based on authentication method
	if authMethod == authBasic {
		// Encode credentials to base64 for the Authorization header
		client = &http.Client{}
		// Add the Authorization header to the request
		req.Header.Add("Authorization", "Basic "+base64Auth)
	} else if authMethod == authx509 {
		// Load client certificate
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			fmt.Println("Error loading certificate:", err)
			return nil, err
		}

		// Create a HTTPS client and supply the created CA pool and certificate
		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			// In a real application, you should adjust the TLS settings according to your security requirements.
		}
		client = &http.Client{Transport: &http.Transport{TLSClientConfig: config}}
	} else {
		fmt.Println("Unsupported authentication method")
		return nil, fmt.Errorf("unsupported authentication method %s", authMethod)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}
	fmt.Println(string(body))
	return body, nil
}

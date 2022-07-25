package web

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// ----- выполнить GET-запрос -----
// 01
func SendGetRequest(url string, contentType string, authToken string) (string, error) {

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New("Error (DB-010101): http request failure")
	}

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	if len(authToken) > 0 {
		req.Header.Set("Auth-token", authToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Error (DB-010102): " + err.Error())
	}
	defer resp.Body.Close()

	//fmt.Printf("StatusCode: %d\n", resp.StatusCode)
	//j, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("resp.Body: %s\n", string(j))

	if resp.StatusCode != 200 {
		return "", errors.New("Error (DB-010103): http request failure: " + resp.Status)
	}

	resp_body, _ := ioutil.ReadAll(resp.Body)

	return string(resp_body), nil
}

// ----- выполнить POST-запрос -----
// 02
func SendPostRequest(url string, requestBody string, contentType string, authToken string) (string, error) {

	var jsonBody = []byte(requestBody)

	client := &http.Client{
		Timeout: time.Second * 300,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", errors.New("Error (DB-010201): http request failure")
	}

	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType)
	}
	if len(authToken) > 0 {
		req.Header.Set("Auth-token", authToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Error (DB-010202): " + err.Error())
	}
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Printf("StatusCode: %d\n", resp.StatusCode)
	//j, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("resp.Body: %s\n", string(resp_body))

	if resp.StatusCode != 200 {
		return "", errors.New("Error (DB-010203): http request failure: " + resp.Status + " : " + string(resp_body))
	}

	return string(resp_body), nil
}

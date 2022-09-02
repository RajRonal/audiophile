package Ginkgo

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func getResponse(url string) ([]byte, error) {
	//url := "http://localhost:8080/api/signup"
	if len(url) == 0 {
		return nil, errors.New("Invalid URL")
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)

	if err == nil && code != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("Server status error: %v", http.StatusText(code))
	}
	return body, nil
}

func CreateRequest(jsonStr []byte) *http.Request {
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	if err != nil {
		logrus.Error("Error in Creating request")
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

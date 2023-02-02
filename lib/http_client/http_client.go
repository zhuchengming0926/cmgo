package http_client

import (
	"bytes"
	"cmgo/lib/logger"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

// http的json请求方式
func PostJson(url string, requestParams interface{}) ([]byte, error) {
	resp := []byte{}
	requestBody := new(bytes.Buffer)
	jsoniter.NewEncoder(requestBody).Encode(requestParams)
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	httpResponse, err := client.Do(req)
	if err != nil {
		logger.Warn("PostJson failed", zap.String("url", url))
		return resp, err
	}
	defer httpResponse.Body.Close()
	body, err := ioutil.ReadAll(httpResponse.Body)

	return body, err
}

// Post请求
func Post(url string, data []byte, a ...int) ([]byte, error) {
	return PostWithHeader(url, nil, data, a...)
}

func PostWithHeader(url string, header map[string]string, data []byte,
	a ...int) ([]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	if len(header) == 0 {
		header = make(map[string]string)
		header["Content-Type"] = "application/json"
	}
	return RetryDoRequest("POST", url, header, data, timeout, retryTimes, interval)
}

// Get请求
func Get(url string, a ...int) ([]byte, error) {
	return GetWithHeader(url, nil, a...)
}

func GetWithHeader(url string, header map[string]string, a ...int) ([]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	return RetryDoRequest("GET", url, header, nil, timeout, retryTimes, interval)
}

func RetryDoRequest(reqType, URL string, headers map[string]string, data []byte,
	timeout, retryTimes, interval int) ([]byte, error) {
	var err1 error
	for i := 0; i < retryTimes+1; i++ {
		_, body, _, err := DoRequest(reqType, URL, headers, data, timeout)
		if err != nil {
			err1 = fmt.Errorf("%s[try %d times]", err, i+1)
			time.Sleep(time.Duration(interval) * time.Millisecond)
			continue
		}
		return body, nil
	}
	return nil, err1
}

func DoRequest(reqType, url string, headers map[string]string, data []byte, timeout int) (
	int, []byte, map[string][]string, error) {
	var reader io.Reader
	if len(data) > 0 {
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(reqType, url, reader)
	if err != nil {
		return 0, nil, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	to := time.Duration(time.Duration(timeout) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()
	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}
	statusCode := resp.StatusCode
	header := resp.Header
	if statusCode != 200 && statusCode != 201 {
		return statusCode, nil, header, fmt.Errorf("response status error: %d", statusCode)
	}
	return statusCode, body, header, nil
}

func parseParameters(a ...int) (timeout, retryTimes, interval int, err error) {
	if len(a) == 2 || len(a) > 3 {
		err = errors.New("http Post retry parameters count error")
		return
	}
	timeout = 10000
	if len(a) > 0 {
		timeout = a[0]
	}
	retryTimes, interval = 0, 10
	if len(a) == 3 {
		retryTimes, interval = a[1], a[2]
	}
	return
}

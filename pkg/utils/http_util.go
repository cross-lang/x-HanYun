package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"x-HanYun/pkg/core/log"

	"go.uber.org/zap"
)

// 内容类型常量
const (
	JsonContentType           = "application/json"
	FormUrlEncodedContentType = "application/x-www-form-urlencoded"
)

// ExtractPath 提权出Url中的路径（Path）
func ExtractPath(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("<<<<<<<< Http Failed to parse url:%s, err:%v", urlStr, err)
	}

	pathWithQuery := parsedURL.Path
	if parsedURL.RawQuery != "" {
		pathWithQuery += "?" + parsedURL.RawQuery
	}

	return pathWithQuery, nil
}

// Get 方法用于发送 HTTP GET 请求，并将响应解析为 JSON 对象
// url: 请求的目标 URL
// header: 请求头信息，以键值对形式存储
// response: 用于存储响应结果的对象指针，调用时需传入对象的指针
// 返回值: 可能出现的错误
func Get(url string, header map[string]string, response interface{}) error {
	log.Info(">>>>>>>> Http Get", zap.String("Url", url))

	// 创建 HTTP GET 请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("<<<<<<<< Http Get, NewRequest error", zap.Error(err))
		return err
	}

	// 设置请求头
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			// 临时跳过证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("<<<<<<<< Http Get, Do error", zap.Error(err))
		return err
	}

	// 确保响应体在函数结束时关闭
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("<<<<<<<< Http Get, Error closing response body", zap.Error(err))
		}
	}()

	// 先读取响应体，再处理状态码检查
	bodyBytes, err := ReadResponseBody(resp)
	if err != nil {
		log.Error("<<<<<<<< Http Get, ReadAll error", zap.Error(err))
		return err
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error("<<<<<<<< Http Get, Status Code", zap.Int("code", resp.StatusCode), zap.String("response", string(bodyBytes)))
		return fmt.Errorf("HTTP Get, Request failed with status %d:%s", resp.StatusCode, string(bodyBytes))
	}

	log.Info(">>>>>>>> Http Get, Successfully", zap.String("Response Body", string(bodyBytes)))

	// 解析响应体为 JSON
	if err := json.Unmarshal(bodyBytes, response); err != nil {
		log.Error("<<<<<<<< Http Get, Failed to Unmarshal json str", zap.Error(err), zap.String("str", string(bodyBytes)))
		return err
	}

	return nil
}

// PostV1 方法用于发送 HTTP POST 请求，请求体为 JSON 对象，响应解析为 JSON 对象
// url: 请求的目标 URL
// headers: 请求头信息，以键值对形式存储
// request: 请求体对象
// response: 用于存储响应结果的对象指针，调用时需传入对象的指针
// 返回值: 可能出现的错误
func PostV1(url string, headers map[string]string, request, response interface{}) error {
	log.Info(">>>>>>>> Http PostV1", zap.String("Url", url), zap.Any("Req Body", request))

	// 将请求对象序列化为 JSON 字符串
	requestStr, err := json.Marshal(request)
	if err != nil {
		log.Error("<<<<<<<< Http PostV1, Failed to Marshal json obj", zap.Error(err), zap.Any("obj", request))
		return err
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestStr))
	if err != nil {
		log.Error("<<<<<<<< Http PostV1, Failed to create HTTP POST request", zap.Error(err))
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", JsonContentType)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			// 临时跳过证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("<<<<<<<< Http PostV1, Do error", zap.Error(err))
		return err
	}

	// 确保响应体在函数结束时关闭
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("<<<<<<<< Http PostV1, Error closing response body", zap.Error(err))
		}
	}()

	// 先读取响应体，再处理状态码检查
	bodyBytes, err := ReadResponseBody(resp)
	if err != nil {
		log.Error("<<<<<<<< Http PostV1, Failed to read response body", zap.Error(err))
		return err
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error("<<<<<<<< Http PostV1, Do Failed", zap.String("Url", url), zap.Int("Code", resp.StatusCode), zap.String("Resp", string(bodyBytes)))
		return fmt.Errorf("HTTP PostV1, Request failed with status %d:%s", resp.StatusCode, string(bodyBytes))
	}

	log.Info(">>>>>>>> Http PostV1, Do Successfully", zap.String("Response Body", string(bodyBytes)))

	// 将响应体反序列化为对象
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		log.Error("<<<<<<<< Http PostV1, Failed to Unmarshal json str", zap.Error(err), zap.String("str", string(bodyBytes)))
		return err
	}

	return nil
}

// PostV2 方法用于发送 HTTP POST 请求，请求体为字符串，返回响应字符串
// url: 请求的目标 URL
// headers: 请求头信息，以键值对形式存储
// request: 请求体字符串
// 返回值: 响应字符串和可能出现的错误
func PostV2(url string, headers map[string]string, request string) (string, error) {
	log.Info(">>>>>>>> Http PostV2", zap.String("Url", url), zap.Any("Req Body", request))

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(request))
	if err != nil {
		log.Error("<<<<<<<< Http PostV2, Failed to create HTTP POST request", zap.Error(err))
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", JsonContentType)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// 创建 HTTP 的客户端
	client := &http.Client{
		Transport: &http.Transport{
			// 临时跳过证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("<<<<<<<< Http PostV2, Do error", zap.Error(err))
		return "", err
	}

	// 确保响应体在函数结束时关闭
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("<<<<<<<< Http PostV2, Error closing response body", zap.Error(err))
		}
	}()

	// 先读取响应体，再处理状态码检查
	bodyBytes, err := ReadResponseBody(resp)
	if err != nil {
		log.Error("<<<<<<<< Http PostV2, Failed to read response body", zap.Error(err))
		return "", err
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error("<<<<<<<< Http PostV2, Do Failed", zap.String("Url", url), zap.Int("Code", resp.StatusCode), zap.String("Resp", string(bodyBytes)))
		return "", fmt.Errorf("HTTP PostV2, Request failed with status %d:%s", resp.StatusCode, string(bodyBytes))
	}

	log.Info(">>>>>>>> Http PostV2, Do Successfully", zap.String("Response Body", string(bodyBytes)))

	return string(bodyBytes), nil
}

// PostV3 方法用于发送 HTTP POST 请求，使用x-www-form-urlencoded格式
// url: 请求的目标 URL
// headers: 请求头信息，以键值对形式存储
// params: 请求参数，以键值对形式存储
// response: 用于存储响应结果的对象指针，调用时需传入对象的指针
// 返回值: 可能出现的错误
func PostV3(xurl string, headers map[string]string, params map[string]string, response interface{}) error {
	log.Info(">>>>>>>> Http PostV3", zap.String("Url", xurl), zap.Any("Params", params))

	// 构建x-www-form-urlencoded格式请求体
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodPost, xurl, strings.NewReader(formData.Encode()))
	if err != nil {
		log.Error("<<<<<<<< Http PostV3, Failed to create HTTP POST request", zap.Error(err))
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", FormUrlEncodedContentType)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			// 临时跳过证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("<<<<<<<< Http PostV3, Do error", zap.Error(err))
		return err
	}

	// 确保响应体在函数结束时关闭
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("<<<<<<<< Http PostV3, Error closing response body", zap.Error(err))
		}
	}()

	// 先读取响应体，再处理状态码检查
	bodyBytes, err := ReadResponseBody(resp)
	if err != nil {
		log.Error("<<<<<<<< Http PostV3, Failed to read response body", zap.Error(err))
		return err
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error("<<<<<<<< Http PostV3, Do Failed", zap.String("Url", xurl), zap.Int("Code", resp.StatusCode), zap.String("Resp", string(bodyBytes)))
		return fmt.Errorf("HTTP PostV3, Request failed with status %d:%s", resp.StatusCode, string(bodyBytes))
	}

	log.Info(">>>>>>>> Http PostV3, Do Successfully", zap.String("Response Body", string(bodyBytes)))

	// 将响应体反序列化为对象
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		log.Error("<<<<<<<< Http PostV3, Failed to Unmarshal json str", zap.Error(err), zap.String("str", string(bodyBytes)))
		return err
	}

	return nil
}

// PostV4 方法用于发送 HTTP POST 请求，使用form-data格式
// url: 请求的目标 URL
// headers: 请求头信息，以键值对形式存储
// params: 请求参数，以键值对形式存储
// response: 用于存储响应结果的对象指针，调用时需传入对象的指针
// 返回值: 可能出现的错误
func PostV4(url string, headers map[string]string, params map[string]string, response interface{}) error {
	log.Info(">>>>>>>> Http PostV4", zap.String("Url", url), zap.Any("Params", params))

	// 构建form-data格式请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			log.Error("<<<<<<<< Http PostV4, Failed to write form field", zap.String("key", k), zap.Error(err))
			return err
		}
	}
	if err := writer.Close(); err != nil {
		log.Error("<<<<<<<< Http PostV4, Failed to close multipart writer", zap.Error(err))
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Error("<<<<<<<< Http PostV4, Failed to create HTTP POST request", zap.Error(err))
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			// 临时跳过证书验证
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("<<<<<<<< Http PostV4, Do error", zap.Error(err))
		return err
	}

	// 确保响应体在函数结束时关闭
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error("<<<<<<<< Http PostV4, Error closing response body", zap.Error(err))
		}
	}()

	// 先读取响应体，再处理状态码检查
	bodyBytes, err := ReadResponseBody(resp)
	if err != nil {
		log.Error("<<<<<<<< Http PostV4, Failed to read response body", zap.Error(err))
		return err
	}

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Error("<<<<<<<< Http PostV4, Do Failed", zap.String("Url", url), zap.Int("Code", resp.StatusCode), zap.String("Resp", string(bodyBytes)))
		return fmt.Errorf("HTTP PostV4, Request failed with status %d:%s", resp.StatusCode, string(bodyBytes))
	}

	log.Info(">>>>>>>> Http PostV4, Do Successfully", zap.String("Response Body", string(bodyBytes)))

	// 将响应体反序列化为对象
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		log.Error("<<<<<<<< Http PostV4, Failed to Unmarshal json str", zap.Error(err), zap.String("str", string(bodyBytes)))
		return err
	}

	return nil
}

// ReadResponseBody 方法用于读取 HTTP 响应体，并处理可能的压缩
// response: HTTP 响应对象
// 返回值: 响应体字节切片和可能出现的错误
func ReadResponseBody(response *http.Response) ([]byte, error) {
	var bodyReader io.Reader
	var err error

	// 处理可能的gzip压缩
	if response.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		bodyReader = response.Body
	}

	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	// 限制日志输出长度，最多显示200个字符
	if len(data) > 200 {
		log.Info(">>>>>>>> Http Resp Body", zap.String("value", string(data[:200])+"..."))
	} else {
		log.Info(">>>>>>>> Http Resp Body", zap.String("value", string(data)))
	}
	return data, nil
}

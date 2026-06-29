package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"x-HanYun/pkg/core/log"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"

	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type DocContent struct {
	Doc         string `json:"docs"`
	DocAsUpsert bool   `json:"doc_as_upsert"`
}

type ElasticsearchClient struct {
	Client *elasticsearch.Client
}

// NewElasticsearchClient 用于初始化 Elasticsearch 客户端
func NewElasticsearchClient(addresses, username, password string) (*ElasticsearchClient, error) {
	//logx.Infof(">>>>>>>> Elastic Search Addresses:%s", addresses)
	//logx.Infof(">>>>>>>> Elastic Search Username:%s", username)
	//logx.Infof(">>>>>>>> Elastic Search Password:%s", password)

	// 创建一个忽略证书验证的 HTTP 客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		// 设置连接超时
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second, // 连接建立超时
			KeepAlive: 30 * time.Second,
		}).DialContext,
		// 设置响应头超时
		ResponseHeaderTimeout: 10 * time.Second,
		// 设置TLS握手超时
		TLSHandshakeTimeout: 10 * time.Second,
	}

	cfg := elasticsearch.Config{
		Addresses: []string{addresses},
		Username:  username,
		Password:  password,
		Transport: tr, // 设置 Transport
	}

	esCli, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Error("<<<<<<<< Error creating Elasticsearch client", zap.Error(err))
		return nil, err
	}

	// 验证连接
	res, err2 := esCli.Info()
	if err2 != nil {
		logMsg := fmt.Errorf("<<<<<<<< Step 1, error getting Elasticsearch info:%v", err2)
		log.Error(logMsg.Error())
		return nil, logMsg
	}
	defer func(Body io.ReadCloser) {
		err4 := Body.Close()
		if err4 != nil {
			log.Error(err4.Error())
		}
	}(res.Body)

	if res.IsError() {
		logMsg := fmt.Errorf("<<<<<<<< Step 2, error getting Elasticsearch info:%s", res.Status())
		log.Error(logMsg.Error())
		return nil, logMsg
	}

	return &ElasticsearchClient{Client: esCli}, nil
}

func Init(addresses, userName, password string) {
	// 初始化客户端
	_, err := NewElasticsearchClient(addresses, userName, password)
	if err != nil {
		log.Error("<<<<<<<< Error Creating Elastic Search Client", zap.Error(err))
		panic(err)
	}

	// 初始化索引
}

// Search 用于执行搜索查询
func (es *ElasticsearchClient) Search(index, query string) (*esapi.Response, error) {
	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(strings.NewReader(query)),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("<<<<<<<< Error executing search query:%w", err)
	}

	return res, nil
}

// Index 用于向指定索引插入数据
func (es *ElasticsearchClient) Index(index, docId string, jsonMessage string) (*esapi.Response, error) {
	// 将 JSON 字符串转换为字节数组
	data := []byte(jsonMessage)

	// 创建插入请求
	req := esapi.IndexRequest{
		Index:      index, // 指定索引
		DocumentID: docId, // 文档 ID
		Body:       bytes.NewReader(data),
		Refresh:    "true", // 刷新索引，以确保立即可搜索
	}

	// 执行插入操作
	resp, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, fmt.Errorf("<<<<<<<< Error Indexing Document:%w", err)
	}

	return resp, nil
}

func (es *ElasticsearchClient) IndexBulk(index string, docs map[string]string) (*esapi.Response, error) {
	// 构建 Bulk 请求体
	var bulkRequest bytes.Buffer
	for docId, jsonMessage := range docs {
		// 使用 "update" 或 "index" 操作，确保支持更新或创建
		meta := fmt.Sprintf(`{ "update" : { "_id" : "%s" } }%s`, docId, "\n")
		bulkRequest.WriteString(meta)

		// 解析 jsonMessage 为 map[string]interface{}
		var doc map[string]interface{}
		err := json.Unmarshal([]byte(jsonMessage), &doc)
		if err != nil {
			return nil, fmt.Errorf("<<<<<<<< Failed to unmarshal JSON message:%w", err)
		}

		// 构建 update 操作的内容，包含 doc 和 upsert
		// 减少重复代码：无需在逻辑层判断是创建还是更新，交由 Elasticsearch 自动处理。
		// 保证幂等性：同一个 _id 多次执行不会重复创建。
		updateDoc := map[string]interface{}{
			"doc":           doc,
			"doc_as_upsert": true,
		}

		// 将 updateDoc 转换为 JSON 字符串
		updateDocJSON, err := json.Marshal(updateDoc)
		if err != nil {
			return nil, fmt.Errorf("<<<<<<<< Failed to marshal update document:%w", err)
		}

		// 将 JSON 字符串写入请求体
		bulkRequest.Write(updateDocJSON)
		bulkRequest.WriteByte('\n')
	}

	// 创建 Bulk 插入请求
	req := esapi.BulkRequest{
		Index:   index,
		Body:    bytes.NewReader(bulkRequest.Bytes()),
		Refresh: "true", // 刷新索引，以确保立即可搜索
	}

	// 执行 Bulk 操作
	resp, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, fmt.Errorf("<<<<<<<< Error bulk indexing documents:%w", err)
	}

	return resp, nil
}

// Update 用于向指定索引更新数据
func (es *ElasticsearchClient) Update(index, docId, jsonMessage string) (*esapi.Response, error) {
	// 将 JSON 字符串转换为字节数组
	data := []byte(jsonMessage)

	// 创建更新请求
	req := esapi.UpdateRequest{
		Index:      index, // 指定索引
		DocumentID: docId, // 指定文档ID
		Body:       bytes.NewReader(data),
		Refresh:    "true", // 刷新索引，以确保立即可搜索
	}

	// 执行更新操作
	resp, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, fmt.Errorf("<<<<<<<< Error updating document:%w", err)
	}

	return resp, nil
}

// CreateIndex 创建带映射的索引
func (es *ElasticsearchClient) CreateIndex(index string, mapping string) error {
	// 检查索引是否已存在
	exists, err := es.Client.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("<<<<<<<< Error checking if index exists:%w", err)
	}
	if exists.StatusCode == 200 {
		log.Info(fmt.Sprintf(">>>>>>>> Index %s already exists", index))
		return nil
	}

	// 创建索引
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  bytes.NewReader([]byte(mapping)),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return fmt.Errorf("<<<<<<<< Error creating index:%w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("<<<<<<<< Error response from Elasticsearch:%s", res.String())
	}

	log.Info(fmt.Sprintf(">>>>>>>> Index %s created successfully", index))
	return nil
}

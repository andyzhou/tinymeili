package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
 * embed http api
 * - base on `MiniLM` embed service
 */

//const variables
const (
	EmbedUri = "http://localhost:8000/embed"
	Timeout  = 5 * time.Second
)

// EmbedResponse 是 /embed 接口的响应结构
type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GetEmbedding 向本地 embedding 服务发送 POST 请求，返回嵌入向量
func GetEmbedding(text string) ([]float32, error) {
	// 请求体结构
	payload := map[string]string{"text": text}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建 HTTP 客户端
	client := &http.Client{Timeout: Timeout}
	req, err := http.NewRequest("POST", EmbedUri, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取并解析响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service error: %s", responseBody)
	}

	var embedResp EmbedResponse
	if err := json.Unmarshal(responseBody, &embedResp); err != nil {
		return nil, fmt.Errorf("failed to parse embedding: %w", err)
	}

	return embedResp.Embedding, nil
}

// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/liuxiaobopro/dashscope-go/common"
)

var (
	sharedHTTPClient     *http.Client
	sharedHTTPClientOnce sync.Once
)

func sharedClient() *http.Client {
	sharedHTTPClientOnce.Do(func() {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxIdleConnsPerHost = 10
		t.IdleConnTimeout = 90 * time.Second
		sharedHTTPClient = &http.Client{Transport: t}
	})
	return sharedHTTPClient
}

// CloseSharedSyncSession close the shared HTTP client idle connections.
func CloseSharedSyncSession() {
	if sharedHTTPClient != nil {
		sharedHTTPClient.CloseIdleConnections()
	}
}

// HttpRequest HTTP / SSE request.
type HttpRequest struct {
	URL             string
	FlattenedOutput bool
	AsyncRequest    bool
	Encryption      *Encryption
	Headers         map[string]string
	Query           bool
	Method          string
	Stream          bool
	Timeout         time.Duration
	Data            *ApiRequestData
	HTTPClient      *http.Client
}

func defaultRequestHeaders(apiKey, userAgent string) map[string]string {
	ua := common.GetUserAgent()
	if userAgent != "" {
		ua += "; " + userAgent
	}
	h := map[string]string{"user-agent": ua}
	for k, v := range common.GetSDKHeaders("") {
		h[k] = v
	}
	h["Accept"] = "application/json; charset=utf-8"
	h["Authorization"] = "Bearer " + apiKey
	if strings.ToLower(os.Getenv(common.DASHSCOPE_DISABLE_DATA_INSPECTION_ENV)) == "false" {
		h["X-DashScope-DataInspection"] = "enable"
	}
	return h
}

// NewHttpRequest create HttpRequest.
func NewHttpRequest(url, apiKey, httpMethod string, stream, asyncRequest, query bool, timeout int, taskID string, flattened bool, encryption *Encryption, userAgent string) *HttpRequest {
	h := defaultRequestHeaders(apiKey, userAgent)
	if encryption != nil && encryption.IsValid() {
		keyJSON, _ := json.Marshal(map[string]string{
			"public_key_id": encryption.GetPubKeyID(),
			"encrypt_key":   encryption.GetEncryptedAESKeyStr(),
			"iv":            encryption.GetBase64IVStr(),
		})
		h["X-DashScope-EncryptionKey"] = string(keyJSON)
	}
	if asyncRequest && !query {
		h["X-DashScope-Async"] = "enable"
	}
	if httpMethod == common.HTTPMethodPOST {
		h["Content-Type"] = "application/json; charset=utf-8"
	}
	if stream {
		h["Accept"] = common.SSE_CONTENT_TYPE
		h["X-Accel-Buffering"] = "no"
		h["X-DashScope-SSE"] = "enable"
	}
	if query {
		url = strings.Replace(url, "/api/", "/api-task/", 1)
		url += taskID
	}
	to := common.DEFAULT_REQUEST_TIMEOUT_SECONDS
	if timeout > 0 {
		to = timeout
	}
	return &HttpRequest{
		URL:             url,
		FlattenedOutput: flattened,
		AsyncRequest:    asyncRequest,
		Encryption:      encryption,
		Headers:         h,
		Query:           query,
		Method:          httpMethod,
		Stream:          stream,
		Timeout:         time.Duration(to) * time.Second,
	}
}

func (r *HttpRequest) AddHeader(key, value string) {
	r.Headers[key] = value
}

func (r *HttpRequest) AddHeaders(headers map[string]string) {
	for k, v := range headers {
		r.Headers[k] = v
	}
}

func (r *HttpRequest) SetData(data *ApiRequestData) {
	r.Data = data
}

func (r *HttpRequest) client() *http.Client {
	if r.HTTPClient != nil {
		return r.HTTPClient
	}
	return sharedClient()
}

// Call send request. If stream, returns a channel; otherwise a single response.
func (r *HttpRequest) Call(ctx context.Context) (*DashScopeAPIResponse, error) {
	ch, err := r.CallStream(ctx)
	if err != nil {
		return nil, err
	}
	var last *DashScopeAPIResponse
	for rsp := range ch {
		last = rsp
		if !r.Stream {
			// drain
			for range ch {
			}
			return last, nil
		}
	}
	return last, nil
}

// CallStream send request and yield responses.
func (r *HttpRequest) CallStream(ctx context.Context) (<-chan *DashScopeAPIResponse, error) {
	out := make(chan *DashScopeAPIResponse)
	go func() {
		defer close(out)
		resp, err := r.doHTTP(ctx)
		if err != nil {
			out <- &DashScopeAPIResponse{StatusCode: -1, Code: "Unknown", Message: err.Error()}
			return
		}
		for rsp := range r.handleResponse(resp) {
			out <- rsp
		}
	}()
	return out, nil
}

func (r *HttpRequest) doHTTP(ctx context.Context) (*http.Response, error) {
	var body io.Reader
	if r.Method == common.HTTPMethodPOST {
		var obj map[string]any
		if r.Data != nil {
			_, _, obj = r.Data.GetHTTPPayload()
		}
		common.Log.Debug("Request body: %v", obj)
		b, err := json.Marshal(obj)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	if r.Method == common.HTTPMethodGET && r.Data != nil {
		q := req.URL.Query()
		for k, v := range r.Data.Parameters {
			q.Set(k, fmtString(v))
		}
		req.URL.RawQuery = q.Encode()
	}
	client := r.client()
	c := *client
	c.Timeout = r.Timeout
	if r.Stream {
		c.Timeout = 0
	}
	common.Log.Debug("Starting request: %s", r.URL)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *HttpRequest) handleResponse(resp *http.Response) <-chan *DashScopeAPIResponse {
	ch := make(chan *DashScopeAPIResponse)
	go func() {
		defer close(ch)
		defer func() { _ = resp.Body.Close() }()
		headers := common.HeaderMap(resp.Header)
		requestID := ""
		ct := resp.Header.Get("content-type")
		if resp.StatusCode == http.StatusOK && r.Stream && strings.Contains(ct, common.SSE_CONTENT_TYPE) {
			scannerCh := scanSSE(resp.Body)
			for item := range scannerCh {
				if item.err != nil {
					ch <- &DashScopeAPIResponse{StatusCode: http.StatusBadRequest, Code: "Unknown", Message: item.err.Error(), Headers: headers}
					return
				}
				data := item.event.Data
				var msg map[string]any
				if err := json.Unmarshal([]byte(data), &msg); err != nil {
					ch <- &DashScopeAPIResponse{RequestID: requestID, StatusCode: http.StatusBadRequest, Code: "Unknown", Message: data, Headers: headers}
					continue
				}
				common.Log.Debug("Stream message: %v", msg)
				if item.isError {
					code, _ := msg["code"].(string)
					message, _ := msg["message"].(string)
					ch <- &DashScopeAPIResponse{RequestID: requestID, StatusCode: item.status, Code: code, Message: message, Headers: headers}
					return
				}
				var output any
				var usage any
				if v, ok := msg["output"]; ok {
					output = v
				}
				if v, ok := msg["usage"]; ok {
					usage = v
				}
				if v, ok := msg["request_id"].(string); ok {
					requestID = v
				}
				if r.FlattenedOutput {
					ch <- &DashScopeAPIResponse{StatusCode: http.StatusOK, Output: msg, Headers: headers}
					continue
				}
				if r.Encryption != nil && r.Encryption.IsValid() {
					output = r.Encryption.Decrypt(output)
				}
				ch <- &DashScopeAPIResponse{RequestID: requestID, StatusCode: http.StatusOK, Output: output, Usage: usage, Headers: headers}
			}
			return
		}
		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			var jsonContent map[string]any
			if err := json.Unmarshal(body, &jsonContent); err != nil {
				ch <- &DashScopeAPIResponse{StatusCode: resp.StatusCode, Code: "Unknown", Message: string(body), Headers: headers}
				return
			}
			common.Log.Debug("Response: %v", jsonContent)
			if r.FlattenedOutput {
				ch <- &DashScopeAPIResponse{StatusCode: http.StatusOK, Output: jsonContent, Headers: headers}
				return
			}
			var output any
			var usage any
			if v, ok := jsonContent["output"]; ok {
				output = v
			} else if _, ok := jsonContent["task_id"]; ok {
				output = map[string]any{"task_id": jsonContent["task_id"]}
			}
			if v, ok := jsonContent["usage"]; ok {
				usage = v
			}
			if v, ok := jsonContent["request_id"].(string); ok {
				requestID = v
			}
			if r.Encryption != nil && r.Encryption.IsValid() {
				output = r.Encryption.Decrypt(output)
			}
			ch <- &DashScopeAPIResponse{RequestID: requestID, StatusCode: http.StatusOK, Output: output, Usage: usage, Headers: headers}
			return
		}
		failed := handleFailed(resp, headers)
		ch <- failed
	}()
	return ch
}

func handleFailed(resp *http.Response, headers map[string]string) *DashScopeAPIResponse {
	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("content-type")
	if strings.Contains(ct, "application/json") {
		var errorBody map[string]any
		_ = json.Unmarshal(body, &errorBody)
		msg, _ := errorBody["message"].(string)
		if msg == "" {
			msg, _ = errorBody["msg"].(string)
		}
		code, _ := errorBody["code"].(string)
		rid, _ := errorBody["request_id"].(string)
		return &DashScopeAPIResponse{RequestID: rid, StatusCode: resp.StatusCode, Code: code, Message: msg, Headers: headers}
	}
	return &DashScopeAPIResponse{StatusCode: resp.StatusCode, Code: "Unknown", Message: string(body), Headers: headers}
}

type sseItem struct {
	isError bool
	status  int
	event   common.SSEEvent
	err     error
}

func scanSSE(r io.Reader) <-chan sseItem {
	ch := make(chan sseItem)
	go func() {
		defer close(ch)
		for item := range common.HandleStream(r) {
			ch <- sseItem{isError: item.IsError, status: item.StatusCode, event: item.Event, err: item.Err}
		}
	}()
	return ch
}

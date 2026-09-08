// Copyright (c) Alibaba, Inc. and its affiliates.

package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/common"
)

// CallParams common call parameters.
type CallParams struct {
	Model              string
	Input              any
	TaskGroup          string
	Task               string
	Function           string
	APIKey             string
	Workspace          string
	Stream             bool
	APIProtocol        string
	HTTPMethod         string
	WSStreamMode       string
	IsBinaryInput      bool
	RequestTimeout     int
	Headers            map[string]string
	Form               map[string]any
	Resources          any
	BaseAddress        string
	FlattenedOutput    bool
	ExtraURLParameters map[string]any
	UserAgent          string
	EnableEncryption   bool
	AsyncRequest       bool
	Query              bool
	TaskID             string
	SDKModule          string
	Parameters         map[string]any
	Extra              map[string]any
}

func mergeParams(p *CallParams) map[string]any {
	out := map[string]any{}
	if p.Parameters != nil {
		for k, v := range p.Parameters {
			out[k] = v
		}
	}
	if p.Extra != nil {
		for k, v := range p.Extra {
			out[k] = v
		}
	}
	return out
}

func applyWorkspace(p *CallParams) {
	if p.Workspace != "" {
		if p.Headers == nil {
			p.Headers = map[string]string{}
		}
		p.Headers["X-DashScope-WorkSpace"] = p.Workspace
	}
}

func validateAPIKeyModel(apiKey, model string) (string, string, error) {
	if apiKey == "" {
		var err error
		apiKey, err = common.GetDefaultAPIKey()
		if err != nil {
			return "", "", err
		}
	}
	if model == "" {
		return "", "", common.NewModelRequired("Model is required!")
	}
	return apiKey, model, nil
}

func toRequestOptions(p *CallParams) apientities.RequestOptions {
	params := mergeParams(p)
	if p.Stream {
		params["stream"] = true
	}
	isService := p.BaseAddress == "" || p.TaskGroup != "" || p.Task != "" || p.Function != ""
	return apientities.RequestOptions{
		Model:              p.Model,
		Input:              p.Input,
		TaskGroup:          p.TaskGroup,
		Task:               p.Task,
		Function:           p.Function,
		APIKey:             p.APIKey,
		IsService:          isService,
		APIProtocol:        p.APIProtocol,
		HTTPMethod:         p.HTTPMethod,
		Stream:             p.Stream,
		AsyncRequest:       p.AsyncRequest,
		RequestTimeout:     p.RequestTimeout,
		WSStreamMode:       p.WSStreamMode,
		IsBinaryInput:      p.IsBinaryInput,
		Query:              p.Query,
		Headers:            p.Headers,
		Form:               p.Form,
		Resources:          p.Resources,
		BaseAddress:        p.BaseAddress,
		FlattenedOutput:    p.FlattenedOutput,
		ExtraURLParameters: p.ExtraURLParameters,
		UserAgent:          p.UserAgent,
		TaskID:             p.TaskID,
		EnableEncryption:   p.EnableEncryption,
		SDKModule:          p.SDKModule,
		Parameters:         params,
	}
}

// BaseApi base API, internal use only.
type BaseApi struct{}

// Call call service and get result.
func (BaseApi) Call(ctx context.Context, p *CallParams) (*apientities.DashScopeAPIResponse, error) {
	apiKey, model, err := validateAPIKeyModel(p.APIKey, p.Model)
	if err != nil {
		return nil, err
	}
	p.APIKey = apiKey
	p.Model = model
	applyWorkspace(p)
	req, err := apientities.BuildAPIRequest(toRequestOptions(p))
	if err != nil {
		return nil, err
	}
	ch, err := req.CallStream(ctx)
	if err != nil {
		return nil, err
	}
	var last *apientities.DashScopeAPIResponse
	for rsp := range ch {
		last = rsp
	}
	return last, nil
}

// CallStream call service and stream results.
func (BaseApi) CallStream(ctx context.Context, p *CallParams) (<-chan *apientities.DashScopeAPIResponse, error) {
	apiKey, model, err := validateAPIKeyModel(p.APIKey, p.Model)
	if err != nil {
		return nil, err
	}
	p.APIKey = apiKey
	p.Model = model
	p.Stream = true
	applyWorkspace(p)
	req, err := apientities.BuildAPIRequest(toRequestOptions(p))
	if err != nil {
		return nil, err
	}
	return req.CallStream(ctx)
}

func isRepeatable(code int) bool {
	for _, s := range common.RepeatableStatus {
		if s == code {
			return true
		}
	}
	return false
}

func normalizationURL(baseAddress string, args ...string) string {
	base := baseAddress
	if base == "" {
		base = common.BaseHTTPAPIURL
	}
	return common.JoinURL(base, args...)
}

func getTaskID(task any) (string, error) {
	switch t := task.(type) {
	case string:
		if t == "" {
			return "", common.NewInvalidParameter("Task id required!")
		}
		return t, nil
	case *apientities.DashScopeAPIResponse:
		if t.StatusCode == http.StatusOK {
			id := t.TaskID()
			if id == "" {
				return "", common.NewInvalidParameter("Task id required!")
			}
			return id, nil
		}
		return "", common.NewInvalidTask("Invalid task, task create failed: " + t.String())
	default:
		return "", common.NewInvalidParameter("Task invalid!")
	}
}

func httpGetJSON(ctx context.Context, rawURL string, headers map[string]string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	common.Log.Debug("Starting request: %s", u.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return parseSimpleResponse(resp)
}

func httpPostJSON(ctx context.Context, rawURL string, headers map[string]string, body io.Reader) (*apientities.DashScopeAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return parseSimpleResponse(resp)
}

func parseSimpleResponse(resp *http.Response) (*apientities.DashScopeAPIResponse, error) {
	headers := common.HeaderMap(resp.Header)
	body, _ := io.ReadAll(resp.Body)
	var jsonContent map[string]any
	if err := json.Unmarshal(body, &jsonContent); err != nil {
		return &apientities.DashScopeAPIResponse{StatusCode: resp.StatusCode, Code: "Unknown", Message: string(body), Headers: headers}, nil
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg, _ := jsonContent["message"].(string)
		code, _ := jsonContent["code"].(string)
		rid, _ := jsonContent["request_id"].(string)
		return &apientities.DashScopeAPIResponse{RequestID: rid, StatusCode: resp.StatusCode, Code: code, Message: msg, Headers: headers}, nil
	}
	rid, _ := jsonContent["request_id"].(string)
	var output any
	var usage any
	if v, ok := jsonContent["output"]; ok {
		output = v
	} else if v, ok := jsonContent["data"]; ok {
		output = v
	} else {
		cp := map[string]any{}
		for k, v := range jsonContent {
			if k != "request_id" {
				cp[k] = v
			}
		}
		output = cp
	}
	if v, ok := jsonContent["usage"]; ok {
		usage = v
	}
	return &apientities.DashScopeAPIResponse{RequestID: rid, StatusCode: resp.StatusCode, Output: output, Usage: usage, Headers: headers}, nil
}

func defaultTaskHeaders(apiKey, workspace string) (map[string]string, error) {
	h, err := common.DefaultHeaders(apiKey, "tasks")
	if err != nil {
		return nil, err
	}
	for k, v := range common.WorkspaceHeader(workspace) {
		h[k] = v
	}
	return h, nil
}

// BaseAsyncApi for async task, internal use only.
type BaseAsyncApi struct{}

func (BaseAsyncApi) get(ctx context.Context, taskID, apiKey, workspace, baseAddress string) (*apientities.DashScopeAPIResponse, error) {
	statusURL := normalizationURL(baseAddress, "tasks", taskID)
	headers, err := defaultTaskHeaders(apiKey, workspace)
	if err != nil {
		return nil, err
	}
	return httpGetJSON(ctx, statusURL, headers, nil)
}

// AsyncCall call async service return async task information.
func (BaseAsyncApi) AsyncCall(ctx context.Context, p *CallParams) (*apientities.DashScopeAPIResponse, error) {
	p.AsyncRequest = true
	p.Query = false
	p.Stream = false
	return (BaseApi{}).Call(ctx, p)
}

// Fetch query async task status.
func (a BaseAsyncApi) Fetch(ctx context.Context, task any, apiKey, workspace, baseAddress string) (*apientities.DashScopeAPIResponse, error) {
	taskID, err := getTaskID(task)
	if err != nil {
		return nil, err
	}
	return a.get(ctx, taskID, apiKey, workspace, baseAddress)
}

// Cancel cancel PENDING task.
func (BaseAsyncApi) Cancel(ctx context.Context, task any, apiKey, workspace, baseAddress string) (*apientities.DashScopeAPIResponse, error) {
	taskID, err := getTaskID(task)
	if err != nil {
		return nil, err
	}
	u := normalizationURL(baseAddress, "tasks", taskID, "cancel")
	headers, err := defaultTaskHeaders(apiKey, workspace)
	if err != nil {
		return nil, err
	}
	return httpPostJSON(ctx, u, headers, nil)
}

// List async tasks.
func (BaseAsyncApi) List(ctx context.Context, params map[string]string, apiKey, workspace, baseAddress string) (*apientities.DashScopeAPIResponse, error) {
	u := normalizationURL(baseAddress, "tasks")
	headers, err := defaultTaskHeaders(apiKey, workspace)
	if err != nil {
		return nil, err
	}
	return httpGetJSON(ctx, u, headers, params)
}

// Wait wait for async task completion and return task result.
func (a BaseAsyncApi) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int, baseAddress string) (*apientities.DashScopeAPIResponse, error) {
	taskID, err := getTaskID(task)
	if err != nil {
		return nil, err
	}
	waitSeconds := 1
	maxWaitSeconds := 5
	incrementSteps := 3
	step := 0
	start := time.Now()
	for {
		step++
		if waitSeconds < maxWaitSeconds && step%incrementSteps == 0 {
			waitSeconds *= 2
			if waitSeconds > maxWaitSeconds {
				waitSeconds = maxWaitSeconds
			}
		}
		if waitTimeout > 0 && time.Since(start) >= time.Duration(waitTimeout)*time.Second {
			common.Log.Warn("Wait task: %s timeout after %d seconds.", taskID, waitTimeout)
			return &apientities.DashScopeAPIResponse{
				RequestID:  taskID,
				StatusCode: http.StatusRequestTimeout,
				Code:       "WaitTaskTimeout",
				Message:    "Wait task: " + taskID + " timeout after specified seconds.",
			}, nil
		}
		rsp, err := a.get(ctx, taskID, apiKey, workspace, baseAddress)
		if err != nil {
			return nil, err
		}
		if rsp.StatusCode == http.StatusOK {
			if rsp.Output == nil {
				return rsp, nil
			}
			st := rsp.TaskStatus()
			switch st {
			case common.TaskStatusFailed, common.TaskStatusCanceled, common.TaskStatusSucceeded, common.TaskStatusUnknown:
				return rsp, nil
			default:
				common.Log.Info("The task %s is  %s", taskID, st)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(time.Duration(waitSeconds) * time.Second):
				}
			}
		} else if isRepeatable(rsp.StatusCode) {
			common.Log.Warn("Get task: %s temporary failure, status_code: %d, code: %s message: %s, will try again.",
				taskID, rsp.StatusCode, rsp.Code, rsp.Message)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(waitSeconds) * time.Second):
			}
		} else {
			return rsp, nil
		}
	}
}

// Call async call then wait.
func (a BaseAsyncApi) Call(ctx context.Context, p *CallParams, waitTimeout int) (*apientities.DashScopeAPIResponse, error) {
	task, err := a.AsyncCall(ctx, p)
	if err != nil {
		return nil, err
	}
	return a.Wait(ctx, task, p.APIKey, p.Workspace, waitTimeout, p.BaseAddress)
}

func getURL(customBase, defaultPath, path string) string {
	base := customBase
	if base == "" {
		base = common.BaseHTTPAPIURL
	}
	if path != "" {
		return common.JoinURL(base, path)
	}
	return common.JoinURL(base, defaultPath)
}

func restHeaders(apiKey, workspace, module string, extra map[string]string) (map[string]string, error) {
	h, err := common.DefaultHeaders(apiKey, module)
	if err != nil {
		return nil, err
	}
	for k, v := range common.WorkspaceHeader(workspace) {
		h[k] = v
	}
	for k, v := range extra {
		h[k] = v
	}
	return h, nil
}

// RESTGet GET helper for mixins.
func RESTGet(ctx context.Context, rawURL string, params map[string]string, apiKey, workspace, module string, extraHeaders map[string]string) (*apientities.DashScopeAPIResponse, error) {
	h, err := restHeaders(apiKey, workspace, module, extraHeaders)
	if err != nil {
		return nil, err
	}
	return httpGetJSON(ctx, rawURL, h, params)
}

// RESTPost POST JSON helper.
func RESTPost(ctx context.Context, rawURL string, payload any, apiKey, workspace, module string, extraHeaders map[string]string) (*apientities.DashScopeAPIResponse, error) {
	h, err := restHeaders(apiKey, workspace, module, extraHeaders)
	if err != nil {
		return nil, err
	}
	h["Content-Type"] = "application/json; charset=utf-8"
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(b))
	}
	return httpPostJSON(ctx, rawURL, h, body)
}

// RESTDelete DELETE helper.
func RESTDelete(ctx context.Context, rawURL string, apiKey, workspace, module string) (*apientities.DashScopeAPIResponse, error) {
	h, err := restHeaders(apiKey, workspace, module, nil)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, rawURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return parseSimpleResponse(resp)
}

// RESTPut PUT JSON helper.
func RESTPut(ctx context.Context, rawURL string, payload any, apiKey, workspace, module string) (*apientities.DashScopeAPIResponse, error) {
	h, err := restHeaders(apiKey, workspace, module, nil)
	if err != nil {
		return nil, err
	}
	h["Content-Type"] = "application/json; charset=utf-8"
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(string(b))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return parseSimpleResponse(resp)
}

// GetURL exported helper.
func GetURL(customBase, defaultPath, path string) string {
	return getURL(customBase, defaultPath, path)
}

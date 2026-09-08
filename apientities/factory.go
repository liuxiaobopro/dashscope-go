// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"context"
	"net/url"
	"strings"

	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/protocol"
)

// RequestOptions options for building an API request.
type RequestOptions struct {
	Model               string
	Input               any
	TaskGroup           string
	Task                string
	Function            string
	APIKey              string
	IsService           bool
	APIProtocol         string
	HTTPMethod          string
	Stream              bool
	AsyncRequest        bool
	RequestTimeout      int
	WSStreamMode        string
	IsBinaryInput       bool
	Query               bool
	Headers             map[string]string
	Form                map[string]any
	Resources           any
	BaseAddress         string
	FlattenedOutput     bool
	ExtraURLParameters  map[string]any
	UserAgent           string
	TaskID              string
	EnableEncryption    bool
	PreTaskID           string
	SDKModule           string
	Parameters          map[string]any
}

// BuildAPIRequest build API request object.
func BuildAPIRequest(opts RequestOptions) (Request, error) {
	if opts.APIProtocol == "" {
		opts.APIProtocol = common.ApiProtocolHTTPS
	}
	if opts.HTTPMethod == "" {
		opts.HTTPMethod = common.HTTPMethodPOST
	}
	if opts.WSStreamMode == "" {
		opts.WSStreamMode = protocol.WebsocketStreamingModeOut
	}
	if !opts.Stream && opts.WSStreamMode == protocol.WebsocketStreamingModeOut {
		opts.WSStreamMode = protocol.WebsocketStreamingModeNone
	}
	if opts.Headers != nil {
		if headerUA, ok := opts.Headers["user-agent"]; ok {
			delete(opts.Headers, "user-agent")
			if opts.UserAgent != "" {
				if headerUA != "" {
					opts.UserAgent = headerUA + "; " + opts.UserAgent
				}
			} else {
				opts.UserAgent = headerUA
			}
		}
	}

	var encryption *Encryption
	var req Request

	switch opts.APIProtocol {
	case common.ApiProtocolHTTP, common.ApiProtocolHTTPS:
		base := opts.BaseAddress
		if base == "" {
			base = common.BaseHTTPAPIURL
		}
		if !strings.HasSuffix(base, "/") {
			base += "/"
		}
		if opts.IsService {
			base += common.SERVICE_API_PATH + "/"
		}
		if opts.TaskGroup != "" {
			base += opts.TaskGroup + "/"
		}
		if opts.Task != "" {
			base += opts.Task + "/"
		}
		if opts.Function != "" {
			base += opts.Function
		}
		if len(opts.ExtraURLParameters) > 0 {
			vals := url.Values{}
			for k, v := range opts.ExtraURLParameters {
				vals.Set(k, fmtString(v))
			}
			base += "?" + vals.Encode()
		}
		if opts.EnableEncryption {
			encryption = &Encryption{}
			encryption.Initialize()
			if encryption.IsValid() {
				common.Log.Debug("encryption enabled")
			}
		}
		httpReq := NewHttpRequest(base, opts.APIKey, opts.HTTPMethod, opts.Stream, opts.AsyncRequest, opts.Query, opts.RequestTimeout, opts.TaskID, opts.FlattenedOutput, encryption, opts.UserAgent)
		req = httpReq
	case common.ApiProtocolWebsocket:
		wsURL := opts.BaseAddress
		if wsURL == "" {
			wsURL = common.BaseWebsocketAPIURL
		}
		wsReq := NewWebSocketRequest(wsURL, opts.APIKey, opts.Stream, opts.WSStreamMode, opts.IsBinaryInput, opts.RequestTimeout, opts.FlattenedOutput, opts.PreTaskID, opts.UserAgent)
		req = wsReq
	default:
		return nil, common.NewUnsupportedApiProtocol(
			"Unsupported protocol: " + opts.APIProtocol + ", support [http, https, websocket]",
		)
	}

	merged := common.GetSDKHeaders(opts.SDKModule)
	if opts.Headers != nil {
		for k, v := range opts.Headers {
			merged[k] = v
		}
	}
	if len(merged) > 0 {
		req.AddHeaders(merged)
	}
	if opts.Input == nil && opts.Form == nil {
		return nil, common.NewInputDataRequired("There is no input data and form data")
	}
	input := opts.Input
	if encryption != nil && encryption.IsValid() {
		input = encryption.Encrypt(input)
	}
	data := NewApiRequestData(opts.Model, opts.TaskGroup, opts.Task, opts.Function, input, opts.Form, opts.IsBinaryInput, opts.APIProtocol)
	data.AddResources(opts.Resources)
	if opts.Parameters != nil {
		data.AddParameters(opts.Parameters)
	}
	req.SetData(data)
	return req, nil
}

// Request unified request interface.
type Request interface {
	AddHeaders(headers map[string]string)
	SetData(data *ApiRequestData)
	CallStream(ctx context.Context) (<-chan *DashScopeAPIResponse, error)
}

// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/protocol"
)

// WebSocketRequest websocket inference request.
type WebSocketRequest struct {
	URL             string
	APIKey          string
	Stream          bool
	WSStreamMode    string
	IsBinaryInput   bool
	Timeout         time.Duration
	FlattenedOutput bool
	PreTaskID       string
	UserAgent       string
	Headers         map[string]string
	Data            *ApiRequestData
}

// NewWebSocketRequest create websocket request.
func NewWebSocketRequest(url, apiKey string, stream bool, wsStreamMode string, isBinary bool, timeout int, flattened bool, preTaskID, userAgent string) *WebSocketRequest {
	to := common.DEFAULT_REQUEST_TIMEOUT_SECONDS
	if timeout > 0 {
		to = timeout
	}
	h := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"user-agent":    common.GetUserAgent(),
	}
	if userAgent != "" {
		h["user-agent"] = h["user-agent"] + "; " + userAgent
	}
	for k, v := range common.GetSDKHeaders("") {
		h[k] = v
	}
	return &WebSocketRequest{
		URL:             url,
		APIKey:          apiKey,
		Stream:          stream,
		WSStreamMode:    wsStreamMode,
		IsBinaryInput:   isBinary,
		Timeout:         time.Duration(to) * time.Second,
		FlattenedOutput: flattened,
		PreTaskID:       preTaskID,
		UserAgent:       userAgent,
		Headers:         h,
	}
}

func (r *WebSocketRequest) AddHeaders(headers map[string]string) {
	for k, v := range headers {
		r.Headers[k] = v
	}
}

func (r *WebSocketRequest) SetData(data *ApiRequestData) {
	r.Data = data
}

func (r *WebSocketRequest) CallStream(ctx context.Context) (<-chan *DashScopeAPIResponse, error) {
	out := make(chan *DashScopeAPIResponse)
	go func() {
		defer close(out)
		if err := r.run(ctx, out); err != nil {
			out <- &DashScopeAPIResponse{StatusCode: -1, Code: "Unknown", Message: err.Error()}
		}
	}()
	return out, nil
}

func (r *WebSocketRequest) run(ctx context.Context, out chan<- *DashScopeAPIResponse) error {
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}
	hdr := http.Header{}
	for k, v := range r.Headers {
		hdr.Set(k, v)
	}
	conn, _, err := websocket.Dial(ctx, r.URL, &websocket.DialOptions{HTTPHeader: hdr})
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	startPayload := map[string]any{}
	if r.Data != nil {
		startPayload = r.Data.GetWebsocketStartData()
	}
	header := map[string]any{
		protocol.ActionKey: protocol.ActionStart,
		"streaming":        r.WSStreamMode,
	}
	if r.PreTaskID != "" {
		header[protocol.TaskID] = r.PreTaskID
	}
	msg := map[string]any{
		protocol.Header: header,
		"payload":       startPayload,
	}
	if err := wsjson.Write(ctx, conn, msg); err != nil {
		return err
	}

	for {
		var raw map[string]any
		if err := wsjson.Read(ctx, conn, &raw); err != nil {
			return err
		}
		h, _ := raw[protocol.Header].(map[string]any)
		event, _ := h[protocol.EventKey].(string)
		taskID, _ := h[protocol.TaskID].(string)
		payload, _ := raw["payload"].(map[string]any)
		switch event {
		case protocol.EventStarted:
			continue
		case protocol.EventGenerated:
			output := payload["output"]
			usage := payload["usage"]
			out <- &DashScopeAPIResponse{RequestID: taskID, StatusCode: http.StatusOK, Output: output, Usage: usage}
			if r.WSStreamMode == protocol.WebsocketStreamingModeNone || r.WSStreamMode == protocol.WebsocketStreamingModeIn {
				return nil
			}
		case protocol.EventFinished:
			output := any(nil)
			usage := any(nil)
			if payload != nil {
				output = payload["output"]
				usage = payload["usage"]
			}
			out <- &DashScopeAPIResponse{RequestID: taskID, StatusCode: http.StatusOK, Output: output, Usage: usage}
			return nil
		case protocol.EventFailed:
			code, _ := h[protocol.ErrorName].(string)
			message, _ := h[protocol.ErrorMessage].(string)
			out <- &DashScopeAPIResponse{RequestID: taskID, StatusCode: common.WEBSOCKET_ERROR_CODE, Code: code, Message: message}
			return nil
		default:
			b, _ := json.Marshal(raw)
			out <- &DashScopeAPIResponse{StatusCode: -1, Code: "Unknown", Message: string(b)}
			return nil
		}
	}
}

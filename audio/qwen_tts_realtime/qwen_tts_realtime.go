// Copyright (c) Alibaba, Inc. and its affiliates.

package qwenttsrealtime

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

type QwenTTSRealtimeCallParams struct {
	Model     string
	Input     map[string]any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type QwenTTSRealtimeService struct{}

func (QwenTTSRealtimeService) Call(ctx context.Context, p *QwenTTSRealtimeCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "audio", Task: "tts", Function: "SpeechSynthesizer",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Extra: p.Extra, SDKModule: "audio",
	})
}

var QwenTTSRealtime = QwenTTSRealtimeService{}

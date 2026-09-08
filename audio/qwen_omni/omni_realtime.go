// Copyright (c) Alibaba, Inc. and its affiliates.

package qwenomni

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

type OmniRealtimeCallParams struct {
	Model     string
	Input     map[string]any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type OmniRealtimeService struct{}

func (OmniRealtimeService) Call(ctx context.Context, p *OmniRealtimeCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "audio", Task: "multimodal-generation", Function: "generation",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Extra: p.Extra, SDKModule: "audio",
	})
}

var OmniRealtime = OmniRealtimeService{}

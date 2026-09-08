// Copyright (c) Alibaba, Inc. and its affiliates.

package multimodal

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const (
	DialogIdle     = "idle"
	DialogListening = "listening"
	DialogThinking = "thinking"
	DialogResponding = "responding"
)

type DialogCallParams struct {
	Model     string
	Input     map[string]any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type MultiModalDialogService struct{}

func (MultiModalDialogService) Call(ctx context.Context, p *DialogCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "multimodal", Task: "dialog", Function: "multimodal-dialog",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Extra: p.Extra, SDKModule: "multimodal",
	})
}

var MultiModalDialog = MultiModalDialogService{}

type TingwuCallParams struct {
	Model     string
	Input     map[string]any
	APIKey    string
	Workspace string
	WaitTimeout int
	Extra     map[string]any
}

type TingwuService struct{}

func (TingwuService) Call(ctx context.Context, p *TingwuCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "multimodal", Task: "tingwu", Function: "tingwu",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "multimodal",
	}, p.WaitTimeout)
}

func (TingwuService) AsyncCall(ctx context.Context, p *TingwuCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).AsyncCall(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "multimodal", Task: "tingwu", Function: "tingwu",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "multimodal",
	})
}

var Tingwu = TingwuService{}

type TingwuRealtimeService struct{}

func (TingwuRealtimeService) Call(ctx context.Context, p *DialogCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "multimodal", Task: "tingwu", Function: "tingwu-realtime",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Extra: p.Extra, SDKModule: "multimodal",
	})
}

var TingwuRealtime = TingwuRealtimeService{}

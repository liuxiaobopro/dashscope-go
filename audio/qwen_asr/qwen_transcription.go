// Copyright (c) Alibaba, Inc. and its affiliates.

package qwenasr

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
)

type QwenTranscriptionCallParams struct {
	Model     string
	Messages  []any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type QwenTranscriptionService struct{}

func (QwenTranscriptionService) Call(ctx context.Context, p *QwenTranscriptionCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: map[string]any{"messages": p.Messages},
		TaskGroup: "audio", Task: "asr", Function: "qwen-transcription",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "audio",
	})
}

var QwenTranscription = QwenTranscriptionService{}

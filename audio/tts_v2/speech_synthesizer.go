// Copyright (c) Alibaba, Inc. and its affiliates.

package ttsv2

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

type SpeechSynthesizerCallParams struct {
	Model     string
	Text      string
	Voice     string
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type SpeechSynthesizerService struct{}

func (SpeechSynthesizerService) Call(ctx context.Context, p *SpeechSynthesizerCallParams) (*apientities.SpeechSynthesisResponse, error) {
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: map[string]any{"text": p.Text},
		TaskGroup: "audio", Task: "tts", Function: "SpeechSynthesizer",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Parameters: map[string]any{"voice": p.Voice}, Extra: p.Extra, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.SpeechSynthesisResponseFromAPI(rsp), nil
}

var SpeechSynthesizer = SpeechSynthesizerService{}

type EnrollmentService struct{}

func (EnrollmentService) Create(ctx context.Context, model, url, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: model, Input: map[string]any{"url": url},
		TaskGroup: "audio", Task: "tts", Function: "enrollment",
		APIKey: apiKey, Workspace: workspace, SDKModule: "audio",
	})
}

var Enrollment = EnrollmentService{}

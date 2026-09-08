// Copyright (c) Alibaba, Inc. and its affiliates.

package tts

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

// SpeechSynthesizerCallParams speech synthesis via websocket.
type SpeechSynthesizerCallParams struct {
	Model      string
	Text       string
	Voice      string
	Format     string
	SampleRate int
	APIKey     string
	Workspace  string
	Extra      map[string]any
}

type SpeechSynthesizerService struct{}

func (SpeechSynthesizerService) Call(ctx context.Context, p *SpeechSynthesizerCallParams) (*apientities.SpeechSynthesisResponse, error) {
	input := map[string]any{"text": p.Text}
	params := map[string]any{}
	if p.Voice != "" {
		params["voice"] = p.Voice
	}
	if p.Format != "" {
		params["format"] = p.Format
	}
	if p.SampleRate != 0 {
		params["sample_rate"] = p.SampleRate
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "audio", Task: "tts", Function: "speech-synthesizer",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		Parameters: params, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.SpeechSynthesisResponseFromAPI(rsp), nil
}

var SpeechSynthesizer = SpeechSynthesizerService{}

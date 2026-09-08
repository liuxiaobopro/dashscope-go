// Copyright (c) Alibaba, Inc. and its affiliates.

package httptts

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
)

// HttpSpeechSynthesizerCallParams HTTP TTS.
type HttpSpeechSynthesizerCallParams struct {
	Model     string
	Input     map[string]any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type HttpSpeechSynthesizerService struct{}

func (HttpSpeechSynthesizerService) Call(ctx context.Context, p *HttpSpeechSynthesizerCallParams) (*apientities.TextToSpeechResponse, error) {
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "audio", Task: "tts", Function: "SpeechSynthesizer",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.TextToSpeechResponseFromAPI(rsp), nil
}

var HttpSpeechSynthesizer = HttpSpeechSynthesizerService{}

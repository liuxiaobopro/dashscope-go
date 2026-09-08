// Copyright (c) Alibaba, Inc. and its affiliates.

package qwentts

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
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

func (SpeechSynthesizerService) Call(ctx context.Context, p *SpeechSynthesizerCallParams) (*apientities.TextToSpeechResponse, error) {
	input := map[string]any{"text": p.Text}
	if p.Voice != "" {
		input["voice"] = p.Voice
	}
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "audio", Task: "tts", Function: "SpeechSynthesizer",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.TextToSpeechResponseFromAPI(rsp), nil
}

var SpeechSynthesizer = SpeechSynthesizerService{}

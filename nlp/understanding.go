// Copyright (c) Alibaba, Inc. and its affiliates.

package nlp

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const OpenNLUV1 = "opennlu-v1"

// UnderstandingCallParams Call generation model service.
type UnderstandingCallParams struct {
	// Model The requested model, such as opennlu-v1
	Model string
	// Sentence The text content entered by the user that needs to be processed supports both Chinese and English.
	Sentence string
	// Labels For the extraction task, label is the name of the type that needs to be extracted. For classification tasks, label is the classification system.
	Labels string
	// Task Task type, optional as extraction or classification, default as extraction.
	Task   string
	APIKey string
	Extra  map[string]any
}

type UnderstandingService struct{}

func (UnderstandingService) Call(ctx context.Context, p *UnderstandingCallParams) (*apientities.DashScopeAPIResponse, error) {
	if p.Sentence == "" || p.Labels == "" {
		return nil, common.NewInputRequired("sentence and labels is required!")
	}
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	if _, ok := p.Extra["stream"]; ok {
		common.Log.Warn("stream option not supported for Understanding.")
		delete(p.Extra, "stream")
	}
	input := map[string]any{"sentence": p.Sentence, "labels": p.Labels}
	if p.Task != "" {
		input["task"] = p.Task
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "nlp", Task: "understanding", Function: "understanding",
		APIKey: p.APIKey, Extra: p.Extra, SDKModule: "nlp",
	})
}

var Understanding = UnderstandingService{}

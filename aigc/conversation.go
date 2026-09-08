// Copyright (c) Alibaba, Inc. and its affiliates.

package aigc

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

// HistoryItem A conversation history item.
type HistoryItem map[string]any

// NewHistoryItem Init a history item.
func NewHistoryItem(role string, text string, extra map[string]any) HistoryItem {
	common.Log.Warn("%s", common.DEPRECATED_MESSAGE)
	item := HistoryItem{role: []any{}}
	if text != "" {
		item[role] = append(item[role].([]any), map[string]any{"text": text})
	}
	for k, v := range extra {
		item[role] = append(item[role].([]any), map[string]any{k: v})
	}
	return item
}

// History Manage the conversation history.
type History []HistoryItem

// ConversationCallParams conversation call params.
type ConversationCallParams struct {
	Model             string
	Prompt            string
	History           History
	APIKey            string
	Workspace         string
	Messages          []apientities.Message
	Stream            bool
	Temperature       *float64
	TopP              *float64
	TopK              *int
	MaxTokens         *int
	ResultFormat      *string
	IncrementalOutput *bool
	EnableSearch      *bool
	Extra             map[string]any
}

// Conversation conversational robot interface.
type ConversationService struct{}

func (ConversationService) Call(ctx context.Context, p *ConversationCallParams) (*apientities.GenerationResponse, error) {
	if p.Prompt == "" && len(p.Messages) == 0 {
		return nil, common.NewInputRequired("prompt or messages is required!")
	}
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	input := map[string]any{}
	params := map[string]any{}
	if len(p.History) > 0 {
		common.Log.Warn("%s", common.DEPRECATED_MESSAGE)
		input[common.HISTORY] = p.History
	}
	if p.Prompt != "" {
		input[common.PROMPT] = p.Prompt
	}
	if len(p.Messages) > 0 {
		input[common.MESSAGES] = p.Messages
	}
	if p.Temperature != nil {
		params["temperature"] = *p.Temperature
	}
	if p.TopP != nil {
		params["top_p"] = *p.TopP
	}
	if p.TopK != nil {
		params["top_k"] = *p.TopK
	}
	if p.MaxTokens != nil {
		params["max_tokens"] = *p.MaxTokens
	}
	if p.ResultFormat != nil {
		params["result_format"] = *p.ResultFormat
	}
	if p.IncrementalOutput != nil {
		params["incremental_output"] = *p.IncrementalOutput
	}
	if p.EnableSearch != nil {
		params["enable_search"] = *p.EnableSearch
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "aigc", Task: "text-generation", Function: "generation",
		APIKey: p.APIKey, Workspace: p.Workspace, Stream: p.Stream, Parameters: params, SDKModule: "aigc",
	})
	if err != nil {
		return nil, err
	}
	return apientities.GenerationResponseFromAPI(rsp), nil
}

// Conversation conversational robot interface.
var Conversation = ConversationService{}

// Copyright (c) Alibaba, Inc. and its affiliates.

package aigc

import (
	"context"
	"iter"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/utils"
)

const mmTask = "multimodal-generation"
const mmFunction = "generation"

// MultiModalConversationCallParams Call the conversation model service.
type MultiModalConversationCallParams struct {
	// Model The requested model, such as 'qwen-vl-max'.
	Model string
	// Messages The generation messages.
	Messages []any
	APIKey   string
	// Workspace The dashscope workspace id.
	Workspace string
	// Text The text to generate.
	Text string
	// Voice The voice name for qwen tts.
	Voice string
	// LanguageType The synthesized language type.
	LanguageType      string
	Stream            bool
	Temperature       *float64
	TopP              *float64
	TopK              *int
	MaxTokens         *int
	Seed              *int
	Stop              any
	RepetitionPenalty *float64
	PresencePenalty   *float64
	ResultFormat      *string
	IncrementalOutput *bool
	EnableSearch      *bool
	Tools             []map[string]any
	ToolChoice        any
	EnableThinking    *bool
	N                 *int
	OCROptions        map[string]any
	Logprobs          *bool
	TopLogprobs       *int
	Extra             map[string]any
	Headers           map[string]string
	// BaseAddress overrides the default HTTP API base (e.g. https://dashscope.aliyuncs.com/api/v1).
	BaseAddress string
}

func (p *MultiModalConversationCallParams) toParams() map[string]any {
	kw := map[string]any{}
	if p.Stream {
		kw["stream"] = true
	}
	if p.Temperature != nil {
		kw["temperature"] = *p.Temperature
	}
	if p.TopP != nil {
		kw["top_p"] = *p.TopP
	}
	if p.TopK != nil {
		kw["top_k"] = *p.TopK
	}
	if p.MaxTokens != nil {
		kw["max_tokens"] = *p.MaxTokens
	}
	if p.Seed != nil {
		kw["seed"] = *p.Seed
	}
	if p.Stop != nil {
		kw["stop"] = p.Stop
	}
	if p.RepetitionPenalty != nil {
		kw["repetition_penalty"] = *p.RepetitionPenalty
	}
	if p.PresencePenalty != nil {
		kw["presence_penalty"] = *p.PresencePenalty
	}
	if p.ResultFormat != nil {
		kw["result_format"] = *p.ResultFormat
	}
	if p.IncrementalOutput != nil {
		kw["incremental_output"] = *p.IncrementalOutput
	}
	if p.EnableSearch != nil {
		kw["enable_search"] = *p.EnableSearch
	}
	if p.Tools != nil {
		kw["tools"] = p.Tools
	}
	if p.ToolChoice != nil {
		kw["tool_choice"] = p.ToolChoice
	}
	if p.EnableThinking != nil {
		kw["enable_thinking"] = *p.EnableThinking
	}
	if p.N != nil {
		kw["n"] = *p.N
	}
	if p.OCROptions != nil {
		kw["ocr_options"] = p.OCROptions
	}
	if p.Logprobs != nil {
		kw["logprobs"] = *p.Logprobs
	}
	if p.TopLogprobs != nil {
		kw["top_logprobs"] = *p.TopLogprobs
	}
	for k, v := range p.Extra {
		kw[k] = v
	}
	return kw
}

func prepareMM(p *MultiModalConversationCallParams) (*client.CallParams, bool, int, error) {
	if p.Model == "" {
		return nil, false, 1, common.NewModelRequired("Model is required!")
	}
	input := map[string]any{}
	if len(p.Messages) > 0 {
		input[common.MESSAGES] = p.Messages
	}
	if p.Text != "" {
		input["text"] = p.Text
	}
	if p.Voice != "" {
		input["voice"] = p.Voice
	}
	if p.LanguageType != "" {
		input["language_type"] = p.LanguageType
	}
	kw := p.toParams()
	isStream := p.Stream
	toMerge := false
	if utils.ShouldModifyIncrementalOutput(p.Model) && isStream && p.IncrementalOutput != nil && !*p.IncrementalOutput {
		toMerge = true
		kw["incremental_output"] = true
	}
	n := 1
	if p.N != nil {
		n = *p.N
	}
	return &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "aigc", Task: mmTask, Function: mmFunction,
		APIKey: p.APIKey, Workspace: p.Workspace, Stream: isStream, Headers: p.Headers, Parameters: kw,
		BaseAddress: p.BaseAddress, SDKModule: "aigc",
	}, toMerge, n, nil
}

// MultiModalConversationService MultiModal conversational robot interface.
type MultiModalConversationService struct{}

const MultiModalQwenVLChatV1 = "qwen-vl-chat-v1"

func (MultiModalConversationService) Call(ctx context.Context, p *MultiModalConversationCallParams) (*apientities.MultiModalConversationResponse, error) {
	cp, _, _, err := prepareMM(p)
	if err != nil {
		return nil, err
	}
	rsp, err := (client.BaseApi{}).Call(ctx, cp)
	if err != nil {
		return nil, err
	}
	return apientities.MultiModalConversationResponseFromAPI(rsp), nil
}

func (MultiModalConversationService) CallStream(ctx context.Context, p *MultiModalConversationCallParams) iter.Seq2[*apientities.MultiModalConversationResponse, error] {
	p.Stream = true
	return func(yield func(*apientities.MultiModalConversationResponse, error) bool) {
		cp, toMerge, n, err := prepareMM(p)
		if err != nil {
			yield(nil, err)
			return
		}
		ch, err := (client.BaseApi{}).CallStream(ctx, cp)
		if err != nil {
			yield(nil, err)
			return
		}
		accumulated := map[any]any{}
		for rsp := range ch {
			parsed := apientities.MultiModalConversationResponseFromAPI(rsp)
			if toMerge {
				ok, _ := utils.MergeMultimodalSingleResponse(parsed, accumulated, n)
				if ok && !yield(parsed, nil) {
					return
				}
			} else if !yield(parsed, nil) {
				return
			}
		}
	}
}

// MultiModalConversation MultiModal conversational robot interface.
var MultiModalConversation = MultiModalConversationService{}

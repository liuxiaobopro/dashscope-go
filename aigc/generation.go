// Copyright (c) Alibaba, Inc. and its affiliates.

package aigc

import (
	"context"
	"encoding/json"
	"iter"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/utils"
)

const generationTask = "text-generation"

// GenerationModels model names.
type GenerationModels struct{}

const (
	GenerationQwenV1     = "qwen-v1"      // deprecated, use qwen_turbo instead
	GenerationQwenPlusV1 = "qwen-plus-v1" // deprecated, use qwen_plus instead
	GenerationBailianV1  = "bailian-v1"
	GenerationDolly12bV2 = "dolly-12b-v2"
	GenerationQwenTurbo  = "qwen-turbo"
	GenerationQwenPlus   = "qwen-plus"
	GenerationQwenMax    = "qwen-max"
)

// GenerationCallParams Call generation model service.
type GenerationCallParams struct {
	// Model The requested model, such as qwen-turbo.
	Model string
	// Prompt The input prompt.
	Prompt any
	// History The user provided history, deprecated.
	History []any
	// APIKey The api api_key, can be None.
	APIKey string
	// Messages The generation messages.
	Messages []apientities.Message
	// Plugins The plugin config, str or dict.
	Plugins any
	// Workspace The dashscope workspace id.
	Workspace string
	// Stream Enable streaming output.
	Stream bool
	// Temperature Controls randomness, range [0, 2).
	Temperature *float64
	// TopP Nucleus sampling, range (0, 1.0].
	TopP *float64
	// TopK Size of candidate token set for sampling.
	TopK *int
	// MaxTokens Maximum output token count.
	MaxTokens *int
	// Seed Random seed for reproducibility.
	Seed *int
	// Stop Stop sequences.
	Stop any
	// RepetitionPenalty Penalizes repeated sequences. 1.0 means no penalty.
	RepetitionPenalty *float64
	// PresencePenalty Controls content repetition, range [-2.0, 2.0].
	PresencePenalty *float64
	// ResultFormat "message" or "text".
	ResultFormat *string
	// IncrementalOutput In streaming mode, output only new tokens (True) vs. cumulative output (False).
	IncrementalOutput *bool
	// EnableSearch Enable web search.
	EnableSearch *bool
	// Tools Tool definitions for function calling.
	Tools []map[string]any
	// ToolChoice Tool selection strategy.
	ToolChoice any
	// EnableThinking Enable thinking mode for hybrid thinking models.
	EnableThinking *bool
	// ThinkingBudget Maximum token budget for thinking mode.
	ThinkingBudget *int
	// N Number of responses to generate (1-4).
	N *int
	// Logprobs Whether to return log probabilities of the output tokens.
	Logprobs *bool
	// TopLogprobs Number of most likely tokens to return at each token position when logprobs is enabled.
	TopLogprobs *int
	// SearchOptions Configuration options for web search feature.
	SearchOptions map[string]any
	// ParallelToolCalls Enable parallel tool calls for function calling.
	ParallelToolCalls *bool
	// ResponseFormat Format constraint for response, e.g., {"type": "json_object"} for JSON mode.
	ResponseFormat map[string]any
	// OutputFormat Output format for qwen-deep-research model.
	OutputFormat *string
	// Extra Additional parameters passed to the API.
	Extra   map[string]any
	Headers map[string]string
	// BaseAddress overrides the default HTTP API base (e.g. https://dashscope.aliyuncs.com/api/v1).
	BaseAddress string
}

func (p *GenerationCallParams) kwargs() map[string]any {
	kw := map[string]any{}
	if p.Stream {
		kw["stream"] = true
	}
	put := func(k string, v any) {
		if v != nil {
			kw[k] = v
		}
	}
	if p.Temperature != nil {
		put("temperature", *p.Temperature)
	}
	if p.TopP != nil {
		put("top_p", *p.TopP)
	}
	if p.TopK != nil {
		put("top_k", *p.TopK)
	}
	if p.MaxTokens != nil {
		put("max_tokens", *p.MaxTokens)
	}
	if p.Seed != nil {
		put("seed", *p.Seed)
	}
	if p.Stop != nil {
		put("stop", p.Stop)
	}
	if p.RepetitionPenalty != nil {
		put("repetition_penalty", *p.RepetitionPenalty)
	}
	if p.PresencePenalty != nil {
		put("presence_penalty", *p.PresencePenalty)
	}
	if p.ResultFormat != nil {
		put("result_format", *p.ResultFormat)
	}
	if p.IncrementalOutput != nil {
		put("incremental_output", *p.IncrementalOutput)
	}
	if p.EnableSearch != nil {
		put("enable_search", *p.EnableSearch)
	}
	if p.Tools != nil {
		put("tools", p.Tools)
	}
	if p.ToolChoice != nil {
		put("tool_choice", p.ToolChoice)
	}
	if p.EnableThinking != nil {
		put("enable_thinking", *p.EnableThinking)
	}
	if p.ThinkingBudget != nil {
		put("thinking_budget", *p.ThinkingBudget)
	}
	if p.N != nil {
		put("n", *p.N)
	}
	if p.Logprobs != nil {
		put("logprobs", *p.Logprobs)
	}
	if p.TopLogprobs != nil {
		put("top_logprobs", *p.TopLogprobs)
	}
	if p.SearchOptions != nil {
		put("search_options", p.SearchOptions)
	}
	if p.ParallelToolCalls != nil {
		put("parallel_tool_calls", *p.ParallelToolCalls)
	}
	if p.ResponseFormat != nil {
		put("response_format", p.ResponseFormat)
	}
	if p.OutputFormat != nil {
		put("output_format", *p.OutputFormat)
	}
	for k, v := range p.Extra {
		kw[k] = v
	}
	return kw
}

func buildGenerationInput(model string, prompt any, history []any, messages []apientities.Message, kw map[string]any) (map[string]any, map[string]any, error) {
	if model == GenerationQwenV1 {
		common.Log.Warn("Model %s is deprecated, use %s instead!", GenerationQwenV1, GenerationQwenTurbo)
	}
	if model == GenerationQwenPlusV1 {
		common.Log.Warn("Model %s is deprecated, use %s instead!", GenerationQwenPlusV1, GenerationQwenPlus)
	}
	parameters := map[string]any{}
	input := map[string]any{}
	if history != nil {
		common.Log.Warn("%s", common.DEPRECATED_MESSAGE)
		input[common.HISTORY] = history
	}
	if prompt != nil && prompt != "" {
		input[common.PROMPT] = prompt
	} else if messages != nil {
		msgs := make([]any, 0, len(messages)+1)
		for _, m := range messages {
			b, _ := json.Marshal(m)
			var mm map[string]any
			_ = json.Unmarshal(b, &mm)
			msgs = append(msgs, mm)
		}
		if prompt != nil && prompt != "" {
			msgs = append(msgs, map[string]any{"role": apientities.RoleUser, "content": prompt})
		}
		input = map[string]any{common.MESSAGES: msgs}
	} else if prompt != nil {
		input[common.PROMPT] = prompt
	}
	if len(model) >= 4 && model[:4] == "qwen" {
		if v, ok := kw["enable_search"]; ok {
			delete(kw, "enable_search")
			if b, ok := v.(bool); ok && b {
				parameters["enable_search"] = true
			}
		}
	} else if len(model) >= 7 && model[:7] == "bailian" {
		customized, ok := kw["customized_model_id"]
		delete(kw, "customized_model_id")
		if !ok || customized == nil {
			return nil, nil, common.NewInputRequired("customized_model_id is required for " + model)
		}
		input[common.CUSTOMIZED_MODEL_ID] = customized
	}
	for k, v := range kw {
		parameters[k] = v
	}
	return input, parameters, nil
}

func hasPromptOrMessages(p *GenerationCallParams) bool {
	if p.Prompt != nil && p.Prompt != "" {
		return true
	}
	return len(p.Messages) > 0
}

func prepareGeneration(p *GenerationCallParams) (*client.CallParams, bool, int, error) {
	if !hasPromptOrMessages(p) {
		return nil, false, 1, common.NewInputRequired("prompt or messages is required!")
	}
	if p.Model == "" {
		return nil, false, 1, common.NewModelRequired("Model is required!")
	}
	kw := p.kwargs()
	headers := map[string]string{}
	for k, v := range p.Headers {
		headers[k] = v
	}
	if p.Plugins != nil {
		switch v := p.Plugins.(type) {
		case string:
			headers["X-DashScope-Plugin"] = v
		default:
			b, _ := json.Marshal(v)
			headers["X-DashScope-Plugin"] = string(b)
		}
	}
	input, parameters, err := buildGenerationInput(p.Model, p.Prompt, p.History, p.Messages, kw)
	if err != nil {
		return nil, false, 1, err
	}
	isStream := false
	if v, ok := parameters["stream"].(bool); ok {
		isStream = v
	}
	var isInc *bool
	if p.IncrementalOutput != nil {
		isInc = p.IncrementalOutput
	}
	toMerge := false
	if utils.ShouldModifyIncrementalOutput(p.Model) && isStream && isInc != nil && !*isInc {
		toMerge = true
		parameters["incremental_output"] = true
	}
	flag := "0"
	if toMerge {
		flag = "1"
	}
	newUA := "incremental_to_full/" + flag
	if existing, ok := parameters["user_agent"].(string); ok && existing != "" {
		parameters["user_agent"] = existing + "; " + newUA
	} else {
		parameters["user_agent"] = newUA
	}
	n := 1
	if v, ok := parameters["n"].(int); ok {
		n = v
	}
	ua, _ := parameters["user_agent"].(string)
	delete(parameters, "user_agent")
	cp := &client.CallParams{
		Model:       p.Model,
		Input:       input,
		TaskGroup:   "aigc",
		Task:        generationTask,
		Function:    "generation",
		APIKey:      p.APIKey,
		Workspace:   p.Workspace,
		Stream:      isStream,
		Headers:     headers,
		UserAgent:   ua,
		Parameters:  parameters,
		BaseAddress: p.BaseAddress,
		SDKModule:   "aigc",
	}
	return cp, toMerge, n, nil
}

// GenerationService API for AI-Generated Content(AIGC) models.
type GenerationService struct{}

// Call Call generation model service.
func (GenerationService) Call(ctx context.Context, p *GenerationCallParams) (*apientities.GenerationResponse, error) {
	if p.Stream {
		return nil, common.NewInvalidParameter("use CallStream when Stream is true")
	}
	cp, _, _, err := prepareGeneration(p)
	if err != nil {
		return nil, err
	}
	rsp, err := (client.BaseApi{}).Call(ctx, cp)
	if err != nil {
		return nil, err
	}
	return apientities.GenerationResponseFromAPI(rsp), nil
}

// CallStream Call generation model service with streaming.
func (GenerationService) CallStream(ctx context.Context, p *GenerationCallParams) iter.Seq2[*apientities.GenerationResponse, error] {
	p.Stream = true
	return func(yield func(*apientities.GenerationResponse, error) bool) {
		cp, toMerge, n, err := prepareGeneration(p)
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
			parsed := apientities.GenerationResponseFromAPI(rsp)
			if toMerge {
				ok, extras := utils.MergeSingleResponse(parsed, accumulated, n)
				if ok {
					if !yield(parsed, nil) {
						return
					}
				}
				for _, e := range extras {
					if !yield(e, nil) {
						return
					}
				}
			} else {
				if !yield(parsed, nil) {
					return
				}
			}
		}
	}
}

// Generation API for AI-Generated Content(AIGC) models.
var Generation = GenerationService{}

// Copyright (c) Alibaba, Inc. and its affiliates.

package aigc

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const TongyiLingmaV1 = "tongyi-lingma-v1"

// CodeGenerationCallParams API for AI-Generated Content(AIGC) models.
type CodeGenerationCallParams struct {
	Model     string
	Messages  []any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type CodeGenerationService struct{}

func (CodeGenerationService) Call(ctx context.Context, p *CodeGenerationCallParams) (*apientities.DashScopeAPIResponse, error) {
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	if len(p.Messages) == 0 {
		return nil, common.NewInputRequired("messages is required!")
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: map[string]any{common.MESSAGES: p.Messages},
		TaskGroup: "aigc", Task: "code-generation", Function: "generation",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "aigc",
	})
}

var CodeGeneration = CodeGenerationService{}

// ChatCompletionCallParams OpenAI-compatible chat completion.
type ChatCompletionCallParams struct {
	Model       string
	Messages    []any
	APIKey      string
	Workspace   string
	Stream      bool
	Temperature *float64
	TopP        *float64
	MaxTokens   *int
	Tools       []map[string]any
	Extra       map[string]any
}

type ChatCompletionService struct{}

func (ChatCompletionService) Call(ctx context.Context, p *ChatCompletionCallParams) (*apientities.DashScopeAPIResponse, error) {
	params := map[string]any{}
	if p.Temperature != nil {
		params["temperature"] = *p.Temperature
	}
	if p.TopP != nil {
		params["top_p"] = *p.TopP
	}
	if p.MaxTokens != nil {
		params["max_tokens"] = *p.MaxTokens
	}
	if p.Tools != nil {
		params["tools"] = p.Tools
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: map[string]any{common.MESSAGES: p.Messages},
		TaskGroup: "aigc", Task: "text-generation", Function: "generation",
		APIKey: p.APIKey, Workspace: p.Workspace, Stream: p.Stream, Parameters: params,
		BaseAddress: common.BaseCompatibleAPIURL, SDKModule: "aigc",
	})
}

var ChatCompletion = ChatCompletionService{}

// ImageGenerationCallParams image generation (qwen-image etc).
type ImageGenerationCallParams struct {
	Model        string
	Messages     []any
	Prompt       any
	APIKey       string
	Workspace    string
	Size         *string
	N            *int
	PromptExtend *bool
	Watermark    *bool
	Extra        map[string]any
	WaitTimeout  int
	BaseAddress  string
}

func (p *ImageGenerationCallParams) toCall() *client.CallParams {
	input := map[string]any{}
	if len(p.Messages) > 0 {
		input[common.MESSAGES] = p.Messages
	}
	if p.Prompt != nil {
		input[common.PROMPT] = p.Prompt
	}
	params := map[string]any{}
	if p.Size != nil {
		params["size"] = *p.Size
	}
	if p.N != nil {
		params["n"] = *p.N
	}
	if p.PromptExtend != nil {
		params["prompt_extend"] = *p.PromptExtend
	}
	if p.Watermark != nil {
		params["watermark"] = *p.Watermark
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	return &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "aigc", Task: "image-generation", Function: "generation",
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: params, BaseAddress: p.BaseAddress, SDKModule: "aigc",
	}
}

type ImageGenerationService struct{}

func (ImageGenerationService) Call(ctx context.Context, p *ImageGenerationCallParams) (*apientities.ImageGenerationResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Call(ctx, p.toCall(), p.WaitTimeout)
	if err != nil {
		return nil, err
	}
	return apientities.ImageGenerationResponseFromAPI(rsp), nil
}

func (ImageGenerationService) AsyncCall(ctx context.Context, p *ImageGenerationCallParams) (*apientities.ImageGenerationResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).AsyncCall(ctx, p.toCall())
	if err != nil {
		return nil, err
	}
	return apientities.ImageGenerationResponseFromAPI(rsp), nil
}

func (ImageGenerationService) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int, baseAddress string) (*apientities.ImageGenerationResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Wait(ctx, task, apiKey, workspace, waitTimeout, baseAddress)
	if err != nil {
		return nil, err
	}
	return apientities.ImageGenerationResponseFromAPI(rsp), nil
}

var ImageGeneration = ImageGenerationService{}

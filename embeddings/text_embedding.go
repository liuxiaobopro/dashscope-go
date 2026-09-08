// Copyright (c) Alibaba, Inc. and its affiliates.

package embeddings

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const (
	TextEmbeddingV1 = "text-embedding-v1"
	TextEmbeddingV2 = "text-embedding-v2"
	TextEmbeddingV3 = "text-embedding-v3"
	TextEmbeddingV4 = "text-embedding-v4"
)

// TextEmbeddingCallParams Get embedding of text input.
type TextEmbeddingCallParams struct {
	// Model The embedding model name.
	Model string
	// Input The text input, can be a text or list of text.
	Input any
	Workspace string
	APIKey    string
	// TextType "query" for search queries, "document" (default) for corpus/symmetric tasks.
	TextType *string
	// Dimension Output vector dimension. Options: 2048 (v4 only), 1536 (v4 only), 1024 (default), 768, 512, 256, 128, 64. Only for v3/v4.
	Dimension *int
	// OutputType Output format: "dense" (default), "sparse", or "dense&sparse". Only for v3/v4.
	OutputType *string
	// Instruct Custom task instruction to guide model understanding of query intent.
	Instruct *string
	Extra    map[string]any
}

// TextEmbeddingService text embedding API.
type TextEmbeddingService struct {
	Task string
}

func (TextEmbeddingService) Call(ctx context.Context, p *TextEmbeddingCallParams) (*apientities.DashScopeAPIResponse, error) {
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	embeddingInput := map[string]any{}
	switch v := p.Input.(type) {
	case string:
		embeddingInput[common.TEXT_EMBEDDING_INPUT_KEY] = []string{v}
	default:
		embeddingInput[common.TEXT_EMBEDDING_INPUT_KEY] = p.Input
	}
	kw := map[string]any{}
	if p.TextType != nil {
		kw["text_type"] = *p.TextType
	}
	if p.Dimension != nil {
		kw["dimension"] = *p.Dimension
	}
	if p.OutputType != nil {
		kw["output_type"] = *p.OutputType
	}
	if p.Instruct != nil {
		kw["instruct"] = *p.Instruct
	}
	for k, v := range p.Extra {
		kw[k] = v
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: embeddingInput, TaskGroup: "embeddings", Task: "text-embedding", Function: "text-embedding",
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: kw, SDKModule: "embeddings",
	})
}

// TextEmbedding text embedding API.
var TextEmbedding = TextEmbeddingService{Task: "text-embedding"}

// MultiModalEmbeddingItemBase multimodal embedding item.
type MultiModalEmbeddingItemBase struct {
	Factor float64 `json:"factor"`
	Text   string  `json:"text,omitempty"`
	Image  string  `json:"image,omitempty"`
	Audio  string  `json:"audio,omitempty"`
}

func NewMultiModalEmbeddingItemText(text string, factor float64) MultiModalEmbeddingItemBase {
	return MultiModalEmbeddingItemBase{Text: text, Factor: factor}
}
func NewMultiModalEmbeddingItemImage(image string, factor float64) MultiModalEmbeddingItemBase {
	return MultiModalEmbeddingItemBase{Image: image, Factor: factor}
}
func NewMultiModalEmbeddingItemAudio(audio string, factor float64) MultiModalEmbeddingItemBase {
	return MultiModalEmbeddingItemBase{Audio: audio, Factor: factor}
}

const (
	MultimodalEmbeddingOnePeaceV1   = "multimodal-embedding-one-peace-v1"
	MultimodalEmbeddingV1           = "multimodal-embedding-v1"
	Qwen3VLEmbedding                = "qwen3-vl-embedding"
	Qwen25VLEmbedding               = "qwen2.5-vl-embedding"
	TongyiEmbeddingVisionPlus       = "tongyi-embedding-vision-plus"
	TongyiEmbeddingVisionFlash      = "tongyi-embedding-vision-flash"
)

// MultiModalEmbeddingCallParams Get embedding multimodal contents.
type MultiModalEmbeddingCallParams struct {
	Model          string
	Input          []MultiModalEmbeddingItemBase
	APIKey         string
	Workspace      string
	Dimension      *int
	OutputType     *string
	FPS            *float64
	Instruct       *string
	EnableFusion   *bool
	ResLevel       *int
	MaxVideoFrames *int
	Extra          map[string]any
}

type MultiModalEmbeddingService struct {
	Task string
}

func (MultiModalEmbeddingService) Call(ctx context.Context, p *MultiModalEmbeddingCallParams) (*apientities.DashScopeAPIResponse, error) {
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	if len(p.Input) == 0 {
		return nil, common.NewInputRequired("input is required!")
	}
	kw := map[string]any{}
	if p.Dimension != nil {
		kw["dimension"] = *p.Dimension
	}
	if p.OutputType != nil {
		kw["output_type"] = *p.OutputType
	}
	if p.FPS != nil {
		kw["fps"] = *p.FPS
	}
	if p.Instruct != nil {
		kw["instruct"] = *p.Instruct
	}
	if p.EnableFusion != nil {
		kw["enable_fusion"] = *p.EnableFusion
	}
	if p.ResLevel != nil {
		kw["res_level"] = *p.ResLevel
	}
	if p.MaxVideoFrames != nil {
		kw["max_video_frames"] = *p.MaxVideoFrames
	}
	for k, v := range p.Extra {
		kw[k] = v
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: map[string]any{"contents": p.Input}, TaskGroup: "embeddings",
		Task: "multimodal-embedding", Function: "multimodal-embedding",
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: kw, SDKModule: "embeddings",
	})
}

var MultiModalEmbedding = MultiModalEmbeddingService{Task: "multimodal-embedding"}

// BatchTextEmbeddingCallParams batch text embedding.
type BatchTextEmbeddingCallParams struct {
	Model     string
	Input     any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type BatchTextEmbeddingService struct{}

func (BatchTextEmbeddingService) Call(ctx context.Context, p *BatchTextEmbeddingCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "embeddings", Task: "text-embedding",
		Function: "batch-text-embedding", APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "embeddings",
	}, -1)
}

func (BatchTextEmbeddingService) AsyncCall(ctx context.Context, p *BatchTextEmbeddingCallParams) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).AsyncCall(ctx, &client.CallParams{
		Model: p.Model, Input: p.Input, TaskGroup: "embeddings", Task: "text-embedding",
		Function: "batch-text-embedding", APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "embeddings",
	})
}

func (BatchTextEmbeddingService) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Wait(ctx, task, apiKey, workspace, waitTimeout, "")
}

func (BatchTextEmbeddingService) Fetch(ctx context.Context, task any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Fetch(ctx, task, apiKey, workspace, "")
}

var BatchTextEmbedding = BatchTextEmbeddingService{}

type BatchTextEmbeddingResponse = apientities.DashScopeAPIResponse

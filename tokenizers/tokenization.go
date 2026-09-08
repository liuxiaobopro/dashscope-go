// Copyright (c) Alibaba, Inc. and its affiliates.

package tokenizers

import (
	"context"
	"unicode/utf8"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const (
	TokenizerQwenTurbo     = "qwen-turbo"
	TokenizerQwenPlus      = "qwen-plus"
	TokenizerQwen7bChat    = "qwen-7b-chat"
	TokenizerQwen14bChat   = "qwen-14b-chat"
	TokenizerLlama27bChatV2 = "llama2-7b-chat-v2"
	TokenizerLlama213bChatV2 = "llama2-13b-chat-v2"
	TokenizerTextEmbeddingV2 = "text-embedding-v2"
	TokenizerQwen72bChat   = "qwen-72b-chat"
)

type TokenizationCallParams struct {
	Model     string
	Prompt    any
	Messages  []apientities.Message
	History   []any
	APIKey    string
	Workspace string
	Extra     map[string]any
}

type TokenizationService struct{}

func (TokenizationService) Call(ctx context.Context, p *TokenizationCallParams) (*apientities.DashScopeAPIResponse, error) {
	if p.Prompt == nil && len(p.Messages) == 0 {
		return nil, common.NewInputRequired("prompt or messages is required!")
	}
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	input := map[string]any{}
	if p.History != nil {
		common.Log.Warn("%s", common.DEPRECATED_MESSAGE)
		input[common.HISTORY] = p.History
	}
	if p.Prompt != nil {
		input[common.PROMPT] = p.Prompt
	}
	if len(p.Messages) > 0 {
		input[common.MESSAGES] = p.Messages
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "tokenizer", Task: "tokenizer", Function: "tokenizer",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "tokenizers",
	})
}

var Tokenization = TokenizationService{}

// Tokenizer local tokenizer fallback (character-based if vocab not loaded).
type Tokenizer struct {
	Name string
}

func (t *Tokenizer) Encode(text string) []int {
	ids := make([]int, 0, utf8.RuneCountInString(text))
	for _, r := range text {
		ids = append(ids, int(r))
	}
	return ids
}

func (t *Tokenizer) Decode(ids []int) string {
	rs := make([]rune, len(ids))
	for i, id := range ids {
		rs[i] = rune(id)
	}
	return string(rs)
}

func GetTokenizer(name string) *Tokenizer {
	if name == "" {
		name = TokenizerQwenTurbo
	}
	return &Tokenizer{Name: name}
}

func ListTokenizers() []string {
	return []string{
		TokenizerQwenTurbo, TokenizerQwenPlus, TokenizerQwen7bChat, TokenizerQwen14bChat,
		TokenizerLlama27bChatV2, TokenizerLlama213bChatV2, TokenizerTextEmbeddingV2, TokenizerQwen72bChat,
	}
}

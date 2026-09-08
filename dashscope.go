// Copyright (c) Alibaba, Inc. and its affiliates.

// Package dashscope provides a Go SDK for Alibaba Cloud Model Studio (Bailian) APIs,
// covering text generation, multi-modal understanding, embeddings, reranking,
// image/video generation, speech synthesis & recognition, and more.
package dashscope

import (
	"github.com/liuxiaobopro/dashscope-go/aigc"
	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/app"
	"github.com/liuxiaobopro/dashscope-go/assistants"
	httptts "github.com/liuxiaobopro/dashscope-go/audio/http_tts"
	"github.com/liuxiaobopro/dashscope-go/audio/asr"
	"github.com/liuxiaobopro/dashscope-go/audio/tts"
	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/embeddings"
	"github.com/liuxiaobopro/dashscope-go/finetune"
	"github.com/liuxiaobopro/dashscope-go/nlp"
	"github.com/liuxiaobopro/dashscope-go/rerank"
	"github.com/liuxiaobopro/dashscope-go/threads"
	"github.com/liuxiaobopro/dashscope-go/tokenizers"
)

// SetAPIKey set the API key via code.
func SetAPIKey(apiKey string) { common.APIKey = apiKey }

// SetAPIKeyFilePath specify the API key file path via code.
func SetAPIKeyFilePath(p string) { common.APIKeyFilePath = p }

// SaveAPIKey save the API key to a file.
// If apiKeyFilePath is empty, save to default location "~/.dashscope/api_key".
func SaveAPIKey(apiKey, apiKeyFilePath string) error {
	return common.SaveAPIKey(apiKey, apiKeyFilePath)
}

func APIKey() string             { return common.APIKey }
func APIKeyFilePath() string     { return common.APIKeyFilePath }
func BaseHTTPAPIURL() string     { return common.BaseHTTPAPIURL }
func BaseWebsocketAPIURL() string { return common.BaseWebsocketAPIURL }
func BaseCompatibleAPIURL() string { return common.BaseCompatibleAPIURL }

func SetBaseHTTPAPIURL(u string)      { common.BaseHTTPAPIURL = u }
func SetBaseWebsocketAPIURL(u string) { common.BaseWebsocketAPIURL = u }
func SetBaseCompatibleAPIURL(u string) { common.BaseCompatibleAPIURL = u }

func CloseSharedSyncSession() { apientities.CloseSharedSyncSession() }

func Ptr[T any](v T) *T { return common.Ptr(v) }

type Message = apientities.Message
type GenerationResponse = apientities.GenerationResponse
type DashScopeAPIResponse = apientities.DashScopeAPIResponse

var (
	Generation              = aigc.Generation
	Conversation            = aigc.Conversation
	ImageSynthesis          = aigc.ImageSynthesis
	MultiModalConversation  = aigc.MultiModalConversation
	VideoSynthesis          = aigc.VideoSynthesis
	CodeGeneration          = aigc.CodeGeneration
	ChatCompletion          = aigc.ChatCompletion
	ImageGeneration         = aigc.ImageGeneration
	Application             = app.Application
	Transcription           = asr.Transcription
	SpeechSynthesizer       = tts.SpeechSynthesizer
	HttpSpeechSynthesizer   = httptts.HttpSpeechSynthesizer
	TextEmbedding           = embeddings.TextEmbedding
	MultiModalEmbedding     = embeddings.MultiModalEmbedding
	BatchTextEmbedding      = embeddings.BatchTextEmbedding
	TextReRank              = rerank.TextReRank
	Understanding           = nlp.Understanding
	Tokenization            = tokenizers.Tokenization
	FineTunes               = finetune.FineTunes
	Deployments             = finetune.Deployments
	Assistants              = assistants.Assistants
	Threads                 = threads.Threads
	Messages                = threads.Messages
	Runs                    = threads.Runs
	Steps                   = threads.Steps
)

func GetTokenizer(name string) *tokenizers.Tokenizer { return tokenizers.GetTokenizer(name) }
func ListTokenizers() []string                         { return tokenizers.ListTokenizers() }

type (
	Assistant          = assistants.Assistant
	AssistantList      = assistants.AssistantList
	AssistantFile      = assistants.AssistantFile
	DeleteResponse     = assistants.DeleteResponse
	Thread             = threads.Thread
	ThreadMessage      = threads.ThreadMessage
	ThreadMessageList  = threads.ThreadMessageList
	MessageFile        = threads.MessageFile
	Run                = threads.Run
	RunList            = threads.RunList
	RunStep            = threads.RunStep
	RunStepList        = threads.RunStepList
	History            = aigc.History
	HistoryItem        = aigc.HistoryItem
	Tokenizer          = tokenizers.Tokenizer
	MultiModalEmbeddingItemText  = embeddings.MultiModalEmbeddingItemBase
	MultiModalEmbeddingItemImage = embeddings.MultiModalEmbeddingItemBase
	MultiModalEmbeddingItemAudio = embeddings.MultiModalEmbeddingItemBase
	BatchTextEmbeddingResponse   = embeddings.BatchTextEmbeddingResponse
)

// Copyright (c) Alibaba, Inc. and its affiliates.

package asr

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

// TranscriptionCallParams speech transcription (async).
type TranscriptionCallParams struct {
	Model          string
	FileURLs       []string
	APIKey         string
	Workspace      string
	WaitTimeout    int
	Extra          map[string]any
}

type TranscriptionService struct{}

func (TranscriptionService) Call(ctx context.Context, p *TranscriptionCallParams) (*apientities.TranscriptionResponse, error) {
	input := map[string]any{"file_urls": p.FileURLs}
	rsp, err := (client.BaseAsyncApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "audio", Task: "asr", Function: "transcription",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "audio",
	}, p.WaitTimeout)
	if err != nil {
		return nil, err
	}
	return apientities.TranscriptionResponseFromAPI(rsp), nil
}

func (TranscriptionService) AsyncCall(ctx context.Context, p *TranscriptionCallParams) (*apientities.TranscriptionResponse, error) {
	input := map[string]any{"file_urls": p.FileURLs}
	rsp, err := (client.BaseAsyncApi{}).AsyncCall(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "audio", Task: "asr", Function: "transcription",
		APIKey: p.APIKey, Workspace: p.Workspace, Extra: p.Extra, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.TranscriptionResponseFromAPI(rsp), nil
}

func (TranscriptionService) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int) (*apientities.TranscriptionResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Wait(ctx, task, apiKey, workspace, waitTimeout, "")
	if err != nil {
		return nil, err
	}
	return apientities.TranscriptionResponseFromAPI(rsp), nil
}

func (TranscriptionService) Fetch(ctx context.Context, task any, apiKey, workspace string) (*apientities.TranscriptionResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Fetch(ctx, task, apiKey, workspace, "")
	if err != nil {
		return nil, err
	}
	return apientities.TranscriptionResponseFromAPI(rsp), nil
}

var Transcription = TranscriptionService{}

// RecognitionCallParams realtime recognition via websocket.
type RecognitionCallParams struct {
	Model         string
	Format        string
	SampleRate    int
	APIKey        string
	Workspace     string
	IsBinaryInput bool
	Extra         map[string]any
}

type RecognitionService struct{}

func (RecognitionService) Call(ctx context.Context, p *RecognitionCallParams) (*apientities.RecognitionResponse, error) {
	input := map[string]any{}
	params := map[string]any{"format": p.Format, "sample_rate": p.SampleRate}
	for k, v := range p.Extra {
		params[k] = v
	}
	rsp, err := (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "audio", Task: "asr", Function: "recognition",
		APIKey: p.APIKey, Workspace: p.Workspace, APIProtocol: common.ApiProtocolWebsocket,
		IsBinaryInput: p.IsBinaryInput, Parameters: params, SDKModule: "audio",
	})
	if err != nil {
		return nil, err
	}
	return apientities.RecognitionResponseFromAPI(rsp), nil
}

var Recognition = RecognitionService{}

type TranslationRecognizerService struct{}

func (TranslationRecognizerService) Call(ctx context.Context, p *RecognitionCallParams) (*apientities.RecognitionResponse, error) {
	p2 := *p
	return (RecognitionService{}).Call(ctx, &p2)
}

var TranslationRecognizer = TranslationRecognizerService{}

type VocabularyService struct{}

func (VocabularyService) Create(ctx context.Context, model string, targetModel string, prefix string, vocabulary []any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: model, Input: map[string]any{"target_model": targetModel, "prefix": prefix, "vocabulary": vocabulary},
		TaskGroup: "audio", Task: "asr", Function: "vocabulary", APIKey: apiKey, Workspace: workspace, SDKModule: "audio",
	})
}

var Vocabulary = VocabularyService{}

type AsrPhraseManagerService struct{}

func (AsrPhraseManagerService) Create(ctx context.Context, model string, phrases []any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: model, Input: map[string]any{"phrases": phrases},
		TaskGroup: "audio", Task: "asr", Function: "asr-phrase-manager", APIKey: apiKey, Workspace: workspace, SDKModule: "audio",
	})
}

var AsrPhraseManager = AsrPhraseManagerService{}

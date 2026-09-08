// Copyright (c) Alibaba, Inc. and its affiliates.

package app

import (
	"context"
	"iter"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
	"github.com/liuxiaobopro/dashscope-go/utils"
)

const (
	DocReferenceSimple  = "simple"
	DocReferenceIndexed = "indexed"
)

// ApplicationCallParams Application completion call.
type ApplicationCallParams struct {
	AppID            string
	Prompt           string
	History          []any
	Workspace        string
	APIKey           string
	Messages         []apientities.Message
	Stream           bool
	Temperature      *float64
	TopP             *float64
	TopK             *int
	Seed             *int
	SessionID        *string
	BizParams        map[string]any
	HasThoughts      *bool
	DocTagCodes      []string
	DocReferenceType *string
	MemoryID         *string
	ImageList        []string
	FileList         []string
	RagOptions       map[string]any
	IncrementalOutput *bool
	Extra            map[string]any
}

func (p *ApplicationCallParams) toInputParams() (map[string]any, map[string]any, error) {
	if p.AppID == "" {
		return nil, nil, common.NewInputRequired("App id is required!")
	}
	if p.Prompt == "" && len(p.Messages) == 0 {
		return nil, nil, common.NewInputRequired("prompt or messages is required!")
	}
	input := map[string]any{}
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
	params := map[string]any{}
	if p.Temperature != nil {
		params["temperature"] = *p.Temperature
	}
	if p.TopP != nil {
		params["top_p"] = *p.TopP
	}
	if p.TopK != nil {
		params["top_k"] = *p.TopK
	}
	if p.Seed != nil {
		params["seed"] = *p.Seed
	}
	if p.SessionID != nil {
		params["session_id"] = *p.SessionID
	}
	if p.BizParams != nil {
		params["biz_params"] = p.BizParams
	}
	if p.HasThoughts != nil {
		params["has_thoughts"] = *p.HasThoughts
	}
	if p.DocTagCodes != nil {
		params["doc_tag_codes"] = p.DocTagCodes
	}
	if p.DocReferenceType != nil {
		params["doc_reference_type"] = *p.DocReferenceType
	}
	if p.MemoryID != nil {
		params["memory_id"] = *p.MemoryID
	}
	if p.ImageList != nil {
		params["image_list"] = p.ImageList
	}
	if p.FileList != nil {
		params["file_list"] = p.FileList
	}
	if p.RagOptions != nil {
		params["rag_options"] = p.RagOptions
	}
	if p.IncrementalOutput != nil {
		params["incremental_output"] = *p.IncrementalOutput
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	return input, params, nil
}

// ApplicationService API for app completion calls.
type ApplicationService struct{}

func (ApplicationService) Call(ctx context.Context, p *ApplicationCallParams) (*apientities.DashScopeAPIResponse, error) {
	input, params, err := p.toInputParams()
	if err != nil {
		return nil, err
	}
	apiKey := p.APIKey
	if apiKey == "" {
		apiKey, err = common.GetDefaultAPIKey()
		if err != nil {
			return nil, err
		}
	}
	return (client.BaseApi{}).Call(ctx, &client.CallParams{
		Model: p.AppID, Input: input, TaskGroup: "apps", Task: p.AppID, Function: "completion",
		APIKey: apiKey, Workspace: p.Workspace, Stream: p.Stream, Parameters: params, SDKModule: "app",
	})
}

func (ApplicationService) CallStream(ctx context.Context, p *ApplicationCallParams) iter.Seq2[*apientities.DashScopeAPIResponse, error] {
	p.Stream = true
	return func(yield func(*apientities.DashScopeAPIResponse, error) bool) {
		input, params, err := p.toInputParams()
		if err != nil {
			yield(nil, err)
			return
		}
		ch, err := (client.BaseApi{}).CallStream(ctx, &client.CallParams{
			Model: p.AppID, Input: input, TaskGroup: "apps", Task: p.AppID, Function: "completion",
			APIKey: p.APIKey, Workspace: p.Workspace, Stream: true, Parameters: params, SDKModule: "app",
		})
		if err != nil {
			yield(nil, err)
			return
		}
		accumulated := map[any]any{}
		for rsp := range ch {
			parsed := apientities.GenerationResponseFromAPI(rsp)
			ok, _ := utils.MergeSingleResponse(parsed, accumulated, 1)
			if ok && !yield(&parsed.DashScopeAPIResponse, nil) {
				return
			}
		}
	}
}

var Application = ApplicationService{}

type ApplicationResponse = apientities.DashScopeAPIResponse

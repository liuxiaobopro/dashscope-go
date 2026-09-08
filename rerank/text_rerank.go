// Copyright (c) Alibaba, Inc. and its affiliates.

package rerank

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const (
	GteRerank         = "gte-rerank"
	GteRerankV2       = "gte-rerank-v2"
	Qwen3Rerank       = "qwen3-rerank"
	Qwen3VLRerank     = "qwen3-vl-rerank"
)

// TextReRankCallParams Calling rerank service.
type TextReRankCallParams struct {
	// Model The model to use.
	Model string
	// Query The query string.
	Query string
	// Documents The documents to rank.
	Documents []string
	// ReturnDocuments enable return origin documents, system default is false.
	ReturnDocuments *bool
	// TopN how many documents to return, default return all the documents.
	TopN *int
	APIKey string
	// Instruct Custom task instruction to guide ranking strategy. English recommended.
	Instruct *string
	Extra    map[string]any
	Workspace string
}

func prepareRerank(p *TextReRankCallParams) (*client.CallParams, error) {
	if p.Query == "" || len(p.Documents) == 0 {
		return nil, common.NewInputRequired("query and documents are required!")
	}
	if p.Model == "" {
		return nil, common.NewModelRequired("Model is required!")
	}
	input := map[string]any{"query": p.Query, "documents": p.Documents}
	parameters := map[string]any{}
	if p.ReturnDocuments != nil {
		parameters["return_documents"] = *p.ReturnDocuments
	}
	if p.TopN != nil {
		parameters["top_n"] = *p.TopN
	}
	if p.Instruct != nil {
		parameters["instruct"] = *p.Instruct
	}
	for k, v := range p.Extra {
		parameters[k] = v
	}
	return &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "rerank", Task: "text-rerank", Function: "text-rerank",
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: parameters, SDKModule: "rerank",
	}, nil
}

// TextReRankService API for rerank models.
type TextReRankService struct {
	Task string
}

func (TextReRankService) Call(ctx context.Context, p *TextReRankCallParams) (*apientities.ReRankResponse, error) {
	cp, err := prepareRerank(p)
	if err != nil {
		return nil, err
	}
	rsp, err := (client.BaseApi{}).Call(ctx, cp)
	if err != nil {
		return nil, err
	}
	return apientities.ReRankResponseFromAPI(rsp), nil
}

// TextReRank API for rerank models.
var TextReRank = TextReRankService{Task: "text-rerank"}

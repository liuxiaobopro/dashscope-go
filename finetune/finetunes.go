// Copyright (c) Alibaba, Inc. and its affiliates.

package finetune

import (
	"context"
	"strconv"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
)

// FineTunesService fine-tune jobs.
type FineTunesService struct{}

type FineTuneCallParams struct {
	Model             string
	TrainingFileIDs   any
	ValidationFileIDs any
	Mode              string
	HyperParameters   map[string]any
	APIKey            string
	Workspace         string
	Extra             map[string]any
}

func (FineTunesService) Call(ctx context.Context, p *FineTuneCallParams) (*apientities.DashScopeAPIResponse, error) {
	payload := map[string]any{"model": p.Model, "training_file_ids": p.TrainingFileIDs}
	if p.ValidationFileIDs != nil {
		payload["validation_file_ids"] = p.ValidationFileIDs
	}
	if p.Mode != "" {
		payload["mode"] = p.Mode
	}
	if p.HyperParameters != nil {
		payload["hyper_parameters"] = p.HyperParameters
	}
	for k, v := range p.Extra {
		payload[k] = v
	}
	return client.RESTPost(ctx, client.GetURL("", "fine-tunes", ""), payload, p.APIKey, p.Workspace, "finetune", nil)
}

func (FineTunesService) Get(ctx context.Context, jobID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "fine-tunes", jobID), nil, apiKey, workspace, "finetune", nil)
}

func (FineTunesService) List(ctx context.Context, pageNo, pageSize int, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "fine-tunes", ""), map[string]string{
		"page_no": strconv.Itoa(pageNo), "page_size": strconv.Itoa(pageSize),
	}, apiKey, workspace, "finetune", nil)
}

func (FineTunesService) Cancel(ctx context.Context, jobID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "fine-tunes", jobID+"/cancel"), nil, apiKey, workspace, "finetune", nil)
}

func (FineTunesService) Delete(ctx context.Context, jobID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "fine-tunes", jobID), apiKey, workspace, "finetune")
}

var FineTunes = FineTunesService{}

type DeploymentsService struct{}

func (DeploymentsService) Call(ctx context.Context, model, suffix, apiKey, workspace string, extra map[string]any) (*apientities.DashScopeAPIResponse, error) {
	payload := map[string]any{"model_name": model}
	if suffix != "" {
		payload["suffix"] = suffix
	}
	for k, v := range extra {
		payload[k] = v
	}
	return client.RESTPost(ctx, client.GetURL("", "deployments", ""), payload, apiKey, workspace, "finetune", nil)
}

func (DeploymentsService) Get(ctx context.Context, deployedModel, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "deployments", deployedModel), nil, apiKey, workspace, "finetune", nil)
}

func (DeploymentsService) List(ctx context.Context, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "deployments", ""), params, apiKey, workspace, "finetune", nil)
}

func (DeploymentsService) Delete(ctx context.Context, deployedModel, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "deployments", deployedModel), apiKey, workspace, "finetune")
}

var Deployments = DeploymentsService{}

// Copyright (c) Alibaba, Inc. and its affiliates.

package assistants

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
)

type Assistant struct {
	ID           string `json:"id"`
	Object       string `json:"object"`
	CreatedAt    int64  `json:"created_at"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Model        string `json:"model"`
	Instructions string `json:"instructions"`
	Tools        []any  `json:"tools"`
	FileIDs      []string `json:"file_ids"`
	Metadata     map[string]any `json:"metadata"`
}

type AssistantList struct {
	Object  string      `json:"object"`
	Data    []Assistant `json:"data"`
	FirstID string      `json:"first_id"`
	LastID  string      `json:"last_id"`
	HasMore bool        `json:"has_more"`
}

type AssistantFile struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	CreatedAt   int64  `json:"created_at"`
	AssistantID string `json:"assistant_id"`
}

type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type AssistantsService struct{}

func (AssistantsService) Create(ctx context.Context, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "assistants", ""), payload, apiKey, workspace, "assistants", nil)
}

func (AssistantsService) Get(ctx context.Context, assistantID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "assistants", assistantID), nil, apiKey, workspace, "assistants", nil)
}

func (AssistantsService) List(ctx context.Context, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "assistants", ""), params, apiKey, workspace, "assistants", nil)
}

func (AssistantsService) Update(ctx context.Context, assistantID string, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "assistants", assistantID), payload, apiKey, workspace, "assistants", nil)
}

func (AssistantsService) Delete(ctx context.Context, assistantID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "assistants", assistantID), apiKey, workspace, "assistants")
}

var Assistants = AssistantsService{}

type AssistantFilesService struct{}

func (AssistantFilesService) Create(ctx context.Context, assistantID, fileID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "assistants", assistantID+"/files"), map[string]any{"file_id": fileID}, apiKey, workspace, "assistants", nil)
}

func (AssistantFilesService) List(ctx context.Context, assistantID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "assistants", assistantID+"/files"), nil, apiKey, workspace, "assistants", nil)
}

func (AssistantFilesService) Get(ctx context.Context, assistantID, fileID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "assistants", assistantID+"/files/"+fileID), nil, apiKey, workspace, "assistants", nil)
}

func (AssistantFilesService) Delete(ctx context.Context, assistantID, fileID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "assistants", assistantID+"/files/"+fileID), apiKey, workspace, "assistants")
}

var AssistantFiles = AssistantFilesService{}

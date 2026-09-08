// Copyright (c) Alibaba, Inc. and its affiliates.

package threads

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
)

type Thread struct {
	ID        string         `json:"id"`
	Object    string         `json:"object"`
	CreatedAt int64          `json:"created_at"`
	Metadata  map[string]any `json:"metadata"`
}

type ThreadMessage struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	CreatedAt   int64  `json:"created_at"`
	ThreadID    string `json:"thread_id"`
	Role        string `json:"role"`
	Content     any    `json:"content"`
	AssistantID string `json:"assistant_id"`
	RunID       string `json:"run_id"`
	FileIDs     []string `json:"file_ids"`
}

type ThreadMessageList struct {
	Object  string          `json:"object"`
	Data    []ThreadMessage `json:"data"`
	HasMore bool            `json:"has_more"`
}

type MessageFile struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	MessageID string `json:"message_id"`
}

type Run struct {
	ID           string `json:"id"`
	Object       string `json:"object"`
	CreatedAt    int64  `json:"created_at"`
	ThreadID     string `json:"thread_id"`
	AssistantID  string `json:"assistant_id"`
	Status       string `json:"status"`
	RequiredAction any  `json:"required_action"`
}

type RunList struct {
	Object  string `json:"object"`
	Data    []Run  `json:"data"`
	HasMore bool   `json:"has_more"`
}

type RunStep struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

type RunStepList struct {
	Object  string    `json:"object"`
	Data    []RunStep `json:"data"`
	HasMore bool      `json:"has_more"`
}

type ThreadsService struct{}

func (ThreadsService) Create(ctx context.Context, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", ""), payload, apiKey, workspace, "assistants", nil)
}
func (ThreadsService) Get(ctx context.Context, threadID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID), nil, apiKey, workspace, "assistants", nil)
}
func (ThreadsService) Update(ctx context.Context, threadID string, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", threadID), payload, apiKey, workspace, "assistants", nil)
}
func (ThreadsService) Delete(ctx context.Context, threadID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "threads", threadID), apiKey, workspace, "assistants")
}

var Threads = ThreadsService{}

type MessagesService struct{}

func (MessagesService) Create(ctx context.Context, threadID string, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", threadID+"/messages"), payload, apiKey, workspace, "assistants", nil)
}
func (MessagesService) List(ctx context.Context, threadID, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/messages"), params, apiKey, workspace, "assistants", nil)
}
func (MessagesService) Get(ctx context.Context, threadID, messageID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/messages/"+messageID), nil, apiKey, workspace, "assistants", nil)
}

var Messages = MessagesService{}

type MessageFilesService struct{}

func (MessageFilesService) List(ctx context.Context, threadID, messageID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/messages/"+messageID+"/files"), nil, apiKey, workspace, "assistants", nil)
}

var MessageFiles = MessageFilesService{}

type RunsService struct{}

func (RunsService) Create(ctx context.Context, threadID string, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", threadID+"/runs"), payload, apiKey, workspace, "assistants", nil)
}
func (RunsService) List(ctx context.Context, threadID, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/runs"), params, apiKey, workspace, "assistants", nil)
}
func (RunsService) Get(ctx context.Context, threadID, runID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/runs/"+runID), nil, apiKey, workspace, "assistants", nil)
}
func (RunsService) Cancel(ctx context.Context, threadID, runID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", threadID+"/runs/"+runID+"/cancel"), nil, apiKey, workspace, "assistants", nil)
}
func (RunsService) SubmitToolOutputs(ctx context.Context, threadID, runID string, payload map[string]any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, client.GetURL("", "threads", threadID+"/runs/"+runID+"/submit_tool_outputs"), payload, apiKey, workspace, "assistants", nil)
}

var Runs = RunsService{}

type StepsService struct{}

func (StepsService) List(ctx context.Context, threadID, runID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/runs/"+runID+"/steps"), nil, apiKey, workspace, "assistants", nil)
}
func (StepsService) Get(ctx context.Context, threadID, runID, stepID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "threads", threadID+"/runs/"+runID+"/steps/"+stepID), nil, apiKey, workspace, "assistants", nil)
}

var Steps = StepsService{}

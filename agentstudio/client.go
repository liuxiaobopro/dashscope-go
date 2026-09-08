// Copyright (c) Alibaba, Inc. and its affiliates.

package agentstudio

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

func base() string {
	return common.JoinURL(common.BaseHTTPAPIURL, "agentstudio")
}

type Client struct {
	APIKey    string
	Workspace string
}

func NewClient(apiKey, workspace string) *Client {
	return &Client{APIKey: apiKey, Workspace: workspace}
}

func (c *Client) AgentsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "agents"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) AgentsGet(ctx context.Context, id string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "agents", id), nil, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) AgentsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "agents"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) AgentsDelete(ctx context.Context, id string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, common.JoinURL(base(), "agents", id), c.APIKey, c.Workspace, "agentstudio")
}

func (c *Client) SessionsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "sessions"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) SessionsGet(ctx context.Context, id string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "sessions", id), nil, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) SessionsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "sessions"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

func (c *Client) SkillsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "skills"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) SkillsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "skills"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

func (c *Client) FilesList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "files"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

func (c *Client) DeploymentsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "deployments"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) DeploymentsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "deployments"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

func (c *Client) VaultsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "vaults"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) VaultsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "vaults"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

func (c *Client) WebhookEndpointsCreate(ctx context.Context, payload map[string]any) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTPost(ctx, common.JoinURL(base(), "webhook-endpoints"), payload, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) EnvironmentsList(ctx context.Context, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "environments"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}
func (c *Client) SessionEventsList(ctx context.Context, sessionID string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(base(), "sessions", sessionID+"/events"), params, c.APIKey, c.Workspace, "agentstudio", nil)
}

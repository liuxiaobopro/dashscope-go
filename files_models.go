// Copyright (c) Alibaba, Inc. and its affiliates.

package dashscope

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

// Files file API.
type FilesService struct{}

// Upload Upload file for model fine-tune or other tasks.
func (FilesService) Upload(ctx context.Context, filePath, purpose, description, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	if purpose == "" {
		purpose = common.FilePurposeFineTune
	}
	if purpose == common.FilePurposeFineTune {
		if !common.IsValidateFineTuneFile(filePath) {
			return nil, common.NewInvalidFileFormat("The file " + filePath + " is not in valid jsonl format")
		}
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}
	_ = w.WriteField("purpose", purpose)
	if description != "" {
		_ = w.WriteField("description", description)
	}
	_ = w.Close()
	u := client.GetURL("", "files", "")
	h, err := common.DefaultHeaders(apiKey, "files")
	if err != nil {
		return nil, err
	}
	for k, v := range common.WorkspaceHeader(workspace) {
		h[k] = v
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return nil, err
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return &apientities.DashScopeAPIResponse{StatusCode: resp.StatusCode, Message: string(body), Output: map[string]any{"raw": string(body)}}, nil
}

func (FilesService) List(ctx context.Context, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "files", ""), params, apiKey, workspace, "files", nil)
}

func (FilesService) Get(ctx context.Context, fileID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, client.GetURL("", "files", fileID), nil, apiKey, workspace, "files", nil)
}

func (FilesService) Delete(ctx context.Context, fileID, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTDelete(ctx, client.GetURL("", "files", fileID), apiKey, workspace, "files")
}

var Files = FilesService{}

// Models model API.
type ModelsService struct{}

func (ModelsService) Get(ctx context.Context, name, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	u := common.JoinURL(common.BaseHTTPAPIURL, "models")
	rsp, err := client.RESTGet(ctx, u, map[string]string{"model": name, "page_no": "1", "page_size": "1"}, apiKey, workspace, "", nil)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != http.StatusOK {
		return rsp, nil
	}
	out := rsp.OutputMap()
	if out == nil {
		rsp.StatusCode = 404
		return rsp, nil
	}
	models, _ := out["models"].([]any)
	if len(models) == 0 {
		rsp.StatusCode = 404
	}
	return rsp, nil
}

func (ModelsService) List(ctx context.Context, apiKey, workspace string, params map[string]string) (*apientities.DashScopeAPIResponse, error) {
	return client.RESTGet(ctx, common.JoinURL(common.BaseHTTPAPIURL, "models"), params, apiKey, workspace, "", nil)
}

var Models = ModelsService{}

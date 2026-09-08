// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
)

var sdkSessionID = newSessionID()
var sdkClient = "go-sdk"

func newSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GetTaskGroupAndTask get task_group and task name based on api file package.
func GetTaskGroupAndTask(pkg, fileBase string) (taskGroup, task string) {
	task = strings.ReplaceAll(fileBase, "_", "-")
	parts := strings.Split(pkg, "/")
	if len(parts) > 0 {
		taskGroup = parts[len(parts)-1]
	}
	return taskGroup, task
}

// IsURL check the input url is valid url.
func IsURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	switch u.Scheme {
	case "http", "https", "oss":
		return true
	default:
		return false
	}
}

// IsPath check the input path is valid local path.
func IsPath(path string) bool {
	u, err := url.Parse(path)
	if err != nil {
		return false
	}
	if u.Scheme == "file" || u.Scheme == "" {
		p := u.Path
		if p == "" {
			p = path
		}
		_, err := os.Stat(p)
		return err == nil
	}
	return false
}

// GetUserAgent user-agent string.
func GetUserAgent() string {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	return fmt.Sprintf("dashscope/%s; go/%s; platform/%s; processor/%s",
		Version, runtime.Version(), platform, runtime.GOARCH)
}

var sdkModulePackageOverrides = map[string]string{
	"threads": "assistants",
}

// GetAPIModule derive the sdk module tag from package path.
func GetAPIModule(pkg string) string {
	const prefix = "github.com/liuxiaobopro/dashscope-go/"
	rest := pkg
	if i := strings.Index(pkg, prefix); i >= 0 {
		rest = pkg[i+len(prefix):]
	}
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	pkgName := parts[0]
	if v, ok := sdkModulePackageOverrides[pkgName]; ok {
		return v
	}
	return pkgName
}

// SetSDKClient override the client identifier in the x-dashscope-sdk-client header.
func SetSDKClient(client string) {
	if client != "" {
		sdkClient = client
	}
}

// GetSDKHeaders structured SDK identification headers for backend statistics.
// Set DASHSCOPE_DISABLE_SDK_HEADERS=1 to omit them.
func GetSDKHeaders(module string) map[string]string {
	if os.Getenv("DASHSCOPE_DISABLE_SDK_HEADERS") != "" {
		return map[string]string{}
	}
	parts := []string{sdkClient, Version}
	if module != "" {
		parts = append(parts, module)
	}
	return map[string]string{
		"x-dashscope-sdk-client":     strings.Join(parts, "/"),
		"x-dashscope-sdk-session-id": sdkSessionID,
	}
}

// DefaultHeaders default request headers.
func DefaultHeaders(apiKey, module string) (map[string]string, error) {
	ua := GetUserAgent()
	headers := map[string]string{"user-agent": ua}
	for k, v := range GetSDKHeaders(module) {
		headers[k] = v
	}
	if apiKey == "" {
		var err error
		apiKey, err = GetDefaultAPIKey()
		if err != nil {
			return nil, err
		}
	}
	headers["Authorization"] = "Bearer " + apiKey
	headers["Accept"] = "application/json; charset=utf-8"
	return headers, nil
}

// JoinURL join base url with path segments.
func JoinURL(baseURL string, args ...string) string {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	u := baseURL
	for _, arg := range args {
		if arg != "" {
			u += arg + "/"
		}
	}
	return strings.TrimRight(u, "/")
}

// WorkspaceHeader workspace header.
func WorkspaceHeader(workspace string) map[string]string {
	if workspace != "" {
		return map[string]string{"X-DashScope-WorkSpace": workspace}
	}
	return map[string]string{}
}

// MergeHeaders merge header maps, later overrides earlier.
func MergeHeaders(ms ...map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// Ptr helper.
func Ptr[T any](v T) *T { return &v }

// IsValidateFineTuneFile check jsonl file.
func IsValidateFineTuneFile(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if !json.Valid([]byte(line)) {
			return false
		}
	}
	return sc.Err() == nil
}

// HeaderMap convert http.Header to map.
func HeaderMap(h map[string][]string) map[string]string {
	m := map[string]string{}
	for k, vs := range h {
		if len(vs) > 0 {
			m[k] = vs[0]
		}
	}
	return m
}

// SSEEvent server-sent event.
type SSEEvent struct {
	ID        string
	EventType string
	Data      string
}

// StreamItem SSE stream item.
type StreamItem struct {
	IsError    bool
	StatusCode int
	Event      SSEEvent
	Err        error
}

// HandleStream parse SSE stream.
func HandleStream(r io.Reader) <-chan StreamItem {
	ch := make(chan StreamItem)
	go func() {
		defer close(ch)
		isError := false
		statusCode := http.StatusBadRequest
		event := SSEEvent{}
		eventType := ""
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), "\r")
			if line == "" {
				continue
			}
			switch {
			case strings.HasPrefix(line, "id:"):
				event.ID = strings.TrimSpace(line[len("id:"):])
			case strings.HasPrefix(line, "event:"):
				eventType = strings.TrimSpace(line[len("event:"):])
				event.EventType = eventType
				if eventType == "error" {
					isError = true
				}
			case strings.HasPrefix(line, "status:"):
				_, _ = fmt.Sscanf(strings.TrimSpace(line[len("status:"):]), "%d", &statusCode)
			case strings.HasPrefix(line, "data:"):
				event.Data = strings.TrimSpace(line[len("data:"):])
				if eventType == "done" {
					continue
				}
				ch <- StreamItem{IsError: isError, StatusCode: statusCode, Event: event}
				if isError {
					return
				}
			default:
				continue
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- StreamItem{Err: err}
		}
	}()
	return ch
}

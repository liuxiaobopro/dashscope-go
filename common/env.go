// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import (
	"os"
	"path/filepath"
)

var (
	// APIRegion API region, default cn-beijing.
	APIRegion = getenv(DASHSCOPE_API_REGION_ENV, "cn-beijing")
	// APIVersion API version, default v1.
	APIVersion = getenv(DASHSCOPE_API_VERSION_ENV, "v1")
	// APIKey API key from env DASHSCOPE_API_KEY.
	APIKey = os.Getenv(DASHSCOPE_API_KEY_ENV)
	// APIKeyFilePath API key file path from env DASHSCOPE_API_KEY_FILE_PATH.
	APIKeyFilePath = os.Getenv(DASHSCOPE_API_KEY_FILE_PATH_ENV)
	// BaseHTTPAPIURL HTTP API base url, ensure end /.
	BaseHTTPAPIURL = getenv("DASHSCOPE_HTTP_BASE_URL", "https://dashscope.aliyuncs.com/api/"+APIVersion)
	// BaseWebsocketAPIURL websocket API base url.
	BaseWebsocketAPIURL = getenv("DASHSCOPE_WEBSOCKET_BASE_URL", "wss://dashscope.aliyuncs.com/api-ws/"+APIVersion+"/inference")
	// BaseCompatibleAPIURL OpenAI-compatible API base url.
	BaseCompatibleAPIURL = getenv("DASHSCOPE_COMPATIBLE_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/"+APIVersion)
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// DefaultCachePath default cache directory ~/.dashscope
func DefaultCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".dashscope"
	}
	return filepath.Join(home, ".dashscope")
}

// DefaultAPIKeyFilePath default api key file ~/.dashscope/api_key
func DefaultAPIKeyFilePath() string {
	return filepath.Join(DefaultCachePath(), "api_key")
}

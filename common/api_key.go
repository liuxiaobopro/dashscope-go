// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import (
	"os"
	"path/filepath"
	"strings"
)

// GetDefaultAPIKey returns the api key from code, env, or file.
func GetDefaultAPIKey() (string, error) {
	if APIKey != "" {
		return APIKey, nil
	}
	if APIKeyFilePath != "" {
		b, err := os.ReadFile(APIKeyFilePath)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	p := DefaultAPIKeyFilePath()
	if _, err := os.Stat(p); err == nil {
		b, err := os.ReadFile(p)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	return "", NewAuthenticationError(
		"No api key provided. You can set by dashscope.SetAPIKey(your_api_key) in code, " +
			"or you can set it via environment variable DASHSCOPE_API_KEY= your_api_key. " +
			"You can store your api key to a file, and use dashscope.SetAPIKeyFilePath(api_key_file_path) in code, " +
			"or you can set api key file path via environment variable DASHSCOPE_API_KEY_FILE_PATH, " +
			"You can call SaveAPIKey to api_key_file_path or default path(~/.dashscope/api_key).",
	)
}

// SaveAPIKey saves the api key to a file.
// If apiKeyFilePath is empty, save to default location "~/.dashscope/api_key".
func SaveAPIKey(apiKey string, apiKeyFilePath string) error {
	if apiKeyFilePath == "" {
		if err := os.MkdirAll(DefaultCachePath(), 0o755); err != nil {
			return err
		}
		return os.WriteFile(DefaultAPIKeyFilePath(), []byte(apiKey), 0o600)
	}
	if err := os.MkdirAll(filepath.Dir(apiKeyFilePath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(apiKeyFilePath, []byte(apiKey), 0o600)
}

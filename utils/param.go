// Copyright (c) Alibaba, Inc. and its affiliates.

package utils

import "strings"

// ParamUtil parameter helpers.
type ParamUtil struct{}

// ShouldModifyIncrementalOutput determine if increment_output parameter needs to be modified based on model name.
// Returns false if model contains 'tts', 'omni', or 'qwen-deep-research', true otherwise.
func ShouldModifyIncrementalOutput(modelName string) bool {
	if modelName == "" {
		return true
	}
	lower := strings.ToLower(modelName)
	if strings.Contains(lower, "tts") {
		return false
	}
	if strings.Contains(lower, "omni") {
		return false
	}
	if strings.Contains(lower, "qwen-deep-research") {
		return false
	}
	return true
}

// Copyright (c) Alibaba, Inc. and its affiliates.

package utils

import (
	"os"
	"strings"

	"github.com/liuxiaobopro/dashscope-go/common"
)

// CheckAndUploadLocal if path is local file, caller should upload to OSS first.
// This helper reports whether the value looks like a local file path.
func CheckAndUploadLocal(model, value string) (uploaded bool, url string, err error) {
	if value == "" {
		return false, value, nil
	}
	if strings.HasPrefix(value, common.FILE_PATH_SCHEMA) {
		p := strings.TrimPrefix(value, common.FILE_PATH_SCHEMA)
		if _, e := os.Stat(p); e != nil {
			return false, "", e
		}
		return true, value, nil
	}
	if common.IsPath(value) {
		return true, value, nil
	}
	return false, value, nil
}

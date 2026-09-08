// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import "net/http"

const (
	DASHSCOPE_API_KEY_ENV              = "DASHSCOPE_API_KEY"
	DASHSCOPE_API_KEY_FILE_PATH_ENV    = "DASHSCOPE_API_KEY_FILE_PATH"
	DASHSCOPE_API_REGION_ENV           = "DASHSCOPE_API_REGION"
	DASHSCOPE_API_VERSION_ENV          = "DASHSCOPE_API_VERSION"
	DASHSCOPE_DISABLE_DATA_INSPECTION_ENV = "DASHSCOPE_DISABLE_DATA_INSPECTION"
	DASHSCOPE_LOGGING_LEVEL_ENV        = "DASHSCOPE_LOGGING_LEVEL"

	DEFAULT_REQUEST_TIMEOUT_SECONDS = 300
	REQUEST_TIMEOUT_KEYWORD         = "request_timeout"
	SERVICE_API_PATH                = "services"

	PROMPT                = "prompt"
	MESSAGES              = "messages"
	NEGATIVE_PROMPT       = "negative_prompt"
	HISTORY               = "history"
	CUSTOMIZED_MODEL_ID   = "customized_model_id"
	IMAGES                = "images"
	REFERENCE_VIDEO_URLS  = "reference_video_urls"
	REFERENCE_URLS        = "reference_urls"
	MEDIA_URLS            = "media"
	TEXT_EMBEDDING_INPUT_KEY = "texts"
	SERVICE_503_MESSAGE   = "Service temporarily unavailable, possibly overloaded or not ready."
	WEBSOCKET_ERROR_CODE  = 44
	SSE_CONTENT_TYPE      = "text/event-stream"
	DEPRECATED_MESSAGE    = "history and auto_history are deprecated for qwen serial models and will be remove in future, use messages"
	SCENE                 = "scene"
	MESSAGE               = "message"
	REQUEST_CONTENT_TEXT  = "text"
	REQUEST_CONTENT_IMAGE = "image"
	REQUEST_CONTENT_AUDIO = "audio"
	FILE_PATH_SCHEMA      = "file://"

	ENCRYPTION_AES_SECRET_KEY_BYTES = 32
	ENCRYPTION_AES_IV_LENGTH        = 12
)

var RepeatableStatus = []int{
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
}

// FilePurpose file upload purpose.
type FilePurpose struct{}

const (
	FilePurposeFineTune   = "fine_tune"
	FilePurposeAssistants = "assistants"
)

// DeploymentStatus deployment status.
const (
	DeploymentStatusDeploying = "DEPLOYING"
	DeploymentStatusServing   = "RUNNING"
	DeploymentStatusDeleting  = "DELETING"
	DeploymentStatusFailed    = "FAILED"
	DeploymentStatusPending   = "PENDING"
)

// ApiProtocol API protocol.
const (
	ApiProtocolWebsocket = "websocket"
	ApiProtocolHTTP      = "http"
	ApiProtocolHTTPS     = "https"
)

// HTTPMethod HTTP methods.
const (
	HTTPMethodGET     = "GET"
	HTTPMethodHEAD    = "HEAD"
	HTTPMethodPOST    = "POST"
	HTTPMethodPUT     = "PUT"
	HTTPMethodDELETE  = "DELETE"
	HTTPMethodCONNECT = "CONNECT"
	HTTPMethodOPTIONS = "OPTIONS"
	HTTPMethodTRACE   = "TRACE"
	HTTPMethodPATCH   = "PATCH"
)

// TaskStatus async task status.
const (
	TaskStatusPending   = "PENDING"
	TaskStatusSuspended = "SUSPENDED"
	TaskStatusSucceeded = "SUCCEEDED"
	TaskStatusCanceled  = "CANCELED"
	TaskStatusRunning   = "RUNNING"
	TaskStatusFailed    = "FAILED"
	TaskStatusUnknown   = "UNKNOWN"
)

// Tasks task names.
const (
	TaskTextGeneration         = "text-generation"
	TaskAutoSpeechRecognition  = "asr"
)

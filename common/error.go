// Copyright (c) Alibaba, Inc. and its affiliates.

package common

import (
	"encoding/json"
	"fmt"
)

// DashScopeException is the base error type.
type DashScopeException struct {
	Msg string
}

func (e *DashScopeException) Error() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

func newErr(msg string) *DashScopeException {
	return &DashScopeException{Msg: msg}
}

// AuthenticationError no api key provided.
type AuthenticationError struct{ DashScopeException }

func NewAuthenticationError(msg string) error {
	return &AuthenticationError{DashScopeException{Msg: msg}}
}

// InvalidParameter invalid parameter.
type InvalidParameter struct{ DashScopeException }

func NewInvalidParameter(msg string) error {
	return &InvalidParameter{DashScopeException{Msg: msg}}
}

// InvalidTask invalid async task.
type InvalidTask struct{ DashScopeException }

func NewInvalidTask(msg string) error {
	return &InvalidTask{DashScopeException{Msg: msg}}
}

// UnsupportedModel unsupported model.
type UnsupportedModel struct{ DashScopeException }

func NewUnsupportedModel(msg string) error {
	return &UnsupportedModel{DashScopeException{Msg: msg}}
}

// UnsupportedTask unsupported task.
type UnsupportedTask struct{ DashScopeException }

func NewUnsupportedTask(msg string) error {
	return &UnsupportedTask{DashScopeException{Msg: msg}}
}

// ModelRequired model is required.
type ModelRequired struct{ DashScopeException }

func NewModelRequired(msg string) error {
	return &ModelRequired{DashScopeException{Msg: msg}}
}

// InvalidModel invalid model.
type InvalidModel struct{ DashScopeException }

func NewInvalidModel(msg string) error {
	return &InvalidModel{DashScopeException{Msg: msg}}
}

// InvalidInput invalid input.
type InvalidInput struct{ DashScopeException }

func NewInvalidInput(msg string) error {
	return &InvalidInput{DashScopeException{Msg: msg}}
}

// InvalidFileFormat invalid file format.
type InvalidFileFormat struct{ DashScopeException }

func NewInvalidFileFormat(msg string) error {
	return &InvalidFileFormat{DashScopeException{Msg: msg}}
}

// UnsupportedApiProtocol unsupported protocol.
type UnsupportedApiProtocol struct{ DashScopeException }

func NewUnsupportedApiProtocol(msg string) error {
	return &UnsupportedApiProtocol{DashScopeException{Msg: msg}}
}

// NotImplementedError not implemented.
type NotImplementedError struct{ DashScopeException }

func NewNotImplemented(msg string) error {
	return &NotImplementedError{DashScopeException{Msg: msg}}
}

// MultiInputsWithBinaryNotSupported binary with multi inputs not supported.
type MultiInputsWithBinaryNotSupported struct{ DashScopeException }

func NewMultiInputsWithBinaryNotSupported(msg string) error {
	return &MultiInputsWithBinaryNotSupported{DashScopeException{Msg: msg}}
}

// UnexpectedMessageReceived unexpected websocket message.
type UnexpectedMessageReceived struct{ DashScopeException }

func NewUnexpectedMessageReceived(msg string) error {
	return &UnexpectedMessageReceived{DashScopeException{Msg: msg}}
}

// UnsupportedData unsupported data.
type UnsupportedData struct{ DashScopeException }

func NewUnsupportedData(msg string) error {
	return &UnsupportedData{DashScopeException{Msg: msg}}
}

// AssistantError assistant API error.
type AssistantError struct {
	Message   string
	Code      string
	RequestID string
}

func NewAssistantError(messageJSON string) *AssistantError {
	e := &AssistantError{}
	var msg map[string]any
	if err := json.Unmarshal([]byte(messageJSON), &msg); err == nil {
		if v, ok := msg["request_id"].(string); ok {
			e.RequestID = v
		}
		if v, ok := msg["code"].(string); ok {
			e.Code = v
		}
		if v, ok := msg["message"].(string); ok {
			e.Message = v
		}
	}
	return e
}

func (e *AssistantError) Error() string {
	return fmt.Sprintf("Request failed, request_id: %s, code: %s, message: %s", e.RequestID, e.Code, e.Message)
}

// RequestFailure server generation or inference error.
type RequestFailure struct {
	RequestID string
	Message   string
	Name      string
	HTTPCode  int
}

func (e *RequestFailure) Error() string {
	return fmt.Sprintf("Request failed, request_id: %s, http_code: %d error_name: %s, error_message: %s",
		e.RequestID, e.HTTPCode, e.Name, e.Message)
}

// UnknownMessageReceived unknown websocket message.
type UnknownMessageReceived struct{ DashScopeException }

func NewUnknownMessageReceived(msg string) error {
	return &UnknownMessageReceived{DashScopeException{Msg: msg}}
}

// InputDataRequired input data required.
type InputDataRequired struct{ DashScopeException }

func NewInputDataRequired(msg string) error {
	return &InputDataRequired{DashScopeException{Msg: msg}}
}

// InputRequired required input missing.
type InputRequired struct{ DashScopeException }

func NewInputRequired(msg string) error {
	return &InputRequired{DashScopeException{Msg: msg}}
}

// UnsupportedDataType unsupported data type.
type UnsupportedDataType struct{ DashScopeException }

func NewUnsupportedDataType(msg string) error {
	return &UnsupportedDataType{DashScopeException{Msg: msg}}
}

// ServiceUnavailableError service unavailable.
type ServiceUnavailableError struct{ DashScopeException }

func NewServiceUnavailableError(msg string) error {
	return &ServiceUnavailableError{DashScopeException{Msg: msg}}
}

// UnsupportedHTTPMethod unsupported HTTP method.
type UnsupportedHTTPMethod struct{ DashScopeException }

func NewUnsupportedHTTPMethod(msg string) error {
	return &UnsupportedHTTPMethod{DashScopeException{Msg: msg}}
}

// AsyncTaskCreateFailed async task create failed.
type AsyncTaskCreateFailed struct{ DashScopeException }

func NewAsyncTaskCreateFailed(msg string) error {
	return &AsyncTaskCreateFailed{DashScopeException{Msg: msg}}
}

// UploadFileException file upload error.
type UploadFileException struct{ DashScopeException }

func NewUploadFileException(msg string) error {
	return &UploadFileException{DashScopeException{Msg: msg}}
}

// TimeoutException timeout.
type TimeoutException struct{ DashScopeException }

func NewTimeoutException(msg string) error {
	return &TimeoutException{DashScopeException{Msg: msg}}
}

var _ = newErr

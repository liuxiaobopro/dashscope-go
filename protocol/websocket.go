// Copyright (c) Alibaba, Inc. and its affiliates.

package protocol

// WebsocketStreamingMode websocket stream mode.
const (
	WebsocketStreamingModeNone   = "none"
	WebsocketStreamingModeIn     = "in"
	WebsocketStreamingModeOut    = "out"
	WebsocketStreamingModeDuplex = "duplex"
)

const (
	ActionKey    = "action"
	EventKey     = "event"
	Header       = "header"
	TaskID       = "task_id"
	ErrorName    = "error_code"
	ErrorMessage = "error_message"
)

// EventType websocket events.
const (
	EventStarted   = "task-started"
	EventGenerated = "result-generated"
	EventFinished  = "task-finished"
	EventFailed    = "task-failed"
)

// ActionType websocket actions.
const (
	ActionStart    = "run-task"
	ActionContinue = "continue-task"
	ActionFinished = "finish-task"
)

// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"encoding/json"
)

// ApiRequestData request payload builder.
type ApiRequestData struct {
	Model         string
	Task          string
	TaskGroup     string
	Function      string
	Input         any
	Parameters    map[string]any
	Form          map[string]any
	Resources     any
	IsBinaryInput bool
	APIProtocol   string
}

// NewApiRequestData create request data.
func NewApiRequestData(model, taskGroup, task, function string, input any, form map[string]any, isBinary bool, apiProtocol string) *ApiRequestData {
	return &ApiRequestData{
		Model:         model,
		Task:          task,
		TaskGroup:     taskGroup,
		Function:      function,
		Input:         input,
		Parameters:    map[string]any{},
		Form:          form,
		IsBinaryInput: isBinary,
		APIProtocol:   apiProtocol,
	}
}

func (d *ApiRequestData) AddParameters(params map[string]any) {
	for k, v := range params {
		d.Parameters[k] = v
	}
}

func (d *ApiRequestData) AddResources(resources any) {
	d.Resources = resources
}

func (d *ApiRequestData) ToRequestObject() map[string]any {
	o := map[string]any{}
	if d.Model != "" {
		o["model"] = d.Model
	}
	if d.Input != nil {
		o["input"] = d.Input
	}
	if len(d.Parameters) > 0 {
		o["parameters"] = d.Parameters
	}
	if d.Resources != nil {
		o["resources"] = d.Resources
	}
	return o
}

func (d *ApiRequestData) GetHTTPPayload() (isForm bool, form map[string]any, data map[string]any) {
	data = d.ToRequestObject()
	if d.Form != nil {
		return true, d.Form, data
	}
	return false, nil, data
}

func (d *ApiRequestData) GetWebsocketStartData() map[string]any {
	if d.IsBinaryInput {
		return d.onlyParameters()
	}
	obj := map[string]any{}
	if d.Model != "" {
		obj["model"] = d.Model
	}
	if d.Input != nil {
		obj["input"] = d.Input
	}
	if len(d.Parameters) > 0 {
		obj["parameters"] = d.Parameters
	}
	if d.Task != "" {
		obj["task"] = d.Task
	}
	if d.TaskGroup != "" {
		obj["task_group"] = d.TaskGroup
	}
	if d.Function != "" {
		obj["function"] = d.Function
	}
	if d.Resources != nil {
		obj["resources"] = d.Resources
	}
	return obj
}

func (d *ApiRequestData) onlyParameters() map[string]any {
	params := map[string]any{}
	for k, v := range d.Parameters {
		params[k] = v
	}
	var tempInput any
	if v, ok := params["raw_input"]; ok {
		tempInput = v
		delete(params, "raw_input")
	}
	obj := map[string]any{"model": d.Model, "parameters": params, "input": map[string]any{}}
	if tempInput != nil {
		obj["input"] = tempInput
	}
	if d.Task != "" {
		obj["task"] = d.Task
	}
	if d.TaskGroup != "" {
		obj["task_group"] = d.TaskGroup
	}
	if d.Function != "" {
		obj["function"] = d.Function
	}
	if d.Resources != nil {
		obj["resources"] = d.Resources
	}
	return obj
}

func fmtString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

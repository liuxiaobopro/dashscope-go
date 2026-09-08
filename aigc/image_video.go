// Copyright (c) Alibaba, Inc. and its affiliates.

package aigc

import (
	"context"

	"github.com/liuxiaobopro/dashscope-go/apientities"
	"github.com/liuxiaobopro/dashscope-go/client"
	"github.com/liuxiaobopro/dashscope-go/common"
)

const (
	WanxV1               = "wanx-v1"
	WanxSketchToImageV1  = "wanx-sketch-to-image-v1"
	Wanx21ImageEdit      = "wanx2.1-imageedit"
)

// ImageSynthesisCallParams Call image(s) synthesis service and get result.
type ImageSynthesisCallParams struct {
	Model           string
	Prompt          any
	NegativePrompt  any
	Images          []string
	APIKey          string
	SketchImageURL  string
	RefImg          string
	Workspace       string
	ExtraInput      map[string]any
	Task            string
	Function        string
	MaskImageURL    string
	BaseImageURL    string
	Size            *string
	N               *int
	Seed            *int
	Style           *string
	RefStrength     *float64
	RefMode         *string
	PromptExtend    *bool
	Watermark       *bool
	BBoxList        []any
	EnableSequential *bool
	ThinkingMode    *string
	ColorPalette    *string
	WaitTimeout     int
	Extra           map[string]any
}

func (p *ImageSynthesisCallParams) toCall() *client.CallParams {
	input := map[string]any{}
	if p.Prompt != nil {
		input[common.PROMPT] = p.Prompt
	}
	if p.NegativePrompt != nil {
		input[common.NEGATIVE_PROMPT] = p.NegativePrompt
	}
	if p.Images != nil {
		input[common.IMAGES] = p.Images
	}
	if p.SketchImageURL != "" {
		input["sketch_image_url"] = p.SketchImageURL
	}
	if p.RefImg != "" {
		input["ref_img"] = p.RefImg
	}
	if p.MaskImageURL != "" {
		input["mask_image_url"] = p.MaskImageURL
	}
	if p.BaseImageURL != "" {
		input["base_image_url"] = p.BaseImageURL
	}
	for k, v := range p.ExtraInput {
		input[k] = v
	}
	params := map[string]any{}
	if p.Size != nil {
		params["size"] = *p.Size
	}
	if p.N != nil {
		params["n"] = *p.N
	}
	if p.Seed != nil {
		params["seed"] = *p.Seed
	}
	if p.Style != nil {
		params["style"] = *p.Style
	}
	if p.RefStrength != nil {
		params["ref_strength"] = *p.RefStrength
	}
	if p.RefMode != nil {
		params["ref_mode"] = *p.RefMode
	}
	if p.PromptExtend != nil {
		params["prompt_extend"] = *p.PromptExtend
	}
	if p.Watermark != nil {
		params["watermark"] = *p.Watermark
	}
	if p.BBoxList != nil {
		params["bbox_list"] = p.BBoxList
	}
	if p.EnableSequential != nil {
		params["enable_sequential"] = *p.EnableSequential
	}
	if p.ThinkingMode != nil {
		params["thinking_mode"] = *p.ThinkingMode
	}
	if p.ColorPalette != nil {
		params["color_palette"] = *p.ColorPalette
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	task := p.Task
	if task == "" {
		task = "text2image"
	}
	fn := p.Function
	if fn == "" {
		fn = "image-synthesis"
	}
	return &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "aigc", Task: task, Function: fn,
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: params, SDKModule: "aigc",
	}
}

// ImageSynthesisService API for image synthesis.
type ImageSynthesisService struct{ Task string }

func (ImageSynthesisService) Call(ctx context.Context, p *ImageSynthesisCallParams) (*apientities.ImageSynthesisResponse, error) {
	if p.Prompt == nil {
		return nil, common.NewInputRequired("prompt is required!")
	}
	rsp, err := (client.BaseAsyncApi{}).Call(ctx, p.toCall(), p.WaitTimeout)
	if err != nil {
		return nil, err
	}
	return apientities.ImageSynthesisResponseFromAPI(rsp), nil
}

func (ImageSynthesisService) AsyncCall(ctx context.Context, p *ImageSynthesisCallParams) (*apientities.ImageSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).AsyncCall(ctx, p.toCall())
	if err != nil {
		return nil, err
	}
	return apientities.ImageSynthesisResponseFromAPI(rsp), nil
}

func (ImageSynthesisService) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int) (*apientities.ImageSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Wait(ctx, task, apiKey, workspace, waitTimeout, "")
	if err != nil {
		return nil, err
	}
	return apientities.ImageSynthesisResponseFromAPI(rsp), nil
}

func (ImageSynthesisService) Fetch(ctx context.Context, task any, apiKey, workspace string) (*apientities.ImageSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Fetch(ctx, task, apiKey, workspace, "")
	if err != nil {
		return nil, err
	}
	return apientities.ImageSynthesisResponseFromAPI(rsp), nil
}

func (ImageSynthesisService) Cancel(ctx context.Context, task any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Cancel(ctx, task, apiKey, workspace, "")
}

var ImageSynthesis = ImageSynthesisService{Task: "text2image"}

const (
	WanxTxt2VideoPro  = "wanx-txt2video-pro"
	WanxImg2VideoPro  = "wanx-img2video-pro"
	Wanx21T2VTurbo    = "wanx2.1-t2v-turbo"
	Wanx21T2VPlus     = "wanx2.1-t2v-plus"
	Wanx21I2VPlus     = "wanx2.1-i2v-plus"
	Wanx21I2VTurbo    = "wanx2.1-i2v-turbo"
	Wanx21Kf2VPlus    = "wanx2.1-kf2v-plus"
	WanxKf2V          = "wanx-kf2v"
)

const (
	MediaFirstFrame     = "first_frame"
	MediaLastFrame      = "last_frame"
	MediaReferenceImage = "reference_image"
	MediaReferenceVideo = "reference_video"
	MediaReferenceVoice = "reference_voice"
	MediaReferenceAudio = "reference_audio"
	MediaVideo          = "video"
	MediaFirstClip      = "first_clip"
	MediaDrivingAudio   = "driving_audio"
	MediaFile           = "file"
	MediaLink           = "link"
)

// VideoSynthesisCallParams API for video synthesis.
type VideoSynthesisCallParams struct {
	Model                      string
	Prompt                     any
	ExtendPrompt               *bool
	NegativePrompt             string
	Template                   string
	ImgURL                     string
	AudioURL                   string
	ReferenceVideoURLs         []string
	ReferenceURLs              []string
	ReferenceURL               string
	ReferenceVideoDescription  []string
	APIKey                     string
	ExtraInput                 map[string]any
	Workspace                  string
	Task                       string
	HeadFrame                  string
	TailFrame                  string
	FirstFrameURL              string
	LastFrameURL               string
	Media                      []map[string]any
	WaitTimeout                int
	Extra                      map[string]any
}

func (p *VideoSynthesisCallParams) toCall() *client.CallParams {
	input := map[string]any{}
	if p.Prompt != nil {
		input[common.PROMPT] = p.Prompt
	}
	if p.NegativePrompt != "" {
		input[common.NEGATIVE_PROMPT] = p.NegativePrompt
	}
	if p.ImgURL != "" {
		input["img_url"] = p.ImgURL
	}
	if p.AudioURL != "" {
		input["audio_url"] = p.AudioURL
	}
	if p.ReferenceVideoURLs != nil {
		input[common.REFERENCE_VIDEO_URLS] = p.ReferenceVideoURLs
	}
	if p.ReferenceURLs != nil {
		input[common.REFERENCE_URLS] = p.ReferenceURLs
	}
	if p.ReferenceURL != "" {
		input["reference_url"] = p.ReferenceURL
	}
	if p.HeadFrame != "" {
		input["head_frame"] = p.HeadFrame
	}
	if p.TailFrame != "" {
		input["tail_frame"] = p.TailFrame
	}
	if p.FirstFrameURL != "" {
		input["first_frame_url"] = p.FirstFrameURL
	}
	if p.LastFrameURL != "" {
		input["last_frame_url"] = p.LastFrameURL
	}
	if p.Media != nil {
		input[common.MEDIA_URLS] = p.Media
	}
	if p.Template != "" {
		input["template"] = p.Template
	}
	if p.ReferenceVideoDescription != nil {
		input["reference_video_description"] = p.ReferenceVideoDescription
	}
	for k, v := range p.ExtraInput {
		input[k] = v
	}
	params := map[string]any{}
	if p.ExtendPrompt != nil {
		params["prompt_extend"] = *p.ExtendPrompt
	}
	for k, v := range p.Extra {
		params[k] = v
	}
	task := p.Task
	if task == "" {
		task = "video-generation"
	}
	return &client.CallParams{
		Model: p.Model, Input: input, TaskGroup: "aigc", Task: task, Function: "video-synthesis",
		APIKey: p.APIKey, Workspace: p.Workspace, Parameters: params, SDKModule: "aigc",
	}
}

type VideoSynthesisService struct{ Task string }

func (VideoSynthesisService) Call(ctx context.Context, p *VideoSynthesisCallParams) (*apientities.VideoSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Call(ctx, p.toCall(), p.WaitTimeout)
	if err != nil {
		return nil, err
	}
	return apientities.VideoSynthesisResponseFromAPI(rsp), nil
}

func (VideoSynthesisService) AsyncCall(ctx context.Context, p *VideoSynthesisCallParams) (*apientities.VideoSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).AsyncCall(ctx, p.toCall())
	if err != nil {
		return nil, err
	}
	return apientities.VideoSynthesisResponseFromAPI(rsp), nil
}

func (VideoSynthesisService) Wait(ctx context.Context, task any, apiKey, workspace string, waitTimeout int) (*apientities.VideoSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Wait(ctx, task, apiKey, workspace, waitTimeout, "")
	if err != nil {
		return nil, err
	}
	return apientities.VideoSynthesisResponseFromAPI(rsp), nil
}

func (VideoSynthesisService) Fetch(ctx context.Context, task any, apiKey, workspace string) (*apientities.VideoSynthesisResponse, error) {
	rsp, err := (client.BaseAsyncApi{}).Fetch(ctx, task, apiKey, workspace, "")
	if err != nil {
		return nil, err
	}
	return apientities.VideoSynthesisResponseFromAPI(rsp), nil
}

func (VideoSynthesisService) Cancel(ctx context.Context, task any, apiKey, workspace string) (*apientities.DashScopeAPIResponse, error) {
	return (client.BaseAsyncApi{}).Cancel(ctx, task, apiKey, workspace, "")
}

var VideoSynthesis = VideoSynthesisService{Task: "video-generation"}

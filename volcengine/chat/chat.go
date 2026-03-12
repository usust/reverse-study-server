package chat

import (
	"context"
	"fmt"
	"os"
	dbmodel "reverse-study-server/volcengine/model"
	"strings"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	sdkmodel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

// PromptChatRequest 是调用模型 API 的请求参数。
type PromptChatRequest struct {
	ModelAPIID string                `json:"model_api_id"`
	Config     dbmodel.ModelAPIModel `json:"config"`
	Prompt     string                `json:"prompt"`
}

// PromptChatResponse 是调用模型 API 的结构化返回。
type PromptChatResponse struct {
	Request  PromptChatRequestInfo `json:"request"`
	Response PromptChatResult      `json:"response"`
	Info     PromptChatInfo        `json:"info"`
	Usage    ExecutePromptUsage    `json:"usage"`
}

/////////////////  PromptChatResponse 内部结构  >>>  /////////////////////////////////

// PromptChatRequestInfo 是请求侧关键信息。
type PromptChatRequestInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Provider     string `json:"provider"`
	Model        string `json:"model"`
	DisplayModel string `json:"display_model"`
}

// PromptChatResult 是模型响应正文。
type PromptChatResult struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content"`
	FinishReason     string `json:"finish_reason"`
}

// PromptChatInfo 是模型响应元信息。
type PromptChatInfo struct {
	RequestID   string `json:"request_id"`
	Object      string `json:"object"`
	Created     int64  `json:"created"`
	ServiceTier string `json:"service_tier"`
}

// ExecutePromptUsage 是 token 消耗信息。
type ExecutePromptUsage struct {
	PromptTokens         int `json:"prompt_tokens"`
	CompletionTokens     int `json:"completion_tokens"`
	TotalTokens          int `json:"total_tokens"`
	ReasoningTokens      int `json:"reasoning_tokens"`
	CachedPromptTokens   int `json:"cached_prompt_tokens"`
	ProvisionedTokensIn  int `json:"provisioned_tokens_in"`
	ProvisionedTokensOut int `json:"provisioned_tokens_out"`
}

/////////////////  <<<  PromptChatResponse 内部结构  /////////////////////////////////

// PromptChat 对外提供的唯一模型调用入口。
func PromptChat(ctx context.Context, chatPromptRequest PromptChatRequest) (PromptChatResponse, error) {
	if err := prepareCheckInput(chatPromptRequest); err != nil {
		return PromptChatResponse{}, err
	}

	client := arkruntime.NewClientWithApiKey(chatPromptRequest.Config.APIKey, arkruntime.WithBaseUrl(chatPromptRequest.Config.BaseURL))
	req := sdkmodel.CreateChatCompletionRequest{
		Model: chatPromptRequest.Config.APIModel,
		Messages: []*sdkmodel.ChatCompletionMessage{
			{
				Role: sdkmodel.ChatMessageRoleUser,
				Content: &sdkmodel.ChatCompletionMessageContent{
					ListValue: []*sdkmodel.ChatCompletionMessageContentPart{
						{
							Type: sdkmodel.ChatCompletionMessageContentPartTypeText,
							Text: chatPromptRequest.Prompt,
						},
					},
				},
			},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return PromptChatResponse{}, err
	}

	if len(resp.Choices) == 0 || resp.Choices[0] == nil {
		return PromptChatResponse{}, fmt.Errorf("empty model response")
	}

	provisionedTokensOut := 0
	if resp.Usage.CompletionTokensDetails.ProvisionedTokens != nil {
		provisionedTokensOut = *resp.Usage.CompletionTokensDetails.ProvisionedTokens
	}
	provisionedTokensIn := 0
	if resp.Usage.PromptTokensDetails.ProvisionedTokens != nil {
		provisionedTokensIn = *resp.Usage.PromptTokensDetails.ProvisionedTokens
	}

	result := PromptChatResponse{
		Request: PromptChatRequestInfo{
			ID:           chatPromptRequest.Config.ID,
			Name:         chatPromptRequest.Config.Name,
			Provider:     chatPromptRequest.Config.Provider,
			Model:        chatPromptRequest.Config.APIModel,
			DisplayModel: chatPromptRequest.Config.Model,
		},
		Info: PromptChatInfo{
			RequestID:   resp.ID,
			Object:      resp.Object,
			Created:     resp.Created,
			ServiceTier: resp.ServiceTier,
		},
		Usage: ExecutePromptUsage{
			PromptTokens:         resp.Usage.PromptTokens,
			CompletionTokens:     resp.Usage.CompletionTokens,
			TotalTokens:          resp.Usage.TotalTokens,
			ReasoningTokens:      resp.Usage.CompletionTokensDetails.ReasoningTokens,
			CachedPromptTokens:   resp.Usage.PromptTokensDetails.CachedTokens,
			ProvisionedTokensOut: provisionedTokensOut,
			ProvisionedTokensIn:  provisionedTokensIn,
		},
	}

	message := resp.Choices[0].Message
	if message.Content != nil && message.Content.StringValue != nil {
		result.Response.Content = *message.Content.StringValue
	}
	if message.ReasoningContent != nil {
		result.Response.ReasoningContent = *message.ReasoningContent
	}
	result.Response.FinishReason = string(resp.Choices[0].FinishReason)

	return result, nil
}

// prepareCheckInput 把参数校验和配置归一化集中在一起。
func prepareCheckInput(req PromptChatRequest) error {
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	config := req.Config
	if !config.Enabled {
		return fmt.Errorf("model api config is disabled")
	}

	model := strings.TrimSpace(config.APIModel)
	if model == "" {
		model = strings.TrimSpace(config.Model)
	}

	displayModel := strings.TrimSpace(config.Model)
	if displayModel == "" {
		displayModel = strings.TrimSpace(config.APIModel)
	}

	apiKey := strings.TrimSpace(config.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("ARK_API_KEY"))
	}
	if apiKey == "" {
		return fmt.Errorf("missing api key")
	}

	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		return fmt.Errorf("baseUrl is required")
	}

	return nil
}

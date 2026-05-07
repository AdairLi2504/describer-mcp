package internal

import (
	"context"
	"describer/internal/config"
	"describer/internal/messages/prompt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func DescribeHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get Parameters
	imageSource, err := req.RequireString("image")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	base64 := req.GetBool("base64", false)
	compress := req.GetBool("compress", false)
	// Encode Image
	img, err := EncodeImage(imageSource, base64, compress)
	if err != nil {
		errStr := err.Error()
		// Specail Type
		if os.IsExist(err) {
			return mcp.NewToolResultError("image not found on local"), nil
		} else if strings.HasPrefix(errStr, "request failed:") {
			return mcp.NewToolResultError(errStr), nil
		}
		switch errStr {
		case "unsupported image source type", "impossbile to compress image url", "reponse type is not image", "path is a directory":
			return mcp.NewToolResultError(errStr), nil
		default:
			return nil, err
		}
	}
	// Construct LLM Client
	llmConfig := config.Get()
	os.Setenv("OPENAI_API_KEY", llmConfig.APIKey)
	llm, err := openai.New(
		openai.WithBaseURL(llmConfig.APIEndpoint),
		openai.WithModel(llmConfig.Model),
	)
	if err != nil {
		return nil, err
	}
	// Construct Messages
	messages := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextPart(prompt.Describe),
				llms.ImageURLPart(img),
			},
		},
	}
	resp, err := llm.GenerateContent(ctx, messages)
	if err != nil {
		return nil, err
	}
	return ConstructTextToCallToolResult(resp.Choices[0].Content), nil
}

package main

import (
	//"context"

	"fmt"
	"log"

	. "github.com/AdairLi2504/describer-mcp/internal"
	"github.com/AdairLi2504/describer-mcp/internal/config"
	"github.com/AdairLi2504/describer-mcp/internal/messages/description"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load Config
	if err := config.Load(); err != nil {
		log.Fatalf("%v", err)
	}
	// Describe
	toolDescribe := mcp.NewTool("describe",
		mcp.WithDescription(description.DescribeTool),
		mcp.WithString("image",
			mcp.Required(),
			mcp.Description(description.DescribeImage)),
		mcp.WithBoolean("base64",
			mcp.DefaultBool(false),
			mcp.Description(description.DescribeBase64)),
		mcp.WithBoolean("compress",
			mcp.DefaultBool(false),
			mcp.Description(description.DescribeCompress)),
	)
	// Define the MCP server
	s := server.NewMCPServer(
		"Describer",
		"0.1.0",
		server.WithToolCapabilities(true),
	)
	s.AddTool(toolDescribe, DescribeHandler)
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

}

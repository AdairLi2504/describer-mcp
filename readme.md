# describer-MCP

[![Golang](https://img.shields.io/badge/golang-1.26+-00ADD8.svg)](https://go.dev)
[![MCP](https://img.shields.io/badge/MCP-Compatible-green.svg)](https://modelcontextprotocol.io/)

**An MCP server use vision-capable LLM to explain images to your agent, bridging the gap between visual input and agent reasoning**

## Overview

The describer is a role provides the audio description, additional commentary that explains what’s happening on screen. Similarly, this MCP can use vision-capable LLM to explain images to the agent, bridging the gap between visual input and agent reasoning. Saluate describers!

## Features

| Tool | Description |
|------|-------------|
| `describe` | Takes an image from a URL, local file path, or base64 input and produces a written description of that image. |

## Quick Start

### Requirements

Please make sure that you have installed golang, and the `go` command is directly available from the terminal. [How to install Golang](https://go.dev/doc/install)

### 1. Download

```bash
git clone https://github.com/AdairLi2504/describer-mcp
```

### 2. Configure Environment Variables

Before running the application, you need to set up the following environment variables to authenticate and connect to the vision-capable LLM API.

| Variable | Default | Description |
|---|---|---|
| `DESCRIBER_API_ENDPOIN` | *(required)* | OpenAI-compatible chat completions endpoint like 'https://your-provider.ai/api/v1/'. The last path segment should be the version, e.g., 'v1'. |
| `DESCRIBER_API_KEY` | *(required)* | Model name |
| `DESCRIBER_MODEL` | *(empty)* | API key (optional for local models) |

### 3. Run

```bash
go run -C /path/to/describer-mcp/ main.go
```

The server communicates over stdio, designed to be launched by an MCP client such as Openclaw/Nanobot.

> Replace `/path/to/describer-mcp/` with the actual install path.

### Alternative: Install

Before you start, confirm that the `bin` folder path under `GOPATH` (check it by `go env GOPATH`) has been written into the PATH environment variable. If you do not know what is that mean, please read [How to configure GOPATH properly](https://labex.io/tutorials/go-how-to-configure-gopath-properly-451553). The command `go install` will install the tool into this path.

```bash
go install github.com/AdairLi2504/describer-mcp@latest # Install directly
```

Then you can use `describer` to run it directly.

## Intergation

Add to your `config.json`:

```json
{
  "mcpServers": {
    "describer-mcp": {
      "command": "go",
      "args": ["run","-C","/path/to/describer-mcp/","main.go"],
	    "env": {
          "DESCRIBER_API_ENDPOIN": "",
          "DESCRIBER_API_KEY":"",
          "DESCRIBER_MODEL": ""
        }
    }
  }
}
```
Replace the empty strings with your actual values.

> Replace `/path/to/describer-mcp` with the actual install path.

If you have already installed that, you can call it directly

```json
{
  "mcpServers": {
    "describer-mcp": {
      "command": "describer",
	    "env": {
          "DESCRIBER_API_ENDPOINT": "",
          "DESCRIBER_API_KEY":"",
          "DESCRIBER_MODEL": ""
        }
    }
  }
}
```
Replace the empty strings with your actual values.



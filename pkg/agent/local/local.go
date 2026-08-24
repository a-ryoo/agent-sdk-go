// Package local exposes the SDK's in-process agent loop without importing the full agent facade.
package local

import (
	"context"
	"fmt"

	sdkruntime "github.com/agenticenv/agent-sdk-go/internal/runtime"
	localruntime "github.com/agenticenv/agent-sdk-go/internal/runtime/local"
	"github.com/agenticenv/agent-sdk-go/internal/types"
	"github.com/agenticenv/agent-sdk-go/pkg/interfaces"
	"github.com/agenticenv/agent-sdk-go/pkg/logger"
)

// Config is the minimal configuration for a sequential, single-attempt local run.
type Config struct {
	Name          string
	SystemPrompt  string
	LLMClient     interfaces.LLMClient
	Tools         []interfaces.Tool
	MaxIterations int
}

// Run executes the local agent loop with retries disabled for LLM and tool calls.
func Run(ctx context.Context, config Config, prompt string) (string, error) {
	runtime, err := localruntime.NewLocalRuntime(
		localruntime.WithLogger(logger.NoopLogger()),
		localruntime.WithToolExecutionMode(types.AgentToolExecutionModeSequential),
		localruntime.WithAgentSpec(sdkruntime.AgentSpec{Name: config.Name, SystemPrompt: config.SystemPrompt}),
		localruntime.WithAgentConfig(sdkruntime.AgentConfig{
			LLM:    sdkruntime.AgentLLM{Client: config.LLMClient},
			Limits: sdkruntime.AgentLimits{MaxIterations: config.MaxIterations},
			ExecutionConfigs: sdkruntime.ExecutionConfigs{
				LLM:         sdkruntime.ExecutionConfig{MaxAttempts: 1},
				ToolExecute: sdkruntime.ExecutionConfig{MaxAttempts: 1},
			},
		}),
	)
	if err != nil {
		return "", err
	}
	defer runtime.Close()
	handle, err := runtime.Run(ctx, &sdkruntime.RunRequest{UserPrompt: prompt, Tools: config.Tools})
	if err != nil {
		return "", err
	}
	result, err := handle.Get(ctx)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", fmt.Errorf("local agent returned no result")
	}
	return result.Content, nil
}

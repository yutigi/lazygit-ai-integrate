package oscommands

import (
	"errors"
	"io"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestAICommitArgs(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			args, err := aiCommitArgs(config.AIConfig{Provider: provider})
			assert.NoError(t, err)
			assert.Equal(t, provider, args[0])
			assert.NotContains(t, args, "--model")

			args, err = aiCommitArgs(config.AIConfig{Provider: provider, Model: "custom-model-1.0"})
			assert.NoError(t, err)
			assert.Contains(t, args, "--model")
			assert.Contains(t, args, "custom-model-1.0")
			if provider == "codex" {
				assert.Equal(t, "-", args[len(args)-1])
				assert.Contains(t, args, "read-only")
			} else {
				assert.Contains(t, args, "--tools")
				assert.Contains(t, args, "--strict-mcp-config")
			}
		})
	}
	_, err := aiCommitArgs(config.AIConfig{Provider: "unknown"})
	assert.Error(t, err)
	for _, model := range []string{"--help", "model & echo injected", `model"`, "%PATH%", "model\nname"} {
		_, err := aiCommitArgs(config.AIConfig{Provider: "claude", Model: model})
		assert.Error(t, err, model)
	}
}

func TestGenerateCommitMessage(t *testing.T) {
	for _, provider := range []string{"claude", "codex"} {
		t.Run(provider, func(t *testing.T) {
			const diff = "diff --git a/file b/file\n+new content"
			runner := NewFakeRunner(t).ExpectFunc("AI receives the diff on stdin", func(cmd *CmdObj) bool {
				assert.Equal(t, provider, cmd.Args()[0])
				input, err := io.ReadAll(cmd.GetCmd().Stdin)
				assert.NoError(t, err)
				assert.Equal(t, commitMessagePrompt+diff, string(input))
				assert.NotContains(t, cmd.ToString(), diff)
				assert.False(t, cmd.ShouldLog())
				return true
			}, "  Explain the purpose\r\n\r\nUseful detail\r\n", nil)
			message, err := NewDummyOSCommandWithRunner(runner).GenerateCommitMessage(config.AIConfig{Provider: provider}, diff)
			assert.NoError(t, err)
			assert.Equal(t, "Explain the purpose\n\nUseful detail", message)
			runner.CheckForMissingCalls()
		})
	}
}

func TestGenerateCommitMessageFailure(t *testing.T) {
	for _, test := range []struct {
		name   string
		output string
		err    error
	}{
		{name: "empty response", output: " \n"},
		{name: "CLI failure", err: errors.New("not authenticated")},
	} {
		t.Run(test.name, func(t *testing.T) {
			runner := NewFakeRunner(t).ExpectFunc("AI CLI", func(*CmdObj) bool { return true }, test.output, test.err)
			message, err := NewDummyOSCommandWithRunner(runner).GenerateCommitMessage(config.AIConfig{Provider: "claude"}, "staged diff")
			assert.Error(t, err)
			assert.Empty(t, message)
			runner.CheckForMissingCalls()
		})
	}
	message, err := NewDummyOSCommandWithRunner(NewFakeRunner(t)).GenerateCommitMessage(config.AIConfig{Provider: "claude"}, "")
	assert.Error(t, err)
	assert.Empty(t, message)
}

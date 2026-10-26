package oscommands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCommitMessageWindowsLauncher(t *testing.T) {
	binDir := t.TempDir()
	launcher := filepath.Join(binDir, "codex.cmd")
	assert.NoError(t, os.WriteFile(launcher, []byte("@echo off\r\necho Draft from Windows launcher\r\n"), 0o600))
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	c := NewDummyOSCommand()
	c.Platform = &Platform{OS: "windows", Shell: "cmd.exe", ShellArg: "/c"}
	c.Cmd.platform = c.Platform
	message, err := c.GenerateCommitMessage(config.AIConfig{Provider: "codex", Model: "test-model"}, "staged changes")
	assert.NoError(t, err)
	assert.Equal(t, "Draft from Windows launcher", message)
}

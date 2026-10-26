package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMethodMenuWithAi = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Offer AI first when configured while keeping manual commit methods available",
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().AI = config.AIConfig{Provider: "claude", Model: "sonnet"}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Press(keys.Files.CommitMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Choose commit method")).
			TopLines(
				Contains("AI commit").IsSelected(),
				Contains("Enter commit message"),
				Contains("Commit changes using git editor"),
				Contains("Commit changes without pre-commit hook"),
			).
			Select(Contains("Enter commit message")).
			Confirm()
		t.ExpectPopup().CommitMessagePanel().Type("Write a manual message").Confirm()
		t.Views().Files().IsEmpty()
		t.Views().Commits().Lines(Contains("Write a manual message"))
	},
})

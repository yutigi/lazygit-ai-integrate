package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMethodMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Choose manual message entry from the commit-method menu without AI configured",
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().Press(keys.Files.CommitMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Choose commit method")).
			TopLines(
				Contains("Enter commit message").IsSelected(),
				Contains("Commit changes using git editor"),
				Contains("Commit changes without pre-commit hook"),
			).
			Confirm()
		t.ExpectPopup().CommitMessagePanel().Type("Explain the change").Confirm()
		t.Views().Files().IsEmpty()
		t.Views().Commits().Lines(Contains("Explain the change"))
	},
})

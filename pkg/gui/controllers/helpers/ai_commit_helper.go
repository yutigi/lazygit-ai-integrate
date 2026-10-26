package helpers

import (
	"errors"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (self *WorkingTreeHelper) OpenCommitMethodMenu() error {
	items := []*types.MenuItem{}
	if self.c.UserConfig().AI.Provider != "" {
		items = append(items, &types.MenuItem{
			Label:   self.c.Tr.AICommit,
			OnPress: self.handleAICommit,
		})
	}
	items = append(items,
		&types.MenuItem{Label: self.c.Tr.EnterCommitMessage, OnPress: self.HandleCommitPress},
		&types.MenuItem{Label: self.c.Tr.CommitChangesWithEditor, OnPress: self.HandleCommitEditorPress},
		&types.MenuItem{Label: self.c.Tr.CommitChangesWithoutHook, OnPress: self.HandleWIPCommitPress},
	)
	return self.c.Menu(types.CreateMenuOptions{
		Title: self.c.Tr.CommitMethodTitle,
		Items: items,
	})
}

func (self *WorkingTreeHelper) handleAICommit() error {
	cfg := self.c.UserConfig().AI
	return self.WithEnsureCommittableFiles(func() error {
		return self.c.WithWaitingStatusBlockingInput(types.WaitingStatusOpts{
			Message: self.c.Tr.GeneratingCommitMessage,
		}, func(gocui.Task) error {
			diff, err := self.c.OS().Cmd.New([]string{
				"git", "diff", "--cached", "--no-ext-diff", "--no-textconv", "--no-color",
			}).DontLog().RunWithOutput()
			if err != nil {
				return err
			}
			message, err := self.c.OS().GenerateCommitMessage(cfg, diff)
			if err != nil {
				return errors.New(utils.ResolvePlaceholderString(self.c.Tr.AICommitFailed, map[string]string{"error": err.Error()}))
			}
			self.c.OnUIThread(func() error {
				return self.HandleCommitPressWithMessage(message, false)
			})
			return nil
		})
	})
}

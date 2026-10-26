# AI commit messages

Press `C` in the files panel or its working-tree diff to choose how to commit.
The menu offers message entry, the Git editor, and committing without hooks.
Lowercase `c` still opens message entry directly.

To add **AI commit** as the first menu option, install and authenticate either
Claude Code or Codex CLI, then configure a provider in `config.yml`:

```yaml
ai:
  provider: claude # claude or codex
  model: sonnet   # model name/alias; empty uses the CLI's default
```

For Codex, set `provider: codex` and set `model` to a model available to your
Codex account, or leave it empty. The selected CLI must be on your PATH.
Lazygit uses its existing authentication; no API key is stored in this config.

Choose **AI commit** to generate a message from the staged diff. If no files
are staged, the usual stage-all prompt applies. The staged diff is sent to the
configured provider through its CLI. Generation has a two-minute timeout.
The generated draft opens in the commit panel, where you can edit or cancel it.
Confirming the panel creates the commit using the usual hooks and signing.

Leave `ai.provider` empty to hide the AI option. The menu shortcut is
`keybinding.files.commitMenu`. Direct editor access remains available through
`keybinding.files.commitChangesWithEditor`, which has no default shortcut.

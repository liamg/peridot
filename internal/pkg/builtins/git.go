package builtin

import (
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func init() {

	git := module.NewFactory("git").
		WithInputs([]config.Variable{
			{
				Name:     "username",
				Required: true,
			},
			{
				Name:     "email",
				Required: true,
			},
			{
				Name:    "editor",
				Default: "vim",
			},
			{
				Name:    "signingkey",
				Default: "",
			},
			{
				Name:    "aliases",
				Default: []interface{}{},
			},
			{
				Name: "extra",
			},
			{
				Name:    "ignores",
				Default: []interface{}{},
			},
		}).
		WithFilesFunc(gitFiles).
		Build()

	module.RegisterBuiltin("git", git)
}

var (
	gitIgnoreTemplate = `{{ range .ignores }}{{ . }}
{{ end }}`
	gitConfigTemplate = `
[user]
	name = {{ .username }}
	email = {{ .email }}
    {{ if .signingkey }}signingkey = {{ .signingkey }}{{ end }}
[commit]
	gpgsign = {{ if .signingkey }}true{{ else }}false{{ end }}
[core]
	excludesfile = ~/.gitignore
{{ if .editor }}	editor = {{ .editor }}
{{ else }}	editor = vim
{{ end }}
[pull]
	rebase = true
[alias]
{{range .aliases}}	{{ . }}
{{end}}
{{ if .extra }}{{ .extra }}{{ end }}
`
)

func gitFiles(vars variable.Collection) []module.File {
	var files []module.File

	home := vars.Get("user_home_dir").AsString()

	files = append(files, module.NewMemoryFile(
		filepath.Join(home, ".gitignore"),
		gitIgnoreTemplate,
		true,
		vars,
	))

	files = append(files, module.NewMemoryFile(
		filepath.Join(home, ".gitconfig"),
		gitConfigTemplate,
		true,
		vars,
	))

	return files
}

package builtin

import (
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/builtins/validate"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func init() {
	module.RegisterBuiltin("git", &gitBuiltin{})
}

type gitBuiltin struct {
	variables variable.Collection
}

func (b *gitBuiltin) Name() string {
	return "builtin:git"
}

func (b *gitBuiltin) Path() string {
	panic("not implemented")
}

func (b *gitBuiltin) Children() []module.Module {
	return nil
}

func (b *gitBuiltin) Files() []module.File {
	var files []module.File

	home := b.variables.Get("user_home_dir").AsString()

	files = append(files, module.NewMemoryFile(
		filepath.Join(home, ".gitignore"),
		`{{ range .ignores }}{{ . }}
{{ end }}`,
		b.variables,
	))

	files = append(files, module.NewMemoryFile(
		filepath.Join(home, ".gitconfig"),
		`
[user]
	name = {{ .username }}
	email = {{ .email }}
[commit]
	gpgsign = true
[core]
{{ if .editor }}	editor = {{ .editor }}
{{ else }}	editor = vim
{{ end }}
[ignore]
	file = ~/.gitignore
[pull]
	rebase = true
[alias]
{{range .aliases}}	{{ . }}
{{end}}
{{ if .extra }}{{ .extra }}{{ end }}
`,
		b.variables,
	))

	return files
}

func (b *gitBuiltin) Validate() error {
	validator := validate.New(b.variables)
	return validator.EnsureDefined(
		"user_home_dir",
		"username",
		"email",
	)
}

func (b *gitBuiltin) RequiresUpdate() bool {
	return false
}

func (b *gitBuiltin) RequiresInstall() bool {
	return false
}

func (b *gitBuiltin) Install() error {
	return nil
}

func (b *gitBuiltin) Update() error {
	return nil
}

func (b *gitBuiltin) AfterFileChange() error {
	return nil
}

func (b *gitBuiltin) ApplyVariables(vars variable.Collection) error {
	b.variables = vars
	return nil
}

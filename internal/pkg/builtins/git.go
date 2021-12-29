package builtin

import (
	"fmt"
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/builtins/validate"
	"github.com/liamg/peridot/internal/pkg/module"
)

func init() {
	module.RegisterBuiltin("git", &gitBuiltin{})
}

type gitBuiltin struct {
	variables map[string]interface{}
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

	home := b.variables["user_home_dir"].(string)

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

	if err := validator.EnsureString("user_home_dir"); err != nil {
		return fmt.Errorf("no home directory available, cannot determine git config location")
	}

	if err := validator.EnsureStrings("username", "email"); err != nil {
		return err
	}

	if err := validator.EnsureStringIfDefined("extra"); err != nil {
		return err
	}

	if err := validator.EnsureStringIfDefined("editor"); err != nil {
		return err
	}

	if err := validator.EnsureStringSliceIfDefined("aliases"); err != nil {
		return err
	}

	return nil
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

func (b *gitBuiltin) ApplyVariables(vars map[string]interface{}) error {
	b.variables = vars
	return nil
}

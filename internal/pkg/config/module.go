package config

import "fmt"

type variableSniff struct {
	Variables []Variable `yaml:"variables"`
}

type Module struct {
	Name      string        `yaml:"-"`
	Dir       string        `yaml:"-"`
	Path      string        `yaml:"-"`
	Variables []Variable    `yaml:"variables"`
	Files     []File        `yaml:"files"`
	Modules   []InnerModule `yaml:"modules"`
	Scripts   Scripts       `yaml:"scripts"`
}

type Scripts struct {
	UpdateRequired  Script `yaml:"should_update"`
	Update          Script `yaml:"update"`
	InstallRequired Script `yaml:"should_install"`
	Install         Script `yaml:"install"`
	AfterFileChange Script `yaml:"after_file_change"`
}

type Script struct {
	Command string `yaml:"command"`
	Sudo    bool   `yaml:"sudo"`
}

type InnerModule struct {
	Name      string                 `yaml:"name"`
	Source    string                 `yaml:"source"`
	DependsOn []string               `yaml:"depends_on"`
	Variables map[string]interface{} `yaml:"variables"`
	Filters   Filters                `yaml:"filters"`
}

type Variable struct {
	Name     string      `yaml:"name"`
	Default  interface{} `yaml:"default"`
	Required bool        `yaml:"required"`
}

type File struct {
	Target            string `yaml:"target"`
	Source            string `yaml:"source"`
	DisableTemplating bool   `yaml:disable_templating`
}

func (m *Module) Validate() error {
	for _, v := range m.Variables {
		if v.Name == "" {
			return fmt.Errorf("module '%s' has a variable with no name", m.Name)
		}
		if !v.Required && v.Default == nil {
			return fmt.Errorf("variable '%s' in module '%s' is not required but has no default value set", v.Name, m.Name)
		}
		if v.Required && v.Default != nil {
			return fmt.Errorf("variable '%s' in module '%s' is always required but has a default value set", v.Name, m.Name)
		}
	}
	for _, child := range m.Modules {
		if child.Name == "" {
			return fmt.Errorf("module '%s' has a child module with no name", m.Name)
		}
		if child.Source == "" {
			return fmt.Errorf("module '%s' has a child module '%s' with no source", m.Name, child.Name)
		}
	}
	if m.Scripts.InstallRequired.Command != "" && m.Scripts.Install.Command == "" {
		return fmt.Errorf("module '%s' has a should_install script defined, but no install script", m.Name)
	} else if m.Scripts.Install.Command != "" && m.Scripts.InstallRequired.Command == "" {
		return fmt.Errorf("module '%s' has an install script defined, but no should_install script", m.Name)
	}
	if m.Scripts.UpdateRequired.Command != "" && m.Scripts.Update.Command == "" {
		return fmt.Errorf("module '%s' has a should_update script defined, but no update script", m.Name)
	} else if m.Scripts.Update.Command != "" && m.Scripts.UpdateRequired.Command == "" {
		return fmt.Errorf("module '%s' has an update script defined, but no should_update script", m.Name)
	}
	if m.Scripts.AfterFileChange.Command != "" && len(m.Files) == 0 {
		return fmt.Errorf("module '%s' has an after_file_change script defined, but has no files configured", m.Name)
	}

	return nil
}

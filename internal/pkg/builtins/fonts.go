package builtin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func getFontsDir(vars variable.Collection) (string, error) {
	// sort out a fonts directory
	var dir string
	if vars.Has("dir") {
		dir = vars.Get("dir").AsString()
	}

	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".local/share/fonts")
	}
	return dir, nil
}

func init() {

	fontsBuiltin := module.NewFactory("fonts").
		WithInputs([]config.Variable{
			{
				Name:     "files",
				Required: true,
			},
			{
				Name:    "dir",
				Default: "",
			},
		}).
		WithRequiresInstallFunc(func(_ *module.Runner, vars variable.Collection) (bool, error) {

			dir, err := getFontsDir(vars)
			if err != nil {
				return false, err
			}
			if err := os.MkdirAll(dir, 0700); err != nil {
				return false, err
			}

			for _, file := range vars.Get("files").AsList().All() {
				var filename string
				if strings.HasPrefix(file.AsString(), "./") {
					filename = filepath.Base(file.AsString())
				} else if parsed, err := url.Parse(file.AsString()); err == nil && parsed.Host != "" {
					filename = path.Base(parsed.Path)
				}

				expectedPath := filepath.Join(dir, filename)
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					return true, nil
				}
			}

			return false, nil
		}).
		WithInstallFunc(func(r *module.Runner, vars variable.Collection) error {

			dir, err := getFontsDir(vars)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(dir, 0700); err != nil {
				return err
			}

			for _, file := range vars.Get("files").AsList().All() {
				var filename string
				if strings.HasPrefix(file.AsString(), "./") {
					filename = filepath.Base(file.AsString())
					targetPath := filepath.Join(dir, filename)
					if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
						continue
					}

					fontData, err := ioutil.ReadFile(file.AsString())
					if err != nil {
						return fmt.Errorf("failed to install font from %s: %s", file.AsString(), err)
					}
					if err := ioutil.WriteFile(targetPath, fontData, 0600); err != nil {
						return err
					}
				} else if parsed, err := url.Parse(file.AsString()); err == nil && parsed.Host != "" {
					filename = path.Base(parsed.Path)
					targetPath := filepath.Join(dir, filename)
					if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
						continue
					}
					if err := func() error {
						resp, err := http.Get(parsed.String())
						if err != nil {
							return err
						}
						defer resp.Body.Close()

						fontData, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							return err
						}

						if err := ioutil.WriteFile(targetPath, fontData, 0600); err != nil {
							return err
						}

						return nil
					}(); err != nil {
						return err
					}
				} else {
					return fmt.Errorf("invalid font file: '%s'", file.AsString())
				}
			}

			return r.Run("fc-cache -f", false)
		}).
		Build()

	module.RegisterBuiltin("fonts", fontsBuiltin)
}

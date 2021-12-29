package config

type Filters struct {
	OperatingSystem []string `yaml:"os"`
	Distribution    []string `yaml:"distro"`
	Architecture    []string `yaml:"arch"`
}

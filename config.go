package main

type config struct {
	Https   bool                   `yaml:"https"`
	Host    string                 `yaml:"host"`
	Name    string                 `yaml:"name"`
	Files   []string               `syaml:"files"`
	Use     map[string][]string    `yaml:"use"`
	Service map[string]interface{} `yaml:"service"`
	Shared  map[string]interface{} `yaml:"shared"`
}

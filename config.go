package main

type config struct {
	Host    string                 `yaml:"host"`
	Name    string                 `yaml:"name"`
	Files   []string               `yaml:"files"`
	Service map[string]interface{} `yaml:"service"`
	Shared  map[string]interface{} `yaml:"shared"`
}

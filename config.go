package main

type config struct {
	Host    string                 `json:"host" yaml:"host"`
	Name    string                 `json:"name" yaml:"name"`
	Files   []string               `json:"files" yaml:"files"`
	Use     map[string][]string    `json:"use" yaml:"use"`
	Service map[string]interface{} `json:"service" yaml:"service"`
	Shared  map[string]interface{} `json:"shared" yaml:"shared"`
}

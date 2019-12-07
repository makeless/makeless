package main

type config struct {
	Host    string      `yaml:"host"`
	Service string      `yaml:"service"`
	Files   []string    `yaml:"files"`
	Compose interface{} `yaml:"compose"`
}

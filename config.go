package main

type config struct {
	Service string      `yaml:"service"`
	Files   []string    `yaml:"files"`
	Compose interface{} `yaml:"compose"`
}

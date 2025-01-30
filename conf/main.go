package conf

import (
	"mangia_nastri/logger"
	"os"

	"gopkg.in/yaml.v3"
)

var log = logger.New("conf")

type ignore struct {
	Headers []string `yaml:"headers"`
	Body    []string `yaml:"body"`
}

type Config struct {
	Ignore ignore `yaml:"ignore"`
	Time   string `yaml:"time"`
}

func New(fileName string) *Config {
	c := &Config{}

	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	log.Printf("Configuration loaded: %v", c)

	return c
}

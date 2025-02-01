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

type DataSource struct {
	Type string `yaml:"type"`
	URI  string `yaml:"uri"`
}

type Proxy struct {
	Port        string     `yaml:"port"`
	Name        string     `yaml:"name"`
	Destination string     `yaml:"destination"`
	Ignore      ignore     `yaml:"ignore"`
	DataSource  DataSource `yaml:"dataSource"`
}

type sharedConf struct {
	Ignore ignore `yaml:"ignore"`
}

type Config struct {
	SharedConf sharedConf `yaml:"sharedConf"`
	Proxy      []Proxy    `yaml:"proxy"`
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

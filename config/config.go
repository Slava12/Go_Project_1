package config

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type MyTime struct {
	Hour int 	`yaml:"hour"`
	Minute int 	`yaml:"minute"`
	Second int `yaml:"second"`
}

type MyRepository struct {
	Directory string 	`yaml:"directory"`
	File string 	`yaml:"file"`
}

type MySystem struct {
	Name string 	`yaml:"name"`
	Release string 	`yaml:"release"`
	Codename string 	`yaml:"codename"`
}

type Config struct {
	Language string 	`yaml:"language"`
    Time MyTime 	`yaml:"time"`
    Repository MyRepository 	`yaml:"repository"`
    System MySystem 	`yaml:"system"`
}

func LoadConfig(configurationPath string) Config {
	

    bytesFile, errorReadFile := ioutil.ReadFile(configurationPath)
    if errorReadFile != nil {
        log.Fatal("Файл конфигурации не был загружен!")
        
    }
    
	config0 := Config{}
	errorUnmarshal := yaml.Unmarshal(bytesFile, &config0)
	if errorUnmarshal != nil {
    	log.Fatalf("error: %v", errorUnmarshal)
    }
    return config0
    //log.Println("Загружен файл конфигурации:")
    //fmt.Printf("%+v\n", config)
}
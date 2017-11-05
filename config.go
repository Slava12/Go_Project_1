package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Port       string `yaml:"port"`
	Mongodb    string `yaml:"mongodb"`
	Tempfolder string `yaml:"tempfolder"`
	Folder     string `yaml:"folder"`
}

func LoadConfigFile() (Config, error) {

	config := Config{}

	configurationPath := flag.String("path", "", "Путь до файла конфигурации.")
	flag.Parse()

	bytesFile, errorReadFile := ioutil.ReadFile(*configurationPath)
	if errorReadFile != nil {
		log.Println("Файл конфигурации не был загружен!")
		return config, errorReadFile
	}

	errorUnmarshal := yaml.Unmarshal(bytesFile, &config)
	if errorUnmarshal != nil {
		log.Println("Ошибка распаковки файла!")
		return config, errorUnmarshal
	}
	log.Println("Загружен файл конфигурации:")
	fmt.Printf("%+v\n", config)
	return config, nil
}

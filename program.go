package main

import (
	"flag"
	"log"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"regexp"
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
    
	config := Config{}
	errorUnmarshal := yaml.Unmarshal(bytesFile, &config)
	if errorUnmarshal != nil {
    	log.Fatalf("error: %v", errorUnmarshal)
    }
    log.Println("Загружен файл конфигурации:")
    fmt.Printf("%+v\n", config)
    return config
}

func ConfigValidation (config Config) bool {

    matched, errMatch := regexp.MatchString(`\S\.go`, config.Repository.File)
    if errMatch != nil {
        log.Println("error: %v", errMatch)
        return false
    }
    if matched == false {
        log.Println("Файл должен иметь расширение .go!")
        return false
    }

    if(config.Time.Hour > 60 || config.Time.Hour < 0) {
        log.Println("Часы вне диапазона!")
    	return false
    }

    if(config.Time.Minute > 60 || config.Time.Minute < 0) {
        log.Println("Минуты вне диапазона!")
        return false
    }

    if(config.Time.Second > 60 || config.Time.Second < 0) {
        log.Println("Секунды вне диапазона!")
        return false
    }
    log.Println("Валидация прошла успешно.")
    return true
}

func main() {
	log.Println("Старт программы.")

    configurationPath := flag.String("path", "", "Путь до файла конфигурации.")
    flag.Parse()

    config := LoadConfig(*configurationPath)

    if ConfigValidation(config) == false {
    	log.Fatal("Валидация закончилась неудачей!")
    }
    
	log.Println("Программа закончила работу.")
}


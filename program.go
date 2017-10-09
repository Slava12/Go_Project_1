package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/Slava12/Go_Project_0/validation"
	"github.com/Slava12/Go_Project_0/config"
)

func main() {
	log.Println("Старт программы.")

    configurationPath := flag.String("path", "", "Путь до файла конфигурации.") // -path /home/svyatoslav/goProjects/Go_Project_0/src/configuration.yaml
    flag.Parse()

    config1 := config.LoadConfig(*configurationPath)
    
    log.Println("Загружен файл конфигурации:")
    fmt.Printf("%+v\n", config1)

    validation.ConfigValidation(config1)
    
	log.Println("Программа закончила работу.")
}
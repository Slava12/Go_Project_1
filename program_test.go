package main

import(
	"testing"
)

func TestConfigValidation(t *testing.T){
	
	config1 := Config{}

	config1.Time.Hour = 18
	config1.Repository.File = "lol.go"
	if ConfigValidation(config1) == false {
		t.Error("Valid Error")
	}
	config1.Time.Hour = 78
	config1.Repository.File = "lol.go"
	if ConfigValidation(config1) == true {
		t.Error("Valid Error")
	}
	config1.Time.Hour = 18
	config1.Repository.File = "lol.go"
	if ConfigValidation(config1) == false {
		t.Error("Valid Error")
	}
	config1.Time.Hour = 18
	config1.Repository.File = "lol.g"
	if ConfigValidation(config1) == true {
		t.Error("Valid Error")
	}
}

func TestLoadConfig(t *testing.T) {
	configurationPath := "/home/svyatoslav/goProjects/src/github.com/Slava12/Go_Project_1/configuration.yaml"
	config1 := Config{}
	config, _ := LoadConfig(configurationPath)
	if config == config1 {
		t.Error("Valid Error")
	}
	configurationPath = "/home/svyatoslav/goProjects/src/github.com/Slava12/Go_Project_1/configuration.yam"
	config, _ = LoadConfig(configurationPath)
	if config != config1 {
		t.Error("Valid Error")
	}
}
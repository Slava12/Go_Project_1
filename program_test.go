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
package test_lol

import(
	//"log"
	//"fmt"
	"testing"
	"github.com/Slava12/Go_Project_1/config"
	"github.com/Slava12/Go_Project_1/validation"
)

/*func TestLoadConfig(t *testing.T){
	configurationPath := "/home/svyatoslav/goProjects/Go_Project_0/src/testConfiguration.yaml";
	config0 := config.LoadConfig(configurationPath)
	log.Println("Загружен файл конфигурации:")
    fmt.Printf("%+v\n", config0)
}*/

func TestValidationConfig(t *testing.T){
	//configurationPath := "/home/svyatoslav/goProjects/Go_Project_0/src/testConfiguration.yaml";
	//config0 := config.LoadConfig(configurationPath)
	config1 := config.Config{}
	//config1.Repository.File = "lol.go"
	config1.Time.Hour = 18
	if validation.ConfigValidation(config1) == false {
		t.Error("Valid Error")
	}
	//config1.Repository.File = "lol.g"
	config1.Time.Hour = 78
	if validation.ConfigValidation(config1) == true {
		t.Error("Valid Error")
	}
}
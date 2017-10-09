package validation

import (
	//"regexp"
	"github.com/Slava12/Go_Project_0/config"
	"log"
)

func ConfigValidation (config1 config.Config) bool {
	//fmt.Printf("%+v\n", config1.Repository.File)
    /*matched, errMatch := regexp.MatchString(`\S\.go`, config1.Repository.File)
    if errMatch != nil {
        log.Fatalf("error: %v", errMatch)
        return false
    }
    if matched == false {
        log.Fatal("Файл должен иметь расширение .go!")
        return false
    }*/

    if(config1.Time.Hour > 60) {
        log.Fatalf("error")
    	return false
    }

    return true
	//fmt.Println(matched, err)
}
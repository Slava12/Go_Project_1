package main

import (
	"flag"
	"fmt"
	"github.com/Slava12/Go_Project_1/loadobj"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	//"encoding/json"
)

var tpl *template.Template
var session *mgo.Session

type Config struct {
	Port    string `yaml:"port"`
	Mongodb string `yaml:"mongodb"`
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

func InitMongo(url string) *mgo.Session {
	session1, errorDial := mgo.Dial(url)
	if errorDial != nil {
		log.Fatal("Не удалось установить соединение с базой данных!")
	}
	log.Println("Соединение с базой данных установлено.")

	session1.SetMode(mgo.Monotonic, true)
	return session1
}

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := []loadobj.ModelRecord{}
		session1 := session.Copy()
		defer session1.Close()
		c := session1.DB("test").C("records")
		//err := c.Find(bson.M{"filename": "Soft_chair_OBJ.obj"}).One(&result)
		err := c.Find(bson.M{}).All(&result)
		if err != nil {
			log.Fatal(err)
		}
		var stringOut string
		for i := 0; i < len(result); i++ {
			stringOut += `<tr><td>` + result[i].Name + `</td>
        	<td>` + strconv.Itoa(result[i].Vertices) + `</td>
        	<td>` + strconv.Itoa(result[i].Normals) + `</td>
        	<td>` + strconv.Itoa(result[i].Textures) + `</td>
        	<td>` + strconv.Itoa(result[i].Faces) + `</td>
        	<td>` + result[i].FileName + `</td>
        	<td>` + strconv.FormatInt(result[i].FileSize, 10) + `</td>
        	<td><button>Delete</button></td></tr>`
		}
		err = tpl.ExecuteTemplate(w, "records.gohtml", template.HTML(stringOut))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		_, fileHeader, err := r.FormFile("f")

		if err != nil {
			fmt.Println(err.Error())
		}
		extension := filepath.Ext(fileHeader.Filename)
		println(extension)
		if extension == ".obj" {
			file, err := fileHeader.Open()

		bytesOfFile, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Файл не был прочитан!")
		}

		fileInServer, err := os.Create(fileHeader.Filename)
		if err != nil {
			log.Println("Файл не был создан!")
		}
		log.Println("Файл был создан.")

		_, err = fileInServer.Write(bytesOfFile)
		if err != nil {
			log.Println("Запись файла не удалась!")
		}
		log.Println("Файл был записан.")

		fileInServer.Close()
		log.Println("Файл был закрыт.")

		err = InsertIntoDB(fileHeader.Filename, session)
		if err != nil {
			log.Println("Вставка записи в базу данных окончилась неудачей!")
		}

		http.Redirect(w, r, "/records", 302)
		} else {
			html := `<html><p>Не obj-файл</p><a href="/records">Вернуться к списку записей</a></html>`
			w.Write([]byte(html))
		}
	}
}

func deleteRecords(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modelName := r.FormValue("name")
		session1 := session.Copy()
		defer session1.Close()
		c := session1.DB("test").C("records")
		err := c.Remove(bson.M{"isbn": &modelName})
		if err != nil {
			log.Fatal(err)
		}

		err = tpl.ExecuteTemplate(w, "delete.gohtml", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", 302)
	}
}

func insertRecord(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "insert.gohtml", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		_, fh, err := r.FormFile("f")

		if err != nil {
			fmt.Println(err.Error())
		}
		log.Println("Файлик пришел")
		log.Println(fh.Filename)
	}

}

func main() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

	log.Println("Старт сервера.")
	config, errorLoadConfig := LoadConfigFile()
	if errorLoadConfig != nil {
		log.Fatalf("error: %v", errorLoadConfig)
	}
	port := config.Port
	session = InitMongo(config.Mongodb)
	defer session.Close()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/records", getAllRecords)
	//http.HandleFunc("/insert", insertRecord)
	//http.HandleFunc(pat.Delete("/delete"), deleteRecords(session))

	http.ListenAndServe(":"+port, nil)
}

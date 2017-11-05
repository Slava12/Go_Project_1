package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
	//"os"
	"path/filepath"
	"strconv"
)

var tpl *template.Template
var tempFolder string

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		result := GetAllRecords(session)
		var stringOut string
		for i := 0; i < len(result); i++ {
			stringOut += `<tr><td>` + result[i].Name + `</td>
        	<td>` + strconv.Itoa(result[i].Vertices) + `</td>
        	<td>` + strconv.Itoa(result[i].Normals) + `</td>
        	<td>` + strconv.Itoa(result[i].Textures) + `</td>
        	<td>` + strconv.Itoa(result[i].Faces) + `</td>
        	<td>` + result[i].FileName + `</td>
        	<td>` + strconv.FormatInt(result[i].FileSize, 10) + `</td>`
			//<td><button>Delete</button></td></tr>`
		}
		err := tpl.ExecuteTemplate(w, "records.gohtml", template.HTML(stringOut))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		_, fileHeader, err := r.FormFile("f")
		if r.FormValue("c") != "" && r.FormValue("n") == "" && err != nil {
			RemoveAllRecords(session)
			http.Redirect(w, r, "/records", 302)
		} else if r.FormValue("c") == "" && r.FormValue("n") != "" && err != nil {
			RemoveOneRecord("name", r.FormValue("n"), session)
			http.Redirect(w, r, "/records", 302)
		} else if r.FormValue("c") == "" && r.FormValue("n") == "" && err == nil {
			extension := filepath.Ext(fileHeader.Filename)
			if extension == ".obj" {
				SaveFile(tempFolder, fileHeader)
				objModel, errorLoadObjFile := LoadObjFileInfo(tempFolder + fileHeader.Filename)
				if errorLoadObjFile != nil {
					log.Println("Не удалось получить информацию о модели из файла!")
				}
				log.Println("Получена информация о модели", objModel.Name)
				if objModel.Name != "" {
					internalName := FindOneResult("name", objModel.Name, session)
					if internalName != "" {

					} else {
						err = InsertIntoDB(tempFolder+fileHeader.Filename, session)
						if err != nil {
							log.Println("Вставка записи в базу данных окончилась неудачей!")
						}
					}
					http.Redirect(w, r, "/records", 302)
				} else {
					log.Println("Файл", fileHeader.Filename, "имеет пустое имя модели!")
					//html := `<html><head><title>LOL</title></head><body><div>Empty name!</div></body></html>`
					html := `<html><head><title>LOL</title></head><body><div>Empty name!</div><div><a href="/records">Вернуться к списку записей</a></div></body></html>`
					w.Write([]byte(html))
				}
				RemoveFile(tempFolder + fileHeader.Filename)
			} else {
				err := tpl.ExecuteTemplate(w, "wrongFormat.gohtml", nil)
				if err != nil {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}
	}
}

/*func deleteRecords(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
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
}*/

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

	tempFolder = config.Tempfolder

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/records", getAllRecords)

	http.ListenAndServe(":"+port, nil)
}

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
var folder string

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
		}
		err := tpl.ExecuteTemplate(w, "records.gohtml", template.HTML(stringOut))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		_, fileHeader, err := r.FormFile("f")
		if r.FormValue("c") != "" && r.FormValue("n") == "" && err != nil {
			RemoveAllRecords(session)
			RemoveAllFiles(folder)
			CreateDirectory(folder)
			http.Redirect(w, r, "/records", 302)
		} else if r.FormValue("c") == "" && r.FormValue("n") != "" && err != nil {
			removedObjModel := FindOneResult("name", r.FormValue("n"), session)
			RemoveFile(folder + removedObjModel.FileName)
			RemoveOneRecord("name", r.FormValue("n"), session)
			http.Redirect(w, r, "/records", 302)
		} else if r.FormValue("c") == "" && r.FormValue("n") == "" && err == nil {
			category := r.FormValue("category")
			log.Println("Категория:", category)
			subcategory := r.FormValue("subcategory")
			if category == "other" {
				subcategory = "other"
			}
			log.Println("Подкатегория:", subcategory)
			if category != "" && subcategory != "" {
				extension := filepath.Ext(fileHeader.Filename)
				if extension == ".obj" {
					SaveFile(tempFolder, fileHeader)
					objModel, errorLoadObjFile := LoadObjFileInfo(tempFolder + fileHeader.Filename)
					if errorLoadObjFile != nil {
						log.Println("Не удалось получить информацию о модели из файла!")
					}
					log.Println("Получена информация о модели", objModel.Name)
					if objModel.Name != "" {
						internalObjModel := FindOneResult("name", objModel.Name, session)
						internalName := internalObjModel.Name
						if internalName != "" {

						} else {
							err = InsertIntoDB(tempFolder+fileHeader.Filename, session)
							if err != nil {
								log.Println("Вставка записи в базу данных окончилась неудачей!")
							}
							SaveFile(folder, fileHeader)
						}
						http.Redirect(w, r, "/records", 302)
					} else {
						log.Println("Файл", fileHeader.Filename, "имеет пустое имя модели!")
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
			} else {
				log.Println("Категория или подкатегория пусты!")
				html := `<html><head><title>Empty kategory</title></head><body><div>Empty kategory or subcategory!</div><div><a href="/records">Вернуться к списку записей</a></div></body></html>`
				w.Write([]byte(html))
			}
		}
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

	tempFolder = config.Tempfolder
	folder = config.Folder
	CreateDirectory(folder)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/records", getAllRecords)

	http.ListenAndServe(":"+port, nil)
}

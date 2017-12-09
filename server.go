package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	tpl        *template.Template
	tempFolder string
	folder     string
	key        = []byte("super-secret-key")
	store      = sessions.NewCookieStore(key)
)

func index(w http.ResponseWriter, r *http.Request) {
	sessionBrowser, _ := store.Get(r, "cookie-name")

	if sessionBrowser.Values["texturesCheck"] == nil {
		sessionBrowser.Values["texturesCheck"] = 0
	}
	if sessionBrowser.Values["normalsCheck"] == nil {
		sessionBrowser.Values["normalsCheck"] = 0
	}
	if r.Method == "GET" {
		showTitle(w, r, sessionBrowser.Values["texturesCheck"].(int), sessionBrowser.Values["normalsCheck"].(int))
		showFilterRecords(w, r)
	}
}

func showTitle(w http.ResponseWriter, r *http.Request, texturesCheck int, normalsCheck int) {
	var one string
	var two string
	if texturesCheck != 0 {
		log.Println("texturesCheck = 1")
		one = "checked"
	} else {
		log.Println("texturesCheck = 0")
		one = ""
	}
	if normalsCheck != 0 {
		log.Println("normalsCheck = 1")
		two = "checked"
	} else {
		log.Println("normalsCheck = 0")
		two = ""
	}
	userAgent := r.UserAgent()
	log.Println("user agent:", userAgent)
	data := struct {
		One string
		Two string
	}{
		One: one,
		Two: two,
	}
	err := tpl.ExecuteTemplate(w, "header", data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func showFilterRecords(w http.ResponseWriter, r *http.Request) {
	result := GetAllRecords(session)
	var str string
	for i := 0; i < len(result); i++ {
		str += result[i].Name + " "
	}
	err := tpl.ExecuteTemplate(w, "index", str)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func setFilters(w http.ResponseWriter, r *http.Request) {
	sessionBrowser, _ := store.Get(r, "cookie-name")
	if r.FormValue("withTextures") != "" {
		sessionBrowser.Values["texturesCheck"] = 1
	} else {
		sessionBrowser.Values["texturesCheck"] = 0
	}

	if r.FormValue("withNormals") != "" {
		sessionBrowser.Values["normalsCheck"] = 1
	} else {
		sessionBrowser.Values["normalsCheck"] = 0
	}
	sessionBrowser.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func GetRecords(w http.ResponseWriter, r *http.Request) {
	result := GetAllRecords(session)

	var stringOut string
	// TODO: templates {{range pipeline}} T1 {{end}}
	// https://golang.org/pkg/text/template/#hdr-Actions
	for i := 0; i < len(result); i++ {
		stringOut += `<tr><td>` + result[i].Name + `</td>
        	<td>` + strconv.Itoa(result[i].Vertices) + `</td>
        	<td>` + strconv.Itoa(result[i].Normals) + `</td>
        	<td>` + strconv.Itoa(result[i].Textures) + `</td>
        	<td>` + strconv.Itoa(result[i].Faces) + `</td>
        	<td>` + result[i].FileName + `</td>
        	<td>` + strconv.FormatInt(result[i].FileSize, 10) + `</td>
        	<td>` + result[i].Category + `</td>
        	<td>` + result[i].Subcategory + `</td>`
	}

	err := tpl.ExecuteTemplate(w, "records.html", template.HTML(stringOut))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ClearRecords(w http.ResponseWriter, r *http.Request) {
	RemoveAllRecords(session)
	RemoveAllFiles(folder)
	CreateDirectory(folder)
	CreateDirectory(tempFolder)
	http.Redirect(w, r, "/records", 302)
}

func RemoveRecords(w http.ResponseWriter, r *http.Request) {
	log.Println("Попытка удаления записи:", r.FormValue("fileName"))
	removedObjModel := FindOneResult("name", r.FormValue("fileName"), session)
	if removedObjModel.Name != "" {
		RemoveFile(folder + removedObjModel.FileName)
		RemoveOneRecord("name", removedObjModel.Name, session)
		log.Println("Удаление записи окончилось успехом.")
	} else {
		log.Println("Удаление записи окончилось неудачей.")
	}
	http.Redirect(w, r, "/records", 302)
}

func ChangeRecords(w http.ResponseWriter, r *http.Request) {
	// Доп.флаг c означает очистку
	// Лучше было бы обработать это отдельным методом /records/clear/
	if r.FormValue("clear") != "" {
		ClearRecords(w, r)
		return
	}

	// Доп. флаг n видимо означает фильтр по имени для удаления (почему??)
	// Кроме того, по REST нужно, чтобы удаление шло методом DELETE
	if r.FormValue("fileName") != "" {
		RemoveRecords(w, r)
		return
	}

	// Остальное - по-видимому добавление
	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Ошибка получения файла: %s", err)
		// В API хорошим тоном считается отдавать корректные коды ошибок
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category := r.FormValue("category")
	log.Println("Категория:", category)
	subcategory := r.FormValue("subcategory")
	if category == "other" {
		subcategory = "other"
	}
	log.Println("Подкатегория:", subcategory)

	if category == "" && subcategory == "" {
		log.Println("Категория или подкатегория пусты!")
		html := `<html><head><title>Empty kategory</title></head><body><div>Empty kategory or subcategory!</div><div><a href="/records">Вернуться к списку записей</a></div></body></html>`
		w.Write([]byte(html))
		return
	}

	extension := filepath.Ext(fileHeader.Filename)
	// При сравнении строк нельзя забывать про регистр
	if strings.ToLower(extension) != ".obj" {
		err := tpl.ExecuteTemplate(w, "wrongFormat.gohtml", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Очистку лучше делать через defer, чтобы она отработала в любом случае
	SaveFile(tempFolder, fileHeader, "")
	defer RemoveFile(tempFolder + fileHeader.Filename)

	objModel, errorLoadObjFile := LoadObjFileInfo(tempFolder + fileHeader.Filename)
	if errorLoadObjFile != nil {
		log.Println("Не удалось получить информацию о модели из файла!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Получена информация о модели", objModel.Name)
	if objModel.Name == "" {
		log.Println("Файл", fileHeader.Filename, "имеет пустое имя модели!")
		html := `<html><head><title>LOL</title></head><body><div>Empty name!</div><div><a href="/records">Вернуться к списку записей</a></div></body></html>`
		w.Write([]byte(html))
		return
	}

	internalObjModel := FindOneResult("name", objModel.Name, session)
	internalName := internalObjModel.Name
	if internalName != "" {
		internalObjModel.Vertices = objModel.Vertices
		internalObjModel.Normals = objModel.Normals
		internalObjModel.Textures = objModel.Textures
		internalObjModel.Faces = objModel.Faces
		internalObjModel.FileSize = objModel.FileSize
		internalObjModel.FileModTime = objModel.FileModTime
		idStr := internalObjModel.Id.String()
		arrStrings := strings.Split(idStr, "\"")
		SaveFile(folder, fileHeader, arrStrings[1]+".obj")
		UpdateOneRecord("name", internalObjModel.Name, internalObjModel, session)
	} else {
		err = InsertModelRecordIntoDB(objModel, session)
		if err != nil {
			log.Println("Вставка записи в базу данных окончилась неудачей!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		objModel = FindOneResult("name", objModel.Name, session)
		idStr := objModel.Id.String()
		arrStrings := strings.Split(idStr, "\"")
		SaveFile(folder, fileHeader, arrStrings[1]+".obj")
		objModel.Category = category
		objModel.Subcategory = subcategory
		objModel.FileName = arrStrings[1] + ".obj"
		UpdateOneRecord("name", objModel.Name, objModel, session)
	}
	http.Redirect(w, r, "/records", 302)
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetRecords(w, r)
		return
	}

	// Добавил, хотя этого не было
	if r.Method == "DELETE" {
		RemoveRecords(w, r)
		return
	}

	// PUT / POST / PATCH
	ChangeRecords(w, r)
}

func main() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

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
	CreateDirectory(tempFolder)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/", index)
	http.HandleFunc("/set_filters", setFilters)
	http.HandleFunc("/records", getAllRecords)

	http.ListenAndServe(":"+port, nil)
}

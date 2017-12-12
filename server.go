package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	tpl        *template.Template
	tempFolder string
	folder     string
	key        = []byte("super-secret-key")
	store      = sessions.NewCookieStore(key)
	path       = "/index"
	fileName   = "111"
)

func CreatePaths() ([]string, []string) {
	result := GetAllRecords(session)
	modelPaths := make([]string, len(result))
	modelFilePaths := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		modelPaths[i] = "/" + result[i].Category + "/" + result[i].Subcategory + "/" + strings.ToLower(result[i].Name)
		modelFilePaths[i] = "/" + result[i].Category + "/" + result[i].Subcategory + "/" + result[i].Name + ".obj"
	}
	return modelPaths, modelFilePaths
}

func index(w http.ResponseWriter, r *http.Request) {
	sessionBrowser, _ := store.Get(r, "cookie-name")

	if sessionBrowser.Values["texturesCheck"] == nil {
		sessionBrowser.Values["texturesCheck"] = 0
	}
	if sessionBrowser.Values["normalsCheck"] == nil {
		sessionBrowser.Values["normalsCheck"] = 0
	}
	if sessionBrowser.Values["sortBySize"] == nil {
		sessionBrowser.Values["sortBySize"] = 0
	}
	if sessionBrowser.Values["searchByName"] == nil {
		sessionBrowser.Values["searchByName"] = ""
	}
	searchByName := sessionBrowser.Values["searchByName"].(string)
	texturesCheck := sessionBrowser.Values["texturesCheck"].(int)
	normalsCheck := sessionBrowser.Values["normalsCheck"].(int)
	sortBySize := sessionBrowser.Values["sortBySize"].(int)
	category := ""
	subcategory := ""
	log.Println("URL:", r.URL.String())
	//log.Println("RequestURI:", r.RequestURI)
	if r.URL.String() == "/index" {
	} else {
		arrStrings := strings.Split(r.URL.String(), "/")
		if len(arrStrings) == 2 {
			category = arrStrings[1]
		} else {
			subcategory = arrStrings[2]
		}
	}
	path = r.URL.String()
	log.Println("category:", category)
	log.Println("subcategory:", subcategory)
	result := filterRecords(searchByName, texturesCheck, normalsCheck, sortBySize, category, subcategory)
	if r.Method == "GET" {
		showTitle(w, r, searchByName, texturesCheck, normalsCheck, sortBySize, category, subcategory)
		showFilterRecords(w, r, result)
	}
}

func filterRecords(searchByName string, texturesCheck int, normalsCheck int, sortBySize int, category string, subcategory string) []ModelRecord {
	var modelRecords []ModelRecord
	if sortBySize == 0 {
		modelRecords = GetAllRecords(session)
	} else if sortBySize == 1 {
		modelRecords = GetSortedRecords("filesize", "", session)
	} else if sortBySize == -1 {
		modelRecords = GetSortedRecords("filesize", "-", session)
	}

	filteredModelRecords := make([]ModelRecord, len(modelRecords))
	newLength := 0
	for i := 0; i < len(modelRecords); i++ {
		if searchByName != "" {
			matched, errMatch := regexp.MatchString(`.*(?i:`+searchByName+`).*`, modelRecords[i].Name)
			if errMatch != nil {
				log.Println("error: ", errMatch)
				return nil
			}
			if matched == false {
				continue
			}
		}
		if texturesCheck != 0 {
			if modelRecords[i].Textures == 0 {
				continue
			}
		}
		if normalsCheck != 0 {
			if modelRecords[i].Normals == 0 {
				continue
			}
		}
		if category != "" {
			if category != modelRecords[i].Category {
				continue
			}
		}
		if subcategory != "" {
			if subcategory != modelRecords[i].Subcategory {
				continue
			}
		}
		filteredModelRecords[i] = modelRecords[i]
		newLength++
	}
	filteredModelRecordsShort := make([]ModelRecord, newLength)
	j := 0
	for i := 0; i < len(modelRecords); i++ {
		if filteredModelRecords[i].Name != "" {
			filteredModelRecordsShort[j] = filteredModelRecords[i]
			j++
		}
	}
	return filteredModelRecordsShort
}

func showTitle(w http.ResponseWriter, r *http.Request, searchByName string, texturesCheck int, normalsCheck int, sortBySize int, category string, subcategory string) {
	if searchByName != "" {
		log.Println("Поиск по слову:", searchByName)
	}
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
	var sortSize1 string
	var sortSize2 string
	var sortSize3 string
	if sortBySize == 0 {
		log.Println("Порядок по умолчанию.")
		sortSize1 = "selected"
		sortSize2 = ""
		sortSize3 = ""
	} else if sortBySize == 1 {
		log.Println("Сортировка по возрастанию размера.")
		sortSize1 = ""
		sortSize2 = "selected"
		sortSize3 = ""
	} else if sortBySize == -1 {
		log.Println("Сортировка по убыванию размера.")
		sortSize1 = ""
		sortSize2 = ""
		sortSize3 = "selected"
	}
	var title string
	if category != "" {
		title = category
	} else if subcategory != "" {
		title = subcategory
	} else {
		title = "Главная страница"
	}
	userAgent := r.UserAgent()
	log.Println("user agent:", userAgent)
	data := struct {
		SearchByName string
		One          string
		Two          string
		SortSize1    string
		SortSize2    string
		SortSize3    string
		Title         string
	}{
		SearchByName: searchByName,
		One:          one,
		Two:          two,
		SortSize1:    sortSize1,
		SortSize2:    sortSize2,
		SortSize3:    sortSize3,
		Title:        title,
	}
	err := tpl.ExecuteTemplate(w, "header.html", data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func showFilterRecords(w http.ResponseWriter, r *http.Request, result []ModelRecord) {
	type SiteData struct {
		Name        string
		Vertices    int
		Normals     int
		Textures    int
		Faces       int
		FileName    string
		FileSize    int64
		Category    string
		Subcategory string
		Link        string
	}
	data := make([]SiteData, len(result))
	for i:= 0; i < len(result); i++ {
		data[i].Name =  result[i].Name
		data[i].Vertices = result[i].Vertices
		data[i].Normals = result[i].Normals
		data[i].Textures = result[i].Textures
		data[i].Faces = result[i].Faces
		data[i].FileName = result[i].FileName
		data[i].FileSize = result[i].FileSize
		data[i].Category = result[i].Category
		data[i].Subcategory = result[i].Subcategory
		data[i].Link = "/" + result[i].Category + "/" + result[i].Subcategory + "/" + strings.ToLower(result[i].Name)
	}
	err := tpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func setFilters(w http.ResponseWriter, r *http.Request) {
	sessionBrowser, _ := store.Get(r, "cookie-name")

	sessionBrowser.Values["searchByName"] = r.FormValue("searchByName")

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

	if r.FormValue("sortBySize") == "noneSize" {
		sessionBrowser.Values["sortBySize"] = 0
	} else if r.FormValue("sortBySize") == "upSize" {
		sessionBrowser.Values["sortBySize"] = 1
	} else if r.FormValue("sortBySize") == "downSize" {
		sessionBrowser.Values["sortBySize"] = -1
	}

	sessionBrowser.Save(r, w)
	http.Redirect(w, r, path, 302)
}

func showModel(w http.ResponseWriter, r *http.Request) {
	log.Println("URL:",r.URL.String())
	arrStrings := strings.Split(r.URL.String(), "/")
	modelName := arrStrings[3]
	modelRecords := GetAllRecords(session)
	var result ModelRecord
	for i := 0; i < len(modelRecords); i++ {
		matched, errMatch := regexp.MatchString(`.*(?i:` + modelName + `).*`, modelRecords[i].Name)
		if errMatch != nil {
			log.Println("error: ", errMatch)
		}
		if matched != false {
			result = modelRecords[i]
			continue
		}
	}
	log.Println("name =", result.Name)
	fileName = result.FileName
	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "model.html", result)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func loadFileFromServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, folder + fileName)
}

//----------------------------------------------------------------------------//
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
	// Лучше было бы обработать это отдельным методом /records/clear/
	if r.FormValue("clear") != "" {
		ClearRecords(w, r)
		return
	}

	if r.FormValue("fileName") != "" {
		RemoveRecords(w, r)
		return
	}

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Ошибка получения файла: %s", err)
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
	if strings.ToLower(extension) != ".obj" {
		err := tpl.ExecuteTemplate(w, "wrongFormat.gohtml", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

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

	if r.Method == "DELETE" {
		RemoveRecords(w, r)
		return
	}

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

	http.HandleFunc("/index", index)
	http.HandleFunc("/set_filters", setFilters)
	http.HandleFunc("/furniture", index)
	http.HandleFunc("/furniture/chair", index)
	http.HandleFunc("/furniture/table", index)
	http.HandleFunc("/furniture/sofa", index)
	http.HandleFunc("/car", index)
	http.HandleFunc("/car/sportcar", index)
	http.HandleFunc("/car/truck", index)
	http.HandleFunc("/car/service", index)
	http.HandleFunc("/weapon", index)
	http.HandleFunc("/weapon/tank", index)
	http.HandleFunc("/weapon/helicopter", index)
	http.HandleFunc("/weapon/plane", index)
	http.HandleFunc("/other", index)

	modelPaths, modelFilePaths := CreatePaths()
	for i:= 0; i < len(modelPaths); i++ {
		http.HandleFunc(modelPaths[i], showModel)
		http.HandleFunc(modelFilePaths[i], loadFileFromServer)
	}

	http.HandleFunc("/records", getAllRecords)

	http.ListenAndServe(":"+port, nil)
}

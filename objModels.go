package main

import (
	"github.com/sheenobu/go-obj/obj"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"time"
)

type ModelRecord struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string
	Vertices    int
	Normals     int
	Textures    int
	Faces       int
	FileName    string
	FileSize    int64
	FileModTime time.Time
	Category    string
	Subcategory string
}

func LoadObjFileInfo(filename string) (ModelRecord, error) {
	modelRecord := ModelRecord{}
	if filename == "" {
		log.Println("Нет файла для загрузки!")
		return modelRecord, nil
	}
	fileObj, errorOpen := os.Open(filename)
	if errorOpen != nil {
		log.Println("Файл с моделью не был открыт!")
		return modelRecord, errorOpen
	}
	model, errorRead := obj.NewReader(fileObj).Read()
	if errorRead != nil {
		log.Println("Файл с моделью не был прочитан!")
		return modelRecord, errorRead
	}
	defer fileObj.Close()

	file, errorStat := os.Stat(filename)
	if errorStat != nil {
		log.Println("Информация о файле с моделью не была получена!")
		return modelRecord, errorStat
	}
	modelRecord.Name = model.Name
	modelRecord.Vertices = len(model.Vertices)
	modelRecord.Normals = len(model.Normals)
	modelRecord.Textures = len(model.Textures)
	modelRecord.Faces = len(model.Faces)
	modelRecord.FileName = file.Name()
	modelRecord.FileSize = file.Size()
	modelRecord.FileModTime = file.ModTime()
	return modelRecord, nil
}

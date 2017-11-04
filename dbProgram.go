package main

import (
	//"fmt"
	//"os"
	"gopkg.in/mgo.v2"
	"log"
	//"gopkg.in/mgo.v2/bson"
	"github.com/Slava12/Go_Project_1/loadobj"
)

func InsertIntoDB(filename string, session *mgo.Session) error {
	modelRecord, errorLoadObjFile := loadobj.LoadObjFileInfo(filename)
	if errorLoadObjFile != nil {
		log.Println("Не удалось загрузить файл!")
		return errorLoadObjFile
	}

	log.Println("Загружены данные о модели:", modelRecord.Name)

	session1 := session.Copy()
	defer session1.Close()
	c := session1.DB("test").C("records")
	errorInsert := c.Insert(&modelRecord)
	if errorInsert != nil {
		log.Println("Невозможно произвести вставку элемента в таблицу - неправильный формат!")
		return errorInsert
	}

	log.Println("Данные о модели", modelRecord.Name, "были добавлены в базу данных.")
	return nil
}

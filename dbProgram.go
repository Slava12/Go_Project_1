package main

import (
	//"fmt"
	//"os"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var session *mgo.Session

func InitMongo(url string) *mgo.Session {
	session1, errorDial := mgo.Dial(url)
	if errorDial != nil {
		log.Fatal("Не удалось установить соединение с базой данных!")
	}
	log.Println("Соединение с базой данных установлено.")

	session1.SetMode(mgo.Monotonic, true)
	return session1
}

func GetAllRecords(session *mgo.Session) []ModelRecord {
	result := []ModelRecord{}
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	err := collection.Find(bson.M{}).All(&result)
	if err != nil {
		log.Println("Не удалось осуществить поиск по базе данных!")
		return nil
	}
	return result
}

func InsertIntoDB(filename string, session *mgo.Session) error {
	modelRecord, errorLoadObjFile := LoadObjFileInfo(filename)
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

func FindOneResult(field string, value string, session *mgo.Session) ModelRecord {
	result := ModelRecord{}
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	err := collection.Find(bson.M{field: value}).One(&result)
	if err != nil {
		log.Println("Нет совпадений в базе данных!")
		return result
	}
	log.Println("Есть совпадение в базе данных по полю:", field, ", значение:", value)
	return result
}

func RemoveAllRecords(session *mgo.Session) {
	internalSession := session.Copy()
	defer internalSession.Close()
	err := internalSession.DB("test").DropDatabase()
	if err != nil {
		log.Println("Не удалось удалить все записи из базы данных!")
		return
	}
	log.Println("Все записи были удалены из базы данных.")
}

func RemoveOneRecord(field string, value string, session *mgo.Session) {
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	err := collection.Remove(bson.M{field: value})
	if err != nil {
		log.Println("Не удалось удалить запись из базы данных!")
		return
	}
	log.Println("Запись", value, "была удалена из базы данных.")
}

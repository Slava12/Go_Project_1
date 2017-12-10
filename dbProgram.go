package main

import (
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

func InsertModelRecordIntoDB(modelRecord ModelRecord, session *mgo.Session) error {
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	errorInsert := collection.Insert(&modelRecord)
	if errorInsert != nil {
		log.Println("Невозможно произвести вставку элемента в таблицу - неправильный формат!")
		return errorInsert
	}

	log.Println("Данные о модели", modelRecord.Name, "были добавлены в базу данных.")
	return nil
}

func InsertFileIntoDB(filename string, session *mgo.Session) error {
	modelRecord, errorLoadObjFile := LoadObjFileInfo(filename)
	if errorLoadObjFile != nil {
		log.Println("Не удалось загрузить файл!")
		return errorLoadObjFile
	}

	log.Println("Загружены данные о модели:", modelRecord.Name)

	err := InsertModelRecordIntoDB(modelRecord, session)
	if err != nil {
		return err
	}
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

func UpdateOneRecord(field string, value string, modelRecord ModelRecord, session *mgo.Session) {
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	err := collection.Update(bson.M{field: value}, &modelRecord)
	if err != nil {
		log.Println("Не удалось обновить запись в базе данных!")
		return
	}
	log.Println("Запись", value, "была обновлена в базе данных.")
}

func GetSortedRecords(field string, order string, session *mgo.Session) []ModelRecord {
	result := []ModelRecord{}
	internalSession := session.Copy()
	defer internalSession.Close()
	collection := internalSession.DB("test").C("records")
	err := collection.Find(nil).Sort(order + field).All(&result)
	if err != nil {
		log.Println("Не удалось отсортировать записи по полю", field, "в базе данных!")
		return nil
	}
	log.Println("Записи в базе данных были отсортированы по полю:", field)
	return result
}

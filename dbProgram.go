package main

import (
	//"fmt"
	//"os"
	"gopkg.in/mgo.v2"
	"log"
	//"gopkg.in/mgo.v2/bson"
	"github.com/Slava12/Go_Project_1/loadobj"
)

func main() {
	session, errorDial := mgo.Dial("localhost")
	if errorDial != nil {
		log.Fatal("Не удалось установить соединение с базой данных!")
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("records")
	modelRecord, errorLoadObjFile := loadobj.LoadObjFileInfo()
	if errorLoadObjFile != nil {
		log.Fatal("Не удалось загрузить файл!")
	}

	errorInsert := c.Insert(&modelRecord)
	if errorInsert != nil {
		log.Fatal("Невозможно произвести вставку элемента в таблицу - неправильный формат!")
	}

	log.Println("Модель", modelRecord.Name, "была добавлена.")

	/*result := loadobj.ModelRecord{}
	  err = c.Find(bson.M{"filename": "untitled.obj"}).One(&result)
	  if err != nil {
	          log.Fatal(err)
	  }

	  fmt.Println("Vertices:", result.Vertices)*/
}

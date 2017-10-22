package main

import (
	"fmt"
	//"os"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ModelRecord struct {
	Vertices string
	Normals string
	Textures string
	Faces string
	FileName string
	FileSize string
	FileModTime string
}

func main() {
	session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

        // Optional. Switch the session to a monotonic behavior.
        session.SetMode(mgo.Monotonic, true)

        c := session.DB("test").C("records")

	    err = c.Insert(&ModelRecord{"12", "13", "14", "15", "lol", "12", "67"})
        if err != nil {
                log.Fatal(err)
        }

        result := ModelRecord{}
        err = c.Find(bson.M{"filename": "lol"}).One(&result)
        if err != nil {
                log.Fatal(err)
        }

        fmt.Println("Vertices:", result.Vertices)
}

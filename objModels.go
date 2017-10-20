package main

import (
	"fmt"
	"github.com/sheenobu/go-obj/obj"
	"os"
)

func main() {
	filename := "untitled.obj"
	if len(os.Args) == 1 {
		fmt.Printf("No argument, using untitled.obj\n")
	} else {
		filename = os.Args[1]
	}
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	model, err := obj.NewReader(f).Read()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println("Name:", model.Name)
	vertices := len(model.Vertices)
	fmt.Println("Vertices:", vertices)
	normals := len(model.Normals)
	fmt.Println("Normals:", normals)
	textures := len(model.Textures)
	fmt.Println("Textures:", textures)
	faces := len(model.Faces)
	fmt.Println("Faces:", faces)
	file, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("File Name:", file.Name())
	fmt.Println("File Size:", file.Size())
	fmt.Println("File Mode:", file.Mode())
	fmt.Println("File ModTime:", file.ModTime())
	fmt.Println("File Sys:", file.Sys())
}

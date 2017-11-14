package main

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

func SaveFile(filePath string, fileHeader *multipart.FileHeader, fileName string) {
	file, err := fileHeader.Open()

	bytesOfFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Файл", fileHeader.Filename, "не был прочитан!")
		return
	}
	fullFilePath := ""
	if fileName != "" {
		fullFilePath = filePath + fileName
	} else {
		fullFilePath = filePath + fileHeader.Filename
	}

	fileInServer, err := os.Create(fullFilePath)
	if err != nil {
		log.Println("Файл", fullFilePath, "не был создан!")
		return
	}
	log.Println("Файл", fullFilePath, "был создан.")

	_, err = fileInServer.Write(bytesOfFile)
	if err != nil {
		log.Println("Запись файла", fullFilePath, "не удалась!")
		return
	}
	log.Println("Файл", fullFilePath, "был записан.")

	fileInServer.Close()
	log.Println("Файл", fullFilePath, "был закрыт.")
}

func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Println("Файл", filePath, "не был удалён!")
	} else {
		log.Println("Файл", filePath, "был удалён.")
	}
}

func RemoveAllFiles(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println("В директории", path, "не удалось произвести удаление всех файлов!")
	} else {
		log.Println("В директории", path, "были удалены все файлы.")
	}
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("Указанный путь", path, " не существует!")
		return false
	}
	log.Println("Указанный путь", path, " существует.")
	return true
}

func CreateDirectory(path string) {
	exist := Exists(path)
	if exist == false {
		err := os.Mkdir(path, 0777)
		if err != nil {
			log.Println("Директория", path, " не была создана!")
		} else {
			log.Println("Директория", path, " была создана.")
		}
	}
}

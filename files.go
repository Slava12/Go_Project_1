package main

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

func SaveFile(filePath string, fileHeader *multipart.FileHeader) {
	file, err := fileHeader.Open()

	bytesOfFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Файл", fileHeader.Filename, "не был прочитан!")
	}
	fullFilePath := filePath + fileHeader.Filename
	fileInServer, err := os.Create(fullFilePath)
	if err != nil {
		log.Println("Файл", fileHeader.Filename, "не был создан!")
	}
	log.Println("Файл", fileHeader.Filename, "был создан.")

	_, err = fileInServer.Write(bytesOfFile)
	if err != nil {
		log.Println("Запись файла", fileHeader.Filename, "не удалась!")
	}
	log.Println("Файл", fileHeader.Filename, "был записан.")

	fileInServer.Close()
	log.Println("Файл", fileHeader.Filename, "был закрыт.")
}

func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		log.Println("Файл", filePath, "не был удалён!")
	}
	log.Println("Файл", filePath, "был удалён.")
}

/*func CreateDirectory (path string) {
	err := os.Mkdir(path, 0777)
	if os.IsExist(err) {
        err = nil // then nullify the error
    }
	if err != nil {
		log.Println("Директория не была создана!")
	}
	log.Println("Директория была создана.")
}*/

package bucket

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func Open(fileori multipart.File, file *multipart.FileHeader, bucketName, fileFolder, file_name string) (url string, err error) {

	bytes, err := ioutil.ReadAll(fileori)
	if err != nil {
		fmt.Println(err)
	}
	dir, _ := os.Getwd()
	dir = dir + "/temp/" + fileFolder
	os.Mkdir(dir, 0777)

	filePath := path.Join(dir, file_name)

	err = ioutil.WriteFile(filePath, bytes, 0666)
	if err != nil {
		fmt.Println("Error")
	}
	ext := filepath.Ext(file.Filename)
	content_type := file.Header.Get("Content-Type")

	file_name = fmt.Sprintf("%s/%s%s", fileFolder, file_name, ext)

	fmt.Println("sampai mau aplod")
	fileUrl, err := UploadFile(file_name, filePath, bucketName, content_type)
	if err != nil {
		return "", fmt.Errorf("Error")
	} else {
		fileori.Close()

		err := os.Remove(filePath)
		if err != nil {
			return "", err
		}
	}
	return fileUrl, nil
}

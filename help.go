package main

import (
	"fmt"
	"io/ioutil"
)

//GetTextFile Получить сообщение из файла
func GetTextFile(fileName string) string {
	content, err := ioutil.ReadFile("messages/" + fileName + ".txt")
	if err != nil {
		fmt.Println("Проблема с вытаскиванием сообщения из файла", fileName, err)
	}

	return string(content)
}

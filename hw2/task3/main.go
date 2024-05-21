package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.WriteFile("some.txt", []byte(""), 0444)
	if err != nil {
		fmt.Println("Ошибка при создании файла и записи:", err)
		return
	}

	file, err := os.OpenFile("some.txt", os.O_RDONLY, 0444)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("Новая строка")
	if err != nil {
		fmt.Println("Ошибка записи:", err)
	} else {
		fmt.Println("Данные записаны")
	}
}

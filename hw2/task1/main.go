package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	inputText := ""
	date := time.Now().Format("2006-01-02 15:04:05")

	fileW, err := os.Create("some.txt")
	if err != nil {
		fmt.Println("Невозможно создать файл", err)
		os.Exit(1)
	}
	defer fileW.Close()

	for inputText != "exit" {
		fmt.Println("Для выхода введите 'exit'\nВедите ваше предложение: ")
		fmt.Fscan(os.Stdin, &inputText)
		if inputText == "exit" {
			break
		}
		fileW.WriteString(date + " " + inputText + "\n")
	}
}

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fileR, err := os.Open("C:/Users/tim67/Desktop/js/GeekB/13 go/Integration-ParallelExe/hw2/task1/some.txt")
	if err != nil {
		fmt.Println("Невозможно открыть файл или он не существует!\n", err)
		os.Exit(1)
	}
	defer fileR.Close()

	readable := make([]byte, 64)
	for {
		n, err := fileR.Read(readable)
		if err == io.EOF {
			fmt.Println("Конец файла.")
			break
		}
		fmt.Print(string(readable[:n]))
	}
}

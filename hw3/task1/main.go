package main

import (
	"fmt"
	"strconv"
)

func main() {
	var inpNum int
	fc := inData()
	sc := madeSquare(&inpNum, fc)
	tc := prodNum(inpNum, sc)
	fmt.Println("Произведение:", <-tc)

}

func inData() chan int {

	inChan := make(chan int)
	for {
		var input string
		fmt.Print("\nДля выходя введите \"stop\"\nВведите число, прогроамма рассчитает его квадрат и произведение квадрата на введенное число: ")
		fmt.Scanln(&input)

		if input == "s" {
			break
		}

		n, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Введите число или команду \"stop\"")
			continue
		}
		go func() {
			inChan <- n
		}()
		break
	}
	return inChan
}

func madeSquare(inpNum *int, inChan chan int) chan int {
	sChan := make(chan int)
	inN := <-inChan
	*inpNum = inN
	fmt.Println("\nВвод:", inN)
	go func() {
		sChan <- inN * inN
	}()
	return sChan
}

func prodNum(inpNum int, sChan chan int) chan int {
	tChan := make(chan int)
	sqN := <-sChan
	fmt.Println("Квадрат:", sqN)
	go func() {
		tChan <- inpNum * sqN
	}()
	return tChan
}

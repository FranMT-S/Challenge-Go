package helpers

import "fmt"

func Helper() {
	fmt.Print("este es un Helper")
}

type MyTest uint32

func (MyTest).IsOne() {
	fmt.Println("test One")
}

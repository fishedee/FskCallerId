package main

import (
	"fmt"
	. "github.com/fishedee/language"
	"github.com/jacobsa/go-serial/serial"
)

func run() {
	defer CatchCrash(func(e Exception) {
		fmt.Println(e)
	})
	options := serial.OpenOptions{
		PortName:              "COM3",
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 200,
		MinimumReadSize:       0,
	}
	port, err := serial.Open(options)
	if err != nil {
		panic(err)
	}
	defer port.Close()

	counter := 0
	for {
		buf := make([]byte, 128)
		n, err := port.Read(buf)
		if err != nil {
			panic("uc")
		}
		if n == 0 {
			continue
		}
		counter++
		fmt.Printf("\n[%d] [%s]", counter, string(buf[0:n]))
	}
}

func main() {
	for {
		run()
		fmt.Println("try again...")
	}

}

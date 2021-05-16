package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		input, _ := reader.ReadString('\n')
		if input == "exit\r\n" {
			return
		}
		result, err := Run("<stdin>", input)
		if err != nil {
			log.Print(err)
		} else {
			fmt.Println(result)
		}
	}
}

func init() {
	log.SetPrefix("Error: ")
	log.SetFlags(0)
}

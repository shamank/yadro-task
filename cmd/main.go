package main

import (
	"fmt"
	"github.com/shamank/yadro-task/internal/app"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Errorf("fail! filename is not specified: %s <filename>", os.Args[0])
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Errorf("fail! unable to open file: %s \n%s", fileName, err.Error())
		return
	}
	defer file.Close()

	a := app.NewApp(file)

	a.Run()

}

package app

import (
	"bufio"
	"fmt"
	usecases "github.com/shamank/yadro-task/internal/usecase"
	"io"
)

type App struct {
	file io.Reader
}

func NewApp(file io.Reader) *App {
	return &App{
		file: file,
	}
}

func (a *App) Run() {

	sc := bufio.NewScanner(a.file)

	club, err := usecases.NewClubFromScanner(sc)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	events, err := usecases.NewEventsFromScanner(sc)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ch := make(chan string)

	go usecases.StartWork(club, events, ch)

	for c := range ch {
		fmt.Println(c)
	}

	return
}

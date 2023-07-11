package app

import (
	"bufio"
	usecases "github.com/shamank/yadro-task/internal/usecase"
	"os"
	"testing"
)

const testDIR = "../../tests"

var filesTests = []string{
	"case_1",
	"case_2",
	"case_3",
	"case_4",
	"case_5",
}

func TestApp_Run(t *testing.T) {
	for _, test := range filesTests {
		_testingFile(t, test)
	}
}

func _testingFile(t *testing.T, test string) {
	file, err := os.Open(testDIR + "/" + test + ".txt")
	defer file.Close()

	if err != nil {
		t.Errorf("error occurred open file in testcase: %s \n - %s", test, err.Error())
		return
	}

	sc := bufio.NewScanner(file)

	ch := make(chan string)

	go usecases.StartWork(sc, ch)

	correctFile, err := os.Open(testDIR + "/" + test + "_out.txt")

	if err != nil {
		t.Errorf("error occurred open file in testcase: %s \n - %s", test, err.Error())
		return
	}

	defer correctFile.Close()

	correctOut := bufio.NewScanner(correctFile)

	i := 1

	for c := range ch {

		correctOut.Scan()

		text := correctOut.Text()

		if c != text {
			t.Errorf("[%s:%d] %s is not equal \"%s\"", test, i, c, text)
		}
		i++
	}

	if correctOut.Scan() {
		t.Errorf("[%s] there are still more out to be given!", test)
	}

}

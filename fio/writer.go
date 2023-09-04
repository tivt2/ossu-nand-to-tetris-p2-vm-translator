package fio

import (
	"fmt"
	"os"
	"regexp"
	"sync"
)

type writer struct {
	filePath string
	Wg       *sync.WaitGroup
}

func NewWriter(filePath string) *writer {
	return &writer{
		filePath: filePath,
		Wg:       &sync.WaitGroup{},
	}
}

func (w *writer) Write(writeChan chan string) {

	regex := regexp.MustCompile(`^(.*?)(\.vm)$`)
	extracted := regex.FindStringSubmatch(w.filePath)

	if len(extracted) < 3 {
		panic("File doesnt match *.vm")
	}

	OutputFilePath := fmt.Sprintf("%v.asm", extracted[1])
	file, err := os.Create(OutputFilePath)
	checkErr(err)
	defer file.Close()

	for parsedText := range writeChan {
		file.Write([]byte(parsedText + "\n"))
	}
	w.Wg.Done()
}

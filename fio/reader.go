package fio

import (
	"bufio"
	"os"
)

type reader struct {
	filePath string
	ReadChan chan string
}

func NewReader(filePath string) *reader {
	readChan := make(chan string, 20)

	return &reader{
		filePath: filePath,
		ReadChan: readChan,
	}
}

func (r *reader) Read() {

	file, err := os.Open(r.filePath)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		r.ReadChan <- scanner.Text()
	}
	close(r.ReadChan)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

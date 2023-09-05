package translator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Translate(path string) {
	if isVMfile(path) {
		translateFile(path)
	}

	info, err := os.Stat(path)
	checkErr(err)
	if info.IsDir() {
		folder, err := os.Open(path)
		checkErr(err)
		files, err2 := folder.Readdirnames(0)
		checkErr(err2)
		for _, file := range files {
			if isVMfile(file) {
				go translateFile(filepath.Join(path, file))
			}
		}
	}
}

func isVMfile(path string) bool {
	length := len(path)
	return length > 3 && path[length-3:] == ".vm"
}

func translateFile(filePath string) {
	lineIn := make(chan string, 40)
	go readFile(filePath, lineIn)

	lineOut := make(chan string, 40)
	go parse(lineIn, lineOut)

	writeFile(filePath, lineOut)
}

func readFile(filePath string, lineIn chan<- string) {

	file, err := os.Open(filePath)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineIn <- scanner.Text()
	}
	close(lineIn)
}

func writeFile(filePath string, lineOut <-chan string) {
	OutputFilePath := fmt.Sprintf("%v.asm", filePath[:len(filePath)-3])
	file, err := os.Create(OutputFilePath)
	checkErr(err)
	defer file.Close()

	for parsedText := range lineOut {
		file.Write([]byte(parsedText + "\n"))
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

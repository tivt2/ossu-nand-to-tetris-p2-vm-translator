package translator

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Translate(path string) {
	outputMap := map[string]string{}
	wg := sync.WaitGroup{}

	if filepath.Ext(path) == ".vm" {
		fileName := strings.ToLower(filepath.Base(path))
		outputMap[fileName] = ""

		translateFile(path, fileName, &outputMap, &wg)

		outputPath := path[:len(path)-2] + "asm"
		wg.Wait()
		writeFile(outputPath, &outputMap, false)
		return
	}

	pathInfo, err := os.Stat(path)
	checkErr(err)
	if pathInfo.IsDir() {
		folder, err := os.Open(path)
		checkErr(err)
		files, err2 := folder.Readdirnames(0)
		checkErr(err2)

		for _, file := range files {
			if filepath.Ext(file) == ".vm" {
				fileName := strings.ToLower(file[:len(file)-3])
				outputMap[fileName] = ""
				filePath := filepath.Join(path, file)

				wg.Add(1)
				go func() {
					translateFile(filePath, fileName, &outputMap, &wg)
					wg.Done()
				}()
			}
		}

		outputPath := filepath.Join(path, filepath.Base(path)+".asm")
		wg.Wait()
		writeFile(outputPath, &outputMap, true)
	}
}

func translateFile(filePath string, fileName string, outputMap *map[string]string, wg *sync.WaitGroup) {
	inputChan := make(chan string, 40)

	wg.Add(2)
	go func() {
		readFile(filePath, inputChan)
		wg.Done()
	}()

	go func() {
		parse(fileName, inputChan, outputMap)
		wg.Done()
	}()
}

func readFile(filePath string, inputChan chan<- string) {
	file, err := os.Open(filePath)
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		inputChan <- scanner.Text()
	}
	close(inputChan)
}

func writeFile(outputPath string, outputMap *map[string]string, bootstrap bool) {
	fmt.Println("\nCreating file ->", outputPath, "<-")
	file, err := os.Create(outputPath)
	checkErr(err)
	defer file.Close()

	if bootstrap {
		bootstrap := "// BOOTSTRAP\n@256\nD=A\n@SP\nM=D\n" + call([]string{"call", "sys.init", "0"}) + "\n"
		file.WriteString(bootstrap)
	}

	for _, text := range *outputMap {
		file.WriteString(text)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

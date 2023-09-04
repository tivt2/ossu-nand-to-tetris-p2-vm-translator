package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tivt2/vm-translator/translator"
)

func main() {
	start := time.Now()
	filePath := os.Args[1]

	translator.Translate(filePath)

	fmt.Printf("Compilation time: %s\n", time.Since(start))
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tivt2/vm-translator/translator"
)

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		log.Fatalf("Usage: VMtranslator <path_to_folder_or_file.vm>")
	}
	arg1 := os.Args[1]

	translator.Translate(arg1)

	fmt.Printf("Compilation time: %v\n", time.Since(start))
}

package translator

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func parse(fileName string, inputChan <-chan string, outputMap *map[string]string) {
	for line := range inputChan {
		cleanLine(&line)
		if line == "" {
			continue
		}
		parsedLine := parseLine(line, fileName)
		(*outputMap)[fileName] += parsedLine + "\n"
	}
}

func cleanLine(line *string) {
	regex := regexp.MustCompile(`//.*`)
	*line = regex.ReplaceAllString(*line, "")
	*line = strings.TrimSpace(*line)
	*line = strings.ToLower(*line)
}

var (
	arithLogical = map[string]bool{
		"add": true,
		"sub": true,
		"neg": true,
		"eq":  true,
		"gt":  true,
		"lt":  true,
		"and": true,
		"or":  true,
		"not": true,
	}

	memoryAccs = map[string]bool{
		"push": true,
		"pop":  true,
	}

	branching = map[string]bool{
		"label":   true,
		"goto":    true,
		"if-goto": true,
	}

	function = map[string]bool{
		"function": true,
		"call":     true,
		"return":   true,
	}
)

func parseLine(line string, fileName string) string {
	commands := strings.Split(line, " ")
	out := ""

	fmt.Println("parsing:", commands)
	if _, ok := arithLogical[commands[0]]; ok {
		out = parseArithLogical(commands)

	} else if _, ok := memoryAccs[commands[0]]; ok {
		out = parseMemoryAccs(commands, fileName)

	} else if _, ok := branching[commands[0]]; ok {
		out = parseBranching(commands)

	} else if _, ok := function[commands[0]]; ok {
		out = parseFunction(commands)

	} else {
		log.Fatalf("Error while parsing line:\n%v", line)
	}

	return fmt.Sprintf("// %s\n%s", line, out)
}

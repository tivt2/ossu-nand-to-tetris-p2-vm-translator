package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type parser struct {
	readChan  chan string
	WriteChan chan string
}

func New(readChan chan string) *parser {
	writeChan := make(chan string, 20)

	return &parser{
		readChan:  readChan,
		WriteChan: writeChan,
	}
}

func (p *parser) Parse() {
	for line := range p.readChan {
		p.cleanLine(&line)
		if line == "" {
			continue
		}
		parsedLine := p.parseLine(line)
		p.WriteChan <- parsedLine
	}
	close(p.WriteChan)
}

func (p *parser) cleanLine(line *string) {
	regex := regexp.MustCompile(`//.*`)
	*line = regex.ReplaceAllString(*line, "")
	*line = strings.TrimSpace(*line)
	*line = strings.ToLower(*line)
}

func (p *parser) parseLine(line string) string {
	regex := regexp.MustCompile(`^(push |pop )`)
	isPushOrPop := regex.MatchString(line)
	parsedLine := ""
	if isPushOrPop {
		parsedLine = p.parsePushPopOperation(line)
	} else {
		parsedLine = p.parseArithLogicalOperation(line)
	}
	return parsedLine
}

var (
	xSubY_D     = "AM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D"
	al_commands = map[string]string{
		"add": "AM=M-1\nD=M\n@SP\nAM=M-1\nM=D+M\n@SP\nM=M+1",
		"sub": "AM=M-1\nD=M\n@SP\nAM=M-1\nM=M-D\n@SP\nM=M+1",
		"neg": "AM=M-1\nM=-M\n@SP\nM=M+1",
		"and": "AM=M-1\nD=M\n@SP\nAM=M-1\nM=D&M\n@SP\nM=M+1",
		"or":  "AM=M-1\nD=M\n@SP\nAM=M-1\nM=D|M\n@SP\nM=M+1",
		"not": "AM=M-1\nM=!M\n@SP\nM=M+1",
	}

	seg_pointer = map[string]string{
		"local":    "@LCL",
		"argument": "@ARG",
		"this":     "@THIS",
		"that":     "@THAT",
	}

	seg_body = map[string]string{
		"push": "D=M\n@SP\nA=M\nM=D\n@SP\nM=M+1",
		"pop":  "@SP\nAM=M-1\nD=D+M\nA=D-M\nM=D-A",
	}
)

func (p *parser) parseArithLogicalOperation(line string) string {
	if line == "eq" || line == "gt" || line == "lt" {
		return logicalDynamic(line)
	}
	command, ok := al_commands[line]
	checkNotOk(ok, line)

	return fmt.Sprintf("// %s\n%s", line, command)
}

var logicalDynamic = func() func(command string) string {
	count := 0
	return func(command string) string {
		count++
		operation := map[string]string{
			"eq": fmt.Sprintf("%s\n@EQ%d\nD;JEQ\n@SP\nA=M\nM=0\n@NOTEQ%d\n0;JMP\n(EQ%d)\n@SP\nA=M\nM=-1\n(NOTEQ%d)\n@SP\nM=M+1", xSubY_D, count, count, count, count),
			"gt": fmt.Sprintf("%s\n@GT%d\nD;JGT\n@SP\nA=M\nM=0\n@NOTGT%d\n0;JMP\n(GT%d)\n@SP\nA=M\nM=-1\n(NOTGT%d)\n@SP\nM=M+1", xSubY_D, count, count, count, count),
			"lt": fmt.Sprintf("%s\n@LT%d\nD;JLT\n@SP\nA=M\nM=0\n@NOTLT%d\n0;JMP\n(LT%d)\n@SP\nA=M\nM=-1\n(NOTLT%d)\n@SP\nM=M+1", xSubY_D, count, count, count, count),
		}
		return operation[command]
	}
}()

func (p *parser) parsePushPopOperation(line string) string {
	operation := strings.Split(line, " ")
	if len(operation) < 3 {
		panic(fmt.Sprintf("Bad operation call:\n%s\n", line))
	}

	printLine := "// " + line

	if operation[1] == "pointer" {
		return fmt.Sprintf("%s\n%s", printLine, pointerCommand(operation))
	}

	if operation[1] == "static" || operation[1] == "temp" {
		return fmt.Sprintf("%s\n%s", printLine, staticTempCommand(operation))
	}

	body, ok := seg_body[operation[0]]
	checkNotOk(ok, line)
	if operation[1] == "constant" {
		return fmt.Sprintf("%s\n@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1", printLine, operation[2])
	}

	pointer, ok := seg_pointer[operation[1]]
	checkNotOk(ok, line)
	if operation[2] == "0" {
		return fmt.Sprintf("%s\n%s\nAD=M\n%s", printLine, pointer, body)
	}

	offset, ok := segOffset(operation[0], operation[2])
	checkNotOk(ok, line)
	return fmt.Sprintf("%s\n%s\nD=M\n%s\n%s", printLine, pointer, offset, body)
}

func segOffset(command string, value string) (string, bool) {
	offset := map[string]string{
		"push": fmt.Sprintf("@%s\nA=D+A\nD=M", value),
		"pop":  fmt.Sprintf("@%s\nD=D+A", value),
	}

	out, ok := offset[command]
	return out, ok
}

func pointerCommand(operations []string) string {
	location := map[string]string{
		"0": "@3",
		"1": "@4",
	}
	pointer := location[operations[2]]
	command := map[string]string{
		"push": fmt.Sprintf("%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1", pointer),
		"pop":  fmt.Sprintf("@SP\nAM=M-1\nD=M\n%s\nM=D", pointer),
	}

	return command[operations[0]]
}

func staticTempCommand(operations []string) string {
	var offset int
	if operations[1] == "static" {
		offset = 16
	} else {
		offset = 5
	}
	strOffset := offsetString(operations[2], offset)
	command := map[string]string{
		"push": fmt.Sprintf("@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1", strOffset),
		"pop":  fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s\nM=D", strOffset),
	}

	return command[operations[0]]
}

func offsetString(str string, offset int) string {
	valInt, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	out := fmt.Sprintf("%v", valInt+offset)
	return out
}

func checkNotOk(ok bool, line string) {
	if !ok {
		panic(fmt.Sprintf("Bad operation call:\n%s\n", line))
	}
}

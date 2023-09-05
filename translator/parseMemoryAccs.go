package translator

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var (
	pointer = map[string]string{
		"local":    "@LCL",
		"argument": "@ARG",
		"this":     "@THIS",
		"that":     "@THAT",
	}

	body = map[string]string{
		"push": "D=M\n@SP\nA=M\nM=D\n@SP\nM=M+1",
		"pop":  "@SP\nAM=M-1\nD=D+M\nA=D-M\nM=D-A",
	}
)

func parseMemoryAccs(commands []string) string {
	if len(commands) < 3 {
		log.Fatalf("Bad operation call:\n%s\n", strings.Join(commands, " "))
	}

	cmd1 := commands[1]
	switch cmd1 {
	case "pointer", "static", "temp":
		return staticTempPointerCommand(commands)
	}

	cmd2 := commands[2]
	if cmd1 == "constant" {
		return fmt.Sprintf("@%s\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1", cmd2)
	}

	cmd0 := commands[0]

	if cmd2 == "0" {
		return fmt.Sprintf("%s\nAD=M\n%s", pointer[cmd1], body[cmd0])
	}

	return fmt.Sprintf("%s\nD=M\n%s\n%s", pointer[cmd1], segOffset(cmd0, cmd2), body[cmd0])
}

func segOffset(command string, memoryIdx string) string {
	offset := map[string]string{
		"push": fmt.Sprintf("@%s\nA=D+A\nD=M", memoryIdx),
		"pop":  fmt.Sprintf("@%s\nD=D+A", memoryIdx),
	}

	out := offset[command]
	return out
}

func staticTempPointerCommand(commands []string) string {
	var offset string
	op1 := commands[1]
	if op1 == "pointer" {
		op2 := commands[2]
		if op2 == "0" {
			offset = "3"
		} else if op2 == "1" {
			offset = "4"
		}
		log.Fatalf("%v <- is not a valid pointer 'i'")
	} else if op1 == "static" {
		offset = offsetString(commands[2], 16)
	} else {
		offset = offsetString(commands[2], 5)
	}

	command := map[string]string{
		"push": fmt.Sprintf("@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1", offset),
		"pop":  fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s\nM=D", offset),
	}

	return command[commands[0]]
}

func offsetString(idx string, offset int) string {
	valInt, err := strconv.Atoi(idx)
	if err != nil {
		log.Fatalf("%v <- is not a valid 'i'", idx)
	}
	out := fmt.Sprintf("%v", valInt+offset)
	return out
}

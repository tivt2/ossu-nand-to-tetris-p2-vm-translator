package translator

import (
	"fmt"
	"strconv"
)

func parseFunction(commands []string) string {
	switch commands[0] {
	case "call":
		return call(commands)
	case "function":
		return functionDef(commands)
	case "return":
		return returnDef
	default:
		return ""
	}
}

// RETURN EXECUTION
const (
	getEndFrameAndRetAddr = "@LCL\nD=M\n@R13\nM=D\n@5\nA=D-A\nD=M\n@R14\nM=D\n"
	popRetVal             = "@SP\nAM=M-1\nD=M\n@ARG\nA=M\nM=D\n"
	repositionSP          = "@ARG\nD=M\n@SP\nM=D+1\n"
	restoreTHAT           = "@R13\nD=M\n@1\nA=D-A\nD=M\n@THAT\nM=D\n"
	restoreTHIS           = "@R13\nD=M\n@2\nA=D-A\nD=M\n@THIS\nM=D\n"
	restoreARG            = "@R13\nD=M\n@3\nA=D-A\nD=M\n@ARG\nM=D\n"
	restoreLCL            = "@R13\nD=M\n@4\nA=D-A\nD=M\n@LCL\nM=D\n"
	restoreFrames         = restoreTHAT + restoreTHIS + restoreARG + restoreLCL
	gotoAddr              = "@R14\nA=M\n0;JMP"
	returnDef             = getEndFrameAndRetAddr + popRetVal + repositionSP + restoreFrames + gotoAddr
)

const (
	saveLCL    = "@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	saveARG    = "@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	saveTHIS   = "@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	saveTHAT   = "@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	saveNewLCL = "@SP\nD=M\n@LCL\nM=D"
	saveFrames = saveLCL + saveARG + saveTHIS + saveTHAT + saveNewLCL
)

var call = func() func(commands []string) string {
	count := -1
	return func(commands []string) string {
		count++
		nArgs, err := strconv.Atoi(commands[2])
		checkErr(err)
		funcName := commands[1]
		saveRetAddr, injectRetAddrLabel := genRetAddr(funcName, count)
		return fmt.Sprintf("%s\n%s\n%s\n@%v\n0;JMP\n%s", saveRetAddr, saveFrames, repositionArg(nArgs), funcName, injectRetAddrLabel)
	}
}()

func genRetAddr(funcName string, n int) (saveRetAddr string, injectRetAddrLabel string) {
	saveRetAddr = fmt.Sprintf("@%v$ret.%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1", funcName, n)
	injectRetAddrLabel = fmt.Sprintf("(%v$ret.%d)", funcName, n)
	return
}

func repositionArg(n int) string {
	n += 5
	return fmt.Sprintf("@SP\nD=M\n@%d\nD=D-A\n@ARG\nM=D", n)
}

func functionDef(commands []string) string {
	nVars, err := strconv.Atoi(commands[2])
	checkErr(err)
	if commands[2] == "0" {
		return fmt.Sprintf("(%v)", commands[1])
	}
	return fmt.Sprintf("(%v)\n%s", commands[1], addLCLVars(nVars))
}

func addLCLVars(n int) string {
	if n == 1 {
		return "@SP\nA=M\nM=0\n@SP\nM=M+1"
	}

	out := "@SP\nA=M\nM=0\n@SP\nAM=M+1\n"
	for n > 1 {
		if n == 2 {
			out += "\nM=0\n@SP\nM=M+1"
			n--
			continue
		}

		out += "\nM=0\n@SP\nAM=M+1"

		n--
	}

	return out
}

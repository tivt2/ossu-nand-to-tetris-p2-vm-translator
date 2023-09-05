package translator

import "fmt"

var (
	xSubY_D     = "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D"
	al_commands = map[string]string{
		"add": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D+M\n@SP\nM=M+1",
		"sub": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M-D\n@SP\nM=M+1",
		"neg": "@SP\nAM=M-1\nM=-M\n@SP\nM=M+1",
		"and": "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D&M\n@SP\nM=M+1",
		"or":  "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D|M\n@SP\nM=M+1",
		"not": "@SP\nAM=M-1\nM=!M\n@SP\nM=M+1",
	}
)

func parseArithLogical(command []string) string {
	cmd0 := command[0]
	switch cmd0 {
	case "eq", "gt", "lt":
		return logicalDynamic(cmd0)
	default:
		return al_commands[cmd0]
	}
}

var logicalDynamic = func() func(command string) string {
	count := -1
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

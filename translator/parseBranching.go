package translator

import "fmt"

func parseBranching(commands []string) string {
	cmd1 := commands[1]
	branching := map[string]string{
		"label":   fmt.Sprintf("(%v)", cmd1),
		"goto":    fmt.Sprintf("@%v\n0;JMP", cmd1),
		"if-goto": fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%v\nD;JNE", cmd1),
	}

	return branching[commands[0]]
}

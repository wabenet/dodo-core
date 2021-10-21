package command

import (
	"strconv"

	"github.com/spf13/cobra"
)

const (
	annotationExitCode = "exit-code"
	exitCodeOK         = 0
)

func SetExitCode(cmd *cobra.Command, code int) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	cmd.Annotations[annotationExitCode] = strconv.Itoa(code)
}

func GetExitCode(cmd *cobra.Command) int {
	if cmd.Annotations != nil {
		if a, ok := cmd.Annotations[annotationExitCode]; ok {
			if code, err := strconv.Atoi(a); err == nil {
				return code
			}
		}
	}

	for _, sub := range cmd.Commands() {
		if code := GetExitCode(sub); code != exitCodeOK {
			return code
		}
	}

	return exitCodeOK
}

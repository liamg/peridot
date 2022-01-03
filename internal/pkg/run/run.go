package run

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/liamg/tml"
)

func Run(command string, dir string, sudo bool) error {
	cmd := createCommand(command, sudo)
	if sudo {
		tml.Printf("\n<bold><blue>This change requires root access. Please enter your password if prompted.</blue></bold>\n")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.Dir = dir
	return cmd.Run()
}

func createCommand(command string, sudo bool) *exec.Cmd {
	if sudo {
		command = fmt.Sprintf("sudo %s", command)
	}
	return exec.Command("sh", "-c", command)
}

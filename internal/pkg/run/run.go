package run

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/liamg/tml"
)

func Run(command string, dir string, sudo bool, interactive bool) error {
	cmd := createCommand(command, sudo)
	errStr := bytes.NewBufferString("")
	if sudo {
		tml.Printf("\n<bold><blue>This change requires root access. Please enter your password if prompted.</blue></bold>\n")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else if interactive {
		errStr.WriteString("see above")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = errStr
	}
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %s", err, errStr.String())
	}
	return nil
}

func createCommand(command string, sudo bool) *exec.Cmd {
	if sudo {
		return exec.Command("sudo", "sh", "-c", command)
	}
	return exec.Command("sh", "-c", command)
}

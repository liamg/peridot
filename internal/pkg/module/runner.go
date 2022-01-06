package module

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/liamg/peridot/internal/pkg/log"
	"github.com/liamg/tml"
)

type Runner struct {
	module Module
}

func NewRunner(module Module) *Runner {
	return &Runner{
		module: module,
	}
}

func (r *Runner) Run(command string, sudo bool) error {

	if sudo {
		tml.Printf("\n<bold><blue>This change requires root access. Please enter your password if prompted.</blue></bold>\n")
	}

	var cmd *exec.Cmd
	if sudo {
		cmd = exec.Command("sudo", "sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd = exec.Command("sh", "-c", command)
		output := log.NewDebugStreamer(
			tml.Sprintf("[<yellow><bold>%s</bold></yellow>]", r.module.Name()),
		)
		defer output.Flush()
		output.Write([]byte(fmt.Sprintf("Running command: %s\n", command)))
		cmd.Stdout = output
		cmd.Stderr = output
	}

	cmd.Dir = r.module.Path()
	return cmd.Run()
}

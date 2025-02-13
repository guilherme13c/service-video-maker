package command

import (
	"bytes"
	"context"
	"os/exec"
)

func (c *commandRepository) Run(ctx context.Context, cmd string, args ...string) error {
	command := exec.CommandContext(ctx, cmd, args...)

	var outb, errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb

	err := command.Run()
	if err != nil {
		return err
	}

	return nil
}

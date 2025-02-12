package command

import (
	"bytes"
	"context"
	"os/exec"

	"log"
)

func (c *commandRepository) Run(ctx context.Context, cmd string, args ...string) error {
	log.Printf("Running command: %s %v\n", cmd, args)

	command := exec.CommandContext(ctx, cmd, args...)

	var outb, errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb

	err := command.Run()
	log.Printf("Command output: %s\terror: %s\n", outb.String(), errb.String())
	if err != nil {
		return err
	}

	return nil
}

package biz

import (
	"context"
	"os"
	"os/exec"
)

type LocalExecutor struct{}

func (e *LocalExecutor) Execute(ctx context.Context, shellCmd string) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", shellCmd)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

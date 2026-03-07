package biz

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/muidea/magicCli/pkg/util"
)

type SSHExecutor struct {
	Host string
	Port int
}

func (e *SSHExecutor) Execute(ctx context.Context, shellCmd string) error {
	args := []string{"-q"}
	if util.IsTTY() {
		args = append(args, "-t")
	}
	if e.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", e.Port))
	}
	args = append(args, e.Host, "--", shellCmd)

	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

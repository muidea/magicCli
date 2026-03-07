package biz

import (
	"context"
	"os"
	"os/exec"

	"github.com/muidea/magicCli/pkg/util"
)

type DockerExecutor struct {
	ContainerID string
	User        string
}

func (e *DockerExecutor) Execute(ctx context.Context, shellCmd string) error {
	args := []string{"exec"}
	if util.IsTTY() {
		args = append(args, "-it")
	}
	if e.User != "" {
		args = append(args, "-u", e.User)
	}
	args = append(args, e.ContainerID, "sh", "-c", shellCmd)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}

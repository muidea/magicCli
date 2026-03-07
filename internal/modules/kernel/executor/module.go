package executor

import (
	"context"
	"strings"

	"github.com/muidea/magicCli/internal/modules/kernel/executor/biz"
)

type Executor interface {
	Execute(ctx context.Context, shellCmd string) error
}

func NewExecutor(target, user string, port int) Executor {
	target = strings.ToLower(target)
	if target == "" || target == "local" || target == "host" {
		return &biz.LocalExecutor{}
	}
	if strings.Contains(target, "@") || strings.Contains(target, ".") {
		return &biz.SSHExecutor{Host: target, Port: port}
	}
	return &biz.DockerExecutor{ContainerID: target, User: user}
}

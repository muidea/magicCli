# magicCli 最终设计文档 (v2.0)

## 1. 项目概述
**magicCli** 是一个基于 Go 开发的统一命令执行工具。它通过抽象层屏蔽了**宿主机 (Local)**、**Docker 容器 (Docker)** 以及**远程服务器 (SSH)** 的底层差异。用户只需指定目标（Target），即可在任何环境下以一致的体验执行 Shell 命令。

## 2. 系统架构设计

### 2.1 目录结构规范
项目采用多模块分层结构，确保核心逻辑（Internal）与应用入口（Application）分离：

```text
magicCli/
├── Makefile                    # 自动化构建脚本
├── go.mod                     # 模块路径: github.com/muidea/magicCli
├── vendor/                    # 依赖vendor目录
├── application/
│   └── magicCli/              # CLI 应用入口
│       └── cmd/
│           └── main.go        # 参数解析与路由分发
├── internal/
│   └── modules/
│       └── kernel/
│           └── executor/      # 执行器内核模块
│               ├── module.go  # 接口定义与工厂
│               └── biz/       # 核心驱动实现
│                   ├── local.go
│                   ├── docker.go
│                   └── ssh.go
└── pkg/
    └── util/
        └── system.go          # 公共系统工具
```

### 2.2 路由识别规则
`Executor` 工厂会根据输入的 `-t` 参数自动选择驱动：
1.  **Local**: 输入为 `local`、`host` 或为空。
2.  **SSH**: 输入包含 `@` 符号（如 `root@10.0.0.1`）或包含 `.` 符号（域名/IP）。
3.  **Docker**: 排除上述情况后的所有字符串，均视为容器名称或 ID。

---

## 3. 核心代码实现

### 3.1 公共包 (pkg)
**`pkg/util/system.go`**
```go
package util

import "os"

// IsTTY 检测当前 Stdout 是否为交互式终端
func IsTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
```

### 3.2 驱动逻辑层 (internal/modules/kernel/executor/biz)

**`local.go` - 宿主机驱动**
```go
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
```

**`docker.go` - 容器驱动**
```go
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
		args = append(args, "-it") // 自动处理交互终端
	}
	if e.User != "" {
		args = append(args, "-u", e.User)
	}
	args = append(args, e.ContainerID, "sh", "-c", shellCmd)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}
```

**`ssh.go` - 远程驱动**
```go
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
		args = append(args, "-t") // 支持交互式如 top/vi
	}
	if e.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", e.Port))
	}
	args = append(args, e.Host, "--", shellCmd)

	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
	return cmd.Run()
}
```

### 3.3 内核模块接口 (internal/modules/kernel/executor)

**`module.go` - 接口与工厂**
```go
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
	// 识别 SSH: 包含 @ 或 . (IP/域名格式)
	if strings.Contains(target, "@") || strings.Contains(target, ".") {
		return &biz.SSHExecutor{Host: target, Port: port}
	}
	// 默认 Docker
	return &biz.DockerExecutor{ContainerID: target, User: user}
}
```

### 3.4 应用入口 (application/magicCli/cmd)

**`main.go`**
```go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/muidea/magicCli/internal/modules/kernel/executor"
	"github.com/spf13/cobra"
)

func main() {
	var target string
	var user string
	var port int

	var rootCmd = &cobra.Command{
		Use:   "magicCli [flags] -- <command>",
		Short: "magicCli: 一键跨环境执行工具",
		Example: `  magicCli -t my-nginx 'ls -l'
  magicCli -t root@10.0.0.1 -p 2222 uptime`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shellCmd := strings.Join(args, " ")
			execInst := executor.NewExecutor(target, user, port)

			err := execInst.Execute(cmd.Context(), shellCmd)
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					os.Exit(exitErr.ExitCode()) // 透传状态码
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&target, "target", "t", "local", "执行目标 (local/容器ID/远程Host)")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "指定用户 (仅限Docker)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH 端口")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

---

## 4. 构建脚本 (Makefile)
```makefile
.PHONY: build install clean

APP_NAME=magicCli
SRC=application/magicCli/cmd/main.go

build:
	go build -o bin/$(APP_NAME) $(SRC)

install:
	cp bin/$(APP_NAME) ~/.local/bin/

clean:
	rm -rf bin/
```

---

## 5. 验收测试说明 (Acceptance)

### 5.1 测试场景与指令
| 测试项 | 执行指令 | 预期结果 |
| :--- | :--- | :--- |
| **本地测试** | `./magicCli "echo $USER"` | 输出当前宿主机用户名 |
| **管道测试** | `./magicCli "ls -l \| wc -l"` | 正确统计文件数量 |
| **Docker测试** | `./magicCli -t my-cnt "cat /etc/hostname"` | 输出容器的 Hostname |
| **SSH测试** | `./magicCli -t user@localhost "uptime"` | 成功通过 SSH 获取运行时间 |
| **SSH端口测试** | `./magicCli -t user@host -p 2222 "uptime"` | 通过指定端口 SSH 连接 |
| **交互式测试** | `./magicCli -t my-cnt top` | 进入 top 实时监控界面，按 Q 退出 |
| **错误透传** | `./magicCli "exit 42" ; echo $?` | 屏幕打印 `42` |

### 5.2 验收合格标准
1.  **一致性**: 不论 target 是什么，命令执行的 Stdout/Stderr 必须实时流式输出。
2.  **隔离性**: Docker 驱动不能影响宿主机文件系统，SSH 驱动必须走加密隧道。
3.  **兼容性**: 必须能正确处理带引号和不带引号的命令参数。
4.  **结构性**: 代码必须通过 `go build` 成功编译，目录结构符合规范。
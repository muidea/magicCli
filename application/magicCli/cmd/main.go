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

	var rootCmd = &cobra.Command{
		Use:   "magicCli [flags] -- <command>",
		Short: "magicCli: 一键跨环境执行工具",
		Example: `  magicCli -t my-nginx 'ls -l'
  magicCli -t root@10.0.0.1 uptime`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shellCmd := strings.Join(args, " ")
			execInst := executor.NewExecutor(target, user)

			err := execInst.Execute(cmd.Context(), shellCmd)
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					os.Exit(exitErr.ExitCode())
				}
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&target, "target", "t", "local", "执行目标 (local/容器ID/远程Host)")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "指定用户 (仅限Docker)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/muidea/magicCli/internal/modules/kernel/executor"
	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
	buildDate string
)

func main() {
	var target string
	var user string
	var port int

	var rootCmd = &cobra.Command{
		Use:   "magicCli [flags] -- <command>",
		Short: "magicCli: 一键跨环境执行工具",
		Example: `  magicCli -t my-nginx 'ls -l'
  magicCli -t root@10.0.0.1 uptime`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shellCmd := strings.Join(args, " ")
			execInst := executor.NewExecutor(target, user, port)

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

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "显示版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:   %s\n", version)
			fmt.Printf("GitCommit: %s\n", gitCommit)
			fmt.Printf("BuildDate: %s\n", buildDate)
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&target, "target", "t", "local", "执行目标 (local/容器ID/远程Host)")
	rootCmd.Flags().StringVarP(&user, "user", "u", "", "指定用户 (仅限Docker)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 22, "SSH 端口")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

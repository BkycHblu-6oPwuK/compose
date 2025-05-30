package cmd

import (
	"docky/config"
	"docky/utils/globalHelper"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:                "share",
	Short:              "Туннелирование локального сайта",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		globalHelper.ValidateWorkDir()
		if err := share(args); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Ошибка: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}

func share(args []string) error {
	hasAuth := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "--auth") {
			hasAuth = true
			break
		}
	}

	if !hasAuth {
		authToken := "e17105f7-e499-470a-bd5b-05c0a579036f"
		args = append(args, "--auth="+authToken)
	}

	cmdArgs := append([]string{
		"run", "--init", "--rm", "-p", "4040:4040", "-t",
		"beyondcodegmbh/expose-server:latest", "share", "http://host.docker.internal:80",
	}, args...)

	cmd := exec.Command("docker", cmdArgs...)
	cmd.Dir = config.GetWorkDirPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}


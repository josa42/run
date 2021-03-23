/*
Copyright © 2021 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/josa42/run/run"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := run.GetTasks()
		tasks.Run(args[0])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

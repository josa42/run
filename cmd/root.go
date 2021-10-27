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
var Version = "development"

var rootCmd = &cobra.Command{
	Use: "run",
	Args: func(cmd *cobra.Command, args []string) error {
		if ok, _ := cmd.Flags().GetBool("version"); ok {
			if len(args) != 0 {
				return fmt.Errorf("accepts no argss, received %d", len(args))
			}
		} else if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		if ok, _ := cmd.Flags().GetBool("version"); ok {
			fmt.Println(Version)
			os.Exit(0)
		}

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
	rootCmd.Flags().BoolP("version", "", false, "Show version")
}

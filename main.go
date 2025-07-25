package main

import (
	"github.com/spf13/cobra"
	"underdog/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "underdog",
		Short: "Underdog: Simple OCI container runtime",
	}

	rootCmd.AddCommand(cmd.BuildCmd)
	rootCmd.AddCommand(cmd.RunCmd)
	rootCmd.Execute()
}

package cmd

import (
	"fmt"
	"underdog/parser"
	"underdog/container"

	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build image from Underdogfile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: underdog build -t <image-name>")
			return
		}
		imageName, _ := cmd.Flags().GetString("tag")
		stages, err := parser.ParseUnderdogfile("Underdogfile")
		if err != nil {
			fmt.Println("Parse error:", err)
			return
		}
		err = container.BuildImage(imageName, stages)
		if err != nil {
			fmt.Println("Build error:", err)
		}
	},
}

func init() {
	BuildCmd.Flags().StringP("tag", "t", "", "Image name")
}

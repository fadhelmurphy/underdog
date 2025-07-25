package cmd

import (
	"fmt"
	"underdog/runtime"

	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a container from image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: underdog run -p 3000:5000 <container-name> <image-name>")
			return
		}
		hostPort, _ := cmd.Flags().GetString("port")
		containerName := args[0]
		imageName := args[1]

		err := runtime.RunContainer(imageName, containerName, hostPort)
		if err != nil {
			fmt.Println("Run error:", err)
		}
	},
}

func init() {
	RunCmd.Flags().StringP("port", "p", "", "Port mapping")
}

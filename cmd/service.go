package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A collection of commands to manage systemd services",
	Long:  `A collection of commands to manage systemd services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please select the server type (--help for available options)")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

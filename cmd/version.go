package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runVersionCmd() {
	fmt.Println("Share cli v0.1.")
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Share cli",
	Long:  `All software has versions. This is share cli's`,
	Run: func(cmd *cobra.Command, args []string) {
		runVersionCmd()
	},
}

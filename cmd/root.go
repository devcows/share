package cmd

import (
	"fmt"

	"github.com/devcows/share/lib"

	"github.com/spf13/cobra"
)

var (
	appSettings lib.SettingsShare
	zipOption   bool
)

func init() {
	RootCmd.AddCommand(VersionCmd, ServerCmd, AddCmd, PsCmd, RmCmd)
	AddCmd.PersistentFlags().BoolVar(&zipOption, "zip", false, "Package with tar")
}

var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "Share is a cli to share quickly a file with http protocol",
	Long:  `Share is a cli to share quickly a file with http protocol with go. Complete documentation is available at https://github.com/devcows/share`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO use add, ps or rm")
	},
}

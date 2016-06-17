package cmd

import (
	"fmt"

	"../lib"

	"github.com/spf13/cobra"
)

var fileNameParam string
var portAPI int
var removeServerID int
var settings lib.SettingsShare

func init() {
	RootCmd.AddCommand(versionCmd, ServerCmd, AddCmd, PsCmd, RmCmd)
	AddCmd.PersistentFlags().StringVar(&fileNameParam, "file", "f", "File for share")
	ServerCmd.PersistentFlags().IntVarP(&portAPI, "port", "p", 7890, "Port APIREST")
	RmCmd.PersistentFlags().IntVarP(&removeServerID, "id", "i", -1, "Id remove server")
}

var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "Share is a cli to share quickly a file with http protocol",
	Long:  `Share is a cli to share quickly a file with http protocol with go. Complete documentation is available at https://github.com/devcows/share`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO use add, ps or rm")
	},
}

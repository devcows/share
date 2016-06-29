package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/devcows/share/api"
	"github.com/devcows/share/lib"
	"github.com/spf13/cobra"
)

func runRmCmd(uuid string, settings lib.SettingsShare) {
	serverEndPoint := fmt.Sprintf("http://%s:%v/rm/%v", settings.ShareDaemon.Host, settings.ShareDaemon.Port, uuid)
	resp, _ := http.Get(serverEndPoint)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	res := api.RmResponse{}
	json.Unmarshal([]byte(body), &res)

	if res.Status {
		fmt.Printf("Removed server{UUID: %s}\n", uuid)
	} else {
		fmt.Printf("Error: %s\n", res.ErrorMessage)
	}
}

var RmCmd = &cobra.Command{
	Use:   "rm UUID [UUID...]",
	Short: "Remove file or folder from server",
	Long:  `Remove file or folder from server`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Empty arguments")
		}

		if err := lib.InitSettings(lib.ConfigFile(), &appSettings); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		for _, arg := range args {
			runRmCmd(arg, appSettings)
		}
	},
}

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/devcows/share/api"
	"github.com/spf13/cobra"
)

func runRmCmd() {
	if len(removeServerUUID) == 0 {
		fmt.Println("Error: UUID empty!")
		return
	}

	url := fmt.Sprintf("http://localhost:7890/rm/%v", removeServerUUID)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	res := api.RmResponse{}
	json.Unmarshal([]byte(body), &res)

	if res.Status {
		fmt.Printf("Ok removed server with id = %s\n", removeServerUUID)
	} else {
		fmt.Printf("Error: %s\n", res.ErrorMessage)
	}
}

var RmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove file or folder from server",
	Long:  `Remove file or folder from server`,
	Run: func(cmd *cobra.Command, args []string) {
		runRmCmd()
	},
}

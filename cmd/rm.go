package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../api"
	"github.com/spf13/cobra"
)

var RmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove file or folder from server",
	Long:  `Remove file or folder from server`,
	Run: func(cmd *cobra.Command, args []string) {
		if removeServerID > 0 {
			url := fmt.Sprintf("http://localhost:7890/rm/%v", removeServerID)
			resp, _ := http.Get(url)
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			res := api.RmResponse{}
			json.Unmarshal([]byte(body), &res)

			if res.Status {
				fmt.Printf("Ok removed server with id = %v\n", removeServerID)
			} else {
				fmt.Printf("Error: %s\n", res.ErrorMessage)
			}
		}
	},
}

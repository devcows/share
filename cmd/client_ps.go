package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../lib"
	"github.com/spf13/cobra"
)

var PsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List files or folders from server",
	Long:  `List files or folders from server`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, _ := http.Get("http://localhost:7890/ps")
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		res := []lib.Server{}
		json.Unmarshal([]byte(body), &res)
		fmt.Println("Id\tFolder\tList Ips")
		for i := 0; i < len(res); i++ {
			fmt.Printf("%v\t%s\t%v\n", res[i].ID, res[i].Path, res[i].ListIps)
		}
	},
}

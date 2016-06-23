package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/devcows/share/api"
	"github.com/ryanuber/columnize"
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

		res := api.PsResponse{}
		json.Unmarshal([]byte(body), &res)

		if res.Status {
			lines := []string{
				"UUID | Folder | List Ips | CreatedAt",
			}

			for i := 0; i < len(res.Servers); i++ {
				server := res.Servers[i]
				line := fmt.Sprintf("%v|%s|%v|%v", server.UUID, server.Path, server.ListIps, server.CreatedAt)
				lines = append(lines, line)
			}

			result := columnize.SimpleFormat(lines)
			fmt.Println(result)
		} else {
			fmt.Printf("%s\n", res.ErrorMessage)
		}
	},
}

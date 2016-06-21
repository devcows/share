package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/devcows/share/lib"
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

		res := []lib.Server{}
		json.Unmarshal([]byte(body), &res)

		lines := []string{
			"UUID | Folder | List Ips | CreatedAt",
		}

		for i := 0; i < len(res); i++ {
			line := fmt.Sprintf("%v|%s|%v|%v", res[i].UUID, res[i].Path, res[i].ListIps, res[i].CreatedAt)
			lines = append(lines, line)
		}

		result := columnize.SimpleFormat(lines)
		fmt.Println(result)
	},
}

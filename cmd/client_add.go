package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/devcows/share/api"
	"github.com/spf13/cobra"
)

func copyClipboard(strToCopy string) {
	fmt.Printf("Copied to clipboard: %s", strToCopy)
	clipboard.WriteAll(strToCopy)
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file or folder to server",
	Long:  `Add file or folder to server`,
	Run: func(cmd *cobra.Command, args []string) {
		absFileNamePath, _ := filepath.Abs(fileNameParam)

		resp, _ := http.PostForm("http://localhost:7890/add", url.Values{"path": {absFileNamePath}})
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		res := api.AddResponse{}
		json.Unmarshal([]byte(body), &res)

		if res.Status {
			if len(res.ListIps) == 1 {
				copyClipboard(res.ListIps[0])
				return
			}

			fmt.Println("Choose option to copy to clipboard:")
			for i := 0; i < len(res.ListIps); i++ {
				fmt.Printf("%v) %s\n", i, res.ListIps[i])
			}

			var option int
			fmt.Scanf("%d", &option)
			copyClipboard(res.ListIps[option])
		} else {
			fmt.Printf("%s\n", res.ErrorMessage)
		}
	},
}

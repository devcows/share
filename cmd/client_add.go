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

func copyClipboard(strToCopy string) error {
	fmt.Printf("Copied to clipboard: %s", strToCopy)
	return clipboard.WriteAll(strToCopy)
}

func runAddCmd() error {
	absFileNamePath, err := filepath.Abs(fileNameParam)
	if err != nil {
		return err
	}

	resp, err := http.PostForm("http://localhost:7890/add", url.Values{"path": {absFileNamePath}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	res := api.AddResponse{}
	json.Unmarshal([]byte(body), &res)

	if res.Status {
		if len(res.ListIps) == 1 {
			return copyClipboard(res.ListIps[0])
		}

		fmt.Println("Choose option to copy to clipboard:")
		for i := 0; i < len(res.ListIps); i++ {
			fmt.Printf("%v) %s\n", i, res.ListIps[i])
		}

		var option int
		fmt.Scanf("%d", &option)
		return copyClipboard(res.ListIps[option])
	}

	fmt.Printf("%s\n", res.ErrorMessage)
	return nil
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file or folder to server",
	Long:  `Add file or folder to server`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAddCmd(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	},
}

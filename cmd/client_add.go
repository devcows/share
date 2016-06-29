package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/devcows/share/api"
	"github.com/devcows/share/lib"
	"github.com/spf13/cobra"
)

func copyClipboard(strToCopy string) error {
	fmt.Printf("Copied to clipboard: %s", strToCopy)
	return clipboard.WriteAll(strToCopy)
}

func selectOptionFrom(options []string) error {
	//TODO: check bad options
	if len(options) == 1 {
		return copyClipboard(options[0])
	}

	fmt.Println("Choose option to copy to clipboard:")
	for i, option := range options {
		fmt.Printf("%v) %s\n", i, option)
	}

	var option int
	fmt.Scanf("%d", &option)
	return copyClipboard(options[option])
}

func runAddCmd(filePath string, settings lib.SettingsShare) error {
	absFileNamePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	serverEndPoint := fmt.Sprintf("http://%s:%v/add", settings.ShareDaemon.Host, settings.ShareDaemon.Port)
	resp, err := http.PostForm(serverEndPoint, url.Values{"path": {absFileNamePath}, "zip": {strconv.FormatBool(zipOption)}})
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
		selectOptionFrom(res.Server.ListIps)
	}

	fmt.Printf("%s\n", res.ErrorMessage)
	return nil
}

var AddCmd = &cobra.Command{
	Use:   "add PATH [PATH...]",
	Short: "Add file or folder to server",
	Long:  `Add file or folder to server`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Empty arguments")
		}

		if err := lib.InitSettings(lib.ConfigFile(), &appSettings); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		for _, arg := range args {
			if err := runAddCmd(arg, appSettings); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
	},
}

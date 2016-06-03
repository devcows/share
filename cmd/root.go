package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/julienschmidt/httprouter"
	"github.com/prestonTao/upnp"
	"github.com/spf13/cobra"
)

var fileNameParam string

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.PersistentFlags().StringVar(&fileNameParam, "file", "f", "File for share")
}

// Ask the kernel for a free open port that is ready to use
func GetPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func upnpPort(port int) bool {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(port, port, "TCP"); err == nil {
		fmt.Println("success !")
		return true
		// remove port mapping in gatway
		// mapping.Reclaim()
	} else {
		fmt.Println("fail !")
		return false
	}
}

func localIps(port int) []string {
	listIps := []string{}

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			listIps = append(listIps, fmt.Sprintf("http://%s:%v", ipv4, port))
		}
	}

	return listIps
}

func publicIps(port int) []string {
	listIps := []string{}
	resp, err := http.Get("http://myexternalip.com/raw")
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			lines := strings.Split(string(body), "\n")
			listIps = append(listIps, fmt.Sprintf("http://%s:%v", lines[0], port))
		}
	}

	return listIps
}

var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "Share is a cli to share quickly a file with http protocol",
	Long:  `Share is a cli to share quickly a file with http protocol with go. Complete documentation is available at https://github.com/devcows/share`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main process
		port := GetPort()
		portOpened := upnpPort(port)

		listIps := []string{}
		if portOpened {
			listIps = append(localIps(port), publicIps(port)...)
		} else {
			listIps = localIps(port)
		}

		if len(listIps) > 0 {
			fmt.Println("Choose option to copy:")
			for i := 0; i < len(listIps); i++ {
				fmt.Printf("%v) %s\n", i, listIps[i])
			}

			var option int
			fmt.Scanf("%d", &option)
			clipboard.WriteAll(listIps[option])

			// Launch as daemon
			fmt.Printf("Server started at: 0.0.0.0:%v\n", port)
			router := httprouter.New()
			absFileNamePath, _ := filepath.Abs(fileNameParam)

			if info, err := os.Stat(absFileNamePath); err == nil && info.IsDir() {
				fmt.Println("Sharing folder: " + absFileNamePath)

				router.ServeFiles("/*filepath", http.Dir(absFileNamePath))
			} else {
				fmt.Println("Sharing file: " + absFileNamePath)

				fileServer := http.FileServer(http.Dir(path.Dir(fileNameParam)))
				router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
					req.URL.Path = "/" + filepath.Base(fileNameParam)
					fmt.Println("GET file: " + req.URL.Path)
					fileServer.ServeHTTP(w, req)
				})
			}

			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
		}
	},
}

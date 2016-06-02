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

func upnpPort(port int) {
	mapping := new(upnp.Upnp)
	if err := mapping.AddPortMapping(port, port, "TCP"); err == nil {
		fmt.Println("success !")
		// remove port mapping in gatway
		// mapping.Reclaim()
	} else {
		fmt.Println("fail !")
	}
}

func localIp(port int) {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Printf("http://%s:%v\n", ipv4, port)
		}
	}
}

func publicIp(port int) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			lines := strings.Split(string(body), "\n")
			fmt.Printf("http://%s:%v\n", lines[0], port)
		}
	}
}

var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "Share is a cli to share quickly a file with http protocol",
	Long:  `Share is a cli to share quickly a file with http protocol with go. Complete documentation is available at https://github.com/devcows/share`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main process
		port := GetPort()

		localIp(port)
		publicIp(port)

		// copy public link

		upnpPort(port)

		// Launch as daemon

		fmt.Printf("Server started at: 0.0.0.0:%v\n", port)
		absFileNamePath, _ := filepath.Abs(fileNameParam)
		fmt.Println("Sharing file: " + absFileNamePath)

		router := httprouter.New()
		fileServer := http.FileServer(http.Dir(path.Dir(fileNameParam)))
		router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			req.URL.Path = "/" + filepath.Base(fileNameParam)
			fmt.Println("GET file: " + req.URL.Path)
			fileServer.ServeHTTP(w, req)
		})

		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
	},
}

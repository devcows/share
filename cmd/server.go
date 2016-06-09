package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"../api"
	"../lib"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/tylerb/graceful"
)

var servers []api.Server
var nextIDServer int

func createHandler(filePathParam string) http.Handler {
	router := httprouter.New()
	absFilePath, _ := filepath.Abs(filePathParam)

	if info, err := os.Stat(absFilePath); err == nil && info.IsDir() {
		fmt.Println("Sharing folder: " + absFilePath)

		router.ServeFiles("/*filepath", http.Dir(absFilePath))
	} else {
		fmt.Println("Sharing file: " + absFilePath)

		fileServer := http.FileServer(http.Dir(path.Dir(filePathParam)))
		router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			req.URL.Path = "/" + filepath.Base(filePathParam)
			fmt.Println("GET file: " + req.URL.Path)
			fileServer.ServeHTTP(w, req)
		})
	}

	return router
}

func serverDaemon(port int, handler http.Handler) {
	fmt.Printf("Server started at: 0.0.0.0:%v\n", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler))
}

func serverDaemon2(port int, srv *graceful.Server) {
	fmt.Printf("Server started at: 0.0.0.0:%v\n", port)
	log.Fatal(srv.ListenAndServe())
}

func getAddPath(path string, c *gin.Context) {
	var msg api.AddResponse
	msg.Path = path

	//TODO: Add port parameter
	//strconv.ParseInt(c.DefaultPostForm("port", "-1"))
	port := lib.GetPort()

	//TODO: check port already taken and return error.
	if false {
		msg.Status = false
		msg.ErrorMessage = "The port is already in use."

		c.JSON(http.StatusOK, msg)
		return
	}
	msg.UpnpOpened = lib.OpenUpnpPort(port)

	msg.ListIps = []string{}
	if msg.UpnpOpened {
		msg.ListIps = append(lib.GetLocalIps(port), lib.GetPublicIps(port)...)
	} else {
		msg.ListIps = lib.GetLocalIps(port)
	}

	if len(msg.ListIps) == 0 {
		msg.Status = false
		msg.ErrorMessage = "No ips, check network connectivity."
		c.JSON(http.StatusOK, msg)
		return
	}

	handler := createHandler(msg.Path)
	//go serverDaemon(port, handler)
	srv := new(graceful.Server)
	srv.Timeout = 0
	srv.Server = new(http.Server)
	srv.Server.Addr = ":" + strconv.Itoa(port)
	srv.Server.Handler = handler

	go serverDaemon2(port, srv)

	var server api.Server
	server.Path = msg.Path
	server.ListIps = msg.ListIps
	server.Srv = srv
	server.ID = nextIDServer
	nextIDServer++
	servers = append(servers, server)

	msg.Status = true

	c.JSON(http.StatusOK, msg)
}

func mainServer() {
	fmt.Printf("Running APIREST 0.0.0.0:%v\n", portAPI)
	r := gin.Default()
	r.POST("/add", func(c *gin.Context) {
		path := c.PostForm("path")

		getAddPath(path, c)
	})

	r.GET("/rm/:id", func(c *gin.Context) {
		var msg api.RmResponse
		id := c.Param("id")

		for i := 0; i < len(servers); i++ {
			if id == strconv.Itoa(servers[i].ID) {
				servers[i].Srv.Stop(0)

				msg.Status = true
				c.JSON(http.StatusOK, msg)
				return
			}
		}

		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Server doesn't found with the id = %s", id)

		c.JSON(http.StatusOK, msg)
	})

	r.GET("/ps", func(c *gin.Context) {
		c.JSON(200, servers)
	})

	r.Run(fmt.Sprintf(":%v", portAPI))
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server APIREST",
	Long:  `Server APIREST`,
	Run: func(cmd *cobra.Command, args []string) {
		servers = []api.Server{}
		nextIDServer = 1
		mainServer()
	},
}

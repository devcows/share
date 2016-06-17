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

func getAddPath(path string, port int, c *gin.Context) {
	var msg api.AddResponse
	var server lib.Server

	msg.Path = path

	if port > 0 {
		server.Port = port
	} else {
		server.Port = lib.GetPort()
	}

	//TODO: check port already taken and return error.
	if false {
		msg.Status = false
		msg.ErrorMessage = "The port is already in use."

		c.JSON(http.StatusOK, msg)
		return
	}
	msg.UpnpOpened = lib.OpenUpnpPort(server.Port)

	msg.ListIps = []string{}
	if msg.UpnpOpened {
		msg.ListIps = append(lib.GetLocalIps(server.Port), lib.GetPublicIps(server.Port)...)
	} else {
		msg.ListIps = lib.GetLocalIps(server.Port)
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
	srv.Server.Addr = ":" + strconv.Itoa(server.Port)
	srv.Server.Handler = handler

	go serverDaemon2(server.Port, srv)

	server.Path = msg.Path
	server.ListIps = msg.ListIps
	server.Srv = srv

	fmt.Printf("iep4")
	_, err := lib.StoreServer(server)
	if err != nil {
		msg.Status = false
		msg.ErrorMessage = "TODO set error message"
		c.JSON(http.StatusOK, msg)
		return
	}

	fmt.Printf("iep3")

	msg.Status = true
	c.JSON(http.StatusOK, msg)
}

func mainServer() {
	r := gin.Default()
	r.POST("/add", func(c *gin.Context) {
		path := c.PostForm("path")
		port := -1 //strconv.ParseInt(c.DefaultPostForm("port", "-1"), "-1")

		getAddPath(path, port, c)
	})

	r.GET("/rm/:id", func(c *gin.Context) {
		var msg api.RmResponse
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			msg.Status = false
			msg.ErrorMessage = fmt.Sprint("Error bad request, cannot parse the parameter id to integer.")
			c.JSON(http.StatusOK, msg)
			return
		}

		err2 := lib.RemoveServer(id)
		if err2 != nil {
			msg.Status = false
			msg.ErrorMessage = fmt.Sprintf("Server doesn't found with the id = %v", id)
			c.JSON(http.StatusOK, msg)
			return
		}

		msg.Status = true
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/ps", func(c *gin.Context) {
		servers, err := lib.ListServers()

		if err != nil {
			c.JSON(http.StatusOK, []lib.Server{})
			return
		}

		c.JSON(http.StatusOK, servers)
	})

	r.Run(lib.ConfigServerEndPoint(settings))
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server APIREST",
	Long:  `Server APIREST`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if err = lib.InitSettings(&settings, portAPI); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if err = lib.InitDB(lib.ConfigFileSQLITE()); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		mainServer()
	},
}

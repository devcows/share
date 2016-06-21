package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devcows/share/api"
	"github.com/devcows/share/lib"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

func getAddPath(path string, port int, c *gin.Context) {
	server := lib.Server{UUID: uuid.NewV4().String(), Port: port, Path: path, CreatedAt: time.Now()}
	msg := api.AddResponse{UpnpOpened: false, ListIps: []string{}, Path: path}

	if port < 0 {
		server.Port = lib.RandomFreePort()
	}

	//TODO: check port already taken and return error.
	if false {
		msg.Status = false
		msg.ErrorMessage = "The port is already in use."

		c.JSON(http.StatusOK, msg)
		return
	}

	if settings.Daemon.EnableUpnp {
		msg.UpnpOpened = lib.OpenUpnpPort(server.Port)
	}

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
	server.ListIps = msg.ListIps

	lib.StartServer(&server)
	err := lib.StoreServer(server)
	if err != nil {
		msg.Status = false
		msg.ErrorMessage = err.Error()
		c.JSON(http.StatusOK, msg)
		return
	}

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

	r.GET("/rm/:uuid", func(c *gin.Context) {
		var msg api.RmResponse
		uuid := c.Param("uuid")

		err := lib.RemoveServer(uuid)
		if err != nil {
			msg.Status = false
			msg.ErrorMessage = fmt.Sprintf("Server doesn't found with the uuid = %s", uuid)
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

func overwriteSettings() {
	//overwrite settings
	if portAPI > 0 {
		settings.Daemon.Port = portAPI
	}
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server APIREST",
	Long:  `Server APIREST`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if err = lib.InitSettings(lib.ConfigFile(), &settings); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if err = lib.InitDB(settings); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		overwriteSettings()
		mainServer()
	},
}

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/devcows/share/api"
	"github.com/devcows/share/lib"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

func prepareAndRunServer(server *lib.Server) (bool, error) {
	upnpOpened := false

	if server.Port < 0 {
		server.Port = lib.RandomFreePort()
	}

	/*
		TODO: check port already taken and return error.
		if false {
			msg.Status = false
			msg.ErrorMessage = "The port is already in use."

			c.JSON(http.StatusOK, msg)
			return
		}
	*/

	if settings.Daemon.EnableUpnp {
		upnpOpened = lib.OpenUpnpPort(server.Port)
	}

	if upnpOpened {
		server.ListIps = append(lib.GetLocalIps(server.Port), lib.GetPublicIps(server.Port)...)
	} else {
		server.ListIps = lib.GetLocalIps(server.Port)
	}

	/*
		TODO: check list ips and if empty return error.
			if len(msg.ListIps) == 0 {
				msg.Status = false
				msg.ErrorMessage = "No ips, check network connectivity."
				c.JSON(http.StatusOK, msg)
				return
			}
	*/
	lib.StartServer(server)

	return upnpOpened, nil
}

func processAddServer(server lib.Server, c *gin.Context) {
	upnpOpened, err := prepareAndRunServer(&server)
	msg := api.AddResponse{UpnpOpened: upnpOpened, ListIps: server.ListIps, Path: server.Path}

	if err != nil {
		msg.Status = false
		msg.ErrorMessage = err.Error()
		c.JSON(http.StatusOK, msg)
		return
	}

	if err = lib.StoreServer(server); err != nil {
		msg.Status = false
		msg.ErrorMessage = err.Error()
		c.JSON(http.StatusOK, msg)
		return
	}

	msg.Status = true
	c.JSON(http.StatusOK, msg)
}

func processRmServer(uuid string, c *gin.Context) {
	var msg api.RmResponse

	if err := lib.RemoveServer(uuid); err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Server doesn't found with the uuid = %s", uuid)
		c.JSON(http.StatusOK, msg)
		return
	}

	msg.Status = true
	c.JSON(http.StatusOK, msg)
}

func processPsServers(c *gin.Context) {
	servers, err := lib.ListServers()

	if err != nil {
		c.JSON(http.StatusOK, []lib.Server{})
		return
	}

	c.JSON(http.StatusOK, servers)
}

func mainServer() {
	r := gin.Default()
	r.POST("/add", func(c *gin.Context) {
		path := c.PostForm("path")
		port := -1 //strconv.ParseInt(c.DefaultPostForm("port", "-1"), "-1")

		server := lib.Server{UUID: uuid.NewV4().String(), Port: port, Path: path, CreatedAt: time.Now()}
		processAddServer(server, c)
	})

	r.GET("/rm/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		processRmServer(uuid, c)
	})

	r.GET("/ps", func(c *gin.Context) {
		processPsServers(c)
	})

	r.Run(lib.ConfigServerEndPoint(settings))
}

func loadInitialServers() error {
	servers, err := lib.ListServers()

	if err != nil {
		return err
	}

	for i := 0; i < len(servers); i++ {
		prepareAndRunServer(&servers[i])
		//TODO: update values
	}

	//TODO: if one server fail return error
	return nil
}

func overwriteSettings() {
	//overwrite settings
	if portAPI > 0 {
		settings.Daemon.Port = portAPI
	}

	if settings.Mode == "release" {
		log.SetFormatter(&log.JSONFormatter{})

		// Only log the warning severity or above.
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)

		//TODO: Output to file instead of stdout.
		//log.SetOutput(os.Stderr)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)

		// Output to stderr instead of stdout, could also be a file.
		log.SetOutput(os.Stderr)
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
		loadInitialServers()
		mainServer()
	},
}

package cmd

import (
	"errors"
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
	var err error
	upnpOpened := false

	if server.Port < 0 {
		server.Port, err = lib.RandomFreePort()

		if err != nil {
			return upnpOpened, err
		}
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

	if appSettings.Daemon.EnableUpnp {
		upnpOpened = lib.OpenUpnpPort(server.Port)
	}

	if upnpOpened {
		server.ListIps = append(lib.GetLocalIps(server.Port), lib.GetPublicIps(server.Port)...)
	} else {
		server.ListIps = lib.GetLocalIps(server.Port)
	}

	if len(server.ListIps) == 0 {
		return upnpOpened, errors.New("No ips, check network connectivity.")
	}

	err = lib.StartServer(server)
	return upnpOpened, err
}

func processAddServer(server lib.Server, c *gin.Context) {
	upnpOpened, err := prepareAndRunServer(&server)
	msg := api.AddResponse{Status: true, UpnpOpened: upnpOpened, ListIps: server.ListIps, Path: server.Path}

	if err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err)
		c.JSON(http.StatusOK, msg)
		return
	}

	if err = lib.StoreServer(server); err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err)
		c.JSON(http.StatusOK, msg)
		return
	}

	c.JSON(http.StatusOK, msg)
}

func processRmServer(uuid string, c *gin.Context) {
	msg := api.RmResponse{Status: true}

	if err := lib.RemoveServer(uuid); err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err) //fmt.Sprintf("Server doesn't found with the uuid = %s", uuid)
		c.JSON(http.StatusOK, msg)
		return
	}

	c.JSON(http.StatusOK, msg)
}

func processPsServers(c *gin.Context) {
	msg := api.PsResponse{Status: true}
	var err error

	msg.Servers, err = lib.ListServers()
	if err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err)

		c.JSON(http.StatusOK, msg)
		return
	}

	c.JSON(http.StatusOK, msg)
}

func mainServer(settings lib.SettingsShare) {
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

func overwriteSettings(settings lib.SettingsShare) {
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

func runServerCmd(configFile string, settings *lib.SettingsShare) error {
	if err := lib.InitSettings(configFile, settings); err != nil {
		return err
	}

	if err := lib.InitDB(*settings); err != nil {
		return err
	}

	overwriteSettings(*settings)
	loadInitialServers()

	return nil
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server APIREST",
	Long:  `Server APIREST`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runServerCmd(lib.ConfigFile(), &appSettings); err != nil {
			log.Error(fmt.Sprintf("Error: %s", err))
			os.Exit(-1)
		}

		mainServer(appSettings)
	},
}

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

func validateServer(server lib.Server) error {
	//TODO: validate flags and properties
	//TODO: file or folder exists
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func newServerParams(c *gin.Context, path string) lib.Server {
	flags := []string{}
	if c.PostForm("zip") == "true" {
		flags = append(flags, "zip")
	}

	return lib.Server{UUID: uuid.NewV4().String(), Source: path, Path: path, CreatedAt: time.Now(), Flags: flags}
}

func processAddServer(c *gin.Context) {
	var err error

	path := c.PostForm("path")
	msg := api.AddResponse{Status: true, UpnpOpened: false}
	msg.Server, err = lib.SearchServerByPath(path)
	if err == nil {
		msg.Server.ListIps = lib.GetServerIps(msg.UpnpOpened, appSettings.FileServerDaemon.Port, msg.Server.UUID)
		c.JSON(http.StatusOK, msg)
		return
	}

	msg.Server = newServerParams(c, path)
	if err = validateServer(msg.Server); err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err)
		c.JSON(http.StatusOK, msg)
		return
	}

	if stringInSlice("zip", msg.Server.Flags) {
		outPutFilePath := lib.CompressedFilePath() + string(os.PathSeparator) + msg.Server.UUID + ".zip"

		if err := lib.CompressFile(msg.Server.Path, outPutFilePath); err != nil {
			msg.Status = false
			msg.ErrorMessage = fmt.Sprintf("Error: %s", err)
			c.JSON(http.StatusOK, msg)
			return
		}

		msg.Server.Path = outPutFilePath
	}

	if appSettings.ShareDaemon.EnableUpnp {
		msg.UpnpOpened = lib.OpenUpnpPort(appSettings.FileServerDaemon.Port)
	}

	msg.Server.ListIps = lib.GetServerIps(msg.UpnpOpened, appSettings.FileServerDaemon.Port, msg.Server.UUID)
	if len(msg.Server.ListIps) == 0 {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: No ips, check network connectivity.")
		c.JSON(http.StatusOK, msg)
		return
	}

	if err = lib.StoreServer(msg.Server); err != nil {
		msg.Status = false
		msg.ErrorMessage = fmt.Sprintf("Error: %s", err)
		c.JSON(http.StatusOK, msg)
		return
	}

	c.JSON(http.StatusOK, msg)
}

// TODO: rm zip file
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

	upnpOpened := false
	if appSettings.ShareDaemon.EnableUpnp {
		upnpOpened = lib.OpenUpnpPort(appSettings.FileServerDaemon.Port)
	}

	for i := 0; i < len(msg.Servers); i++ {
		msg.Servers[i].ListIps = lib.GetServerIps(upnpOpened, appSettings.FileServerDaemon.Port, msg.Servers[i].UUID)
	}

	fmt.Printf("%v\n", msg.Servers)

	c.JSON(http.StatusOK, msg)
}

func runShareServer(settings lib.SettingsShare) {
	r := gin.Default()
	r.POST("/add", func(c *gin.Context) {
		processAddServer(c)
	})

	r.GET("/rm/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		processRmServer(uuid, c)
	})

	r.GET("/ps", func(c *gin.Context) {
		processPsServers(c)
	})

	r.Run(lib.ConfigServerEndPoint(settings))
	log.WithFields(log.Fields{"ip": settings.ShareDaemon.Host, "port": settings.ShareDaemon.Port}).Info("Share server started.")
}

func overwriteSettings(settings lib.SettingsShare) error {
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

	return nil
}

func runServerCmd(configFile string, settings *lib.SettingsShare) error {
	if err := lib.InitSettings(configFile, settings); err != nil {
		return err
	}

	if err := lib.InitDB(*settings); err != nil {
		return err
	}

	if err := overwriteSettings(*settings); err != nil {
		return err
	}

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

		log.Info("Running servers...")
		go lib.RunFileServer(appSettings)
		runShareServer(appSettings) //TODO: launch as go routine and wait
	},
}

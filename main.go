package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/devcows/share/cmd"
	"github.com/gin-gonic/gin"
)

func main() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	if os.Getenv("SHARE_MODE") == "release" {
		// Only log the warning severity or above.
		log.SetLevel(log.WarnLevel)
		gin.SetMode(gin.ReleaseMode)

		//TODO: Output to file instead of stdout.
		//log.SetOutput(os.Stderr)
	} else {
		log.SetLevel(log.DebugLevel)

		// Output to stderr instead of stdout, could also be a file.
		log.SetOutput(os.Stderr)
	}

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

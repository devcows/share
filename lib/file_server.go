package lib

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

func FileServerGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filePathURL := ps.ByName("filePathURL")
	//TODO: check empty filepath

	paths := strings.Split(filePathURL, "/")
	//TODO: check length == 2 minimum

	server, err := SearchServerByUUID(paths[1])
	if err != nil {
		// TODO: respond not found
		fmt.Fprintf(w, "no server %v !\n", err)
		return
	}

	if info, err := os.Stat(server.Path); err == nil && info.IsDir() {
		log.Debug("Returning to folder: " + server.Path)

		if len(paths) == 2 {
			http.ServeFile(w, r, server.Path)
		} else {
			filePath := server.Path + string(os.PathSeparator) + strings.Join(paths[2:len(paths)], string(os.PathSeparator))
			log.Debug("Returning to file: " + filePath)
			http.ServeFile(w, r, filePath)
		}
	} else {
		log.Debug("Returning to file: " + server.Path)
		http.ServeFile(w, r, server.Path)
	}
}

func createHandler() http.Handler {
	router := httprouter.New()
	router.GET("/*filePathURL", FileServerGET)

	return router
}

func RunFileServer(settings SettingsShare) {
	handler := createHandler()
	log.Fatal(http.ListenAndServe(ConfigFileServerEndPoint(settings), handler))
	log.WithFields(log.Fields{"ip": settings.FileServerDaemon.Host, "port": settings.FileServerDaemon.Port}).Info("File server started.")
}

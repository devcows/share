package lib

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
)

func CreateHandler(filePathParam string) http.Handler {
	router := httprouter.New()
	absFilePath, _ := filepath.Abs(filePathParam)

	if info, err := os.Stat(absFilePath); err == nil && info.IsDir() {
		log.Info("Creating handler for folder: " + absFilePath)

		router.ServeFiles("/*filepath", http.Dir(absFilePath))
	} else {
		log.Info("Creating handler for file: " + absFilePath)

		fileServer := http.FileServer(http.Dir(path.Dir(filePathParam)))
		router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			req.URL.Path = "/" + filepath.Base(filePathParam)
			log.Debug("GET file: " + req.URL.Path)
			fileServer.ServeHTTP(w, req)
		})
	}

	return router
}

func ServerDaemon(port int, handler http.Handler) {
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler))
}

func ServerDaemon2(port int, srv *graceful.Server) {
	log.Fatal(srv.ListenAndServe())
}

func StartServer(server *Server) {
	log.WithFields(log.Fields{"ip": "0.0.0.0", "port": server.Port, "path": server.Path}).Info("Server started.")
	handler := CreateHandler(server.Path)
	//go serverDaemon(port, handler)
	srv := new(graceful.Server)
	srv.Timeout = 0
	srv.Server = new(http.Server)
	srv.Server.Addr = ":" + strconv.Itoa(server.Port)
	srv.Server.Handler = handler

	go ServerDaemon2(server.Port, srv)
	server.Srv = srv
}

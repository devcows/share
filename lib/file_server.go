package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func dirList(w http.ResponseWriter, folderPath string, prefix string) error {
	f, err := os.Open(folderPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dirs, err := f.Readdir(-1)
	if err != nil {
		// TODO: log err.Error() to the Server.ErrorLog, once it's possible
		// for a handler to get at its Server via the ResponseWriter. See
		// Issue 12438.
		//log.Error(w, "Error reading directory", http.StatusInternalServerError)
		return err
	}
	sort.Sort(byName(dirs))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", prefix+"/"+url.String(), name)
	}
	fmt.Fprintf(w, "</pre>\n")

	return nil
}

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
		if len(paths) == 2 {
			log.Debug("Returning to folder: " + server.Path)
			dirList(w, server.Path, server.UUID)
		} else {
			filePath := server.Path + string(os.PathSeparator) + strings.Join(paths[2:len(paths)], string(os.PathSeparator))

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

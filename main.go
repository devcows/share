package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("file_name"))
}

func main() {
	router := httprouter.New()
	router.GET("/:file_name", ServeFile)

	log.Fatal(http.ListenAndServe(":8080", router))
}

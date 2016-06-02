package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
)

func ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("file_name"))
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
}

var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "Share is a cli to share quickly a file with http protocol",
	Long: `Share is a cli to share quickly a file with http protocol
								with go. Complete documentation is available at https://github.com/devcows/share`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main process
		fmt.Println("Server started at: 0.0.0.0:8080")

		// GET file_name
		// GET relative path

		// print links
		// copy public link
		// configure upnp
		// Launch as daemon

		router := httprouter.New()
		router.GET("/:file_name", ServeFile)

		log.Fatal(http.ListenAndServe(":8080", router))
	},
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/soulshake/pawd/controllers"
)

func main() {
	print("RemoteServer started \n")
	// Replace "/" with e.g. "/volume" for a path
	port := os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}
	log.Println("listening=" + port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), controllers.NewRouter())
}

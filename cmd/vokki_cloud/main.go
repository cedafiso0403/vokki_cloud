package main

import (
	"log"
	"net/http"
	"vokki_cloud/internal/router"

	_ "github.com/lib/pq"
)

func main() {

	r := router.SetupRouter()

	log.Fatal(http.ListenAndServe(":8000", r))
}

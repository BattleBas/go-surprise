package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BattleBas/go-surprise/pkg/restapi"
)

func main() {
	port := os.Getenv("PORT")
	r := restapi.Handler()
	log.Fatal(http.ListenAndServe(":"+port, r))
}

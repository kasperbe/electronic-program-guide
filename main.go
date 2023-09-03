package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kasperbe/electronic-program-guide/router"
)

const HTTP_PORT = "8080"

func main() {
	fmt.Println("initializing electronic program service")

	mux := http.NewServeMux()
	mux.HandleFunc("/", router.TranslateEPG)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", HTTP_PORT),
		Handler: mux,
	}

	fmt.Printf("listening on port %s\n", HTTP_PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

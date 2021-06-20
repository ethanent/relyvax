package main

import (
	"embed"
	"fmt"
	"github.com/ethanent/calronavaxverif"
	"golang.org/x/crypto/acme/autocert"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var webFS embed.FS

func main() {
	// Get CA key set

	var err error

	caKeySet, err = calronavaxverif.FetchKeys("https://myvaccinerecord.cdph.ca.gov/creds/.well-known/jwks.json")

	if err != nil {
		panic(err)
	}

	// Set up server

	s := http.NewServeMux()

	rt, err := fs.Sub(webFS, "static")

	if err != nil {
		panic(err)
	}

	s.Handle("/", http.FileServer(http.FS(rt)))

	s.HandleFunc("/check", checkHandler)

	fmt.Println("Listening...")

	certMgr := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("relyvax.com", "www.relyvax.com"),
		Cache:      autocert.DirCache("cert-cache"),
	}

	server := &http.Server{
		Addr:      ":https",
		TLSConfig: certMgr.TLSConfig(),
		Handler:   s,
	}

	// Start autocert HTTP server

	go func() {
		useH := certMgr.HTTPHandler(nil)

		err := http.ListenAndServe(":http", useH)

		log.Fatal(err)
	}()

	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

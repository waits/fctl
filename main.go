// Command fctl hosts a web page for administering a Factorio server.
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

var (
	host    = flag.String("host", "", "server hostname")
	dir     = flag.String("certs", "", "directory to store certificates in")
	devMode = flag.Bool("development", false, "run in development mode")
)

func main() {
	flag.Parse()

	t := LoadTemplates("templates/")
	env := &Env{Templates: t, TemplatePath: "templates/"}
	http.Handle("/", Handler{Env: env, Fn: Static})

	if *devMode {
		log.Printf("Starting server at http://0.0.0.0:8080\n")
		http.ListenAndServe("0.0.0.0:8080", nil)
		return
	}

	redir := "https://" + *host
	log.Printf("Starting server at http://" + *host)
	go http.ListenAndServe("0.0.0.0:80", http.RedirectHandler(redir, 301))

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(*host),
		Cache:      autocert.DirCache(*dir),
	}
	s := &http.Server{
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}
	log.Printf("Starting server at https://" + *host)
	log.Fatal(s.ListenAndServeTLS("", ""))
}

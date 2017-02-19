package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// Env holds a map of cached templates.
type Env struct {
	Templates    map[string]*template.Template
	TemplatePath string
}

// Handler wraps a route handler with an Env.
type Handler struct {
	*Env
	Fn func(*Env, http.ResponseWriter, *http.Request) (int, error)
}

// ServeHTTP is called on a reqHandler by net/http; satisfies http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("strict-transport-security", "max-age=31536000")

	status, err := h.Fn(h.Env, w, r)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusBadRequest:
			http.Error(w, err.Error(), status)
		default:
			status = http.StatusInternalServerError
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	log.Printf("%s %s %s %d", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.Path, status)
}

// Static responds to static routes not covered by another handler.
func Static(env *Env, w http.ResponseWriter, r *http.Request) (int, error) {
	switch r.URL.Path {
	case "/":
		return http.StatusOK, RenderTemplate(env, w, "home", gameInfo())
	default:
		return http.StatusNotFound, errors.New("handler: page not found")
	}
}

func gameInfo() string {
	out, err := exec.Command("fctl", "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/GeertJohan/go.rice"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"
)

// Default is the default handler. It searches for the URL and if it doesn't exist or there's an error, it redirects to
// the main page
func Default(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]

		if path == "" {
			if r.Method == http.MethodPost {
				addURL(db, w, r)
				return
			}

			mainPage(w)
			return
		}

		var toURL string
		var err error
		if toURL, err = db.ReadURL(path); err != nil {
			errorPage(err, w)
			return
		}

		if len(strings.Split(toURL, "://")) == 1 {
			toURL = "http://" + toURL
		}

		http.Redirect(w, r, toURL, http.StatusFound)
	}
}

// mainPage renders the main page
func mainPage(w io.Writer) {
	if _, err := fmt.Fprint(w, rice.MustFindBox("static").MustString("index.html")); err != nil {
		log.Printf("error writting the HTTP response at mainPage: %v", err)
	}
}

// errorPage renders an error page with the error provided
func errorPage(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	if _, writeErr := fmt.Fprintf(w, "There was an error processing your request: %v\n", err); writeErr != nil {
		log.Printf("error writting the HTTP response at errorPage: %v", writeErr)
	}
}

// addURL adds a new URL to the DB
func addURL(db *db.DB, w http.ResponseWriter, r *http.Request) {
	if err := db.AddURL(r.FormValue("shortURL"), r.FormValue("longURL")); err != nil {
		errorPage(err, w)
		return
	}

	url := r.FormValue("longURL")

	if len(strings.Split(url, "://")) == 1 {
		url = "http://" + url
	}

	http.Redirect(w, r, url, http.StatusFound)
}

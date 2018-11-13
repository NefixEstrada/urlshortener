package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"
	bolt "go.etcd.io/bbolt"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() http.Header {
	return nil
}
func (mockResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("testing error")
}
func (mockResponseWriter) WriteHeader(int) {}

// Should work as expected
func TestMainPage(t *testing.T) {
	w := httptest.NewRecorder()

	expected, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		t.Errorf("unexpected error when reading the index.html file: %v", err)
	}

	mainPage(w)

	if w.Code != http.StatusOK {
		t.Errorf("expecting %d, but got %d", http.StatusOK, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("expecting %v, but got %v", expected, w.Body.Bytes())
	}
}

// Should fail when writting the response
func TestMainPageErr(t *testing.T) {
	mainPage(mockResponseWriter{})
}

// Should work as expected
func TestErrorPage(t *testing.T) {
	w := httptest.NewRecorder()

	expected := []byte("There was an error processing your request: testing error\n")

	errorPage(errors.New("testing error"), w)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("expecting %v, but got %v", expected, w.Body.Bytes())
	}
}

// Should fail when writting the response
func TestErrorPageErr(t *testing.T) {
	errorPage(fmt.Errorf("%b", []byte("a")), mockResponseWriter{})
}

// Should work as expected
func TestAddURL(t *testing.T) {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		t.Fatalf("error creating the testing DB: %v", err)

	}

	db := &db.DB{
		DB: boltDB,
	}
	err = db.Initialize()
	if err != nil {
		t.Fatalf("error initializing the DB: %v", err)
	}

	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("shortURL", "go")
	form.Add("longURL", "https://golang.org")

	r, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("error preparing the HTTP request: %v", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var expected []byte

	addURL(db, w, r)

	if w.Code != http.StatusFound {
		t.Errorf("expecting %d, but got %d", http.StatusFound, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("expecting %s, but got %s", expected, w.Body.Bytes())
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// There should be an error when adding the URL to the DB
func TestAddURLErr(t *testing.T) {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		t.Fatalf("error creating the testing DB: %v", err)

	}

	db := &db.DB{
		DB: boltDB,
	}
	err = db.Initialize()
	if err != nil {
		t.Fatalf("error initializing the DB: %v", err)
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "/", bytes.NewReader([]byte("")))
	if err != nil {
		t.Fatalf("error preparing the HTTP request: %v", err)
	}

	expected := []byte("There was an error processing your request: the short URL can't be empty\n")

	addURL(db, w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("expecting %s, but got %s", expected, w.Body.Bytes())
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

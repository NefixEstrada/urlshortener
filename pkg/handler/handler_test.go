package handler_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"
	"gitea.nefixestrada.com/nefix/urlshortener/pkg/handler"
	bolt "go.etcd.io/bbolt"
)

// Should return the main page
func TestDefaultHandler(t *testing.T) {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		t.Fatalf("error creating the testing DB: %v", err)
	}

	db := &db.DB{
		DB: boltDB,
	}

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	w := httptest.NewRecorder()

	expected, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		t.Errorf("unexpected error when reading the index.html file: %v", err)
	}

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expecting %d, but got %d", http.StatusOK, w.Code)
	}

	if !bytes.Equal(w.Body.Bytes(), expected) {
		t.Errorf("expecting %s, but got %s", expected, w.Body.Bytes())
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// Should redirect to the requested page
func TestDefaultHandlerRedirect(t *testing.T) {
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

	if err = boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))

		return b.Put([]byte("test"), []byte("https://nefixestrada.com"))
	}); err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	w := httptest.NewRecorder()

	expected := "https://nefixestrada.com"

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("expecting %d, but got %d", http.StatusFound, w.Code)
	}

	if expected != w.Header().Get("Location") {
		t.Errorf("expecting %s, but got %s", expected, w.Header().Get("Location"))
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// Should append http:// when redirecting to an url that doesn't contain http://
func TestDefaultHandlerRedirectNoHttp(t *testing.T) {
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

	if err = boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))

		return b.Put([]byte("test"), []byte("nefixestrada.com"))
	}); err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	w := httptest.NewRecorder()

	expected := "http://nefixestrada.com"

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("expecting %d, but got %d", http.StatusFound, w.Code)
	}

	if expected != w.Header().Get("Location") {
		t.Errorf("expecting %s, but got %s", expected, w.Header().Get("Location"))
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// Should return a not found page
func TestDefaultHandlerNotFound(t *testing.T) {
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

	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}

	w := httptest.NewRecorder()

	expected := []byte("There was an error processing your request: the shortened URL wasn't found in the DB\n")

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expecting %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if !bytes.Equal(expected, w.Body.Bytes()) {
		t.Errorf("expecting %s, but got %s", expected, w.Body.Bytes())
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// Should redirect to the new URL when adding an URL
func TestDefaultHandlerNew(t *testing.T) {
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

	r, err := http.NewRequest("POST", "/", strings.NewReader("shortURL=test&longURL=https://nefixestrada.com"))
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	expected := "https://nefixestrada.com"

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("expecting %d, but got %d", http.StatusFound, w.Code)
	}

	if expected != w.Header().Get("Location") {
		t.Errorf("expecting %s, but got %s", expected, w.Header().Get("Location"))
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// Should append http:// when redirecting to an url that doesn't contain http://
func TestDefaultHandlerNewNoHttp(t *testing.T) {
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

	r, err := http.NewRequest("POST", "/", strings.NewReader("shortURL=test&longURL=nefixestrada.com"))
	if err != nil {
		t.Fatalf("error preparing the test: %v", err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	expected := "http://nefixestrada.com"

	handler := handler.Default(db)
	handler(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("expecting %d, but got %d", http.StatusFound, w.Code)
	}

	if expected != w.Header().Get("Location") {
		t.Errorf("expecting %s, but got %s", expected, w.Header().Get("Location"))
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

package db_test

import (
	"os"
	"testing"

	"gitea.nefixestrada.com/nefix/urlshortener/pkg/db"

	bolt "go.etcd.io/bbolt"
)

var tests = []struct {
	shortURL string
	longURL  string
}{
	{
		shortURL: "git",
		longURL:  "https://gitea.nefixestrada.com",
	},
}

// Should work as expected
func TestReadURL(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			var b *bolt.Bucket
			b, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			if err = b.Put([]byte(tt.shortURL), []byte(tt.longURL)); err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		rsp, err := db.ReadURL(tt.shortURL)
		if err != nil {
			t.Errorf("unexpected error when reading the URL in the DB: %v", err)
		}

		if rsp != tt.longURL {
			t.Errorf("expecting %s, but got %s", tt.shortURL, tt.longURL)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// Should return a not found error
func TestReadURLNotFound(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "the shortened URL wasn't found in the DB"

		rsp, err := db.ReadURL(tt.shortURL)
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if rsp != "" {
			t.Errorf("expecting %s, but got %s", "", rsp)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// Should work as expected
func TestAddURL(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		err = db.AddURL(tt.shortURL, tt.longURL)
		if err != nil {
			t.Errorf("unexpected error adding the URL: %v", err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// The short URL can't be empty
func TestAddURLShortNoEmpty(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "the short URL can't be empty"

		err = db.AddURL("", tt.longURL)
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// The long URL can't be empty
func TestAddURLLongNoEmpty(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "the long URL can't be empty"

		err = db.AddURL(tt.shortURL, "")
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// The long URL needs to be an URL
func TestAddURLLongIsURL(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "the long URL needs to be a valid URL"

		err = db.AddURL(tt.shortURL, "https://notanurl!")
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// The bucket 'urls' doesn't exist
func TestAddURLBucketNotExist(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "the bucket urls doesn't exist"

		err = db.AddURL(tt.shortURL, tt.longURL)
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// The shortened URL already exists
func TestAddURLAlreadyExists(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			var b *bolt.Bucket
			b, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			if err = b.Put([]byte(tt.shortURL), []byte(tt.longURL)); err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "there's already an shortened URL with that URL"

		err = db.AddURL(tt.shortURL, tt.longURL)
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// There should be an error when inserting the URL in the DB
func TestAddURLErrInserting(t *testing.T) {
	for _, tt := range tests {
		boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
		if err != nil {
			t.Fatalf("error creating the testing DB: %v", err)
		}

		if err = boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucket([]byte("urls"))
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Fatalf("error inserting test data to the DB: %v", err)
		}

		db := db.DB{
			DB: boltDB,
		}

		expectedErr := "key too large"

		longKey := make([]byte, bolt.MaxKeySize+1)

		err = db.AddURL(string(longKey), tt.longURL)
		if err.Error() != expectedErr {
			t.Errorf("expecting %s, but got %v", expectedErr, err)
		}

		if err := os.Remove("urlshortener.db"); err != nil {
			t.Fatalf("error finishing the test: %v", err)
		}
	}
}

// Should work as expected
func TestInitialize(t *testing.T) {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		t.Fatalf("error creating the testing DB: %v", err)
	}

	db := db.DB{
		DB: boltDB,
	}

	err = db.Initialize()
	if err != nil {
		t.Errorf("unexpected error initializing the DB: %v", err)
	}

	if err = boltDB.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucket([]byte("urls"))

		return err
	}); err != bolt.ErrBucketExists {
		t.Errorf("expecting %v, but got %v", bolt.ErrBucketExists, err)
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

// There should be an error creating the bucket
func TestInitializeErrBucket(t *testing.T) {
	boltDB, err := bolt.Open("urlshortener.db", 0600, nil)
	if err != nil {
		t.Fatalf("error creating the testing DB: %v", err)
	}

	db := db.DB{
		DB: boltDB,
	}

	db.DB.Close()

	expectedErr := "database not open"

	err = db.Initialize()
	if err.Error() != expectedErr {
		t.Errorf("expecting %s, but got %v", expectedErr, err)
	}

	if err := os.Remove("urlshortener.db"); err != nil {
		t.Fatalf("error finishing the test: %v", err)
	}
}

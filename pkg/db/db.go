package db

import (
	"errors"

	"github.com/asaskevich/govalidator"

	bolt "go.etcd.io/bbolt"
)

// DB is the struct that contains the connection with the Bold DB
type DB struct {
	DB *bolt.DB
}

// ReadURL reads a shortened URL from the DB and returns the target URL for it
func (d *DB) ReadURL(shortURL string) (fullURL string, err error) {
	if err := d.DB.View(func(tx *bolt.Tx) error {
		fullURL = string(tx.Bucket([]byte("urls")).Get([]byte(shortURL)))

		if fullURL == "" {
			return errors.New("the shortened URL wasn't found in the DB")
		}

		return nil
	}); err != nil {
		return "", err
	}

	return fullURL, nil
}

// AddURL adds a new URL to the DB
func (d *DB) AddURL(shortURL string, longURL string) error {
	if shortURL == "" {
		return errors.New("the short URL can't be empty")
	}

	if longURL == "" {
		return errors.New("the long URL can't be empty")
	}

	if !govalidator.IsURL(longURL) {
		return errors.New("the long URL needs to be a valid URL")
	}

	return d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("urls"))
		if b == nil {
			return errors.New("the bucket urls doesn't exist")
		}

		if content := b.Get([]byte(shortURL)); content != nil {
			return errors.New("there's already an shortened URL with that URL")
		}

		if err := b.Put([]byte(shortURL), []byte(longURL)); err != nil {
			return err
		}

		return nil
	})
}

// Initialize creates the required bucket
func (d *DB) Initialize() error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("urls"))

		return err
	})
}

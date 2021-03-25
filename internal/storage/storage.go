package storage

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/vimoppa/turl.to/internal/config"
)

// Accessor is the storage accessor contract.
type Accessor interface {
	WriteOnce(shortURL string, longURL string) error
	ReadOne(shortURL string) (string, error)
	ReadAll() ([]string, error)
	LookUp(longURL string) bool
}

// Store ...
type Store struct {
	file string
	mu   sync.Mutex
}

// WriteOnce writes a new record to the store.
func (s *Store) WriteOnce(shortURL string, longURL string) error {
	input := shortURL + " " + longURL

	errChan := make(chan error, 1)
	go func() {
		s.mu.Lock()
		file, err := os.Open(s.file)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()
		defer s.mu.Unlock()

		if err := ioutil.WriteFile(s.file, []byte(input), fs.ModeAppend); err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	return <-errChan
}

// ReadOne a single record from the store.
func (s *Store) ReadOne(URL string) (string, error) {
	strChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		s.mu.Lock()
		file, err := os.Open(s.file)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()
		defer s.mu.Unlock()

		var (
			longURL  string
			shortURL string
		)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text := scanner.Text()
			urls := strings.Split(text, " ")
			shortURL, longURL = urls[0], urls[1]

			if shortURL == URL {
				break
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}

		strChan <- longURL
	}()

	select {
	case out := <-strChan:
		return out, nil
	case err := <-errChan:
		return "", err
	}
}

// ReadAll reads all the record from the store.
func (s *Store) ReadAll() ([]string, error) {
	outputChan := make(chan []string, 1)
	errChan := make(chan error, 1)

	go func() {
		s.mu.Lock()
		file, err := os.Open(s.file)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()
		defer s.mu.Unlock()

		records := make([]string, 0)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text := scanner.Text()
			records = append(records, text)
		}

		if err := scanner.Err(); err != nil {
			errChan <- err
			return
		}

		outputChan <- records
	}()

	select {
	case out := <-outputChan:
		return out, nil
	case err := <-errChan:
		return nil, err
	}
}

// LookUp searches for a record by the longURL is any.
func (s *Store) LookUp(longURL string) bool {
	boolChan := make(chan bool, 1)

	go func() {
		s.mu.Lock()
		file, err := os.Open(s.file)
		if err != nil {
			log.Fatal("storage.LookUp: failed to open file", err)
			boolChan <- false
			return
		}
		defer file.Close()
		defer s.mu.Unlock()

		exists := false

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text := scanner.Text()

			if strings.Contains(text, longURL) {
				exists = true
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal("storage.LookUp: failed to scan file", err)
		}

		boolChan <- exists
	}()

	return <-boolChan
}

// New creates a new store.
func New(cfg *config.StorageConfiguration) (*Store, error) {
	if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
		_, err := os.Create(cfg.File)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		file: cfg.File,
	}, nil
}

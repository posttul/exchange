package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Response handle the expected response of a storage.
type Response struct {
	Rates map[string]Rate `json:"rates"`
}

// Rate is use to hold the rate information
type Rate struct {
	Value      string    `json:"value"`
	LastUpdate time.Time `json:"last_update"`
}

// Storage is an interface to handle the storage.
type Storage interface {
	Write([]byte) error
	Read() (*Response, error)
	Init() error
}

// NewFileStorage return a new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	f := &FileStorage{
		fileName: file,
	}
	if err := f.Init(); err != nil {
		return nil, err
	}
	return f, nil
}

// FileStorage is use to storage information into a file.
type FileStorage struct {
	Data     []byte `json:"data"`
	fileName string
	file     *os.File
}

// Init the file storage
func (f *FileStorage) Init() error {
	fle, err := os.OpenFile(f.fileName, os.O_CREATE, 777)
	f.file = fle
	if err != nil {
		return err
	}
	return f.update()
}

// Close the file.
func (f *FileStorage) Close() error {
	return f.file.Close()
}

// Read information from the storage.
func (f *FileStorage) Read() (*Response, error) {
	if err := f.update(); err != nil {
		return nil, err
	}
	r := &Response{}
	if err := json.Unmarshal(f.Data, &r); err != nil {
		return nil, err
	}
	return r, nil
}

func (f *FileStorage) update() error {
	// Read file to init with previews data
	dta, err := ioutil.ReadAll(f.file)
	if err != nil {
		return err
	}
	f.Data = dta
	return nil
}

// Write to the data storage.
func (f *FileStorage) Write(data []byte) error {
	f.Data = data
	return nil
}

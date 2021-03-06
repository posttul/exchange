package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Responder is something a store can respond.
type Responder interface {
	String() string
}

// Response handle the expected response of a storage.
type Response struct {
	Rates map[string]Rate `json:"rates"`
}

func (r *Response) String() string {
	return fmt.Sprintf("%+v", *r)
}

// Rate is use to hold the rate information
type Rate struct {
	Value      float64   `json:"value"`
	LastUpdate time.Time `json:"last_update"`
}

// Storage is an interface to handle the storage.
type Storage interface {
	Write([]byte) error
	Read(Responder) (Responder, error)
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
}

// Init the file storage
func (f *FileStorage) Init() error {
	return f.update()
}

// Read information from the storage.
func (f *FileStorage) Read(r Responder) (Responder, error) {
	if err := f.update(); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(f.Data, &r); err != nil {
		log.Printf("Ups JSON string -> \n %s \n err -> %s", string(f.Data), err.Error())
		return nil, err
	}
	return r, nil
}

// Write to the data storage.
func (f *FileStorage) Write(data []byte) error {
	f.Data = data
	return f.save()
}

func (f *FileStorage) save() error {
	fle, err := os.OpenFile(f.fileName, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer fle.Close()
	_, err = fle.Write(f.Data)
	return err
}

// update is use to check file information when information is required.
func (f *FileStorage) update() error {
	fle, err := os.OpenFile(f.fileName, os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer fle.Close()
	// Read file to init with previews data
	dta, err := ioutil.ReadAll(fle)
	if err != nil {
		return err
	}
	return f.Write(dta)
}

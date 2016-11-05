package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Storage interface {
	GetMD5(string) []byte
	Update(editionID string, newData io.Reader) error
}

type FileSystemStorage struct {
	path string
}

func NewFileSystemStorage(path string) FileSystemStorage {
	return FileSystemStorage{
		path: path,
	}
}

func (s FileSystemStorage) getPath(editionID string) string {
	return path.Join(s.path, fmt.Sprintf("%v.mmdb", editionID))
}

func (s FileSystemStorage) GetMD5(editionID string) []byte {
	path := s.getPath(editionID)
	file, err := os.Open(path)
	hasher := md5.New()
	if err == nil {
		io.Copy(hasher, file)
		file.Close()
	}
	return hasher.Sum(nil)
}

func (s FileSystemStorage) Update(editionID string, newData io.Reader) error {
	tmpfile, err := ioutil.TempFile(s.path, "updates")
	defer os.Remove(tmpfile.Name())
	io.Copy(tmpfile, newData)
	tmpfile.Close()
	err = os.Rename(tmpfile.Name(), s.getPath(editionID))
	return err
}

package files

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// RealWebClient is a real web client.
type RealWebClient struct{}

// WebClient gets a response given an url.
type WebClient interface {
	Get(string) ([]byte, error)
}

// Client embedds the WebClient
type Client struct {
	HTTP WebClient
}

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}

// HTTPDownload downloads stuff from the internet.
func HTTPDownload(webClient Client, url string) ([]byte, error) {
	body, err := webClient.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	return body, err
}

// Get gets the response given an url.
func (r *RealWebClient) Get(url string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// WriteFile writes the file to the disk. This is already tested by the Go team. No need to test.
func WriteFile(dst string, d []byte) error {
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		return err
	}
	return nil
}

// Unzip unzips a zip archive.
func Unzip(archive, target string) (string, error) {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return "", err
	}
	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			return "", err
		}
		defer fileReader.Close()
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return "", err
		}
	}
	return filepath.Join(target, "/", reader.File[0].Name), nil
}

// DownloadToFile downloads an http source to a file
func DownloadToFile(client Client, url string, folder string, dst string, deleteAfterUnzip bool) (string, error) {
	d, err := HTTPDownload(client, url)
	if err != nil {
		return "", err
	}
	err = WriteFile(dst, d)
	if err != nil {
		return "", err
	}
	sourcePath, err := Unzip(dst, folder)
	if err != nil {
		return "", err
	}
	if deleteAfterUnzip {
		err = os.Remove(dst)
		if err != nil {
			return "", err
		}
	}
	return sourcePath, nil
}

// IsDirectory checks if a path is a directory.
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}
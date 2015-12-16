/*
This is a utility to extract a tar.gz file
*/
package gunzip

import (
	"compress/gzip"
	"io"
	"net/http"
	"errors"
	"archive/tar"
	"path/filepath"
	"os"
)

func DownloadUrl(url *string) (*http.Response, error) {
	if url == nil || len(*url) == 0 {
		return nil, errors.New("Please provide a valid URL")
	}
	
	resp,error := http.Get(*url)
	
	if error != nil {
		return nil, error
	}
	
	return resp, nil
}

func CloseResponse(response *http.Response) {
	if (response != nil && response.Body != nil) {
		response.Body.Close()	
	}
}

func UnGzip (reader *io.ReadCloser) (*gzip.Reader, error) {
	if reader == nil {
		return nil, errors.New("No reader stream provided")
	}
	
	tgz, error := gzip.NewReader(*reader)
	
	if error != nil {
		return nil, error
	}
	
	return tgz, nil
}

func UnTar (reader *gzip.Reader) (*tar.Reader, error) {
	if reader == nil {
		return nil, errors.New("GZip Reader is nil")
	}
	
	treader := tar.NewReader(reader)
	if treader == nil {
		return nil, errors.New("Could not create reader from .tar file")
	}
	
	return treader, nil
}

func SaveTarToDisk(reader *tar.Reader, dest *string) error {
	if reader == nil {
		return errors.New("tar.Reader is nil")
	}
	
	if dest == nil || len(*dest) == 0 {
		return errors.New("Please provide a valid destination directory")
	}
	
	for {
		header, error := reader.Next()
		if error == io.EOF {
			break
		} else if error != nil {
			return error
		}
		
		info := header.FileInfo()		
		path := filepath.Join(*dest, header.Name)
		if info.IsDir() {
			os.MkdirAll(path, info.Mode())
		} else {
			file, error := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, info.Mode())			
			defer file.Close()
			if error != nil {
				return error
			}
			_, error = io.Copy(file, reader)
			if error != nil {
				return error
			}
		}
	}
	return nil
}
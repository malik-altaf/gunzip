package gunzip

import (
	"net/http"
	"compress/gzip"
	"archive/tar"
	"path/filepath"
	"errors"
	"io"
	"os"
)

type Tarzip struct {
	Url, Dest string
	resp *http.Response
	gzreader *gzip.Reader
	treader *tar.Reader
}

func (t *Tarzip) download() error {
	if t.Url == "" {
		return errors.New("Please provide a valid URL")
	}
	var err error;
	t.resp,err = http.Get(t.Url)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (t *Tarzip) close() {
	if (t.resp != nil && t.resp.Body != nil) {
		t.resp.Body.Close()	
	}
}

func (t *Tarzip) unzip () error {
	if t.resp == nil || t.resp.Body == nil {
		return errors.New("Nothing has been downloaded yet")
	}
	
	var err error;
	t.gzreader,err = gzip.NewReader(t.resp.Body)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (t *Tarzip) untar () error {
	if t.gzreader == nil {
		return errors.New("Archive have not been uncompressed")
	}
	
	t.treader = tar.NewReader(t.gzreader)
	if t.treader == nil {
		return errors.New("Could not create reader from .tar file")
	}
	
	return nil
}

func (t *Tarzip) saveToDisk() error {
	if t.treader == nil {
		return errors.New("Archive have not been extracted.")
	}
	
	if t.Dest == "" {
		return errors.New("Please provide a valid destination directory")
	}
	
	for {
		header, error := t.treader.Next()
		if error == io.EOF {
			break
		} else if error != nil {
			return error
		}
		
		info := header.FileInfo()		
		path := filepath.Join(t.Dest, header.Name)
		if info.IsDir() {
			os.MkdirAll(path, info.Mode())
		} else {
			file, error := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, info.Mode())			
			defer file.Close()
			if error != nil {
				return error
			}
			_, error = io.Copy(file, t.treader)
			if error != nil {
				return error
			}
		}
	}
	return nil
}

func (t *Tarzip) Extract() error {
	defer t.close();
	
	if error := t.download(); error != nil {
		return error;
	}
	
	if error := t.unzip(); error != nil {
		return error;
	}
	
	if error := t.untar(); error != nil {
		return error;
	}
	
	if error := t.saveToDisk(); error != nil {
		return error;
	}
	
	return nil;
}
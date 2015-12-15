package gunzip

import (
	"testing"
)

var (
	dest = "D:/test"
	tgz = "https://jdbc.postgresql.org/download/postgresql-jdbc-9.4-1206.src.tar.gz"
) 

func TestDownloadUrl(t *testing.T) {	
	resp, error := DownloadUrl(nil)
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if resp != nil {
		t.Errorf("Should have returned nil response")
	}
	
	url := ""
	resp, error = DownloadUrl(&url)
	if error == nil {
		t.Errorf("Should have returned error for empty URL")
	}
	if resp != nil {
		t.Errorf("Should have returned nil response for empty URL")
	}
	url = "invalidUrl"
	resp, error = DownloadUrl(&url)
	if error == nil {
		t.Errorf("Should have returned error for invalidurl URL")
	}
	if resp != nil {
		t.Errorf("Should have returned nil response for invalidurl URL")
	}
	url = "http://www.google.com"
	resp, error = DownloadUrl(&url)
	if error != nil {
		t.Errorf("Should not have returned error")
	}
	if resp == nil {
		t.Errorf("Should have returned valid response")
	}
}

func TestCloseResponse(t *testing.T) {
	CloseResponse(nil)
	url := "http://www.google.com"
	resp,_ := DownloadUrl(&url)
	resp.Body.Close()
	CloseResponse(resp);
	b := make([]byte, 5, 5)
	n, error := resp.Body.Read(b)
	if error == nil || n > 0 {
		t.Errorf("Should have returned error")
	}
	
	
	resp,_ = DownloadUrl(&url)	
	CloseResponse(resp);
	n, error = resp.Body.Read(b)
	if error == nil || n > 0 {
		t.Errorf("Should have returned error")
	}
}

func TestUnGzip(t *testing.T) {
	reader, error := UnGzip(nil)
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if reader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	// now test with valid reader for not to a zip stream
	url := "http://www.google.com"
	resp,_ := DownloadUrl(&url)
	reader, error = UnGzip(&resp.Body)
	if error == nil {
		t.Errorf("Should have returned error as the stream was not gz")
	}
	if reader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	url = "http://www.google.com"
	resp,_ = DownloadUrl(&url)
	reader, error = UnGzip(&resp.Body)	
	if error == nil {
		t.Errorf("Should have returned error as the stream was not gz")
	}
	if reader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	url = tgz 
	resp,_ = DownloadUrl(&url)
	reader, error = UnGzip(&resp.Body)		
	if error != nil {
		t.Errorf("Thrown error : ",error)
	}
	
	if reader == nil {
		t.Errorf("Should have returned valid reader")
	}
}

func TestUnTar(t *testing.T) {
	reader, error := UnTar(nil)
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if reader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	url := tgz 
	resp,_ := DownloadUrl(&url)
	greader, error := UnGzip(&resp.Body)
	reader, error = UnTar(greader);	
	CloseResponse(resp)
	
	if error != nil {
		t.Errorf("Thrown error : ",error)
	}
	
	if reader == nil {
		t.Errorf("Should have returned valid tar reader")
	}
}

func TestSaveTarToDisk(t *testing.T) {
	error := SaveTarToDisk(nil, nil)
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	path := "something"
	
	error = SaveTarToDisk(nil, &path)
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	url := tgz 
	resp,_ := DownloadUrl(&url)
	greader, _ := UnGzip(&resp.Body)
	reader, _ := UnTar(greader);
	
	error = SaveTarToDisk(reader, nil)
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	error = SaveTarToDisk(reader, &path)
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	path = dest
	
	error = SaveTarToDisk(reader, &path)
	if error != nil {
		t.Errorf("Got error while extracting the tar.gz file")
	}
	
}


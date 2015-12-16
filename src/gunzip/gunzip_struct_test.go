package gunzip

import (
	"testing"
)

func TestDownload(t *testing.T) {
	tz := new(Tarzip) 	
	error := tz.download()
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if tz.resp != nil {
		t.Errorf("Should have returned nil response")
	}
	
	tz.Url = ""
	error = tz.download()
	if error == nil {
		t.Errorf("Should have returned error for empty URL")
	}
	if tz.resp != nil {
		t.Errorf("Should have returned nil response for empty URL")
	}
	tz.Url = "invalidUrl"
	error = tz.download()
	if error == nil {
		t.Errorf("Should have returned error for invalidurl URL")
	}
	if tz.resp != nil {
		t.Errorf("Should have returned nil response for invalidurl URL")
	}
	
	tz.Url = "http://www.google.com"
	error = tz.download()
	if error != nil {
		t.Errorf("Should not have returned error")
	}
	if tz.resp == nil {
		t.Errorf("Should have returned valid response")
	}
	
	tz.close()
}

func TestClose(t *testing.T) {
	tz := new(Tarzip)
	tz.close()
	tz.Url = "http://www.google.com"
	error := tz.download()
	if (error != nil) {
		t.Errorf("Download failed");
	}
	tz.resp.Body.Close()
	tz.close()
	b := make([]byte, 5, 5)
	n, error := tz.resp.Body.Read(b)
	if error == nil || n > 0 {
		t.Errorf("Should have returned error")
	}
	
	
	error = tz.download()
	if (error != nil) {
		t.Errorf("Download failed");
	}	
	tz.close()
	n, error = tz.resp.Body.Read(b)
	if error == nil || n > 0 {
		t.Errorf("Should have returned error")
	}
}

func TestUnzip(t *testing.T) {
	tz := new(Tarzip)
	error := tz.unzip()
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if tz.gzreader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	// now test with valid reader for not to a zip stream
	tz.Url = "http://www.google.com"
	error = tz.download()
	if (error != nil) {
		t.Errorf("Download failed");
	}
	
	error = tz.unzip()
	if error == nil {
		t.Errorf("Should have returned error as the stream was not gz")
	}
	
	if tz.gzreader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	tz.close()
	
	tz.Url = tgz 
	error = tz.download()	
	if (error != nil) {
		t.Errorf("Download failed")
	}
	
	error = tz.unzip()		
	if error != nil {
		t.Errorf("Thrown error : ",error)
	}
	
	if tz.gzreader == nil {
		t.Errorf("Should have returned valid reader")
	}
	tz.close()
}

func TestUntar(t *testing.T) {
	tz := new(Tarzip)

	error := tz.untar()
	if error == nil {
		t.Errorf("Should have returned error")
	}
	if tz.treader != nil {
		t.Errorf("Should have returned nil response")
	}
	
	tz.Url = tgz 
	error = tz.download()
	if (error != nil) {
		t.Errorf("Download failed")
	}	
	error = tz.unzip()
	if (error != nil) {
		t.Errorf("Unzip failed")
	}
	
	error = tz.untar()
	
	if error != nil {
		t.Errorf("Thrown error : ",error)
	}
	
	if tz.treader == nil {
		t.Errorf("Should have returned valid tar reader")
	}
	
	tz.close()
}

func TestSaveToDisk(t *testing.T) {
	tz := new(Tarzip)
	error := tz.saveToDisk()
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	tz.Dest = "something"
	
	error = tz.saveToDisk()
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	tz.Url = tgz 
	error = tz.download()
	if error != nil {
		t.Errorf("Download error: ",error)
	}
	
	error = tz.unzip()
	if error != nil {
		t.Errorf("Unzip error: ",error)
	}
	
	error = tz.untar()
	if error != nil {
		t.Errorf("Untar error: ",error)
	}
	
	error = tz.saveToDisk()
	if error == nil {
		t.Errorf("Expected error, found nil")
	}
	
	tz.Dest = dest
	
	error = tz.saveToDisk()
	if error != nil {
		t.Errorf("Got error while extracting the tar.gz file")
	}
	
}


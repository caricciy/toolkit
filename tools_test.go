package toolkit

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var (
		testTools Tools
	)

	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Errorf("RandomString() = %v, want %v", len(s), 10)
	}
}

func TestTools_UploadFiles(t *testing.T) {
	var (
		testTools Tools
	)

	// Create a mock HTTP request with multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "test.txt")
	file.Write([]byte("test file content"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call the UploadFiles function
	uploadedFiles, err := testTools.UploadFiles(req, os.TempDir(), true)

	// Assert the results
	if err != nil {
		t.Errorf("UploadFiles() returned an error: %v", err)
	}

	if len(uploadedFiles) != 1 {
		t.Errorf("UploadFiles() returned %d uploaded files, want %d", len(uploadedFiles), 1)
	}

	uploadedFile := uploadedFiles[0]
	if uploadedFile.NewFileName == "" {
		t.Errorf("UploadFiles() did not set the new file name")
	}

	if uploadedFile.FileSize != int64(len("test file content")) {
		t.Errorf("UploadFiles() set incorrect file size, got %d, want %d", uploadedFile.FileSize, len("test file content"))
	}
}

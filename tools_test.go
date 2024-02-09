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
	fileOne, _ := writer.CreateFormFile("fileOne", "test-1.txt")
	fileOneContent := "test file-1 content"
	_, _ = fileOne.Write([]byte(fileOneContent))

	fileTwo, _ := writer.CreateFormFile("fileTwo", "test-2.txt")
	fileTwoContent := "test file-2 content"
	_, _ = fileTwo.Write([]byte(fileTwoContent))

	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call the UploadFiles function
	uploadedFiles, err := testTools.UploadFiles(req, os.TempDir(), true)

	// Assert the results
	if err != nil {
		t.Errorf("UploadFiles() returned an error: %v", err)
	}

	if len(uploadedFiles) != 2 {
		t.Errorf("UploadFiles() returned %d uploaded files, want %d", len(uploadedFiles), 1)
	}

	uploadedFileOne := uploadedFiles[0]
	uploadedFileTwo := uploadedFiles[1]
	if uploadedFileOne.NewFileName == "" && uploadedFileTwo.NewFileName == "" {
		t.Errorf("UploadFiles() did not set the new file name")
	}

	if uploadedFileOne.FileSize != int64(len(fileOneContent)) {
		t.Errorf("UploadFiles() set incorrect file size, got %d, want %d", uploadedFileOne.FileSize, len(fileOneContent))
	}

	if uploadedFileTwo.FileSize != int64(len(fileTwoContent)) {
		t.Errorf("UploadFiles() set incorrect file size, got %d, want %d", uploadedFileTwo.FileSize, len(fileTwoContent))
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	var (
		testTools Tools
	)

	// Create a mock HTTP request with multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "test.txt")
	fileContent := "test file content"
	_, _ = file.Write([]byte(fileContent))
	_ = writer.Close()

	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Call the UploadOneFile function
	uploadedFile, err := testTools.UploadOneFile(req, os.TempDir(), true)

	// Assert the results
	if err != nil {
		t.Errorf("UploadOneFile() returned an error: %v", err)
	}

	if uploadedFile.NewFileName == "" {
		t.Errorf("UploadOneFile() did not set the new file name")
	}

	if uploadedFile.FileSize != int64(len(fileContent)) {
		t.Errorf("UploadOneFile() set incorrect file size, got %d, want %d", uploadedFile.FileSize, len(fileContent))
	}
}

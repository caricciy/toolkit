package toolkit

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"
const oneGiB = 1024 * 1024 * 1024 // 1 GiB

// Tools is the type used to instantiate the toolkit module.
type Tools struct {
	MaxFileSize      int
	AllowedFileTypes []string
}

// NewTools returns a new Tools instance.
func NewTools() *Tools {
	return &Tools{
		MaxFileSize:      0,
		AllowedFileTypes: []string{"image/gif"},
	}
}

// RandomString returns a random string of the given length.
func (t *Tools) RandomString(length int) string {
	s, r := make([]rune, length), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}

// UploadedFile is a struct used to save file upload information.
type UploadedFile struct {
	NewFileName      string
	OriginalFileName string
	FileSize         int64
}

// UploadOneFile uploads a single file from a multipart form request to the given directory.
func (t *Tools) UploadOneFile(r *http.Request, uploadDir string, rename ...bool) (*UploadedFile, error) {

	files, err := t.UploadFiles(r, uploadDir, rename...)
	if err != nil {
		return nil, err
	}

	return files[0], nil
}

// UploadFiles uploads files from a multipart form request to the given directory.
func (t *Tools) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {
	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if t.MaxFileSize == 0 {
		t.MaxFileSize = oneGiB
	}

	err := r.ParseMultipartForm(int64(t.MaxFileSize))

	if err != nil {
		return nil, errors.New("file size too large")
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {

			uploadedFiles, err = func(UploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				var uploadedFile UploadedFile
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer func() {
					_ = infile.Close()
				}()

				allowed, _ := t.getFileContentType(infile)

				if !allowed {
					return nil, errors.New("file type not allowed")
				}

				// Go to the beginning of the file to read it again because we read the first 512 bytes
				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", t.RandomString(25), filepath.Ext(hdr.Filename))
				} else {
					uploadedFile.NewFileName = hdr.Filename
				}

				uploadedFile.OriginalFileName = hdr.Filename

				var outfile *os.File
				defer func() {
					_ = outfile.Close()
				}()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}
				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil
			}(uploadedFiles)

			if err != nil {
				// May some files have been uploaded before the error occurred (should we delete them?)
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}

// getFileContentType reads the first 512 bytes of the file and returns true if the content type is allowed.
func (t *Tools) getFileContentType(infile multipart.File) (bool, error) {
	if len(t.AllowedFileTypes) > 0 {
		// Get the first 512 bytes to sniff the content type
		buff := make([]byte, 512)
		if _, err := infile.Read(buff); err != nil {
			return false, err
		}

		fileType := http.DetectContentType(buff)

		// Check if the detected content type is allowed
		for _, allowedType := range t.AllowedFileTypes {
			if strings.EqualFold(fileType, allowedType) {
				return true, nil
			}
		}
	} else {
		// When no allowed types are specified, allow all
		return true, nil
	}

	return false, nil
}

// CreateDirIfNotExist creates a directory, and all necessary parents,  if it does not exist.
func (t *Tools) CreateDirIfNotExist(path string) error {
	const mode = 0755 // rwxr-xr-x
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

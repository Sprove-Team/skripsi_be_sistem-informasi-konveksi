package helper

import (
	"mime/multipart"
	"net/url"

	"github.com/h2non/filetype"
)

func IsGoogleStorageURL(inputURL string) bool {
	// Parse the URL

	u, err := url.Parse(inputURL)
	if err != nil {
		return false
	}

	if u.Hostname() != "storage.googleapis.com" {
		return false
	}
	return true
}

func IsImg(fh *multipart.FileHeader) bool {
	file, err := fh.Open()
	if err != nil {
		return false
	}

	if fh.Size > 10*1024*1024 { // 10 MB
		return false
	}

	defer file.Close()

	buff := make([]byte, 512)

	if _, err := file.Read(buff); err != nil {
		return false
	}

	kind, _ := filetype.Match(buff)

	if kind == filetype.Unknown {
		return false
	}

	if kind.Extension != "jpg" && kind.Extension != "jpeg" && kind.Extension != "png" {
		return false
	}

	if kind.MIME.Type != "image" {
		return false
	}

	return true
}

package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/be-sistem-informasi-konveksi/common/message"
	"github.com/be-sistem-informasi-konveksi/pkg"
	"golang.org/x/sync/errgroup"
)

func CountUniqueElements(slice []string) int {
	uniqueMap := make(map[string]bool)

	// Iterate over the slice and add each element to the map
	for _, element := range slice {
		uniqueMap[element] = true
	}

	// Return the number of unique elements (length of the map)
	return len(uniqueMap)
}

func GetFiles(fieldNameAllowed string, mform func() (*multipart.Form, error)) ([]*multipart.FileHeader, error) {
	form, err := mform()
	if err != nil {
		return nil, err
	}
	if len(form.File) == 0 {
		return nil, nil
	}
	var r []*multipart.FileHeader
	for fieldName, fileHeaders := range form.File {
		if fieldName == fieldNameAllowed {
			r = append(r, fileHeaders...)
		}
	}
	return r, nil
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateUniqueFileName(fileName string) string {
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	// Use a combination of base name, timestamp, and a random string to ensure uniqueness
	randomString := RandomString(6) // Change the length of the random string as needed
	return fmt.Sprintf("%s_%d_%s", fileName, time.Now().UnixNano(), randomString)
}

func DataParser(data string, out any) error {
	return json.Unmarshal([]byte(data), out)
}

func SaveMultiFileInLocal(files []*multipart.FileHeader) ([]string, error) {
	g := new(errgroup.Group)
	g.SetLimit(10)
	fileNames := make([]string, len(files))
	for i, file := range files {
		file := file
		i := i
		g.Go(func() error {
			if file.Size > 5*1024*1024 { // 5 MB
				return errors.New(message.InvalidImageSize)
			}

			fi, err := file.Open()
			if err != nil {
				return err
			}
			defer fi.Close() // Close the file

			fReader, err := pkg.Decode(fi)

			if err != nil {
				if err == image.ErrFormat {
					return errors.New(message.InvalidImageFormat)
				}
				return err
			}

			fileName := GenerateUniqueFileName(file.Filename) + ".webp"
			filePath := filepath.Join("./public/img/", fileName)
			fo, err := os.Create(filePath)
			if err != nil {
				return err
			}

			defer fo.Close()

			if _, err := io.Copy(fo, fReader); err != nil {
				return err
			}
			fileNames[i] = fileName

			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return fileNames, nil
}

func DeleteMultiFileInLocal(fileNames []string) error {
	g := new(errgroup.Group)
	g.SetLimit(10)
	for _, fileName := range fileNames {
		fileName := fileName
		g.Go(func() error {
			return os.Remove("./public/img/" + fileName)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

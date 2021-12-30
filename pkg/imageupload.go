package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//ImageUpload saves image to the local machine
func ImageUpload(w http.ResponseWriter, r *http.Request) (string, error) {
	var err error
	var maxmem int64 = 1024 * 1024 * 30
	r.Body = http.MaxBytesReader(w, r.Body, maxmem)
	if err = r.ParseMultipartForm(maxmem); err != nil {
		return "", err
	}

	//Save file
	file, fileHeader, err := r.FormFile("uploaded")
	if err != nil {
		Info("No image is uploaded")
		return "", nil
	}

	defer file.Close()

	if fileHeader.Size > 1024*1024*20 {
		Warning("The size of image is bigger than 20 MB")
		return "", errors.New("The size of the image is bigger than 20 MB")
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		Warning("The provided image format is not allowed - 512 bytes")
		return "", errors.New("The provided image format is not allowed. Please upload a JPEG or PNG or GIF image")
	}

	typeofimage := http.DetectContentType(buff)
	if typeofimage != "image/jpeg" && typeofimage != "image/png" && typeofimage != "image/gif" {
		Warning("The provided image format is not allowed")
		return "", errors.New("The provided image format is not allowed. Please upload a JPEG or PNG or GIF image")
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		Warning("The provided image format is not allowed - seekstart")
		return "", errors.New("The provided image format is not allowed. Please upload a JPEG or PNG or GIF image")
	}

	//Create a new file in the uploads directory
	image := fmt.Sprintf("./static/img_posts/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	dst, err := os.Create(image)
	if err != nil {
		Danger("Cannot create image file")
		return "", errors.New("Internal server error. Please try to upload the image later")
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		Danger("Cannot save image to the local machine")
		return "", errors.New("Internal server error. Please try to upload the image later")
	}

	return image[1:], nil
}

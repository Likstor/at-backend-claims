package image

import (
	"net/http"
	"bytes"
	"image/jpeg"
	_ "image/png"
	"image"
)

func IsImage(fileBytes []byte) (bool, string) {
	imageType := http.DetectContentType(fileBytes)
	
	switch imageType {
	case "image/jpeg":
		fallthrough
	case "image/png":
		return true, imageType
	}

	return false, imageType
}

func ConvertToJpeg(file []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(file))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}

	if err := jpeg.Encode(buf, img, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

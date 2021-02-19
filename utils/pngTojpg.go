package utils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"image/jpeg"
	"image/png"
)

func ToPNG(imageBytes []byte, contentType string) ([]byte, error) {
	switch contentType {
	case ".png":
		return imageBytes, nil
	case ".jpg":
		img, err := jpeg.Decode(bytes.NewBuffer(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode jpeg")
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, errors.Wrap(err, "unable to encode png")
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("unable to convert %#v to png", contentType)

}

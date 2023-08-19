package cleanhttp

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Convert file image and return it as base64 string to use it in request.
func ImageToBase64(filePath string) (string, error) {
	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return "data:image/" + filepath.Ext(filePath)[1:] + ";base64," + base64.StdEncoding.EncodeToString(imageData), nil
}

func CalculateContentLength(reader io.Reader) (int64, error) {
	seeker, ok := reader.(io.Seeker)
	if ok {
		currentPos, err := seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, err
		}

		endPos, err := seeker.Seek(0, io.SeekEnd)
		if err != nil {
			return 0, err
		}

		_, err = seeker.Seek(currentPos, io.SeekStart)
		if err != nil {
			return 0, err
		}

		contentLength := endPos - currentPos
		return contentLength, nil
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return 0, err
	}

	contentLength := int64(len(content))
	return contentLength, nil
}

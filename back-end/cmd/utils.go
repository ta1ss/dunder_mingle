package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (app *application) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) readJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJson(w, statusCode, payload, nil)
}

func saveBase64Image(base64String, path string) (imgName string, err error) {
	var imageType string
	if strings.Contains(base64String, "data:image/jpeg;base64,") {
		imageType = "jpg"
		base64String = strings.Replace(base64String, "data:image/jpeg;base64,", "", 1)
	} else if strings.Contains(base64String, "data:image/png;base64,") {
		imageType = "png"
		base64String = strings.Replace(base64String, "data:image/png;base64,", "", 1)
	} else if strings.Contains(base64String, "data:image/jpg;base64,") {
		imageType = "jpg"
		base64String = strings.Replace(base64String, "data:image/jpg;base64,", "", 1)
	} else if strings.Contains(base64String, "data:image/gif;base64,") {
		imageType = "gif"
		base64String = strings.Replace(base64String, "data:image/gif;base64,", "", 1)
	} else {
		return "", fmt.Errorf("unsupported image type or too large")
	}

	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}
	imgName = randSeq(12) + "." + imageType
	err = os.WriteFile(path+imgName, decoded, 0644)
	if err != nil {
		return "", err
	}

	return imgName, nil
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func deleteImage(image, path string) error {
	err := os.Remove(path + image)
	if err != nil {
		return err
	}
	return nil
}

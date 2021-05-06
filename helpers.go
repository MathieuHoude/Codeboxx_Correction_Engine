package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/joho/godotenv"
)

//SSHKeys contains the pointers to the ssh keys
type SSHKeys struct {
	PublicKey  *string
	PrivateKey *string
}

type malformedRequest struct {
	status int
	msg    string
}

func getKeys() SSHKeys {
	cmd := exec.Command("whoami")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	username := string(stdout)
	username = strings.TrimSuffix(username, "\n")

	cmd = exec.Command("cat", "/home/"+username+"/.ssh/id_rsa")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	privateKey := string(stdout)

	cmd = exec.Command("cat", "/home/"+username+"/.ssh/id_rsa.pub")
	stdout, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	publicKey := string(stdout)

	keys := SSHKeys{
		PublicKey:  &publicKey,
		PrivateKey: &privateKey,
	}

	return keys

}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			return &malformedRequest{status: http.StatusBadRequest, msg: "Request body contains badly-formed JSON"}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func deleteFile(filePath string) {
	err := os.Remove(filePath) // remove a single file
	if err != nil {
		fmt.Println(err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func startWorkers(correctionWorkerAmount, gradingWorkerAmount int) {
	for i := 1; i <= correctionWorkerAmount; i++ {
		go worker(i, "correction")
	}
	for j := 1; j <= gradingWorkerAmount; j++ {
		go worker(j, "grading")
	}
}

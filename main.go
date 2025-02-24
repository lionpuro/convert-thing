package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lionpuro/convert-thing/files"
)

const (
	convertError = "error converting file"
)

func main() {
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /files", handleConvert)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("error getting file: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	defer file.Close()

	outputFormat := r.FormValue("outputFormat")

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error reading uploaded file: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	fileType, format, err := files.DetectType(fileBytes)
	formats, ok := files.Formats()[fileType]
	if !ok {
		http.Error(w, "filetype not supported", http.StatusBadRequest)
		return
	}
	supported := slices.Contains(formats, format)
	if !supported {
		http.Error(w, "filetype not supported", http.StatusBadRequest)
		return
	}

	converted, err := files.ConvertTo(outputFormat, bytes.NewReader(fileBytes))
	if err != nil {
		log.Printf("error converting file: %v\n", err)
		http.Error(w, convertError, http.StatusInternalServerError)
		return
	}

	filename := files.ChangeExt(fileHeader.Filename, outputFormat)
	resultPath := filepath.Join("/tmp", filename)
	if err := files.WriteToFile(resultPath, converted); err != nil {
		http.Error(w, convertError, http.StatusInternalServerError)
		return
	}

	result, err := os.ReadFile(resultPath)
	if err != nil {
		log.Printf("error reading output file: %v\n", err)
		http.Error(w, convertError, http.StatusInternalServerError)
		return
	}
	newMime := mimetype.Detect(result).String()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", newMime)
	if _, err := w.Write(result); err != nil {
		log.Printf("error writing file to response writer: %v\n", err)
		http.Error(w, convertError, http.StatusInternalServerError)
		return
	}
}

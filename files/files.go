package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ConvertTo(format string, file io.Reader) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	inputBytes := buf.Bytes()
	inputFile, err := os.CreateTemp("/tmp", "*.tmp")
	if err != nil {
		return nil, err
	}
	defer os.Remove(inputFile.Name())

	if _, err := inputFile.Write(inputBytes); err != nil {
		return nil, err
	}

	input := inputFile.Name()
	output := fmt.Sprintf("/tmp/%s", randomFileName(format))
	err = ffmpeg.
		Input(input).
		Output(output).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
	if err != nil {
		return nil, fmt.Errorf("error converting file: %v\n", err)
	}

	outputFile, err := os.Open(output)
	if err != nil {
		return nil, err
	}
	defer os.Remove(outputFile.Name())

	outputBytes, err := io.ReadAll(outputFile)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(outputBytes), nil
}

func WriteToFile(path string, file io.Reader) error {
	nf, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer nf.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	bytes := buf.Bytes()
	if _, err := nf.Write(bytes); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func DetectType(file []byte) (string, string, error) {
	mime := mimetype.Detect(file).String()
	types := strings.Split(mime, "/")
	if len(types) != 2 {
		return "", "", fmt.Errorf("%s is not a valid mimetype", mime)
	}
	return types[0], types[1], nil
}

func ChangeExt(name, ext string) string {
	n := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
	return fmt.Sprintf("%s.%s", n, ext)
}

package flac

import (
	"errors"
	"io"
)

func readBytes(reader io.Reader, size int) ([]byte, error) {
	bytes := make([]byte, size)
	n, err := reader.Read(bytes)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, errors.New("incorrect size")
	}
	return bytes, nil
}

func readString(reader io.Reader, size int) (string, error) {
	bytes, err := readBytes(reader, size)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

package meta

import "io"

type Padding struct {
	Data []byte
}

func readPadding(reader io.Reader, size int) (*Padding, error) {
	padding := &Padding{
		Data: make([]byte, size),
	}
	_, err := reader.Read(padding.Data)
	return padding, err
}

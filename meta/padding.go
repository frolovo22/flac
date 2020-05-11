package meta

import (
	"github.com/icza/bitio"
)

type Padding struct {
	Data []byte
}

func readPadding(reader *bitio.Reader, size int) (*Padding, error) {
	padding := &Padding{
		Data: make([]byte, size),
	}
	_, err := reader.Read(padding.Data)
	return padding, err
}

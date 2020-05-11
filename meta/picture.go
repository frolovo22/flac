package meta

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/icza/bitio"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

type Picture struct {
	Type           int32
	MIME           string
	Description    string
	Width          int32
	Height         int32
	BitsPerPixel   int32
	NumberOfColors int32
	PictureData    []byte
}

func readPicture(reader *bitio.Reader) (*Picture, error) {
	var picture Picture

	// Picture type
	err := binary.Read(reader, binary.BigEndian, &picture.Type)
	if err != nil {
		return nil, err
	}

	// MIME
	MIMEBytes, err := readLengthData(reader, binary.BigEndian)
	if err != nil {
		return nil, err
	}
	picture.MIME = string(MIMEBytes)

	// Description
	DescriptionBytes, err := readLengthData(reader, binary.BigEndian)
	if err != nil {
		return nil, err
	}
	picture.Description = string(DescriptionBytes)

	// Width
	err = binary.Read(reader, binary.BigEndian, &picture.Width)
	if err != nil {
		return nil, err
	}

	// Height
	err = binary.Read(reader, binary.BigEndian, &picture.Height)
	if err != nil {
		return nil, err
	}

	// Bits per pixel
	err = binary.Read(reader, binary.BigEndian, &picture.BitsPerPixel)
	if err != nil {
		return nil, err
	}

	// Number of colors
	err = binary.Read(reader, binary.BigEndian, &picture.NumberOfColors)
	if err != nil {
		return nil, err
	}

	// Picture data
	picture.PictureData, err = readLengthData(reader, binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return &picture, nil
}

func (p *Picture) GetImage() (image.Image, error) {
	switch p.MIME {
	case "image/jpeg":
		return jpeg.Decode(bytes.NewReader(p.PictureData))
	case "image/png":
		return png.Decode(bytes.NewReader(p.PictureData))
	case "-->":
		return downloadImage(string(p.PictureData))
	}
	return nil, errors.New("incorrect picture type")
}

// Read format:
// [length, data]
func readLengthData(reader *bitio.Reader, order binary.ByteOrder) ([]byte, error) {
	// length
	var length uint32
	err := binary.Read(reader, order, &length)
	if err != nil {
		return nil, err
	}

	// data
	data := make([]byte, length)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(resp.Body)
	return img, err
}

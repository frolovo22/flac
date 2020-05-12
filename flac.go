package flac

import (
	"errors"
	"frolovo22/flac/frame"
	"frolovo22/flac/meta"
	"github.com/icza/bitio"
	"io"
	"os"
)

const StreamMarker = "fLaC"

type FLAC struct {
	Marker         string // always "fLaC"
	MetadataBlocks []meta.MetadataBlock
	Frame          frame.Frame
}

func ReadFile(path string) (*FLAC, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Read(file)
}

func Read(reader io.Reader) (*FLAC, error) {
	flac := FLAC{}

	bits := bitio.NewReader(reader)

	// check format
	err := flac.readMarker(bits)
	if err != nil {
		return &flac, err
	}

	// read metadata
	err = flac.readMetadata(bits)
	if err != nil {
		return &flac, err
	}

	// read frames
	err = flac.readFrame(bits)
	if err != nil {
		return &flac, err
	}

	return &flac, nil
}

func (f *FLAC) readMarker(reader *bitio.Reader) error {
	marker := make([]byte, 4)
	_, err := reader.Read(marker)
	if err != nil {
		return err
	}
	if string(marker) != StreamMarker {
		return errors.New("incorrect marker")
	}
	f.Marker = string(marker)
	return nil
}

func (f *FLAC) readMetadata(reader *bitio.Reader) error {
	isLast := false
	for !isLast {
		metadata, err := meta.ReadMetadataBlock(reader)
		if err != nil {
			return err
		}
		f.MetadataBlocks = append(f.MetadataBlocks, *metadata)
		isLast = metadata.Header.IsLast
	}
	return nil
}

func (f *FLAC) readFrame(reader *bitio.Reader) error {
	frame, err := frame.ReadFrame(reader)
	if err != nil {
		return err
	}
	f.Frame = *frame
	return nil
}

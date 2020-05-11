package flac

import (
	"errors"
	"frolovo22/flac/frame"
	"frolovo22/flac/meta"
	"io"
	"os"
)

const StreamMarker = "fLaC"

type FLAC struct {
	Marker         string // always "fLaC"
	MetadataBlocks []meta.MetadataBlock
	Frames         []frame.Frame
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

	// check format
	err := flac.readMarker(reader)
	if err != nil {
		return &flac, err
	}

	// read metadata
	err = flac.readMetadata(reader)
	if err != nil {
		return &flac, err
	}

	// read frames

	return &flac, nil
}

func (f *FLAC) readMarker(reader io.Reader) error {
	marker, err := readString(reader, 4)
	if err != nil {
		return err
	}
	if marker != StreamMarker {
		return errors.New("incorrect marker")
	}
	f.Marker = marker
	return nil
}

func (f *FLAC) readMetadata(reader io.Reader) error {
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

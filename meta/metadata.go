package meta

import (
	"errors"
	"github.com/icza/bitio"
	"io"
)

type MetadataBlock struct {
	Header MetadataBlockHeader
	Data   MetadataBlockData
}

func ReadMetadataBlock(reader io.Reader) (*MetadataBlock, error) {
	metadata := &MetadataBlock{}

	header, err := readMetadataBlockHeader(reader)
	if err != nil {
		return metadata, err
	}
	metadata.Header = *header

	switch metadata.Header.Type {
	case StreamInfoBlockType:
		metadata.Data, err = readStreamInfo(reader)
	case PaddingBlockType:
		metadata.Data, err = readPadding(reader, metadata.Header.Length)
	case ApplicationBlockType:
		metadata.Data, err = readApplication(reader, metadata.Header.Length)
	case SeekTableBlockType:
		metadata.Data, err = readSeekTable(reader, metadata.Header.Length)
	case VorbisCommentBlockType:
		metadata.Data, err = readVorbisComment(reader)
	case CurSheetBlockType:
		metadata.Data = nil
	case PictureBlockType:
		metadata.Data = nil
	case InvalidBlockType:
		err = errors.New("invalid block type")
	default:
		err = errors.New("error block type")
	}

	return metadata, err
}

type MetadataBlockHeader struct {
	IsLast bool      // Last-metadata-block flag: '1' if this block is the last metadata block before the audio blocks, '0' otherwise.
	Type   BlockType // Block type. 127 - invalid, to avoid confusion with a frame sync code
	Length int       // Length (in bytes) of metadata to follow (does not include the size of the MetadataBlockHeader)
}

type MetadataBlockData interface{}

func readMetadataBlockHeader(reader io.Reader) (*MetadataBlockHeader, error) {
	header := MetadataBlockHeader{}

	bits := bitio.NewReader(reader)

	// IsLast: 1 bit
	isLast, err := bits.ReadBool()
	if err != nil {
		return &header, err
	}
	header.IsLast = isLast

	// Type: bits 2-8
	blockType, err := bits.ReadBits(7)
	if err != nil {
		return &header, err
	}
	header.Type = BlockType(blockType)

	// Size: 3 bytes
	length, err := bits.ReadBits(24)
	if err != nil {
		return &header, err
	}
	header.Length = int(length)

	return &header, nil
}

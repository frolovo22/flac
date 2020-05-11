package meta

import (
	"errors"
	"github.com/icza/bitio"
)

type MetadataBlock struct {
	Header MetadataBlockHeader
	Data   MetadataBlockData
}

func ReadMetadataBlock(reader *bitio.Reader) (*MetadataBlock, error) {
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
		metadata.Data, err = readPadding(reader, header.Length)
	case ApplicationBlockType:
		metadata.Data, err = readApplication(reader, header.Length)
	case SeekTableBlockType:
		metadata.Data, err = readSeekTable(reader, header.Length)
	case VorbisCommentBlockType:
		metadata.Data, err = readVorbisComment(reader)
	case CueSheetBlockType:
		metadata.Data, err = readCueSheet(reader)
	case PictureBlockType:
		metadata.Data, err = readPicture(reader)
	case InvalidBlockType:
		err = errors.New("invalid block type")
	default:
		data := make([]byte, header.Length)
		_, err = reader.Read(data)
	}

	return metadata, err
}

type MetadataBlockHeader struct {
	IsLast bool      // Last-metadata-block flag: '1' if this block is the last metadata block before the audio blocks, '0' otherwise.
	Type   BlockType // Block type. 127 - invalid, to avoid confusion with a frame sync code
	Length int       // Length (in bytes) of metadata to follow (does not include the size of the MetadataBlockHeader)
}

type MetadataBlockData interface{}

func readMetadataBlockHeader(reader *bitio.Reader) (*MetadataBlockHeader, error) {
	header := MetadataBlockHeader{}

	// IsLast: 1 bit
	isLast, err := reader.ReadBool()
	if err != nil {
		return &header, err
	}
	header.IsLast = isLast

	// Type: bits 2-8
	blockType, err := reader.ReadBits(7)
	if err != nil {
		return &header, err
	}
	header.Type = BlockType(blockType)

	// Size: 3 bytes
	length, err := reader.ReadBits(24)
	if err != nil {
		return &header, err
	}
	header.Length = int(length)

	return &header, nil
}

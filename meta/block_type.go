package meta

type BlockType uint8

const (
	StreamInfoBlockType    BlockType = 0
	PaddingBlockType       BlockType = 1
	ApplicationBlockType   BlockType = 2
	SeekTableBlockType     BlockType = 3
	VorbisCommentBlockType BlockType = 4
	CurSheetBlockType      BlockType = 5
	PictureBlockType       BlockType = 6
	InvalidBlockType       BlockType = 127
)

func (b *BlockType) String() string {
	switch *b {
	case StreamInfoBlockType:
		return "STREAMINFO"
	case PaddingBlockType:
		return "PADDING"
	case ApplicationBlockType:
		return "APPLICATION"
	case SeekTableBlockType:
		return "SEEKTABLE"
	case VorbisCommentBlockType:
		return "VORBIS_COMMENT"
	case CurSheetBlockType:
		return "CURSHEET"
	case PictureBlockType:
		return "PICTURE"
	case InvalidBlockType:
		return "INVALID"
	default:
		return ""
	}
}

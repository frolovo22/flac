package meta

import "github.com/icza/bitio"

type CueSheet struct {
	MediaCatalogNumber    string
	NumberOfLeadInSamples uint64
	CompactDisc           bool
	Reserved              []byte
	NumberOfTracks        uint8
	CueSheetTracks        []CueSheetTrack
}

type CueSheetTrack struct {
	OffsetInSamples         uint64
	TrackNumber             uint8
	ISRC                    string
	NonAudioType            bool
	PreEmphasis             bool
	Reserved                []byte
	NumberOfTrackIndexPoint uint8
	CueSheetTrackIndexes    []CueSheetTrackIndex
}

type CueSheetTrackIndex struct {
	OffsetInSamples  uint64
	IndexPointNumber uint8
	Reserved         []byte
}

func readCueSheet(reader *bitio.Reader) (*CueSheet, error) {
	cueSheet := &CueSheet{}

	// media catalog number
	mediaCatalogNumber := make([]byte, 128)
	_, err := reader.Read(mediaCatalogNumber)
	if err != nil {
		return cueSheet, err
	}
	cueSheet.MediaCatalogNumber = string(mediaCatalogNumber)

	// number of Lead-In samples
	cueSheet.NumberOfLeadInSamples, err = reader.ReadBits(64)
	if err != nil {
		return cueSheet, err
	}

	// is compact disc
	cueSheet.CompactDisc, err = reader.ReadBool()
	if err != nil {
		return cueSheet, err
	}

	// reserved
	cueSheet.Reserved = make([]byte, 7+258*8)
	_, err = reader.Read(cueSheet.Reserved)
	if err != nil {
		return cueSheet, err
	}

	// number of tracks
	numberOfTracks, err := reader.ReadBits(8)
	if err != nil {
		return cueSheet, err
	}
	cueSheet.NumberOfTracks = uint8(numberOfTracks)

	for i := uint8(0); i < cueSheet.NumberOfTracks; i++ {
		cueSheetTrack, err := readCueSheetTrack(reader)
		if err != nil {
			return cueSheet, err
		}
		cueSheet.CueSheetTracks = append(cueSheet.CueSheetTracks, *cueSheetTrack)
	}

	return cueSheet, nil
}

func readCueSheetTrack(reader *bitio.Reader) (*CueSheetTrack, error) {
	cueSheetTrack := &CueSheetTrack{}
	var err error

	// Track offset in samples
	cueSheetTrack.OffsetInSamples, err = reader.ReadBits(64)
	if err != nil {
		return cueSheetTrack, err
	}

	// track number
	trackNumber, err := reader.ReadBits(8)
	if err != nil {
		return cueSheetTrack, err
	}
	cueSheetTrack.TrackNumber = uint8(trackNumber)

	// ISRC
	isrc := make([]byte, 12)
	_, err = reader.Read(isrc)
	if err != nil {
		return cueSheetTrack, nil
	}
	cueSheetTrack.ISRC = string(isrc)

	// non audio type
	cueSheetTrack.NonAudioType, err = reader.ReadBool()
	if err != nil {
		return cueSheetTrack, nil
	}

	// Pre-emphasis
	cueSheetTrack.PreEmphasis, err = reader.ReadBool()
	if err != nil {
		return cueSheetTrack, nil
	}

	// reserved
	cueSheetTrack.Reserved = make([]byte, 6+13*8)
	_, err = reader.Read(cueSheetTrack.Reserved)
	if err != nil {
		return cueSheetTrack, err
	}

	// The number of track index points
	indexPoints, err := reader.ReadBits(8)
	if err != nil {
		return cueSheetTrack, err
	}
	cueSheetTrack.NumberOfTrackIndexPoint = uint8(indexPoints)

	for i := uint8(0); i < cueSheetTrack.NumberOfTrackIndexPoint; i++ {
		cueSheetTrackIndex, err := readCueSheetTrackIndex(reader)
		if err != nil {
			return cueSheetTrack, err
		}
		cueSheetTrack.CueSheetTrackIndexes = append(cueSheetTrack.CueSheetTrackIndexes, *cueSheetTrackIndex)
	}

	return cueSheetTrack, nil
}

func readCueSheetTrackIndex(reader *bitio.Reader) (*CueSheetTrackIndex, error) {
	cueSheetTrackIndex := &CueSheetTrackIndex{}
	var err error

	// offset in samples
	cueSheetTrackIndex.OffsetInSamples, err = reader.ReadBits(64)
	if err != nil {
		return cueSheetTrackIndex, err
	}

	// index Point Number
	indexPointNumber, err := reader.ReadBits(8)
	if err != nil {
		return cueSheetTrackIndex, err
	}
	cueSheetTrackIndex.IndexPointNumber = uint8(indexPointNumber)

	// reserved
	cueSheetTrackIndex.Reserved = make([]byte, 3*8)
	_, err = reader.Read(cueSheetTrackIndex.Reserved)
	if err != nil {
		return cueSheetTrackIndex, err
	}

	return cueSheetTrackIndex, nil
}

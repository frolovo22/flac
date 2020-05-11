package meta

import (
	"errors"
	"github.com/icza/bitio"
)

type SeekTable struct {
	SeekPoints []SeekPoint
}

type SeekPoint struct {
	SampleNumberOfFirstSample uint64
	Offset                    uint64
	NumberOfSamples           uint16
}

func readSeekTable(reader *bitio.Reader, size int) (*SeekTable, error) {
	seekTable := &SeekTable{}

	// each seekpoint 18 bytes
	if size%18 != 0 {
		return seekTable, errors.New("incorrect SEEKTABLE size")
	}
	numberSeekPoint := size / 18
	for i := 0; i < numberSeekPoint; i++ {
		seekPoint, err := readSeekPoint(reader)
		if err != nil {
			return seekTable, err
		}
		seekTable.SeekPoints = append(seekTable.SeekPoints, *seekPoint)
	}

	return seekTable, nil
}

func readSeekPoint(reader *bitio.Reader) (*SeekPoint, error) {
	seekPoint := &SeekPoint{}

	// 64 bits per sample numbers of first sample
	sampleNumbers, err := reader.ReadBits(64)
	if err != nil {
		return seekPoint, err
	}
	seekPoint.SampleNumberOfFirstSample = sampleNumbers

	// 64 bits per offset
	offset, err := reader.ReadBits(64)
	if err != nil {
		return seekPoint, err
	}
	seekPoint.Offset = offset

	// 16 bits per number of samples
	numberOfSamples, err := reader.ReadBits(16)
	if err != nil {
		return seekPoint, err
	}
	seekPoint.NumberOfSamples = uint16(numberOfSamples)

	return seekPoint, nil
}

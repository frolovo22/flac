package meta

import (
	"errors"
	"github.com/icza/bitio"
	"io"
)

/*
	FLAC specifies a minimum block size of 16 and a maximum block size of 65535,
	meaning the bit patterns corresponding to the numbers 0-15 in the minimum blocksize and maximum blocksize fields are invalid
*/
type StreamInfo struct {
	MinimumBlockSize     uint16 // The minimum block size (in samples) used in the stream.
	MaximumBlockSize     uint16 // The maximum block size (in samples) used in the stream. (Minimum blocksize == maximum blocksize) implies a fixed-blocksize stream.
	MinimumFrameSize     uint32 // The minimum frame size (in bytes) used in the stream. May be 0 to imply the value is not known.
	MaximumFrameSize     uint32 // The maximum frame size (in bytes) used in the stream. May be 0 to imply the value is not known.
	SampleRate           uint32 // Sample rate in Hz. Though 20 bits are available, the maximum sample rate is limited by the structure of frame headers to 655350Hz. Also, a value of 0 is invalid.
	NumberOfChannels     uint8  // (number of channels)-1. FLAC supports from 1 to 8 channels
	BitsPerSample        uint8  // (bits per sample)-1. FLAC supports from 4 to 32 bits per sample. Currently the reference encoder and decoders only support up to 24 bits per sample.
	TotalSamplesInStream uint32 // Total samples in stream. 'Samples' means inter-channel sample, i.e. one second of 44.1Khz audio will have 44100 samples regardless of the number of channels. A value of zero here means the number of total samples is unknown.
	MD5                  []byte // MD5 signature of the unencoded audio data. This allows the decoder to determine if an error exists in the audio data even when the error does not result in an invalid bitstream.
}

func (si *StreamInfo) check() error {
	if si.MinimumBlockSize > si.MaximumBlockSize {
		return errors.New("minimum block size grater than maximum block size")
	}
	if si.SampleRate == 0 {
		return errors.New("invalid sample rate")
	}
	if si.NumberOfChannels < 1 || si.NumberOfChannels > 8 {
		return errors.New("invalid number of channels")
	}
	if si.BitsPerSample < 4 || si.BitsPerSample > 32 {
		return errors.New("invalid bits per sample")
	}

	return nil
}

func readStreamInfo(reader io.Reader) (*StreamInfo, error) {
	si := &StreamInfo{}

	bits := bitio.NewReader(reader)

	// 16 bits per minimum block size
	minimumBlockSize, err := bits.ReadBits(16)
	if err != nil {
		return si, err
	}
	si.MinimumBlockSize = uint16(minimumBlockSize)

	// 16 bits per maximum block size
	maximumBlockSize, err := bits.ReadBits(16)
	if err != nil {
		return si, err
	}
	si.MaximumBlockSize = uint16(maximumBlockSize)

	// 24 bits per minimum frame size
	minimumFrameSize, err := bits.ReadBits(24)
	if err != nil {
		return si, err
	}
	si.MinimumFrameSize = uint32(minimumFrameSize)

	// 24 bits per maximum frame size
	maximumFrameSize, err := bits.ReadBits(24)
	if err != nil {
		return si, err
	}
	si.MaximumFrameSize = uint32(maximumFrameSize)

	// 20 bit per SampleRate
	sampleRate, err := bits.ReadBits(20)
	if err != nil {
		return si, err
	}
	si.SampleRate = uint32(sampleRate)

	// 3 bits per number of channels
	numberOfChannels, err := bits.ReadBits(3)
	if err != nil {
		return si, err
	}
	si.NumberOfChannels = uint8(numberOfChannels) + 1

	// 5 bits per bits per sample
	bitsPerSample, err := bits.ReadBits(5)
	if err != nil {
		return si, err
	}
	si.BitsPerSample = uint8(bitsPerSample) + 1

	// 36 bits per total samples in stream
	totalSamplesInStream, err := bits.ReadBits(36)
	if err != nil {
		return si, err
	}
	si.TotalSamplesInStream = uint32(totalSamplesInStream)

	// 128 bits (16 bytes) per MD5 signature
	si.MD5 = make([]byte, 16)
	_, err = reader.Read(si.MD5)
	if err != nil {
		return si, err
	}

	return si, si.check()
}

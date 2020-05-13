package frame

import (
	"github.com/icza/bitio"
)

type FrameHeader struct {
	// sync code
	// (14 bits) always 1111 1111 1111 10 (16382)
	SyncCode uint16

	// Reserved:
	// 0 : mandatory value
	// 1 : reserved for future use
	Reserved Reserved

	// Blocking strategy:
	// 0 : fixed-blocksize stream; frame header encodes the frame number
	// 1 : variable-blocksize stream; frame header encodes the sample number
	BlockingStrategy BlockingStrategy

	// Block size in inter-channel samples:
	// 0000 : reserved
	// 0001 : 192 samples
	// 0010-0101 : 576 * (2^(n-2)) samples, i.e. 576/1152/2304/4608
	// 0110 : get 8 bit (blocksize-1) from end of header
	// 0111 : get 16 bit (blocksize-1) from end of header
	// 1000-1111 : 256 * (2^(n-8)) samples, i.e. 256/512/1024/2048/4096/8192/16384/32768
	BlockSize BlockSize

	// Sample rate:
	// 0000 : get from STREAMINFO metadata block
	// 0001 : 88.2kHz
	// 0010 : 176.4kHz
	// 0011 : 192kHz
	// 0100 : 8kHz
	// 0101 : 16kHz
	// 0110 : 22.05kHz
	// 0111 : 24kHz
	// 1000 : 32kHz
	// 1001 : 44.1kHz
	// 1010 : 48kHz
	// 1011 : 96kHz
	// 1100 : get 8 bit sample rate (in kHz) from end of header
	// 1101 : get 16 bit sample rate (in Hz) from end of header
	// 1110 : get 16 bit sample rate (in tens of Hz) from end of header
	// 1111 : invalid, to prevent sync-fooling string of 1s
	SampleRate SampleRate

	// Channel assignment
	// 0000-0111 : (number of independent channels)-1. Where defined, the channel order follows SMPTE/ITU-R recommendations. The assignments are as follows:
	// 1 channel: mono
	// 2 channels: left, right
	// 3 channels: left, right, center
	// 4 channels: front left, front right, back left, back right
	// 5 channels: front left, front right, front center, back/surround left, back/surround right
	// 6 channels: front left, front right, front center, LFE, back/surround left, back/surround right
	// 7 channels: front left, front right, front center, LFE, back center, side left, side right
	// 8 channels: front left, front right, front center, LFE, back left, back right, side left, side right
	// 1000 : left/side stereo: channel 0 is the left channel, channel 1 is the side(difference) channel
	// 1001 : right/side stereo: channel 0 is the side(difference) channel, channel 1 is the right channel
	// 1010 : mid/side stereo: channel 0 is the mid(average) channel, channel 1 is the side(difference) channel
	// 1011-1111 : reserved
	ChannelAssigment ChannelAssigment

	// Sample size in bits:
	// 000 : get from STREAMINFO metadata block
	// 001 : 8 bits per sample
	// 010 : 12 bits per sample
	// 011 : reserved
	// 100 : 16 bits per sample
	// 101 : 20 bits per sample
	// 110 : 24 bits per sample
	// 111 : reserved
	SampleSize SampleSize

	// Reserved:
	// 0 : mandatory value
	// 1 : reserved for future use
	Reserved2 Reserved

	// if(variable blocksize)
	//   <8-56>:"UTF-8" coded sample number (decoded number is 36 bits)
	// else
	//   <8-48>:"UTF-8" coded frame number (decoded number is 31 bits)
	VariableBlockSize uint64

	// if(blocksize bits == 011x)
	//   8/16 bit (blocksize-1)
	BlockSizeEnd uint16

	// if(sample rate bits == 11xx)
	//   8/16 bit sample rate
	SampleRateEnd uint16

	// CRC-8 (polynomial = x^8 + x^2 + x^1 + x^0, initialized with 0) of everything before the crc, including the sync code
	CRC8 uint8
}

type Reserved bool
type BlockingStrategy bool
type BlockSize uint8
type SampleRate uint8
type ChannelAssigment uint8
type SampleSize uint8

const (
	MandatoryValue       Reserved = false
	ReservedForFutureUse Reserved = true

	FixedBlockSizeStream    BlockingStrategy = false
	VariableBlockSizeStream BlockingStrategy = true
)

// return block size
func (bs BlockSize) BlockSize() uint32 {
	switch bs {
	case 0:
		return 0
	case 1:
		return 192
	case 2:
		return 576
	case 3:
		return 1152
	case 4:
		return 2304
	case 5:
		return 4608
	case 6:
		return 0
	case 7:
		return 0
	case 8:
		return 256
	case 9:
		return 512
	case 10:
		return 1024
	case 11:
		return 2048
	case 12:
		return 4096
	case 13:
		return 8192
	case 14:
		return 16384
	case 15:
		return 32768
	}
	return 0
}

// return sample rate in kHz
func (sr SampleRate) SampleRate() float32 {
	switch sr {
	case 0:
		return 0
	case 1:
		return 88.2
	case 2:
		return 176.4
	case 3:
		return 192
	case 4:
		return 8
	case 5:
		return 16
	case 6:
		return 22.05
	case 7:
		return 24
	case 8:
		return 32
	case 9:
		return 44.1
	case 10:
		return 48
	case 11:
		return 96
	}
	return 0
}

func readFrameHeader(reader *bitio.Reader) (*FrameHeader, error) {
	header := &FrameHeader{}

	// always 11 1111 1111 1110
	syncCode, err := reader.ReadBits(14)
	if err != nil {
		return header, err
	}
	header.SyncCode = uint16(syncCode)

	// reserved
	reserved, err := reader.ReadBool()
	if err != nil {
		return header, err
	}
	header.Reserved = Reserved(reserved)

	// blocking strategy
	blockingStrategy, err := reader.ReadBool()
	if err != nil {
		return header, err
	}
	header.BlockingStrategy = BlockingStrategy(blockingStrategy)

	// block size
	blockSize, err := reader.ReadBits(4)
	if err != nil {
		return header, err
	}
	header.BlockSize = BlockSize(blockSize)

	// sample rate
	sampleRate, err := reader.ReadBits(4)
	if err != nil {
		return header, err
	}
	header.SampleRate = SampleRate(sampleRate)

	// channel assigment
	channelAssigment, err := reader.ReadBits(4)
	if err != nil {
		return header, err
	}
	header.ChannelAssigment = ChannelAssigment(channelAssigment)

	// sample size
	sampleSize, err := reader.ReadBits(3)
	if err != nil {
		return header, err
	}
	header.SampleSize = SampleSize(sampleSize)

	// reserved
	reserved2, err := reader.ReadBool()
	if err != nil {
		return header, err
	}
	header.Reserved2 = Reserved(reserved2)

	// variable block size
	//   <8-56>:"UTF-8" coded sample number (decoded number is 36 bits)
	// else
	//   <8-48>:"UTF-8" coded frame number (decoded number is 31 bits)
	if header.BlockingStrategy == VariableBlockSizeStream {
		header.VariableBlockSize, err = reader.ReadBits(36)
		if err != nil {
			return header, err
		}
	}

	if header.BlockingStrategy == FixedBlockSizeStream {
		header.VariableBlockSize, err = reader.ReadBits(31)
		if err != nil {
			return header, err
		}
	}

	// block size
	// 0110 : get 8 bit (blocksize-1) from end of header
	// 0111 : get 16 bit (blocksize-1) from end of header
	var sizeBlockSize uint8
	switch header.BlockSize {
	case 6:
		sizeBlockSize = 8
	case 7:
		sizeBlockSize = 16
	}
	if sizeBlockSize > 0 {
		blockSizeEnd, err := reader.ReadBits(sizeBlockSize)
		if err != nil {
			return header, err
		}
		header.BlockSizeEnd = uint16(blockSizeEnd)
	}

	// sample rate
	// 1100 : get 8 bit sample rate (in kHz) from end of header
	// 1101 : get 16 bit sample rate (in Hz) from end of header
	// 1110 : get 16 bit sample rate (in tens of Hz) from end of header
	var sizeSampleRate uint8
	switch header.SampleRate {
	case 12:
		sizeSampleRate = 8
	case 13:
		sizeSampleRate = 16
	case 14:
		sizeSampleRate = 16
	}

	if sizeSampleRate > 0 {
		sampleRateEnd, err := reader.ReadBits(sizeSampleRate)
		if err != nil {
			return header, err
		}
		header.SampleRateEnd = uint16(sampleRateEnd)
	}

	// CRC-8
	crc, err := reader.ReadBits(8)
	if err != nil {
		return header, err
	}
	header.CRC8 = uint8(crc)

	return header, nil
}

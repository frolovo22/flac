package frame

import "github.com/icza/bitio"

type Frame struct {
	Header FrameHeader
}

func ReadFrame(reader *bitio.Reader) (*Frame, error) {
	frame := &Frame{}

	// header
	header, err := readFrameHeader(reader)
	if err != nil {
		return frame, err
	}
	frame.Header = *header

	return frame, nil
}

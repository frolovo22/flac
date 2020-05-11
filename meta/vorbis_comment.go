package meta

import (
	"encoding/binary"
	"errors"
	"github.com/icza/bitio"
	"strings"
)

type VorbisComment struct {
	VendorLength       uint32
	VendorString       string
	UserCommentsLength uint32
	UserComments       []UserComment
}

type UserComment struct {
	Length uint32
	Key    string
	Value  string
}

// The comment header is decoded as follows:
//
//	1) [vendor_length] = read an unsigned integer of 32 bits
//	2) [vendor_string] = read a UTF-8 vector as [vendor_length] octets
//	3) [user_comment_list_length] = read an unsigned integer of 32 bits
//	4) iterate [user_comment_list_length] times {
//
//		5) [length] = read an unsigned integer of 32 bits
//		6) this iteration's user comment = read a UTF-8 vector as [length] octets
//
//	}
//
//	7) [framing_bit] = read a single bit as boolean
func readVorbisComment(reader *bitio.Reader) (*VorbisComment, error) {
	vorbisComment := &VorbisComment{}

	err := binary.Read(reader, binary.LittleEndian, &vorbisComment.VendorLength)
	if err != nil {
		return vorbisComment, err
	}

	vendorString := make([]byte, vorbisComment.VendorLength)
	_, err = reader.Read(vendorString)
	if err != nil {
		return vorbisComment, err
	}
	vorbisComment.VendorString = string(vendorString)

	err = binary.Read(reader, binary.LittleEndian, &vorbisComment.UserCommentsLength)
	if err != nil {
		return vorbisComment, err
	}

	for i := uint32(0); i < vorbisComment.UserCommentsLength; i++ {
		userComment, err := readUserComment(reader)
		if err != nil {
			return vorbisComment, err
		}
		vorbisComment.UserComments = append(vorbisComment.UserComments, *userComment)
	}

	return vorbisComment, nil
}

func readUserComment(reader *bitio.Reader) (*UserComment, error) {
	userComment := &UserComment{}

	err := binary.Read(reader, binary.LittleEndian, &userComment.Length)
	if err != nil {
		return userComment, err
	}

	userString := make([]byte, userComment.Length)
	_, err = reader.Read(userString)
	if err != nil {
		return userComment, err
	}

	comment := strings.SplitN(string(userString), "=", 2)
	if len(comment) != 2 {
		return userComment, errors.New("error vorbis comment format")
	}

	userComment.Key = comment[0]
	userComment.Value = comment[1]
	return userComment, nil
}

package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"
)

type Charset []byte

var (
	lowerCharset   = Charset("abcdefghijklmnopqrstuvwxyz")
	upperCharset   = Charset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialCharset = Charset("~=+%^*/()[]{}!@#$?|")
	numberSet      = Charset("0123456789")
)

var (
	ErrTooShortLength = errors.New("too short length")
)

const (
	partUpperChars   = 15
	partSpecialChars = 15
	partNumberChars  = 15
	minLength        = 12
	chunkSize        = 6
)

// CreatePassword generates password with a given length
func CreatePassword(length int) (string, error) {
	var (
		err         error
		b           []byte
		password    string
		chosenBytes []byte
	)

	if length <= minLength {
		return "", ErrTooShortLength
	}
	chosenBytes = make([]byte, 0, length)

	lengthUpperChars := length * partUpperChars / 100
	lengthSpecialChars := length * partSpecialChars / 100
	lengthNumberChars := length * partNumberChars / 100
	lengthLowerCharset := length - lengthUpperChars - lengthSpecialChars - lengthNumberChars

	if b, err = chooseCharsFromCharset(lengthLowerCharset, lowerCharset); err != nil {
		return "", err
	}
	chosenBytes = append(chosenBytes, b...)
	if b, err = chooseCharsFromCharset(lengthUpperChars, upperCharset); err != nil {
		return "", err
	}
	chosenBytes = append(chosenBytes, b...)
	if b, err = chooseCharsFromCharset(lengthSpecialChars, specialCharset); err != nil {
		return "", err
	}
	chosenBytes = append(chosenBytes, b...)
	if b, err = chooseCharsFromCharset(lengthNumberChars, numberSet); err != nil {
		return "", err
	}
	chosenBytes = append(chosenBytes, b...)
	shuffleBytes(chosenBytes)

	if password, err = chunkString(chosenBytes, length, chunkSize); err != nil {
		return "", err
	}

	return password, nil
}

func chooseCharsFromCharset(length int, chars Charset) ([]byte, error) {
	var (
		b   byte
		res = make([]byte, 0, length)
		err error
	)
	for i := 0; i < length; i++ {
		if b, err = randomlyChooseChar(chars); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, nil
}

func randomlyChooseChar(chars Charset) (byte, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, fmt.Errorf("failed to choose %v from charset %v", r, chars)
	}
	return chars[byte(r.Int64())], nil
}

func shuffleBytes(in []byte) {
	s := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(s)

	r.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})
}

func chunkString(in []byte, length, chunkSize int) (string, error) {
	var (
		err error
		res strings.Builder
	)

	if length <= chunkSize {
		return string(in), nil
	}

	for i, v := range in {
		if i > 0 && (i)%chunkSize == 0 {
			if err = res.WriteByte('-'); err != nil {
				return "", err
			}
		}
		if err = res.WriteByte(v); err != nil {
			return "", err
		}
	}
	return res.String(), nil
}

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
	partUpperChars   = 15
	partSpecialChars = 15
	partNumberChars  = 15

	minLength = 12
	chunkSize = 6
)

var (
	ErrTooShortLength = errors.New("too short length")
)

// CreatePassword generates password with a given length
func CreatePassword(length int) (string, error) {
	var (
		err      error
		password strings.Builder
	)

	if length <= minLength {
		return "", ErrTooShortLength
	}

	lengthUpperChars := length * partUpperChars / 100
	lengthSpecialChars := length * partSpecialChars / 100
	lengthNumberChars := length * partNumberChars / 100
	lengthLowerCharset := length - lengthUpperChars - lengthSpecialChars - lengthNumberChars

	if err = chooseCharsFromCharset(lengthLowerCharset, lowerCharset, &password); err != nil {
		return "", err
	}

	if err = chooseCharsFromCharset(lengthUpperChars, upperCharset, &password); err != nil {
		return "", err
	}

	if err = chooseCharsFromCharset(lengthSpecialChars, specialCharset, &password); err != nil {
		return "", err
	}

	if err = chooseCharsFromCharset(lengthNumberChars, numberSet, &password); err != nil {
		return "", err
	}

	if password, err = shuffleString([]byte(password.String())); err != nil {
		return "", err
	}

	password, err = chunkString(length, password)
	if err != nil {
		return "", err
	}

	return password.String(), nil
}

func chooseCharsFromCharset(length int, chars Charset, password *strings.Builder) error {
	for i := 0; i < length; i++ {
		b, err := randomlyChooseChar(chars)
		if err != nil {
			return err
		}
		password.WriteByte(b)
	}
	return nil
}

func randomlyChooseChar(chars Charset) (byte, error) {
	r, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		return 0, fmt.Errorf("failed to choose %v from charset %v", r, chars)
	}
	return chars[byte(r.Int64())], nil
}

func shuffleString(in []byte) (strings.Builder, error) {
	var (
		err error
		out strings.Builder
	)
	s := mathRand.NewSource(time.Now().UnixNano())
	r := mathRand.New(s)

	r.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})

	if _, err = out.WriteString(string(in)); err != nil {
		return out, err
	}
	return out, nil
}

func chunkString(lengthAll int, in strings.Builder) (strings.Builder, error) {
	var (
		err error
		res strings.Builder
	)

	if lengthAll <= chunkSize {
		return in, nil
	}

	for i, v := range in.String() {
		if i > 0 && (i)%chunkSize == 0 {
			if err = res.WriteByte('-'); err != nil {
				return res, err
			}
		}
		if err = res.WriteByte(byte(v)); err != nil {
			return res, err
		}
	}
	return res, nil
}

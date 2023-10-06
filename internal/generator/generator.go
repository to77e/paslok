package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type Charset []byte

var (
	lowerCharset   = Charset("abcdefghijklmnopqrstuvwxyz")
	upperCharset   = Charset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialCharset = Charset("~=+%^*/()[]{}!@#$?|")
	numberSet      = Charset("0123456789")
)

const (
	partAllChars     = 100
	partUpperChars   = 15
	partSpecialChars = 15
	partNumberChars  = 15
	chunkSize        = 6
)

func CreatePassword(length int) (string, error) {
	var (
		err         error
		b           []byte
		password    string
		chosenBytes []byte
	)

	chosenBytes = make([]byte, 0, length)

	lengthUpperChars := length * partUpperChars / partAllChars
	lengthSpecialChars := length * partSpecialChars / partAllChars
	lengthNumberChars := length * partNumberChars / partAllChars
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
	err = shuffleBytes(chosenBytes)
	if err != nil {
		return "", err
	}

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

func shuffleBytes(in []byte) error {
	for i := range in {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(in)-i)))
		if err != nil {
			return fmt.Errorf("choose %v from charset %v", j, in)
		}
		jInt := j.Int64()
		in[i], in[i+int(jInt)] = in[i+int(jInt)], in[i]
	}
	return nil
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
				return "", fmt.Errorf("write \"-\": %w", err)
			}
		}
		if err = res.WriteByte(v); err != nil {
			return "", fmt.Errorf("write byte: %w", err)
		}
	}
	return res.String(), nil
}

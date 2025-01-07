package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/to77e/paslok/internal/models"
)

type Charset []byte

var (
	lowerCharset   = Charset("abcdefghijklmnopqrstuvwxyz")
	upperCharset   = Charset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialCharset = Charset("~=+%^*/()[]{}!@#$?|")
	numberSet      = Charset("0123456789")
)

const (
	percentTotal   = 100
	percentUpper   = 15
	percentSpecial = 15
	percentNumber  = 15
)

func CreatePassword(pswd *models.Password) (string, error) {
	var (
		lengthUpperChars   int
		lengthSpecialChars int
		lengthNumberChars  int
		lengthLowerCharset int
		bytes              []byte
		err                error
	)
	chosenBytes := make([]byte, 0, pswd.Length)

	if pswd.Uppercase {
		lengthUpperChars = pswd.Length * percentUpper / percentTotal
		bytes, err = chooseCharsFromCharset(lengthUpperChars, upperCharset)
		if err != nil {
			return "", fmt.Errorf("choose upper charset: %w", err)
		}
		chosenBytes = append(chosenBytes, bytes...)
	}

	if pswd.Special {
		lengthSpecialChars = pswd.Length * percentSpecial / percentTotal
		bytes, err = chooseCharsFromCharset(lengthSpecialChars, specialCharset)
		if err != nil {
			return "", fmt.Errorf("choose special charset: %w", err)
		}
		chosenBytes = append(chosenBytes, bytes...)
	}

	if pswd.Number {
		lengthNumberChars = pswd.Length * percentNumber / percentTotal
		bytes, err = chooseCharsFromCharset(lengthNumberChars, numberSet)
		if err != nil {
			return "", fmt.Errorf("choose number charset: %w", err)
		}
		chosenBytes = append(chosenBytes, bytes...)
	}

	lengthLowerCharset = pswd.Length - lengthUpperChars - lengthSpecialChars - lengthNumberChars
	bytes, err = chooseCharsFromCharset(lengthLowerCharset, lowerCharset)
	if err != nil {
		return "", fmt.Errorf("choose lower charset: %w", err)
	}
	chosenBytes = append(chosenBytes, bytes...)

	err = shuffleBytes(chosenBytes)
	if err != nil {
		return "", fmt.Errorf("shuffle bytes: %w", err)
	}

	password, err := chunkString(chosenBytes, pswd)
	if err != nil {
		return "", fmt.Errorf("chunk string: %w", err)
	}

	return password, nil
}

func chooseCharsFromCharset(length int, chars Charset) ([]byte, error) {
	var res []byte
	for i := 0; i < length; i++ {
		b, err := randomlyChooseChar(chars)
		if err != nil {
			return nil, fmt.Errorf("choose char from charset: %w", err)
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

func chunkString(in []byte, pswd *models.Password) (string, error) {
	var res strings.Builder
	if pswd.Length <= pswd.ChunkSize {
		return string(in), nil
	}
	for i, v := range in {
		if pswd.Dash && (i > 0 && (i)%pswd.ChunkSize == 0) {
			if err := res.WriteByte('-'); err != nil {
				return "", fmt.Errorf("write \"-\": %w", err)
			}
		}
		if err := res.WriteByte(v); err != nil {
			return "", fmt.Errorf("write byte: %w", err)
		}
	}
	return res.String(), nil
}

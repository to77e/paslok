package generator

import (
	//"crypto/rand"
	"errors"
	"math/rand"
	"strings"
	"time"
)

type Charset rune

var (
	lowerCharset   = []Charset("abcdefghijklmnopqrstuvwxyz")
	upperCharset   = []Charset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialCharset = []Charset("~=+%^*/()[]{}/!@#$?|")
	numberSet      = []Charset("0123456789")
)

var (
	partUpperChars   = 15
	partSpecialChars = 15
	partNumberChars  = 15

	minLenght = 12
	chunkSize = 6
)

func CreatePassword(lengthAll int) (string, error) {

	if lengthAll <= minLenght {
		return "", errors.New("too short length")
	}

	var password strings.Builder
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	lengthUpperChars := lengthAll * partUpperChars / 100
	lengthspecialChars := lengthAll * partSpecialChars / 100
	lengthNunberChars := lengthAll * partNumberChars / 100
	lengthLowerCharset := lengthAll - lengthUpperChars - lengthspecialChars - lengthNunberChars

	for i := 0; i < lengthLowerCharset; i++ {
		v := lowerCharset[r.Intn(len(lowerCharset))]
		if _, err := password.WriteRune(rune(v)); err != nil {
			return "", err
		}
	}

	for i := 0; i < lengthUpperChars; i++ {
		v := upperCharset[r.Intn(len(upperCharset))]
		if _, err := password.WriteRune(rune(v)); err != nil {
			return "", err
		}
	}

	for i := 0; i < lengthspecialChars; i++ {
		v := specialCharset[r.Intn(len(specialCharset))]
		if _, err := password.WriteRune(rune(v)); err != nil {
			return "", err
		}
	}

	for i := 0; i < lengthNunberChars; i++ {
		v := numberSet[r.Intn(len(numberSet))]
		if _, err := password.WriteRune(rune(v)); err != nil {
			return "", err
		}
	}

	password, err := shuffleString([]rune(password.String()), r)
	if err != nil {
		return "", err
	}

	password, err = chukString(lengthAll, password)
	if err != nil {
		return "", err
	}

	return password.String(), nil
}

func shuffleString(in []rune, r *rand.Rand) (strings.Builder, error) {
	var out strings.Builder
	r.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})

	if _, err := out.WriteString(string(in)); err != nil {
		return out, err
	}
	return out, nil
}

func chukString(lengthAll int, in strings.Builder) (strings.Builder, error) {
	if lengthAll <= chunkSize {
		return in, nil
	}

	var res strings.Builder
	for i, v := range in.String() {
		if i > 0 && (i)%chunkSize == 0 {
			if _, err := res.WriteRune(rune('-')); err != nil {
				return res, err
			}
		}
		if _, err := res.WriteRune(rune(v)); err != nil {
			return res, err
		}
	}

	return res, nil
}

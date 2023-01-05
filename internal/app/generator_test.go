package app

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreatePassword(t *testing.T) {
	t.Run("generating passwords with different lengths", func(t *testing.T) {
		type testCase struct {
			length int
		}

		testCases := []testCase{{length: 18}, {length: 12}, {length: 0}, {length: -1}}
		for _, v := range testCases {
			actual, err := CreatePassword(v.length)
			if err != nil {
				if err == ErrTooShortLength {
					assert.Empty(t, actual)
					continue
				}
				t.Fatalf("test for checking for creating password is failed")
			}
			assert.NotEmpty(t, actual)
			assert.Len(t, actual, v.length+((v.length/chunkSize)-1))
			assert.True(t, strings.ContainsAny(actual, string(lowerCharset)))
			assert.True(t, strings.ContainsAny(actual, string(upperCharset)))
			assert.True(t, strings.ContainsAny(actual, string(specialCharset)))
			assert.True(t, strings.ContainsAny(actual, string(numberSet)))

		}
	})

	t.Run("check for uniqueness", func(t *testing.T) {
		var (
			length = 1000
			result = make(map[string]struct{}, length)
		)
		for i := 0; i < length; i++ {
			k, err := CreatePassword(18)
			if err != nil {
				t.Fatalf("test for checking for uniqueness is failed")
			}
			result[k] = struct{}{}
		}
		assert.Len(t, result, length)
	})
}

package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func Test_shuffleBytes(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "empty slice",
			args: args{
				in: []byte{},
			},
			want: []byte{},
		},
		{
			name: "single element slice",
			args: args{
				in: []byte{1},
			},
			want: []byte{1},
		},
		{
			name: "multiple element slice",
			args: args{in: []byte{1, 2, 3, 4, 5}},
			want: []byte{5, 1, 4, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]byte, len(tt.args.in))
			copy(inputCopy, tt.args.in)

			err := shuffleBytes(inputCopy)
			assert.NoError(t, err)

			if len(inputCopy) <= 1 {
				assert.EqualValues(t, tt.args.in, inputCopy)
			} else {
				assert.NotEqualValues(t, tt.args.in, inputCopy)
			}
			for _, v := range inputCopy {
				assert.Contains(t, tt.args.in, v)
			}
		})
	}
}

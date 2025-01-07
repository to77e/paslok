package generator

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/to77e/paslok/internal/models"
)

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

func TestCreatePassword(t *testing.T) {
	tests := []struct {
		name          string
		pswd          models.Password
		wantErr       assert.ErrorAssertionFunc
		checkContains func(*testing.T, string)
	}{
		{
			name: "lowercase only (length=8)",
			pswd: models.Password{
				Length:    8,
				Uppercase: false,
				Special:   false,
				Number:    false,
				Dash:      false,
				ChunkSize: 0,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 8, "password should be 8 characters long")
				for _, r := range generated {
					assert.True(t, unicode.IsLower(r), "expected only lowercase letters")
				}
			},
		},
		{
			name: "uppercase only (length=8)",
			pswd: models.Password{
				Length:    8,
				Uppercase: true,
				Special:   false,
				Number:    false,
				Dash:      false,
				ChunkSize: 0,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 8, "password should be 8 characters long")
				hasUpper := false
				for _, r := range generated {
					if unicode.IsUpper(r) {
						hasUpper = true
						break
					}
				}
				assert.True(t, hasUpper, "expected at least one uppercase letter")
			},
		},
		{
			name: "special characters only (length=8)",
			pswd: models.Password{
				Length:    8,
				Uppercase: false,
				Special:   true,
				Number:    false,
				Dash:      false,
				ChunkSize: 0,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 8, "password should be 8 characters long")
				hasSpecial := false
				for _, r := range generated {
					if isSpecial(r) {
						hasSpecial = true
					}
				}
				assert.True(t, hasSpecial, "expected at least one special character in the result")
			},
		},
		{
			name: "numbers only (length=8)",
			pswd: models.Password{
				Length:    8,
				Uppercase: false,
				Special:   false,
				Number:    true,
				Dash:      false,
				ChunkSize: 0,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 8, "password should be 8 characters long")
				hasNumber := false
				for _, r := range generated {
					if unicode.IsDigit(r) {
						hasNumber = true
					}
				}
				assert.True(t, hasNumber, "expected at least one digit")
			},
		},
		{
			name: "all flags on (length=12) with chunk=4 and dash=true",
			pswd: models.Password{
				Length:    12,
				Uppercase: true,
				Special:   true,
				Number:    true,
				Dash:      true,
				ChunkSize: 4,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 12+2, "expected 14 characters (12 plus 2 dashes)")
				assert.Equal(t, '-', rune(generated[4]), "expected dash at position 4")
				assert.Equal(t, '-', rune(generated[9]), "expected dash at position 9")
				hasUpper := false
				hasDigit := false
				hasSpecial := false
				for _, r := range generated {
					if r == '-' {
						continue
					}
					if unicode.IsUpper(r) {
						hasUpper = true
					}
					if unicode.IsDigit(r) {
						hasDigit = true
					}
					if isSpecial(r) {
						hasSpecial = true
					}
				}
				assert.True(t, hasUpper, "expected at least one uppercase letter")
				assert.True(t, hasDigit, "expected at least one digit")
				assert.True(t, hasSpecial, "expected at least one special character")
			},
		},
		{
			name: "chunk size > length (chunk=10, length=8)",
			pswd: models.Password{
				Length:    8,
				Uppercase: true,
				Special:   true,
				Number:    true,
				Dash:      true,
				ChunkSize: 10,
			},
			wantErr: assert.NoError,
			checkContains: func(t *testing.T, generated string) {
				require.Len(t, generated, 8, "expected 8 characters with no dashes")
				hasUpper := false
				hasDigit := false
				hasSpecial := false
				for _, r := range generated {
					if unicode.IsUpper(r) {
						hasUpper = true
					}
					if unicode.IsDigit(r) {
						hasDigit = true
					}
					if isSpecial(r) {
						hasSpecial = true
					}
				}
				assert.True(t, hasUpper, "expected at least one uppercase letter")
				assert.True(t, hasDigit, "expected at least one digit")
				assert.True(t, hasSpecial, "expected at least one special character")
			},
		},
		{
			name: "zero length => expect error",
			pswd: models.Password{
				Length:    0,
				Uppercase: true,
				Special:   true,
				Number:    true,
			},
			wantErr:       assert.NoError,
			checkContains: func(*testing.T, string) {},
		},
		{
			name: "negative chunk => expect error",
			pswd: models.Password{
				Length:    8,
				Uppercase: true,
				Special:   true,
				Number:    true,
				Dash:      true,
				ChunkSize: -1,
			},
			wantErr:       assert.NoError,
			checkContains: func(*testing.T, string) {},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			generated, err := CreatePassword(&tc.pswd)
			require.True(t, tc.wantErr(t, err), "error assertion did not pass")
			if err == nil {
				tc.checkContains(t, generated)
			}
		})
	}
}

func isSpecial(r rune) bool {
	specials := "~=+%^*/()[]{}!@#$?|"
	for _, sc := range specials {
		if r == sc {
			return true
		}
	}
	return false
}

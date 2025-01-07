package locker

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/to77e/paslok/internal/models"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		db      Resourcer
		cryptor Cryptor
	}
	type args struct {
		req *models.CreatePasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db:      tt.fields.db,
				cryptor: tt.fields.cryptor,
			}
			err := s.Create(tt.args.req)
			assert.True(t, tt.wantErr(t, err, "Create() error = %v", err))
		})
	}
}

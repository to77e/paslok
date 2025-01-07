package locker

import (
	"fmt"

	"github.com/to77e/paslok/internal/generator"
	"github.com/to77e/paslok/internal/models"
)

type Cryptor interface {
	Encrypt(password string) (string, error)
	Decrypt(encryptedPassword string) (string, error)
}

type Resourcer interface {
	Create(req *models.CreatePasswordRequest) error
	Read(req *models.ReadPasswordRequest) (string, error)
	List(req *models.ListPasswordsRequest) ([]models.Resource, error)
	Delete(req *models.DeletePasswordRequest) error
	Close() error
}

type Service struct {
	db      Resourcer
	cryptor Cryptor
}

func New(db Resourcer, cryptor Cryptor) *Service {
	return &Service{db: db, cryptor: cryptor}
}

func (s *Service) Create(req *models.CreatePasswordRequest) error {
	var err error
	if len(req.Password) == 0 {
		pswd := &models.Password{
			Length:    req.Length,
			ChunkSize: models.DefaultChunk,
			Uppercase: req.Uppercase,
			Special:   req.Special,
			Number:    req.Number,
			Dash:      req.Dash,
		}
		req.Password, err = generator.CreatePassword(pswd)
		if err != nil {
			return fmt.Errorf("generate password: %v\n", err)
		}
	}
	req.Password, err = s.cryptor.Encrypt(req.Password)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}
	req.Username, err = s.cryptor.Encrypt(req.Username)
	if err != nil {
		return fmt.Errorf("encrypt username: %w", err)
	}
	if err = s.db.Create(req); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	return nil
}

func (s *Service) Read(req *models.ReadPasswordRequest) (string, error) {
	encryptedPassword, err := s.db.Read(req)
	if err != nil {
		return "", fmt.Errorf("database: %w", err)
	}
	pass, err := s.cryptor.Decrypt(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return pass, nil
}

func (s *Service) List(req *models.ListPasswordsRequest) ([]models.Resource, error) {
	resources, err := s.db.List(req)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	for i := range resources {
		resources[i].Username, err = s.cryptor.Decrypt(resources[i].Username)
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}
	}
	return resources, nil
}

func (s *Service) Delete(req *models.DeletePasswordRequest) error {
	if err := s.db.Delete(req); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	return nil
}

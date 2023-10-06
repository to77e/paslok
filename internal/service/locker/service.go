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
	Create(name, password, comment string) error
	Read(name string) (string, error)
	List() ([]models.Resource, error)
	Update(name, password, comment string) error
	Delete(name string) error
	Close() error
}

type Service struct {
	db      Resourcer
	cryptor Cryptor
}

func New(db Resourcer, cryptor Cryptor) *Service {
	return &Service{db: db, cryptor: cryptor}
}

func (s *Service) Create(name, comment string) error {
	const (
		length = 18
	)

	password, err := generator.CreatePassword(length)
	if err != nil {
		return fmt.Errorf("failed to create password: %v\n", err)
	}

	encryptedPassword, err := s.cryptor.Encrypt(password)
	if err != nil {
		return fmt.Errorf("call encrypt: %w", err)
	}
	if err = s.db.Create(name, encryptedPassword, comment); err != nil {
		return fmt.Errorf("call create: %w", err)
	}
	return nil
}

func (s *Service) Read(name string) (string, error) {
	encryptedPassword, err := s.db.Read(name)
	if err != nil {
		return "", fmt.Errorf("call read: %w", err)
	}
	password, err := s.cryptor.Decrypt(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("call decrypt: %w", err)
	}
	return password, nil
}

func (s *Service) List() ([]models.Resource, error) {
	resources, err := s.db.List()
	if err != nil {
		return nil, fmt.Errorf("call list: %w", err)
	}
	return resources, nil
}

func (s *Service) Update(name, password, comment string) error {
	encryptedPassword, err := s.cryptor.Encrypt(password)
	if err != nil {
		return fmt.Errorf("call encrypt: %w", err)
	}
	if err = s.db.Update(name, encryptedPassword, comment); err != nil {
		return fmt.Errorf("call update: %w", err)
	}
	return nil
}

func (s *Service) Delete(name string) error {
	if err := s.db.Delete(name); err != nil {
		return fmt.Errorf("call delete: %w", err)
	}
	return nil
}

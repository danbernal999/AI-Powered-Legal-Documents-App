package services

import (
	"github.com/danbernal999/kiradoc-backend/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Services struct {
	repo *repository.Repository
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{repo: repo}
}

func (s *Services) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *Services) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *Services) CreateUser(id, email, password, name string) error {
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return err
	}
	return s.repo.CreateUser(id, email, hashedPassword, name)
}

func (s *Services) GetUserByEmail(email string) (*repository.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *Services) GetUserByID(id string) (*repository.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Services) GetTemplates() ([]*repository.Template, error) {
	return s.repo.GetTemplates()
}

func (s *Services) GetTemplateByID(id string) (*repository.Template, error) {
	return s.repo.GetTemplateByID(id)
}

func (s *Services) CreateDocument(doc *repository.Document) error {
	return s.repo.CreateDocument(doc)
}

func (s *Services) GetDocumentByID(id string) (*repository.Document, error) {
	return s.repo.GetDocumentByID(id)
}

func (s *Services) GetDocumentsByUserID(userID string) ([]*repository.Document, error) {
	return s.repo.GetDocumentsByUserID(userID)
}

func (s *Services) UpdateDocument(doc *repository.Document) error {
	return s.repo.UpdateDocument(doc)
}

func (s *Services) DeleteDocument(id string) error {
	return s.repo.DeleteDocument(id)
}

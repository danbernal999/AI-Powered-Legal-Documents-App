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

func (s *Services) CreateTemplate(t *repository.Template) error {
	return s.repo.CreateTemplate(t)
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

func (s *Services) CreateSignature(sig *repository.Signature) error {
	return s.repo.CreateSignature(sig)
}

func (s *Services) GetSignaturesByDocumentID(documentID string) ([]*repository.Signature, error) {
	return s.repo.GetSignaturesByDocumentID(documentID)
}

func (s *Services) GetSignatureByID(id string) (*repository.Signature, error) {
	return s.repo.GetSignatureByID(id)
}

func (s *Services) DeleteSignature(id string) error {
	return s.repo.DeleteSignature(id)
}

func (s *Services) CreateDocumentShare(share *repository.DocumentShare) error {
	return s.repo.CreateDocumentShare(share)
}

func (s *Services) GetDocumentShares(documentID string) ([]*repository.DocumentShare, error) {
	return s.repo.GetDocumentShares(documentID)
}

func (s *Services) GetDocumentsSharedWithUser(userID string) ([]*repository.DocumentShare, error) {
	return s.repo.GetDocumentsSharedWithUser(userID)
}

func (s *Services) UpdateDocumentShare(share *repository.DocumentShare) error {
	return s.repo.UpdateDocumentShare(share)
}

func (s *Services) DeleteDocumentShare(id string) error {
	return s.repo.DeleteDocumentShare(id)
}

func (s *Services) CreateShareLink(link *repository.ShareLink) error {
	return s.repo.CreateShareLink(link)
}

func (s *Services) GetShareLinkByToken(token string) (*repository.ShareLink, error) {
	return s.repo.GetShareLinkByToken(token)
}

func (s *Services) GetShareLinks(documentID string) ([]*repository.ShareLink, error) {
	return s.repo.GetShareLinks(documentID)
}

func (s *Services) DeleteShareLink(id string) error {
	return s.repo.DeleteShareLink(id)
}

func (s *Services) CreateDocumentComment(comment *repository.DocumentComment) error {
	return s.repo.CreateDocumentComment(comment)
}

func (s *Services) GetDocumentComments(documentID string) ([]*repository.DocumentComment, error) {
	return s.repo.GetDocumentComments(documentID)
}

func (s *Services) DeleteDocumentComment(id string) error {
	return s.repo.DeleteDocumentComment(id)
}

func (s *Services) CreateDocumentVersion(version *repository.DocumentVersion) error {
	return s.repo.CreateDocumentVersion(version)
}

func (s *Services) GetDocumentVersions(documentID string) ([]*repository.DocumentVersion, error) {
	return s.repo.GetDocumentVersions(documentID)
}

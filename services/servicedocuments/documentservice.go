package servicedocuments

import (
	"goapi/models"
	"goapi/repositories/repodocuments"
)

type DocumentService interface {
	Get(id string) (models.Document, error)
	GetAll() ([]models.Document, error)
	CreateOrUpdate(document models.Document) (bool, error)
	Delete(id string) (bool, error)
}

// DocumentServiceImpl Default implementation for DocumentService
type DocumentServiceImpl struct {
	documentRepo repodocuments.DocumentRepository
}

func NewDocumentServiceImpl(documentRepo repodocuments.DocumentRepository) *DocumentServiceImpl {
	return &DocumentServiceImpl{documentRepo}
}

// Get returns the document with ID.
func (s *DocumentServiceImpl) Get(id string) (models.Document, error) {
	return s.documentRepo.GetById(id)
}

// GetAll return all documents sorted by ID
func (s *DocumentServiceImpl) GetAll() ([]models.Document, error) {
	return s.documentRepo.GetAll()
}

// CreateOrUpdate creates or update given document
func (s *DocumentServiceImpl) CreateOrUpdate(documentToCreate models.Document) (bool, error) {
	return s.documentRepo.CreateOrUpdate(documentToCreate)
}

// Delete delete document id
func (s *DocumentServiceImpl) Delete(idToDelete string) (bool, error) {
	return s.documentRepo.Delete(idToDelete)
}

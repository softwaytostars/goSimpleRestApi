package services

import (
	"goapi/models"
	"log"
)

type DocumentService interface {
	Get(id string) (models.Document, error)
	GetAll() ([]models.Document, error)
	CreateOrUpdate(document models.Document) (bool, error)
	Delete(id string) (bool, error)
}

type DocumentServiceImpl struct {
	documentsById map[string]models.Document
}

func NewDocumentServiceImpl() *DocumentServiceImpl {
	return &DocumentServiceImpl{documentsById: make(map[string]models.Document)}
}

// Get returns the document with ID.
func (s DocumentServiceImpl) Get(id string) (models.Document, error) {
	document, found := s.documentsById[id]
	if found {
		return document, nil
	}
	return models.Document{}, nil
}

// GetAll return all documents
func (s DocumentServiceImpl) GetAll() ([]models.Document, error) {
	values := make([]models.Document, 0, len(s.documentsById))
	for _, doc := range s.documentsById {
		values = append(values, doc)
	}
	return values, nil
}

// CreateOrUpdate creates or update given document
func (s DocumentServiceImpl) CreateOrUpdate(documentToCreate models.Document) (bool, error) {
	_, found := s.documentsById[documentToCreate.ID]
	if found {
		log.Print("document " + documentToCreate.ID + " already exists")
	}
	s.documentsById[documentToCreate.ID] = documentToCreate
	return found, nil
}

// Delete delete document id
func (s DocumentServiceImpl) Delete(idToDelete string) (bool, error) {
	_, found := s.documentsById[idToDelete]
	if !found {
		log.Print("document " + idToDelete + " doesn't exists")
	}
	delete(s.documentsById, idToDelete)
	return found, nil
}

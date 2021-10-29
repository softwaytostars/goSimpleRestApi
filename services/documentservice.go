package services

import (
	"goapi/models"
	"log"
	"sort"
	"sync"
)

type DocumentService interface {
	Get(id string) (models.Document, error)
	GetAll() ([]models.Document, error)
	CreateOrUpdate(document models.Document) (bool, error)
	Delete(id string) (bool, error)
}

// DocumentServiceImpl Default implementation for DocumentService
type DocumentServiceImpl struct {
	documentsById sync.Map
}

func NewDocumentServiceImpl() *DocumentServiceImpl {
	return &DocumentServiceImpl{}
}

// Get returns the document with ID.
func (s *DocumentServiceImpl) Get(id string) (models.Document, error) {
	document, found := s.documentsById.Load(id)
	if found {
		return document.(models.Document), nil
	}
	return models.Document{}, nil
}

// GetAll return all documents sorted by ID
func (s *DocumentServiceImpl) GetAll() ([]models.Document, error) {
	ids := make([]string, 0)
	//retrieve ids and sort them
	s.documentsById.Range(func(id, value interface{}) bool {
		ids = append(ids, id.(string))
		return true
	})
	sort.Strings(ids)

	//fill values for sorted ids
	values := make([]models.Document, 0, len(ids))
	for _, id := range ids {
		doc, _ := s.documentsById.Load(id)
		values = append(values, doc.(models.Document))
	}

	return values, nil
}

// CreateOrUpdate creates or update given document
func (s *DocumentServiceImpl) CreateOrUpdate(documentToCreate models.Document) (bool, error) {
	_, found := s.documentsById.Load(documentToCreate.ID)
	if found {
		log.Print("document " + documentToCreate.ID + " already exists")
	}
	s.documentsById.Store(documentToCreate.ID, documentToCreate)
	return found, nil
}

// Delete delete document id
func (s *DocumentServiceImpl) Delete(idToDelete string) (bool, error) {
	_, found := s.documentsById.Load(idToDelete)
	if !found {
		log.Print("document " + idToDelete + " doesn't exists")
	}
	s.documentsById.Delete(idToDelete)
	return found, nil
}

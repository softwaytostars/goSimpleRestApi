package repositories

import (
	log "github.com/sirupsen/logrus"
	"goapi/models"
	"sort"
	"sync"
)

type InMemoryDocumentRepo struct {
	DocumentsById sync.Map
}

func (r *InMemoryDocumentRepo) GetById(id string) (models.Document, error) {
	document, found := r.DocumentsById.Load(id)
	if found {
		return document.(models.Document), nil
	}
	return models.Document{}, nil
}

func (r *InMemoryDocumentRepo) GetAll() ([]models.Document, error) {
	ids := make([]string, 0)
	//retrieve ids and sort them
	r.DocumentsById.Range(func(id, value interface{}) bool {
		ids = append(ids, id.(string))
		return true
	})
	sort.Strings(ids)

	//fill values for sorted ids
	values := make([]models.Document, 0, len(ids))
	for _, id := range ids {
		doc, _ := r.DocumentsById.Load(id)
		values = append(values, doc.(models.Document))
	}

	return values, nil
}

func (r *InMemoryDocumentRepo) CreateOrUpdate(documentToCreate models.Document) (bool, error) {
	_, found := r.DocumentsById.Load(documentToCreate.ID)
	if found {
		log.Info("document " + documentToCreate.ID + " already exists")
	}
	r.DocumentsById.Store(documentToCreate.ID, documentToCreate)
	return found, nil
}

func (r *InMemoryDocumentRepo) Delete(idToDelete string) (bool, error) {
	_, found := r.DocumentsById.Load(idToDelete)
	if !found {
		log.Info("document " + idToDelete + " doesn't exists")
	}
	r.DocumentsById.Delete(idToDelete)
	return found, nil
}

package repodocuments

import "goapi/models"

type DocumentRepository interface {
	GetById(id string) (models.Document, error)
	GetAll() ([]models.Document, error)
	CreateOrUpdate(document models.Document) (bool, error)
	Delete(id string) (bool, error)
}

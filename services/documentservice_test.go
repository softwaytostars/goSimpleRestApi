package services

import (
	"github.com/stretchr/testify/assert"
	"goapi/models"
	"goapi/repositories"
	"testing"
)

func TestDocumentServiceImpl_CreateOrUpdateWithCreation(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	updated, err := documentServiceImpl.CreateOrUpdate(models.Document{ID: "toto", Description: "descToto", Name: "nameToto"})
	assert.Nil(t, err)
	length := 0
	repo.DocumentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 1, length)
	assert.False(t, updated)
}

func TestDocumentServiceImpl_CreateOrUpdateWithUpdate(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	repo.DocumentsById.Store("toto", doc)

	docUpdate := models.Document{ID: doc.ID, Description: "descUpdateToto", Name: "nameUpdateToto"}
	updated, err := documentServiceImpl.CreateOrUpdate(docUpdate)
	assert.Nil(t, err)
	length := 0
	repo.DocumentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 1, length)
	assert.True(t, updated)
	docFound, _ := repo.DocumentsById.Load(doc.ID)
	assert.Equal(t, docFound.(models.Document), docUpdate)
}

func TestDocumentServiceImpl_DeleteExisting(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	repo.DocumentsById.Store("toto", doc)

	found, err := documentServiceImpl.Delete(doc.ID)
	assert.Nil(t, err)
	assert.True(t, found)
	length := 0
	repo.DocumentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

func TestDocumentServiceImpl_DeleteNonExisting(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	found, err := documentServiceImpl.Delete("toto")
	assert.Nil(t, err)
	assert.False(t, found)
	length := 0
	repo.DocumentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

func TestDocumentServiceImpl_GetFound(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	repo.DocumentsById.Store("toto", doc)
	res, err := documentServiceImpl.Get("toto")
	assert.Nil(t, err)
	assert.Equal(t, doc, res)
}

func TestDocumentServiceImpl_GetNotFound(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	res, err := documentServiceImpl.Get("toto")
	assert.Nil(t, err)
	assert.Equal(t, models.Document{}, res)
}

func TestDocumentServiceImpl_GetAll(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	docToto := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	repo.DocumentsById.Store(docToto.ID, docToto)

	docTata := models.Document{ID: "tata", Description: "descTata", Name: "nameTata"}
	repo.DocumentsById.Store(docTata.ID, docTata)

	res, err := documentServiceImpl.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, docTata, res[0]) //values should be sorted by ID
	assert.Equal(t, docToto, res[1])
}

func TestNewDocumentServiceImpl(t *testing.T) {
	repo := repositories.InMemoryDocumentRepo{}
	documentServiceImpl := NewDocumentServiceImpl(&repo)

	assert.NotNil(t, documentServiceImpl)
	assert.NotNil(t, &repo.DocumentsById)
	length := 0
	repo.DocumentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

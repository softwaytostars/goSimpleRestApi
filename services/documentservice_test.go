package services

import (
	"github.com/stretchr/testify/assert"
	"goapi/models"
	"testing"
)

func TestDocumentServiceImpl_CreateOrUpdateWithCreation(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()

	updated, err := documentServiceImpl.CreateOrUpdate(models.Document{ID: "toto", Description: "descToto", Name: "nameToto"})
	assert.Nil(t, err)
	length := 0
	documentServiceImpl.documentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 1, length)
	assert.False(t, updated)
}

func TestDocumentServiceImpl_CreateOrUpdateWithUpdate(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById.Store("toto", doc)

	docUpdate := models.Document{ID: doc.ID, Description: "descUpdateToto", Name: "nameUpdateToto"}
	updated, err := documentServiceImpl.CreateOrUpdate(docUpdate)
	assert.Nil(t, err)
	length := 0
	documentServiceImpl.documentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 1, length)
	assert.True(t, updated)
	docFound, _ := documentServiceImpl.documentsById.Load(doc.ID)
	assert.Equal(t, docFound.(models.Document), docUpdate)
}

func TestDocumentServiceImpl_DeleteExisting(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById.Store("toto", doc)

	found, err := documentServiceImpl.Delete(doc.ID)
	assert.Nil(t, err)
	assert.True(t, found)
	length := 0
	documentServiceImpl.documentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

func TestDocumentServiceImpl_DeleteNonExisting(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	found, err := documentServiceImpl.Delete("toto")
	assert.Nil(t, err)
	assert.False(t, found)
	length := 0
	documentServiceImpl.documentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

func TestDocumentServiceImpl_GetFound(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById.Store("toto", doc)
	res, err := documentServiceImpl.Get("toto")
	assert.Nil(t, err)
	assert.Equal(t, doc, res)
}

func TestDocumentServiceImpl_GetNotFound(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	res, err := documentServiceImpl.Get("toto")
	assert.Nil(t, err)
	assert.Equal(t, models.Document{}, res)
}

func TestDocumentServiceImpl_GetAll(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()

	docToto := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById.Store(docToto.ID, docToto)

	docTata := models.Document{ID: "tata", Description: "descTata", Name: "nameTata"}
	documentServiceImpl.documentsById.Store(docTata.ID, docTata)

	res, err := documentServiceImpl.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, docTata, res[0]) //values should be sorted by ID
	assert.Equal(t, docToto, res[1])
}

func TestNewDocumentServiceImpl(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	assert.NotNil(t, documentServiceImpl)
	assert.NotNil(t, &documentServiceImpl.documentsById)
	length := 0
	documentServiceImpl.documentsById.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)
}

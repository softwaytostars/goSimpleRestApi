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
	assert.Equal(t, 1, len(documentServiceImpl.documentsById))
	assert.False(t, updated)
}

func TestDocumentServiceImpl_CreateOrUpdateWithUpdate(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById["toto"] = doc

	docUpdate := models.Document{ID: doc.ID, Description: "descUpdateToto", Name: "nameUpdateToto"}
	updated, err := documentServiceImpl.CreateOrUpdate(docUpdate)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(documentServiceImpl.documentsById))
	assert.True(t, updated)
	assert.Equal(t, documentServiceImpl.documentsById[doc.ID], docUpdate)
}

func TestDocumentServiceImpl_DeleteExisting(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById["toto"] = doc

	found, err := documentServiceImpl.Delete(doc.ID)
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, 0, len(documentServiceImpl.documentsById))
}

func TestDocumentServiceImpl_DeleteNonExisting(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	found, err := documentServiceImpl.Delete("toto")
	assert.Nil(t, err)
	assert.False(t, found)
	assert.Equal(t, 0, len(documentServiceImpl.documentsById))
}

func TestDocumentServiceImpl_GetFound(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById["toto"] = doc
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
	doc := models.Document{ID: "toto", Description: "descToto", Name: "nameToto"}
	documentServiceImpl.documentsById["toto"] = doc
	res, err := documentServiceImpl.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, doc, res[0])
}

func TestNewDocumentServiceImpl(t *testing.T) {
	documentServiceImpl := NewDocumentServiceImpl()
	assert.NotNil(t, documentServiceImpl)
	assert.NotNil(t, documentServiceImpl.documentsById)
	assert.Equal(t, 0, len(documentServiceImpl.documentsById))
}

package resources

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goapi/models"
	"goapi/services"
	"net/http"
)

type ResourceDocument struct {
	documentService services.DocumentService
}

func (resource ResourceDocument) validationID(id string) error {
	if len(id) == 0 {
		return errors.New("id must be defined")
	}
	return nil
}

//endpoint to retrieve all documents
func (resource ResourceDocument) getAllDocuments(c *gin.Context) {
	docs, err := resource.documentService.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot get documents [err=%s]", err)})
		return
	}
	c.IndentedJSON(http.StatusOK, docs)
}

//endpoint to retrieve a given document from the path param id
func (resource ResourceDocument) getDocument(c *gin.Context) {
	id := c.Param("id")
	doc, err := resource.documentService.Get(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot get document id %s [err=%s]", id, err)})
		return
	}
	if (models.Document{}) == doc {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("document id %s not found", id)})
		return
	}
	c.IndentedJSON(http.StatusOK, doc)
}

// endpoint to create or update a document
func (resource ResourceDocument) createOrUpdateDocument(c *gin.Context) {
	id := c.Param("id")

	err := resource.validationID(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Validation failed [err=%s]", err)})
		return
	}

	var docToCreateOrUpdate models.Document
	if err := c.BindJSON(&docToCreateOrUpdate); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Cannot deserialize document [err=%s]", err)})
		return
	}
	docToCreateOrUpdate.ID = id

	docUpdated, err := resource.documentService.CreateOrUpdate(docToCreateOrUpdate)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot create or update document [err=%s]", err)})
		return
	}

	if docUpdated {
		c.IndentedJSON(http.StatusOK, docToCreateOrUpdate)
	} else {
		c.IndentedJSON(http.StatusCreated, docToCreateOrUpdate)
	}
}

// endpoint to Delete a given document id
func (resource ResourceDocument) deleteDocument(c *gin.Context) {
	idToDelete := c.Param("id")

	err := resource.validationID(idToDelete)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Validation failed [err=%s]", err)})
		return
	}

	found, err := resource.documentService.Delete(idToDelete)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot delete documents [err=%s]", err)})
		return
	}

	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("document id %s not found", idToDelete)})
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

// RegisterHandlers register all handlers for a router
func RegisterHandlers(r *gin.Engine, documentService services.DocumentService) {
	resource := ResourceDocument{documentService}

	r.GET("/documents", resource.getAllDocuments)
	r.GET("/documents/:id", resource.getDocument)
	r.PUT("/documents/:id", resource.createOrUpdateDocument)
	r.DELETE("/documents/:id", resource.deleteDocument)
}

package resources

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag/example/celler/httputil"
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

// Endpoint to retrieve all documents
// @Summary Retrieve all documents
// @Description Retrieve all documents
// @Produce  json
// @Success 200 {array} models.Document
// @Failure 500 {object} httputil.HTTPError
// @Router /documents [get]
func (resource ResourceDocument) GetAllDocuments(c *gin.Context) {
	docs, err := resource.documentService.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot get documents [err=%s]", err)})
		return
	}
	c.IndentedJSON(http.StatusOK, docs)
}

// Endpoint to retrieve a given document from the path param id
// @Summary Retrieve a given document
// @Description Retrieve  a given document from the path param id
// @Produce  json
// @Param id path int true "Document ID"
// @Success 200 {object} models.Document
// @Failure 500 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /documents/{id} [get]
func (resource ResourceDocument) GetDocument(c *gin.Context) {
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

// Endpoint to create or update a document
// @Summary Create or update a document
// @Description Create or update a document
// @Accept  json
// @Produce  json
// @Param id path int true "Document ID"
// @Param data body models.Document true "The document struct"
// @Success 200 {object} models.Document "update"
// @Success 201 {object} models.Document "creation"
// @Failure 500 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /documents/{id} [put]
func (resource ResourceDocument) CreateOrUpdateDocument(c *gin.Context) {
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

// Endpoint to Delete a given document id
// @Summary  Delete a given document id
// @Description  Delete a given document id
// @Param id path int true "Document ID"
// @Success 200 "OK"
// @Failure 500 {object} httputil.HTTPError
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /documents/{id} [delete]
func (resource ResourceDocument) DeleteDocument(c *gin.Context) {
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

	r.GET("/documents", resource.GetAllDocuments)
	r.GET("/documents/:id", resource.GetDocument)
	r.PUT("/documents/:id", resource.CreateOrUpdateDocument)
	r.DELETE("/documents/:id", resource.DeleteDocument)
}

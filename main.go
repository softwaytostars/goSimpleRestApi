package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type document struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var documents = make([]document, 0)

func getAllDocuments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, documents)
}

//retrieve a given document from the query param id
func getDocument(c *gin.Context) {
	id := c.Param("id")
	for _, document := range documents {
		if document.ID == id {
			c.IndentedJSON(http.StatusOK, document)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("document id %s not found", id)})
}

//create a list of documents
func createDocuments(c *gin.Context) {
	var newDocuments []document
	if err := c.BindJSON(&newDocuments); err != nil {
		return
	}
	for _, newDocument := range newDocuments {
		documents = append(documents, newDocument)
	}
	c.IndentedJSON(http.StatusCreated, newDocuments)
}

func main() {
	router := gin.Default()

	router.GET("/documents", getAllDocuments)
	router.GET("/documents/:id", getDocument)
	router.PUT("/documents", createDocuments)

	err := router.Run("localhost:8040")
	if err != nil {
		return
	}
}

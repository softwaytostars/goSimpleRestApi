package resources

import (
	"bytes"
	"encoding/json"
	_ "encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"goapi/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DocumentServiceMock struct {
	mock.Mock
}

func (s *DocumentServiceMock) Get(id string) (models.Document, error) {
	args := s.Called(id)
	return args.Get(0).(models.Document), args.Error(1)
}

func (s *DocumentServiceMock) GetAll() ([]models.Document, error) {
	args := s.Called()
	return args.Get(0).([]models.Document), args.Error(1)
}

func (s *DocumentServiceMock) CreateOrUpdate(documentToCreate models.Document) (bool, error) {
	args := s.Called(documentToCreate)
	return args.Get(0).(bool), args.Error(1)
}

func (s *DocumentServiceMock) Delete(idToDelete string) (bool, error) {
	args := s.Called(idToDelete)
	return args.Get(0).(bool), args.Error(1)
}

func configureRouter(service *DocumentServiceMock) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	RegisterHandlers(router, service)
	return router
}

func createTestServer(service *DocumentServiceMock) *httptest.Server {
	router := configureRouter(service)
	// Start a local HTTP server
	testServer := httptest.NewServer(router)
	return testServer
}

func executeRequest(t *testing.T, client *http.Client, req *http.Request, documentServiceMock *DocumentServiceMock, expectedBody string, expectedCodeStatus int) {

	//execute the request
	resp, err := client.Do(req)
	assert.Nil(t, err, fmt.Sprintf("Couldn't send request: %v", err))

	// assert that the expectations were met
	documentServiceMock.AssertExpectations(t)

	//check the status
	assert.Equal(t, expectedCodeStatus, resp.StatusCode)

	//read the body response
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err, fmt.Sprintf("Couldn't read body response: %v", err))

	//check the body value
	require.JSONEq(t, expectedBody, string(body[:]))
}

func TestResourceDocument_getAllDocuments(t *testing.T) {

	expected := []models.Document{
		{
			ID:          "toto",
			Name:        "nameOfToto",
			Description: "descOfToto",
		},
		{
			ID:          "titi",
			Name:        "nameOfTiti",
			Description: "descOfTiti",
		},
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("GetAll").Return(expected, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()

	//create request
	req, err := http.NewRequest("GET", testServer.URL+"/documents", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	expectedBody, err := json.Marshal(expected)
	executeRequest(t, testServer.Client(), req, documentServiceMock, string(expectedBody[:]), http.StatusOK)
}

func TestResourceDocument_getAllDocumentsErrorService(t *testing.T) {

	expected := []models.Document{
		{
			ID:          "toto",
			Name:        "nameOfToto",
			Description: "descOfToto",
		},
		{
			ID:          "titi",
			Name:        "nameOfTiti",
			Description: "descOfTiti",
		},
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("GetAll").Return(expected, errors.New("error_service_getAll"))

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()

	//create request
	req, err := http.NewRequest("GET", testServer.URL+"/documents", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot get documents [err=error_service_getAll]"}`
	executeRequest(t, testServer.Client(), req, documentServiceMock, expectedError, http.StatusInternalServerError)
}

func TestResourceDocument_getDocumentOK(t *testing.T) {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Get", "toto").Return(expected, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request should be OK
	req, err := http.NewRequest("GET", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result is OK
	expectedBody, err := json.Marshal(expected)
	executeRequest(t, client, req, documentServiceMock, string(expectedBody[:]), http.StatusOK)
}

func TestResourceDocument_getDocumentErrorService(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Get", "toto").Return(models.Document{}, errors.New("error_service_get"))

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create the request
	req, err := http.NewRequest("GET", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	var expectedError = `{"message":"Cannot get document id toto [err=error_service_get]"}`
	//check result is error
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusInternalServerError)
}

func TestResourceDocument_getDocumentNotFound(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Get", "toto").Return(models.Document{}, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create the request
	req, err := http.NewRequest("GET", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	var expectedError = `{"message":"document id toto not found"}`
	//check result is error
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusNotFound)
}

func TestResourceDocument_createOrUpdateDocumentCreation(t *testing.T) {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("CreateOrUpdate", expected).Return(false, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	payload, err := json.Marshal(expected)
	req, err := http.NewRequest("PUT", testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	executeRequest(t, client, req, documentServiceMock, string(payload), http.StatusCreated)
}

func TestResourceDocument_createOrUpdateDocumentUpdate(t *testing.T) {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("CreateOrUpdate", expected).Return(true, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	payload, err := json.Marshal(models.Document{Name: expected.Name, Description: expected.Description})
	req, err := http.NewRequest("PUT", testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	expectedBody, err := json.Marshal(expected)
	executeRequest(t, client, req, documentServiceMock, string(expectedBody), http.StatusOK)
}

func TestResourceDocument_createOrUpdateDocumentErrorService(t *testing.T) {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("CreateOrUpdate", expected).Return(false, errors.New("error_service_create"))

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	payload, err := json.Marshal(expected)
	req, err := http.NewRequest("PUT", testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot create or update document [err=error_service_create]"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusInternalServerError)
}

func TestResourceDocument_createOrUpdateDocumentNoPayload(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("PUT", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot deserialize document [err=EOF]"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusBadRequest)
}

func TestResourceDocument_createOrUpdateDocumentWrongID(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("PUT", testServer.URL+"/documents/", bytes.NewBuffer([]byte("{\"name\":\"nameToto\"}")))
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Validation failed [err=id must be defined]"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusBadRequest)
}

func TestResourceDocument_deleteDocument(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Delete", "toto").Return(true, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("DELETE", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	executeRequest(t, client, req, documentServiceMock, "null", http.StatusOK)
}

func TestResourceDocument_deleteDocumentErrorService(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Delete", "toto").Return(false, errors.New("error_service_delete"))

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("DELETE", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot delete documents [err=error_service_delete]"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusInternalServerError)
}

func TestResourceDocument_deleteDocumentsWrongID(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("DELETE", testServer.URL+"/documents/", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Validation failed [err=id must be defined]"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusBadRequest)
}

func TestResourceDocument_deleteDocumentNotFound(t *testing.T) {

	//create mock service
	documentServiceMock := new(DocumentServiceMock)
	documentServiceMock.On("Delete", "toto").Return(false, nil)

	//create http server
	testServer := createTestServer(documentServiceMock)
	// Close the server when test finishes
	defer testServer.Close()
	client := testServer.Client()

	//create request
	req, err := http.NewRequest("DELETE", testServer.URL+"/documents/toto", nil)
	assert.Nil(t, err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"document id toto not found"}`
	executeRequest(t, client, req, documentServiceMock, expectedError, http.StatusNotFound)
}

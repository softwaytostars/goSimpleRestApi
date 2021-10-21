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
	"github.com/stretchr/testify/suite"
	"goapi/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Mock  implementation for documentService
*/
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

/*
	Test suite definition
*/

type DocumentResourceTestSuite struct {
	suite.Suite
	documentServiceMock *DocumentServiceMock
	testServer          *httptest.Server
}

//beforeAll method
func (suite *DocumentResourceTestSuite) SetupSuite() {
	//create the mock service
	suite.documentServiceMock = new(DocumentServiceMock)
	//create http server
	suite.testServer = createTestServer(suite.documentServiceMock)
}

//the suite runner
func TestDocumentResourceTestSuite(t *testing.T) {
	suite.Run(t, new(DocumentResourceTestSuite))
}

// The SetupTest method will be run before every test in the suite.
func (suite *DocumentResourceTestSuite) SetupTest() {
	//not yet available suite.documentServiceMock.Off
	//for removing handlers on the mock
	suite.documentServiceMock.ExpectedCalls = nil
}

// The TearDownTest method will be run after every test in the suite.
func (suite *DocumentResourceTestSuite) TearDownTest() {
}

func (suite *DocumentResourceTestSuite) TearDownSuite() {
	//close the http server test
	suite.testServer.Close()
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_getAllDocuments() {

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

	//add handler mock service
	suite.documentServiceMock.On("GetAll").Return(expected, nil)

	//create request
	req, err := http.NewRequest("GET", suite.testServer.URL+"/documents", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	expectedBody, err := json.Marshal(expected)
	executeRequest(suite, req, string(expectedBody[:]), http.StatusOK)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_getAllDocumentsErrorService() {

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

	//add handler mock service
	suite.documentServiceMock.On("GetAll").Return(expected, errors.New("error_service_getAll"))

	//create request
	req, err := http.NewRequest("GET", suite.testServer.URL+"/documents", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot get documents [err=error_service_getAll]"}`
	executeRequest(suite, req, expectedError, http.StatusInternalServerError)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_getDocumentOK() {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//add handler mock service
	suite.documentServiceMock.On("Get", "toto").Return(expected, nil)

	//create request should be OK
	req, err := http.NewRequest("GET", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result is OK
	expectedBody, err := json.Marshal(expected)
	executeRequest(suite, req, string(expectedBody[:]), http.StatusOK)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_getDocumentErrorService() {

	//add handler mock service
	suite.documentServiceMock.On("Get", "toto").Return(models.Document{}, errors.New("error_service_get"))

	//create the request
	req, err := http.NewRequest("GET", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	var expectedError = `{"message":"Cannot get document id toto [err=error_service_get]"}`
	//check result is error
	executeRequest(suite, req, expectedError, http.StatusInternalServerError)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_getDocumentNotFound() {

	//add handler mock service
	suite.documentServiceMock.On("Get", "toto").Return(models.Document{}, nil)

	//create the request
	req, err := http.NewRequest("GET", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	var expectedError = `{"message":"document id toto not found"}`
	//check result is error
	executeRequest(suite, req, expectedError, http.StatusNotFound)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_createOrUpdateDocumentCreation() {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//add handler mock service
	suite.documentServiceMock.On("CreateOrUpdate", expected).Return(false, nil)

	//create request
	payload, err := json.Marshal(expected)
	req, err := http.NewRequest("PUT", suite.testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	executeRequest(suite, req, string(payload), http.StatusCreated)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_createOrUpdateDocumentUpdate() {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//add handler mock service
	suite.documentServiceMock.On("CreateOrUpdate", expected).Return(true, nil)

	//create request
	payload, err := json.Marshal(models.Document{Name: expected.Name, Description: expected.Description})
	req, err := http.NewRequest("PUT", suite.testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	expectedBody, err := json.Marshal(expected)
	executeRequest(suite, req, string(expectedBody), http.StatusOK)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_createOrUpdateDocumentErrorService() {

	expected := models.Document{
		ID:          "toto",
		Name:        "nameOfToto",
		Description: "descOfToto",
	}

	//add handler mock service
	suite.documentServiceMock.On("CreateOrUpdate", expected).Return(false, errors.New("error_service_create"))

	//create request
	payload, err := json.Marshal(expected)
	req, err := http.NewRequest("PUT", suite.testServer.URL+"/documents/"+expected.ID, bytes.NewBuffer(payload))
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot create or update document [err=error_service_create]"}`
	executeRequest(suite, req, expectedError, http.StatusInternalServerError)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_createOrUpdateDocumentNoPayload() {

	//create request
	req, err := http.NewRequest("PUT", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot deserialize document [err=EOF]"}`
	executeRequest(suite, req, expectedError, http.StatusBadRequest)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_createOrUpdateDocumentWrongID() {

	//create request
	req, err := http.NewRequest("PUT", suite.testServer.URL+"/documents/", bytes.NewBuffer([]byte("{\"name\":\"nameToto\"}")))
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Validation failed [err=id must be defined]"}`
	executeRequest(suite, req, expectedError, http.StatusBadRequest)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_deleteDocument() {

	//add handler mock service
	suite.documentServiceMock.On("Delete", "toto").Return(true, nil)

	//create request
	req, err := http.NewRequest("DELETE", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	executeRequest(suite, req, "null", http.StatusOK)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_deleteDocumentErrorService() {

	//add handler mock service
	suite.documentServiceMock.On("Delete", "toto").Return(false, errors.New("error_service_delete"))

	//create request
	req, err := http.NewRequest("DELETE", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Cannot delete documents [err=error_service_delete]"}`
	executeRequest(suite, req, expectedError, http.StatusInternalServerError)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_deleteDocumentsWrongID() {

	//create request
	req, err := http.NewRequest("DELETE", suite.testServer.URL+"/documents/", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"Validation failed [err=id must be defined]"}`
	executeRequest(suite, req, expectedError, http.StatusBadRequest)
}

func (suite *DocumentResourceTestSuite) TestResourceDocument_deleteDocumentNotFound() {

	//add handler mock service
	suite.documentServiceMock.On("Delete", "toto").Return(false, nil)

	//create request
	req, err := http.NewRequest("DELETE", suite.testServer.URL+"/documents/toto", nil)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't create request: %v", err))

	//check result
	var expectedError = `{"message":"document id toto not found"}`
	executeRequest(suite, req, expectedError, http.StatusNotFound)
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

func executeRequest(suite *DocumentResourceTestSuite, req *http.Request, expectedBody string, expectedCodeStatus int) {

	//execute the request
	resp, err := suite.testServer.Client().Do(req)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't send request: %v", err))

	// assert that the expectations were met
	suite.documentServiceMock.AssertExpectations(suite.T())

	//check the status
	assert.Equal(suite.T(), expectedCodeStatus, resp.StatusCode)

	//read the body response
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(suite.T(), err, fmt.Sprintf("Couldn't read body response: %v", err))

	//check the body value
	require.JSONEq(suite.T(), expectedBody, string(body[:]))
}

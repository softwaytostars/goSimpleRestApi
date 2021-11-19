package goconvey

import (
	"bytes"
	"encoding/json"
	"goapi/models"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func getServerHost() string {
	host := os.Getenv("SERVER_HOST")
	if len(host) <= 0 {
		host = "localhost"
	}
	return host
}

func doCall(method string, addr string, body io.Reader) *http.Response {

	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		logrus.Error(err.Error())
	}
	So(err, ShouldBeNil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)

	So(err, ShouldBeNil)

	return resp
}

func getBodyResponseFrom(resp *http.Response) []byte {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err.Error())
	}

	So(err, ShouldBeNil)

	return bodyBytes
}

func TestApiDocument(t *testing.T) {

	Convey("Given api server uses database and database is available", t, func() {
		Convey("Then an empty document list should be available", func() {

			resp := doCall("GET", "http://"+getServerHost()+":8040/documents", nil)
			defer resp.Body.Close()

			bodyBytes := getBodyResponseFrom(resp)

			var documents []models.Document
			json.Unmarshal(bodyBytes, &documents)
			So(documents, ShouldBeEmpty)

		})

		Convey("When I add a document with id toto and name nametoto", func() {

			postBody, _ := json.Marshal(models.Document{Name: "nametoto"})

			resp := doCall("PUT", "http://"+getServerHost()+":8040/documents/toto", bytes.NewBuffer(postBody))
			So(resp.StatusCode, ShouldEqual, 201)
			resp.Body.Close()

			Convey("Then I can retrieve it with name nametoto", func() {

				resp := doCall("GET", "http://"+getServerHost()+":8040/documents/toto", nil)
				defer resp.Body.Close()

				bodyBytes := getBodyResponseFrom(resp)
				var document models.Document
				json.Unmarshal(bodyBytes, &document)
				So(document, ShouldResemble, models.Document{ID: "toto", Name: "nametoto"})

			})

		})

		Convey("When I add a document with same id toto and name nametoto2", func() {

			postBody, _ := json.Marshal(models.Document{Name: "nametoto2"})

			resp := doCall("PUT", "http://"+getServerHost()+":8040/documents/toto", bytes.NewBuffer(postBody))
			So(resp.StatusCode, ShouldEqual, 200)
			resp.Body.Close()

			Convey("Then I can retrieve it with name nametoto2", func() {
				resp := doCall("GET", "http://"+getServerHost()+":8040/documents/toto", nil)
				defer resp.Body.Close()

				bodyBytes := getBodyResponseFrom(resp)
				var document models.Document
				json.Unmarshal(bodyBytes, &document)
				So(document, ShouldResemble, models.Document{ID: "toto", Name: "nametoto2"})
			})

		})

		Convey("When I look fors documents, I find only toto with name toto2", func() {

			resp := doCall("GET", "http://"+getServerHost()+":8040/documents", nil)
			defer resp.Body.Close()

			bodyBytes := getBodyResponseFrom(resp)

			var documents []models.Document
			json.Unmarshal(bodyBytes, &documents)
			So(len(documents), ShouldEqual, 1)
			So(documents[0], ShouldResemble, models.Document{ID: "toto", Name: "nametoto2"})

		})

		Convey("When I delete document id toto", func() {

			resp := doCall("DELETE", "http://"+getServerHost()+":8040/documents/toto", nil)
			defer resp.Body.Close()

			Convey("Then the list of documents is empty", func() {
				resp := doCall("GET", "http://"+getServerHost()+":8040/documents", nil)
				defer resp.Body.Close()

				bodyBytes := getBodyResponseFrom(resp)

				var documents []models.Document
				json.Unmarshal(bodyBytes, &documents)
				So(documents, ShouldBeEmpty)
			})

		})

	})

}

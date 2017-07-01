package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jdelobel/go-api/cmd/apid/handlers"
	"github.com/jdelobel/go-api/internal/image"
	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/web"
	"github.com/jdelobel/go-api/logger"

	"strings"

	"github.com/jdelobel/go-api/config"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Succeed is the Unicode codepoint for a check mark.
	Succeed = "\u2713"

	// Failed is the Unicode codepoint for an X mark.
	Failed = "\u2717"
)

// The web application state for tests
var a *web.App

// init is called before main. We are using init to customize logging output.
func init() {
	err := logger.Init(logger.Conf{Level: "EMERGENCY", App: "go-api-testing"})
	if err != nil {
		log.Fatalf("main: Failed to init config: %v", err)
	}
}

// TestMain gives you a chance to setup / tear down before tests in this package.
func TestMain(m *testing.M) {
	os.Exit(runTest(m))
}

// runTest initializes the environment for the tests and allows for
// the proper return code if the test fails or succeeds.
func runTest(m *testing.M) int {

	// Check the environment for a configured port value.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres://go-api-postgres:go-api-postgres@localhost:5432/go-api-postgres?sslmode=disable"
	}

	// Register the Master Session for the database.
	masterDB, err := db.NewPSQL("postgres", dbHost)
	c := config.Config{}
	if err != nil {
		return 1
	}
	a = handlers.API(masterDB, logger.Log, c, nil).(*web.App)

	return m.Run()
}

// TestImages is the entry point for the images
func TestImages(t *testing.T) {
	t.Run("getImages200Empty", getImages200Empty)
	t.Run("postImage400", postImage400)
	t.Run("getImage404", getImage404)
	t.Run("getImage400", getImage400)
	t.Run("putImage404", putImage404)
	t.Run("crudImages", crudImage)
}

// getImages200Empty validates an empty images list can be retrieved with the endpoint.
func getImages200Empty(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/images", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to fetch an empty list of images with the images endpoint.")
	{
		t.Log("\tTest 0:\tWhen fetching an empty image list.")
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the response.", Succeed)

			recv := w.Body.String()
			resp := `[]`
			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// postImage400 validates an image can't be created with the endpoint
// unless a valid image document is submitted.
func postImage400(t *testing.T) {
	m := image.CreateImage{
		ID:    "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title: "Image Elijah Baley",
		URL:   "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:  "/images/1280/720/test-2260-b1396d-1@1x",
	}

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("POST", "/v1/images", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate a new image can't be created without publisher")
	{
		t.Log("\tTest 0:\tWhen using an incomplete image value.")
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tShould receive a status code of 400 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 400 for the response.", Succeed)

			recv := w.Body.String()
			resps := []string{
				`{
  "error": "field validation failure",
  "fields": [
    {
      "field_name": "Publisher",
      "error": "required"
    }
  ]
}`,
			}

			var found bool
			for _, resp := range resps {
				if resp == recv {
					found = true
					break
				}
			}

			if !found {
				t.Log("Got :", recv)
				t.Log("Want:", resps[0])
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// getImage400 validates an image request for a malformed imageid.
func getImage400(t *testing.T) {
	imageID := "12345"

	r := httptest.NewRequest("GET", "/v1/images/"+imageID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an image with a malformed imageid.")
	{
		t.Logf("\tTest 0:\tWhen using the new image %s.", imageID)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tShould receive a status code of 400 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 400 for the response.", Succeed)

			recv := w.Body.String()
			resp := `{
  "error": "ID is not in it's proper form"
}`
			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// getImage404 validates an image request for an image that does not exist with the endpoint.
func getImage404(t *testing.T) {
	imageID := bson.NewObjectId().Hex()

	r := httptest.NewRequest("GET", "/v1/images/"+imageID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an image with an unknown id.")
	{
		t.Logf("\tTest 0:\tWhen using the new image %s.", imageID)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tShould receive a status code of 404 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 404 for the response.", Succeed)

			recv := w.Body.String()
			resp := "Entity not found"
			if !strings.Contains(recv, resp) {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// putImage404 validates updating an image that does not exist.
func putImage404(t *testing.T) {
	m := image.CreateImage{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Image Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	imageID := bson.NewObjectId().Hex()

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("PUT", "/v1/images/"+imageID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate updating an image that does not exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new image %s.", imageID)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tShould receive a status code of 404 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 404 for the response.", Succeed)

			recv := w.Body.String()
			resp := "Entity not found"
			if !strings.Contains(recv, resp) {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// crudImage performs a complete test of CRUD against the api.
func crudImage(t *testing.T) {
	nm := postImage201(t)
	defer deleteImage204(t, nm.ID)

	getImage200(t, nm.ID)
	putImage204(t, nm)
}

// postImage201 validates an image can be created with the endpoint.
func postImage201(t *testing.T) image.CreateImage {
	m := image.CreateImage{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Image Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	var newImage image.CreateImage

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("POST", "/v1/images", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to create a new image with the images endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the declared image value.")
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tShould receive a status code of 201 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 201 for the response.", Succeed)

			var u image.Image
			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			newImage = m

			m.ID = "47c658e0-68d7-4d79-9f9f-25ece8a1fb04"
			m.Title = "Image Elijah Baley11"
			m.URL = "/images/1280/720/test-3360-b1396d-1@1x.jpeg"
			m.Slug = "/images/1280/720/test-3360-b1396d-1@1x"
			m.Publisher = "etf1"

			doc, err := json.Marshal(&m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal the response : %v", Failed, err)
			}

			recv := string(doc)
			resp := `{
    "id": "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
    "title": "Image Elijah Baley",
    "url": "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
    "slug": "/images/1280/720/test-2260-b1396d-1@1x",
    "publisher": "etf1",
    "published_at": "2017-06-16T13:42:09.500916Z",
    "expired_at": null,
    "specific": null,
    "created_at": "2017-06-16T13:42:09.500916Z",
    "updated_at": null,
    "restored_at": null,
    "deleted_at": null
  }`

			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}

	return newImage
}

// deleteImage204 validates deleting an image that does exist.
func deleteImage204(t *testing.T, imageID string) {
	r := httptest.NewRequest("DELETE", "/v1/images/"+imageID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting an image that does exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new image %s.", imageID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", Succeed)
		}
	}
}

// getImage200 validates an image request for an existing userid.
func getImage200(t *testing.T, imageID string) {
	r := httptest.NewRequest("GET", "/v1/images/"+imageID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an image that exsits.")
	{
		t.Logf("\tTest 0:\tWhen using the new image %s.", imageID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the response.", Succeed)

			var u image.Image
			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			m := image.Image{}
			*m.ID = "47c658e0-68d7-4d79-9f9f-25ece8a1fb03"
			*m.Title = "Image Elijah Baley"
			*m.URL = "/images/1280/720/test-2260-b1396d-1@1x.jpeg"
			*m.Slug = "/images/1280/720/test-2260-b1396d-1@1x"
			*m.Publisher = "etf1"

			doc, err := json.Marshal(&m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal the response : %v", Failed, err)
			}

			recv := string(doc)
			resp := `{"user_id":"1234","type":1,"first_name":"Bill","last_name":"Kennedy","email":"bill@ardanlabs.com","company":"Ardan Labs","addresses":[{"type":1,"line_one":"12973 SW 112th ST","line_two":"Suite 153","city":"Miami","state":"FL","zipcode":"FL","phone":"305-527-3353","date_modified":null,"date_created":null}],"date_modified":null,"date_created":null}`

			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// putImage204 validates updating an image that does exist.
func putImage204(t *testing.T, m image.CreateImage) {
	m.Title = "Image Elijah Baley"
	m.URL = "/images/1280/720/test-2260-b1396d-1@1x.jpeg"

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("PUT", "/v1/images/"+m.ID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to update an image with the images endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the modified image value.")
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", Succeed)

			r = httptest.NewRequest("GET", "/v1/images/"+m.ID, nil)
			w = httptest.NewRecorder()
			a.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the retrieve : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the retrieve.", Succeed)

			var ru image.Image
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			m.Title = "Image Elijah Baley"
			m.URL = "/images/1280/720/test-2260-b1396d-1@1x.jpeg"

			doc, err := json.Marshal(&m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal the response : %v", Failed, err)
			}

			recv := string(doc)
			resp := `{"user_id":"1234","type":1,"first_name":"Lisa","last_name":"Kennedy","email":"lisa@email.com","company":"Ardan Labs","addresses":[{"type":1,"line_one":"12973 SW 112th ST","line_two":"Suite 153","city":"Miami","state":"NY","zipcode":"FL","phone":"305-527-3353","date_modified":null,"date_created":null}],"date_modified":null,"date_created":null}`

			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

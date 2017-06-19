package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-api/cmd/apid/handlers"
	"github.com/go-api/internal/media"
	"github.com/go-api/internal/platform/db"
	"github.com/go-api/internal/platform/web"

	"strings"

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
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
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
		dbHost = "postgres://images:images@127.0.0.1:54321/images?sslmode=disable"
	}

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")
	masterDB, err := db.NewPSQL("postgres", dbHost)
	if err != nil {
		return 1
	}

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	a = handlers.API(masterDB).(*web.App)

	return m.Run()
}

// TestMedias is the entry point for the medias
func TestMedias(t *testing.T) {
	t.Run("getMedias200Empty", getMedias200Empty)
	t.Run("postMedia400", postMedia400)
	t.Run("getMedia404", getMedia404)
	t.Run("getMedia400", getMedia400)
	t.Run("putMedia404", putMedia404)
	t.Run("crudMedias", crudMedia)
}

// getMedias200Empty validates an empty medias list can be retrieved with the endpoint.
func getMedias200Empty(t *testing.T) {
	r := httptest.NewRequest("GET", "/v1/medias", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to fetch an empty list of medias with the medias endpoint.")
	{
		t.Log("\tTest 0:\tWhen fetching an empty media list.")
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

// postMedia400 validates a media can't be created with the endpoint
// unless a valid media document is submitted.
func postMedia400(t *testing.T) {
	m := media.CreateMedia{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Media Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("POST", "/v1/medias", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate a new media can't be created with an invalid document.")
	{
		t.Log("\tTest 0:\tWhen using an incomplete media value.")
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
      "field_name": "Addresses",
      "error": "required"
    },
    {
      "field_name": "FirstName",
      "error": "required"
    }
  ]
}`,
				`{
  "error": "field validation failure",
  "fields": [
    {
      "field_name": "FirstName",
      "error": "required"
    },
    {
      "field_name": "Addresses",
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
				t.Log("Want:", resps[1])
				t.Fatalf("\t%s\tShould get the expected result.", Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", Succeed)
		}
	}
}

// getMedia400 validates a media request for a malformed mediaid.
func getMedia400(t *testing.T) {
	mediaID := "12345"

	r := httptest.NewRequest("GET", "/v1/medias/"+mediaID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a media with a malformed mediaid.")
	{
		t.Logf("\tTest 0:\tWhen using the new media %s.", mediaID)
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

// getMedia404 validates a media request for a media that does not exist with the endpoint.
func getMedia404(t *testing.T) {
	mediaID := bson.NewObjectId().Hex()

	r := httptest.NewRequest("GET", "/v1/medias/"+mediaID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a media with an unknown id.")
	{
		t.Logf("\tTest 0:\tWhen using the new media %s.", mediaID)
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

// putMedia404 validates updating a media that does not exist.
func putMedia404(t *testing.T) {
	m := media.CreateMedia{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Media Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	mediaID := bson.NewObjectId().Hex()

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("PUT", "/v1/medias/"+mediaID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate updating a media that does not exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new media %s.", mediaID)
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

// crudMedia performs a complete test of CRUD against the api.
func crudMedia(t *testing.T) {
	nm := postMedia201(t)
	defer deleteMedia204(t, nm.ID)

	getMedia200(t, nm.ID)
	putMedia204(t, nm)
}

// postMedia201 validates a media can be created with the endpoint.
func postMedia201(t *testing.T) media.CreateMedia {
	m := media.CreateMedia{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Media Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	var newMedia media.CreateMedia

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("POST", "/v1/medias", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to create a new media with the medias endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the declared media value.")
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tShould receive a status code of 201 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 201 for the response.", Succeed)

			var u media.Media
			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			newMedia = m

			m.ID = "47c658e0-68d7-4d79-9f9f-25ece8a1fb04"
			m.Title = "Media Elijah Baley11"
			m.URL = "/images/1280/720/test-3360-b1396d-1@1x.jpeg"
			m.Slug = "/images/1280/720/test-3360-b1396d-1@1x"
			m.Publisher = "etf1"

			doc, err := json.Marshal(&m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to marshal the response : %v", Failed, err)
			}

			recv := string(doc)
			resp := `{
    "media_id": "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
    "title": "Media Elijah Baley",
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

	return newMedia
}

// deleteMedia204 validates deleting a media that does exist.
func deleteMedia204(t *testing.T, mediaID string) {
	r := httptest.NewRequest("DELETE", "/v1/medias/"+mediaID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting a media that does exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new media %s.", mediaID)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", Succeed)
		}
	}
}

// getMedia200 validates a media request for an existing userid.
func getMedia200(t *testing.T, mediaID string) {
	r := httptest.NewRequest("GET", "/v1/medias/"+mediaID, nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to validate getting a media that exsits.")
	{
		t.Logf("\tTest 0:\tWhen using the new media %s.", mediaID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the response.", Succeed)

			var u media.Media
			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			m := media.Media{}
			*m.ID = "47c658e0-68d7-4d79-9f9f-25ece8a1fb03"
			*m.Title = "Media Elijah Baley"
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

// putMedia204 validates updating a media that does exist.
func putMedia204(t *testing.T, m media.CreateMedia) {
	m.Title = "Media Elijah Baley"
	m.URL = "/images/1280/720/test-2260-b1396d-1@1x.jpeg"

	body, _ := json.Marshal(&m)
	r := httptest.NewRequest("PUT", "/v1/medias/"+m.ID, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	t.Log("Given the need to update a media with the medias endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the modified media value.")
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", Succeed)

			r = httptest.NewRequest("GET", "/v1/medias/"+m.ID, nil)
			w = httptest.NewRecorder()
			a.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the retrieve : %v", Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the retrieve.", Succeed)

			var ru media.Media
			if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", Failed, err)
			}

			m.Title = "Media Elijah Baley"
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

package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Swagger structure
type Swagger struct {
	URL string
}

// GetAPIDocs swagger api documentation
func (s *Swagger) GetAPIDocs(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {

	buff, err := ioutil.ReadFile("./cmd/apid/statics/swagger/swagger.yaml")
	if err != nil {
		return errors.Wrap(err, "GetApiDocs")
	}

	buff = []byte(strings.Replace(string(buff), "{{url}}", s.URL, -1))
	if _, err := w.Write(buff); err != nil {
		return errors.Wrap(err, "GetApiDocs")
	}
	return nil
}

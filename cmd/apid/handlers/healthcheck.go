package handlers

import (
	"context"
	"net/http"

	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/web"
)

// Healthzcheck represents the Healthzcheck API method handler set.
type Healthzcheck struct {
	MasterDB *db.DB
	// ADD OTHER STATE LIKE THE LOGGER AND CONFIG HERE.
}

// HealthCheckResp structure
type HealthCheckResp struct {
	Result  bool     `json:"result"`
	Errors  []string `json:"errors"`
	Version string   `json:"version"`
}

// Healthz returns HealthCheckResp.
// 200 Success
func (h *Healthzcheck) Healthz(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	res := HealthCheckResp{Result: true, Errors: []string{}, Version: "1.0.0 "}
	web.Respond(ctx, w, res, http.StatusOK)
	return nil
}

// Readiness returns HealthCheckResp.
// 200 Success, 500 Internal
func (h *Healthzcheck) Readiness(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := h.MasterDB

	if err := reqDB.Ping(); err != nil {
		res := HealthCheckResp{Result: false, Errors: []string{err.Error()}, Version: "1.0.0"}
		web.Respond(ctx, w, res, http.StatusInternalServerError)
		return nil
	}
	res := HealthCheckResp{Result: true, Errors: []string{}, Version: "1.0.0"}
	web.Respond(ctx, w, res, http.StatusOK)
	return nil
}

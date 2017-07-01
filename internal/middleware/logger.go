package middleware

import (
	"context"

	"net/http"
	"time"

	"github.com/jdelobel/go-api/internal/platform/web"
)

// RequestLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func RequestLogger(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
		v := ctx.Value(web.KeyValues).(*web.Values)
		if err := next(ctx, w, r, params); err != nil {
			v.Log.Errorf("RequestLogger: middleware next failed: %v", err)
		}

		v.Log.Infof("%s : (%d) : %s %s -> %s (%s)",
			v.TraceID,
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now),
		)

		// This is the top of the food chain. At this point all error
		// handling has been done including logging.
		return nil
	}
}

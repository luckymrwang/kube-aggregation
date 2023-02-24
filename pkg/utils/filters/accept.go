package filters

import (
	"net/http"

	"github.com/inspur/pkg/utils/request"
)

func WithAcceptHeader(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(request.WithAcceptHeader(req.Context(), req.Header.Get("Accept")))
		handler.ServeHTTP(w, req)
	})
}

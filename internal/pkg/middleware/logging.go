package middleware

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/b0pof/avito-internship/pkg/logger"
)

var ErrHijackAssertion = errors.New("type assertion to http.Hijacker failed")

var currRequestID uint64

type responseWriterInterceptor struct {
	w          http.ResponseWriter
	statusCode int
}

func newResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor {
	return &responseWriterInterceptor{
		w: w,
	}
}

func (wi *responseWriterInterceptor) WriteHeader(statusCode int) {
	wi.w.WriteHeader(statusCode)
	wi.statusCode = statusCode
}

func (wi *responseWriterInterceptor) Header() http.Header {
	return wi.w.Header()
}

func (wi *responseWriterInterceptor) Write(d []byte) (int, error) {
	return wi.w.Write(d)
}

func (wi *responseWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := wi.w.(http.Hijacker)
	if !ok {
		return nil, nil, ErrHijackAssertion
	}
	return h.Hijack()
}

func (wi *responseWriterInterceptor) GetStatusCode() int {
	return wi.statusCode
}

func NewLoggingMiddleware(l *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddUint64(&currRequestID, 1)
			requestLogger := l.With(slog.Uint64("requestID", currRequestID))
			requestLogger.Info("new",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI))

			wi := newResponseWriterInterceptor(w)
			ctx := logger.WithContext(r.Context(), requestLogger)
			start := time.Now()
			next.ServeHTTP(wi, r.Clone(ctx))
			dur := time.Since(start)
			statusCode := wi.GetStatusCode()
			requestLogger.Info("response",
				slog.Int("statusCode", statusCode),
				slog.String("duration", dur.String()))
		})
	}
}

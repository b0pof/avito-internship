package helper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/b0pof/avito-internship/pkg/logger"
)

func Respond(ctx context.Context, w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error(ctx, "marshall error: "+err.Error())
	}
}

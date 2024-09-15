package helper

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/b0pof/avito-internship/internal/model"
)

const (
	_defaultLimit  = 5
	_defaultOffset = 0
)

func ParseLimitOffset(r *http.Request) (int, int, error) {
	var err error
	var limit, offset int

	limitStr, ok := r.URL.Query()["limit"]
	if !ok {
		limit = _defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr[0])
		if err != nil {
			return 0, 0, model.ErrInvalidQueryParam
		}
	}
	offsetStr, ok := r.URL.Query()["offset"]
	if !ok {
		offset = _defaultOffset
	} else {
		offset, err = strconv.Atoi(offsetStr[0])
		if err != nil {
			return 0, 0, model.ErrInvalidQueryParam
		}
	}
	return limit, offset, nil
}

func ParseServiseTypes(r *http.Request) []string {
	types, _ := r.URL.Query()["service_type"]
	return types
}

func ParseUsername(r *http.Request) string {
	username, _ := r.URL.Query()["username"]
	if len(username) == 0 {
		return ""
	}
	return username[0]
}

func ParseTenderID(r *http.Request) string {
	tenderID, _ := mux.Vars(r)["tenderId"]
	return tenderID
}

func ParseVersion(r *http.Request) (int, error) {
	version, err := strconv.Atoi(mux.Vars(r)["version"])
	if err != nil {
		return 0, model.ErrInvalidPathParam
	}
	return version, nil
}

func ParseBidID(r *http.Request) string {
	bidID, _ := mux.Vars(r)["bidId"]
	return bidID
}

func ParseStatus(r *http.Request) string {
	status, _ := r.URL.Query()["status"]
	if len(status) == 0 {
		return ""
	}
	return status[0]
}

func ParseDecision(r *http.Request) string {
	dec, _ := r.URL.Query()["decision"]
	if len(dec) == 0 {
		return ""
	}
	return dec[0]
}

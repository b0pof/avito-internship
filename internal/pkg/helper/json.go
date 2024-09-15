package helper

import (
	"encoding/json"
	"net/http"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/internal/repository"
	"github.com/b0pof/avito-internship/internal/usecase"
)

func ParseTenderFromBody(r *http.Request) (usecase.CreateTenderInput, error) {
	var tender usecase.CreateTenderInput
	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		return usecase.CreateTenderInput{}, model.ErrInvalidBody
	}
	return tender, nil
}

func ParseBidFromBody(r *http.Request) (repository.CreateBidInput, error) {
	var bid repository.CreateBidInput
	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		return repository.CreateBidInput{}, model.ErrInvalidBody
	}
	return bid, nil
}

type UpdateTenderInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}

func ParseUpdateTenderInfo(r *http.Request) (UpdateTenderInfo, error) {
	var info UpdateTenderInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return UpdateTenderInfo{}, model.ErrInvalidBody
	}
	return info, nil
}

type UpdateBidInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ParseUpdateBidInfo(r *http.Request) (UpdateBidInfo, error) {
	var info UpdateBidInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return UpdateBidInfo{}, model.ErrInvalidBody
	}
	return info, nil
}

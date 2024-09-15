package usecase

import (
	"context"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/internal/repository"
)

func (u *Usecase) GetTenders(ctx context.Context, input repository.GetTendersInput) ([]model.Tender, error) {
	return u.repo.GetTenders(ctx, input)
}

type CreateTenderInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationID  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

func (u *Usecase) CreateTender(ctx context.Context, input CreateTenderInput) (model.Tender, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.CreatorUsername)
	if err != nil {
		return model.Tender{}, err
	}
	if !u.repo.IsUserOrganizationResponsible(ctx, userID, input.OrganizationID) {
		return model.Tender{}, model.ErrNoRights
	}
	return u.repo.CreateTender(ctx, repository.CreateTenderInput{
		Name:           input.Name,
		Description:    input.Description,
		ServiceType:    input.ServiceType,
		OrganizationID: input.OrganizationID,
		CreatorID:      userID,
	})
}

type GetMyTendersInput struct {
	Limit    int
	Offset   int
	Username string
}

func (u *Usecase) GetMyTenders(ctx context.Context, input GetMyTendersInput) ([]model.Tender, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	return u.repo.GetMyTenders(ctx, repository.GetMyTendersInput{
		Limit:  input.Limit,
		Offset: input.Offset,
		UserID: userID,
	})
}

type GetTenderStatusInput struct {
	TenderID string
	Username string
}

func (u *Usecase) GetTenderStatus(ctx context.Context, input GetTenderStatusInput) (string, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return "", err
	}
	status, err := u.repo.GetTenderStatus(ctx, input.TenderID)
	if err != nil {
		return "", err
	}
	if !u.repo.IsUserResponsibleForTender(ctx, input.TenderID, userID) {
		return "", model.ErrNoRights
	}
	return status, nil
}

type UpdateTenderStatusInput struct {
	TenderID string `json:"tenderId"`
	Status   string `json:"status"`
	Username string `json:"username"`
}

func (u *Usecase) UpdateTenderStatus(ctx context.Context, input UpdateTenderStatusInput) (model.Tender, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Tender{}, err
	}
	if !u.repo.TenderExists(ctx, input.TenderID) {
		return model.Tender{}, model.ErrTenderNotFound
	}
	if !u.repo.IsUserResponsibleForTender(ctx, input.TenderID, userID) {
		return model.Tender{}, model.ErrNoRights
	}
	return u.repo.UpdateTenderStatus(ctx, repository.UpdateTenderStatusInput{
		Status:   input.Status,
		TenderID: input.TenderID,
	})
}

type UpdateTenderInput struct {
	TenderID    string
	Username    string
	Name        string
	Description string
	ServiceType string
}

func (u *Usecase) UpdateTender(ctx context.Context, input UpdateTenderInput) (model.Tender, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Tender{}, err
	}
	if !u.repo.TenderExists(ctx, input.TenderID) {
		return model.Tender{}, model.ErrTenderNotFound
	}
	if !u.repo.IsUserResponsibleForTender(ctx, input.TenderID, userID) {
		return model.Tender{}, model.ErrNoRights
	}
	return u.repo.UpdateTender(ctx, repository.EditTenderInput{
		TenderID:    input.TenderID,
		Name:        input.Name,
		Description: input.Description,
		ServiceType: input.ServiceType,
	})
}

type RollbackTenderInput struct {
	TenderID string
	Version  int
	Username string
}

func (u *Usecase) RollbackTender(ctx context.Context, input RollbackTenderInput) (model.Tender, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Tender{}, err
	}
	if !u.repo.TenderExists(ctx, input.TenderID) {
		return model.Tender{}, model.ErrTenderNotFound
	}
	hasVersion, err := u.repo.TenderHasVersion(ctx, repository.TenderHasVersionInput{
		TenderID: input.TenderID,
		Version:  input.Version,
	})
	if err != nil {
		return model.Tender{}, err
	}
	if !hasVersion {
		return model.Tender{}, model.ErrNoSuchVersion
	}
	if !u.repo.IsUserResponsibleForTender(ctx, input.TenderID, userID) {
		return model.Tender{}, model.ErrNoRights
	}
	return u.repo.RollbackTender(ctx, repository.RollbackTenderInput{
		TenderID: input.TenderID,
		Version:  input.Version,
	})
}

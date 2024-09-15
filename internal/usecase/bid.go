package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/internal/repository"
)

func (u *Usecase) CreateBid(ctx context.Context, input repository.CreateBidInput) (model.Bid, error) {
	if !u.repo.IsUserExist(ctx, input.AuthorID) {
		return model.Bid{}, model.ErrUserNotFound
	}
	if input.AuthorType == "Organization" {
		_, err := u.repo.GetOrganizationIDByEmployeeID(ctx, input.AuthorID)
		if err != nil {
			return model.Bid{}, errors.Wrap(err, "невозможно созодать предложение от имени организации")
		}
	}
	status, err := u.repo.GetTenderStatus(ctx, input.TenderID)
	if err != nil {
		return model.Bid{}, err
	}
	if status != "Published" {
		return model.Bid{}, errors.Wrap(model.ErrNoRights, "тендер не опубликован")
	}
	return u.repo.CreateBid(ctx, input)
}

type GetMyBidsInput struct {
	Limit    int
	Offset   int
	Username string
}

func (u *Usecase) GetMyBids(ctx context.Context, input GetMyBidsInput) ([]model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	bids, err := u.repo.GetMyBids(ctx, repository.GetMyBidsInput{
		Limit:  input.Limit,
		Offset: input.Offset,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (u *Usecase) GetTenderBids(ctx context.Context, input repository.GetTenderBidsInput) ([]model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if !u.repo.IsTenderExist(ctx, input.TenderID) {
		return nil, model.ErrTenderNotFound
	}
	bids, err := u.repo.GetTenderBids(ctx, input)
	if err != nil {
		return nil, err
	}
	if !u.repo.IsUserResponsibleForTender(ctx, input.TenderID, userID) {
		return nil, model.ErrNoRights
	}
	return bids, nil
}

type GetBidStatusInput struct {
	BidID    string
	Username string
}

func (u *Usecase) GetBidStatus(ctx context.Context, input GetBidStatusInput) (string, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return "", err
	}
	status, err := u.repo.GetBidStatus(ctx, input.BidID)
	if err != nil {
		return "", err
	}
	hasAccess, err := u.repo.IsBidVisibleForUser(ctx, userID, input.BidID)
	if err != nil {
		return "", err
	}
	if !hasAccess {
		return "", model.ErrNoRights
	}
	return status, nil
}

type UpdateBidStatusInput struct {
	BidID    string
	Status   string
	Username string
}

func (u *Usecase) UpdateBidStatus(ctx context.Context, input UpdateBidStatusInput) (model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Bid{}, err
	}
	if !u.repo.BidExists(ctx, input.BidID) {
		return model.Bid{}, model.ErrNoBidFound
	}
	hasAccess, err := u.repo.IsBidVisibleForUser(ctx, userID, input.BidID)
	if err != nil {
		return model.Bid{}, err
	}
	if !hasAccess {
		return model.Bid{}, model.ErrNoRights
	}
	return u.repo.UpdateBidStatus(ctx, repository.UpdateBidStatusInput{
		Status: input.Status,
		BidID:  input.BidID,
	})
}

type SubmitDecisionInput struct {
	BidID    string
	Username string
	Decision string
}

func (u *Usecase) SubmitDecision(ctx context.Context, input SubmitDecisionInput) (model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Bid{}, err
	}
	if !u.repo.UserCanSubmitDecision(ctx, input.BidID, userID) {
		return model.Bid{}, model.ErrNoRights
	}
	if input.Decision == "Approved" {
		if err := u.repo.CloseTenderByBidID(ctx, input.BidID); err != nil {
			return model.Bid{}, err
		}
		return u.repo.GetBidByID(ctx, input.BidID)
	} else if input.Decision == "Rejected" {
		return u.repo.UpdateBidStatus(ctx, repository.UpdateBidStatusInput{
			Status: "Canceled",
			BidID:  input.BidID,
		})
	} else {
		return model.Bid{}, model.ErrWrongDecision
	}
}

type UpdateBidInput struct {
	BidID       string
	Username    string
	Name        string
	Description string
}

func (u *Usecase) UpdateBid(ctx context.Context, input UpdateBidInput) (model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Bid{}, err
	}
	if !u.repo.BidExists(ctx, input.BidID) {
		return model.Bid{}, model.ErrNoBidFound
	}
	authorID, err := u.repo.GetUserIDByBidID(ctx, input.BidID)
	if err != nil {
		return model.Bid{}, err
	}
	if userID != authorID {
		return model.Bid{}, model.ErrNoRights
	}
	return u.repo.UpdateBid(ctx, repository.EditBidInput{
		BidID:       input.BidID,
		Name:        input.Name,
		Description: input.Description,
	})
}

type RollbackBidInput struct {
	BidID    string
	Version  int
	Username string
}

func (u *Usecase) RollbackBid(ctx context.Context, input RollbackBidInput) (model.Bid, error) {
	userID, err := u.repo.GetUserIDByUsername(ctx, input.Username)
	if err != nil {
		return model.Bid{}, err
	}
	if !u.repo.BidExists(ctx, input.BidID) {
		return model.Bid{}, model.ErrNoBidFound
	}
	hasVersion, err := u.repo.BidHasVersion(ctx, repository.BidHasVersionInput{
		BidID:   input.BidID,
		Version: input.Version,
	})
	if err != nil {
		return model.Bid{}, err
	}
	if !hasVersion {
		return model.Bid{}, model.ErrNoSuchVersion
	}
	authorID, err := u.repo.GetUserIDByBidID(ctx, input.BidID)
	if err != nil {
		return model.Bid{}, err
	}
	if userID != authorID {
		return model.Bid{}, model.ErrNoRights
	}
	return u.repo.RollbackBid(ctx, repository.RollbackBidInput{
		BidID:   input.BidID,
		Version: input.Version,
	})
}

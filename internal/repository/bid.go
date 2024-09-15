package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/pkg/logger"
)

func (r *Repository) GetBidByID(ctx context.Context, bidID string) (model.Bid, error) {
	q := `SELECT b.id, bv.name, b.status, b.author_type, b.author_id, bv.version, b.created_at
		FROM bid_version bv
			INNER JOIN bid b ON bv.bid_id = b.id
		WHERE bid_id = $1
    		AND bv.version = (SELECT MAX(version) FROM bid_version WHERE bid_id = $1);`

	var bid model.Bid
	if err := r.db.GetContext(ctx, &bid, q, bidID); err != nil {
		return model.Bid{}, model.ErrNoBidFound
	}
	return bid, nil
}

type CreateBidInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderID    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorID    string `json:"authorId"`
}

func (r *Repository) CreateBid(ctx context.Context, input CreateBidInput) (model.Bid, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return model.Bid{}, model.ErrInternal
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `INSERT INTO bid (tender_id, author_type, author_id)
			VALUES ($1, $2, $3)
			RETURNING id;`

	var bidID string
	if err = tx.GetContext(ctx, &bidID, q, input.TenderID, input.AuthorType, input.AuthorID); err != nil {
		logger.Error(ctx, err.Error())
		return model.Bid{}, model.ErrInternal
	}

	q = `INSERT INTO bid_version (bid_id, name, description, version)
			SELECT b.id, $1, $2, COALESCE(bv.version, 0) + 1
			FROM bid b
				LEFT JOIN bid_version bv ON b.id = bv.bid_id
			WHERE b.id = $3
			ORDER BY version DESC
			LIMIT 1;`

	if _, err = tx.ExecContext(ctx, q, input.Name, input.Description, bidID); err != nil {
		logger.Error(ctx, err.Error())
		return model.Bid{}, model.ErrInternal
	}
	_ = tx.Commit()
	return r.GetBidByID(ctx, bidID)
}

type GetMyBidsInput struct {
	Limit  int
	Offset int
	UserID string
}

func (r *Repository) GetMyBids(ctx context.Context, input GetMyBidsInput) ([]model.Bid, error) {
	q := `SELECT b.id, bv.name, b.status, b.author_type, b.author_id, bv.version, b.created_at
			FROM bid b
				 INNER JOIN bid_version bv
					ON b.id = bv.bid_id
						AND bv.version = (
							SELECT MAX(version)
							FROM bid_version
							WHERE bid_id = b.id
						)
			WHERE b.author_id = $1
			ORDER BY bv.name
			LIMIT $2
			OFFSET $3;`

	bids := make([]model.Bid, 0)
	if err := r.db.SelectContext(ctx, &bids, q, input.UserID, input.Limit, input.Offset); err != nil {
		logger.Error(ctx, err.Error())
		return nil, model.ErrInternal
	}
	return bids, nil
}

func (r *Repository) BidExists(ctx context.Context, bidID string) bool {
	q := `SELECT id FROM bid WHERE id = $1`

	var foundBidID string
	if err := r.db.GetContext(ctx, &foundBidID, q, bidID); err != nil {
		logger.Error(ctx, err.Error())
		return false
	}
	return foundBidID == bidID
}

type GetTenderBidsInput struct {
	TenderID string
	Username string
	Limit    int
	Offset   int
}

func (r *Repository) GetTenderBids(ctx context.Context, input GetTenderBidsInput) ([]model.Bid, error) {
	q := `SELECT b.id, bv.name, b.status, b.author_type, b.author_id, bv.version, b.created_at
			FROM bid b
				INNER JOIN bid_version bv
					ON b.id = bv.bid_id
						AND bv.version = (
							SELECT MAX(version)
							FROM bid_version
							WHERE bid_id = b.id
						)
			WHERE b.tender_id = $1
			ORDER BY bv.name
			LIMIT $2
			OFFSET $3;`

	bids := make([]model.Bid, 0)
	if err := r.db.SelectContext(ctx, &bids, q, input.TenderID, input.Limit, input.Offset); err != nil {
		return nil, model.ErrNoBidsFound
	}
	return bids, nil
}

func (r *Repository) GetBidStatus(ctx context.Context, bidID string) (string, error) {
	q := `SELECT status FROM bid WHERE id = $1`

	var status string
	if err := r.db.GetContext(ctx, &status, q, bidID); err != nil {
		return "", model.ErrNoBidFound
	}
	return status, nil
}

type UpdateBidStatusInput struct {
	BidID  string
	Status string
}

func (r *Repository) UpdateBidStatus(ctx context.Context, input UpdateBidStatusInput) (model.Bid, error) {
	q := `UPDATE bid
			SET status = $1
			WHERE id = $2;`

	if _, err := r.db.ExecContext(ctx, q, input.Status, input.BidID); err != nil {
		return model.Bid{}, model.ErrInvalidAttributeValue
	}
	return r.GetBidByID(ctx, input.BidID)
}

type EditBidInput struct {
	BidID       string
	Name        string
	Description string
}

func (r *Repository) UpdateBid(ctx context.Context, input EditBidInput) (model.Bid, error) {
	q := `INSERT INTO bid_version (bid_id, name, description, version)
			SELECT b.id, %s, %s, COALESCE(bv.version, 0) + 1
			FROM bid b
				LEFT JOIN bid_version bv ON b.id = bv.bid_id
			WHERE b.id = $1
			ORDER BY version DESC
			LIMIT 1;`

	format := "'%s'"

	var name = "name"
	if input.Name != "" {
		name = fmt.Sprintf(format, input.Name)
	}
	var description = "description"
	if input.Description != "" {
		description = fmt.Sprintf(format, input.Description)
	}
	q = fmt.Sprintf(q, name, description)

	if _, err := r.db.ExecContext(ctx, q, input.BidID); err != nil {
		if strings.Contains(err.Error(), "invalid value") {
			return model.Bid{}, model.ErrInvalidAttributeValue
		}
		logger.Error(ctx, err.Error())
		return model.Bid{}, model.ErrInternal
	}
	return r.GetBidByID(ctx, input.BidID)
}

type BidHasVersionInput struct {
	BidID   string
	Version int
}

func (r *Repository) BidHasVersion(ctx context.Context, input BidHasVersionInput) (bool, error) {
	q := `SELECT (SELECT bid_id
			FROM bid_version
			WHERE bid_id = $1
			  AND version = $2
		) IS NOT NULL;`

	var hasVersion bool
	if err := r.db.GetContext(ctx, &hasVersion, q, input.BidID, input.Version); err != nil {
		logger.Error(ctx, err.Error())
		return false, model.ErrInternal
	}
	return hasVersion, nil
}

type RollbackBidInput struct {
	BidID   string
	Version int
}

func (r *Repository) RollbackBid(ctx context.Context, input RollbackBidInput) (model.Bid, error) {
	q := `INSERT INTO bid_version (bid_id, name, description, version)
			SELECT b.id, bv.name, bv.description, COALESCE(
				(SELECT MAX(version) FROM bid_version WHERE bid_id = $1), 0
			) + 1
			FROM bid b
				LEFT JOIN bid_version bv ON b.id = bv.bid_id
			WHERE b.id = $1 AND bv.version = $2;`

	if _, err := r.db.ExecContext(ctx, q, input.BidID, input.Version); err != nil {
		logger.Error(ctx, err.Error())
		return model.Bid{}, model.ErrInternal
	}
	return r.GetBidByID(ctx, input.BidID)
}

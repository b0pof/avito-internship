package repository

import (
	"context"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/pkg/logger"
)

func (r *Repository) IsUserOrganizationResponsible(ctx context.Context, userID, orgID string) bool {
	q := `SELECT user_id
		FROM organization_responsible
		WHERE user_id = $1 AND organization_id = $2;`

	var foundUserID string
	if err := r.db.GetContext(ctx, &foundUserID, q, userID, orgID); err != nil {
		return false
	}
	if foundUserID == userID {
		return true
	}
	return false
}

func (r *Repository) GetUserIDByUsername(ctx context.Context, username string) (string, error) {
	q := `SELECT id FROM employee WHERE username = $1`

	var id string
	if err := r.db.GetContext(ctx, &id, q, username); err != nil {
		logger.Error(ctx, err.Error())
		return "", model.ErrUserNotFound
	}
	return id, nil
}

func (r *Repository) IsUserResponsibleForTender(ctx context.Context, tenderID, userID string) bool {
	q := `SELECT r.user_id
		FROM tender t 
			INNER JOIN organization o ON t.organization_id = o.id
			INNER JOIN organization_responsible r ON o.id = r.organization_id
		WHERE t.id = $1 AND r.user_id = $2;`

	var foundUserID string
	if err := r.db.GetContext(ctx, &foundUserID, q, tenderID, userID); err != nil {
		return false
	}
	if foundUserID == "" {
		return false
	}
	return true
}

func (r *Repository) IsBidVisibleForUser(ctx context.Context, userID, bidID string) (bool, error) {
	q := `SELECT $1 IN (
			SELECT r.user_id
			FROM bid b
				INNER JOIN tender t ON b.tender_id = t.id
				INNER JOIN organization o ON t.organization_id = o.id
				INNER JOIN organization_responsible r ON o.id = r.organization_id
			WHERE b.id = $2
			UNION DISTINCT
			SELECT author_id
			FROM bid
			WHERE id = $2
		);`
	var visible bool
	if err := r.db.GetContext(ctx, &visible, q, userID, bidID); err != nil {
		logger.Error(ctx, err.Error())
		return false, model.ErrInternal
	}
	return visible, nil
}

func (r *Repository) UserCanSubmitDecision(ctx context.Context, bidID, userID string) bool {
	q := `SELECT r.user_id
			FROM bid b
				INNER JOIN tender t ON b.tender_id = t.id
				INNER JOIN organization o ON t.organization_id = o.id
				INNER JOIN organization_responsible r ON o.id = r.organization_id
			WHERE b.id = $1 AND r.user_id = $2;`

	var foundUserID string
	if err := r.db.GetContext(ctx, &foundUserID, q, bidID, userID); err != nil {
		return false
	}
	if foundUserID == userID {
		return true
	}
	return false
}

func (r *Repository) GetUserIDByBidID(ctx context.Context, bidID string) (string, error) {
	q := `SELECT author_id
			FROM bid
			WHERE id = $1;`

	var authorID string
	if err := r.db.GetContext(ctx, &authorID, q, bidID); err != nil {
		return "", model.ErrUserNotFound
	}
	return authorID, nil
}

func (r *Repository) IsUserExist(ctx context.Context, userID string) bool {
	q := `SELECT id FROM employee WHERE id = $1;`

	var foundUserID string
	if err := r.db.GetContext(ctx, &foundUserID, q, userID); err != nil {
		logger.Error(ctx, err.Error())
		return false
	}
	return foundUserID == userID
}

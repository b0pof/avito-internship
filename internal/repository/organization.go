package repository

import (
	"context"

	"github.com/b0pof/avito-internship/internal/model"
)

func (r *Repository) GetOrganizationIDByEmployeeID(ctx context.Context, employeeID string) (string, error) {
	q := `SELECT r.organization_id
		FROM organization_responsible r
		WHERE r.user_id = $1;`

	var orgID string
	if err := r.db.Get(&orgID, q, employeeID); err != nil {
		return "", model.ErrNoOrganizationFound
	}
	return orgID, nil
}

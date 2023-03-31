package users

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetByUserName(ctx context.Context, username string) (*User, error)
	GetUserInfo(ctx context.Context, username string) (*UserInfo, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetByUserName(ctx context.Context, username string) (*User, error) {
	query := `SELECT
			username,
			password
		FROM users
		WHERE username = ?`

	query = r.db.Rebind(query)
	var user User
	err := r.db.GetContext(ctx, &user, query, username)
	return &user, err
}

func (r *repo) GetUserInfo(ctx context.Context, username string) (*UserInfo, error) {
	query := `SELECT
			u.username,
			a.balance,
			a.currency
		FROM users AS u
		JOIN accounts AS a ON u.id = a.user_id
		WHERE u.username = ?  LIMIT 1`

	query = r.db.Rebind(query)
	var userAcc UserInfo
	err := r.db.GetContext(ctx, &userAcc, query, username)
	return &userAcc, err
}

package transfer

import (
	"context"
	"fmt"
	"thunes-api/errors"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	GetAccountById(ctx context.Context, id int64) (*Account, error)
	GetAccountByUserName(ctx context.Context, username string) (*UserAccount, error)
	TransferTx(ctx context.Context, arg TransferTxParams) (*UserAccount, error)
	GetAllBeneficiaries(ctx context.Context, username string) ([]*Beneficiary, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) (Repo, error) {
	return &repo{
		db: db,
	}, nil
}

func (r *repo) GetAccountById(ctx context.Context, id int64) (*Account, error) {
	query := `SELECT
			 balance,
			 currency
		FROM accounts 
		WHERE id = ?`

	query = r.db.Rebind(query)
	var account Account
	err := r.db.GetContext(ctx, &account, query, id)
	fmt.Println(err)
	return &account, err
}

func (r *repo) GetAccountByUserName(ctx context.Context, username string) (*UserAccount, error) {
	query := `SELECT
			u.id,
			u.username,
			a.balance,
			a.currency
		FROM users AS u
		JOIN accounts AS a ON u.id = a.user_id
		WHERE u.username = ?  LIMIT 1`

	query = r.db.Rebind(query)
	var userAcc UserAccount
	err := r.db.GetContext(ctx, &userAcc, query, username)
	fmt.Println(err)
	return &userAcc, err
}

func (r *repo) TransferTx(ctx context.Context, arg TransferTxParams) (*UserAccount, error) {
	// Create a new context, and begin a transaction
	ctxBg := context.Background()
	tx, err := r.db.BeginTx(ctxBg, nil)
	if err != nil {
		fmt.Println(err)
		return nil, errors.SystemError
	}

	_, err = tx.ExecContext(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	if err != nil {
		fmt.Println(err)
		// Incase we find any error in the query execution, rollback the transaction
		tx.Rollback()
		return nil, errors.SystemError
	}

	_, err = tx.ExecContext(ctx, createEntry, arg.FromAccountID, -arg.Amount)
	if err != nil {
		fmt.Println(err)
		//rollback the transaction
		tx.Rollback()
		return nil, errors.SystemError
	}

	_, err = tx.ExecContext(ctx, createEntry, arg.ToAccountID, arg.Amount)
	if err != nil {
		fmt.Println(err)
		//rollback the transaction
		tx.Rollback()
		return nil, errors.SystemError
	}

	_, err = tx.ExecContext(ctx, addAccountBalance, -arg.Amount, arg.FromAccountID)
	if err != nil {
		fmt.Println(err)
		//rollback the transaction
		tx.Rollback()
		return nil, errors.SystemError
	}
	_, err = tx.ExecContext(ctx, addAccountBalance, +arg.Amount, arg.ToAccountID)
	if err != nil {
		fmt.Println(err)
		//rollback the transaction
		tx.Rollback()
		return nil, errors.SystemError
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return nil, errors.SystemError
	}

	accDetails, _ := r.GetAccountByUserName(ctx, arg.Username)
	return accDetails, nil
}

func (r *repo) GetAllBeneficiaries(ctx context.Context, username string) ([]*Beneficiary, error) {
	query := `SELECT
			a.id,
			u.username
		FROM users AS u
		JOIN accounts AS a ON u.id = a.user_id
		WHERE u.username != ? `
	query = r.db.Rebind(query)
	var userAcc []*Beneficiary
	err := r.db.SelectContext(ctx, &userAcc, query, username)
	return userAcc, err
}

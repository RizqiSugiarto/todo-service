package shared

import (
	"database/sql"
	"errors"
)

const TRANSACTION_ERR string = "Transaction is nil"

type TransactionManager interface {
	StartTransaction() error
	SaveTransaction() error
	CancelTransaction() error
}

type Database struct {
	*sql.DB
}

type Transaction struct {
	*sql.Tx
}

type SqlTransactionManager struct {
	Db *Database
	Tx *Transaction
}

func (r *SqlTransactionManager) StartTransaction() error {
	tx, err := r.Db.DB.Begin()
	if err != nil {
		return err
	}

	r.Tx = &Transaction{
		Tx: tx,
	}

	return nil
}

func (r *SqlTransactionManager) SaveTransaction() error {
	if r.Tx == nil {
		return errors.New(TRANSACTION_ERR)
	}

	err := r.Tx.Commit()
	if err != nil {
		errRB := r.CancelTransaction()
		if errRB != nil {
			return errRB
		}

		return err
	}

	r.Tx = nil

	return nil
}

func (r *SqlTransactionManager) CancelTransaction() error {
	if r.Tx == nil {
		return errors.New(TRANSACTION_ERR)
	}

	err := r.Tx.Rollback()
	if err != nil {
		return err
	}

	r.Tx = nil

	return nil
}

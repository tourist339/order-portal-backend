package model

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TxKey struct{}
type Tx struct {
	Tx *sqlx.Tx
}

type Query interface {
	Get(ctx context.Context, q *SelectQuery, u any) error
	GetByID(ctx context.Context, id, tableName string, fields []string, u any) error
}

type Write interface {
	Insert(ctx context.Context, tableName string, data any) (string, error)
	Update(ctx context.Context, tableName string, data any) error
}

func (s *Service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	sqtx, err := s.GetTransaction(ctx)
	childTransaction := false
	if err != nil {
		sqtx, err = s.db.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}
	} else {
		childTransaction = true
	}
	ctxTx := ctx.Value(TxKey{}).(*Tx)
	ctxTx.Tx = sqtx
	defer func() {
		ctxTx := ctx.Value(TxKey{}).(*Tx)
		ctxTx.Tx = nil
	}()
	err = fn(ctx)
	if childTransaction {
		return err
	}
	if err != nil {
		fmt.Println(err)
		rbErr := sqtx.Rollback()
		if rbErr != nil {
			fmt.Println("Error rolling back transaction:", rbErr)
			return rbErr
		}
		return fmt.Errorf("Error executing transaction: %s", err.Error())
	}
	err = sqtx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction: %s", err.Error())
	}
	return nil
}

func (s *Service) GetTransaction(ctx context.Context) (*sqlx.Tx, error) {
	tx := ctx.Value(TxKey{}).(*Tx)
	if tx == nil {
		return nil, fmt.Errorf("transaction not found in context")
	}
	if tx.Tx != nil {
		return tx.Tx, nil
	}
	return nil, fmt.Errorf("transaction not found in context")
}

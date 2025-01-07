package model

import "context"

// Transactional is a `Persistence` that supports transactions.
//
// Always `Release` the transaction even if `Abort` is called. The abort or commit
// is performed in the `Release` method.
//
// [source,go]
// ----
// tx, err := p.Begin(ctx)
//
//	if err != nil {
//	    return err
//	}
//
// defer p.Release(ctx, tx) // <1>
//
// // read and write using `Persistence`
// // ... when out of scope, the `Release` will be called -> commit
// ----
// <1> Always release the transaction to ensure it will either commit or rollback the transaction.
type Transactional interface {
	Persistence
	// Begin will start a new transaction.
	Begin(ctx context.Context) (*Transaction, error)
	// Release will either commit or rollback the transaction.
	Release(ctx context.Context, tx *Transaction) error
	// Abort will mark the the transaction as aborted.
	Abort(ctx context.Context, tx *Transaction) error
}

type Transaction struct {
	// ID is the unique identifier of the transaction.
	ID string
	// Custom is where the `Persistence` can store custom data such as
	// sessions or other data that is needed to be shared.
	//
	// NOTE: A client shall never assume anything or interact with this
	// data directly. It is for the `Persistence` to use.
	Custom map[string]any
}

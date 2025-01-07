package persistencemodel

import "context"

// Transactional is a `Persistence` that supports transactions.
//
// Always `Release` the transaction even if `Abort` is called. The abort or commit
// is performed in the `Release` method.
//
// .Example Update Model Transaction
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
// res := p.Read(
//
//						ctx,
//						ReadOptions{Tx:tx},
//						ReadOperation{
//							ID: PersistenceID{ID:"Östra Nygatan 11B", "temperatures"},
//							ModelType: ModelTypeReported,
//							Model: typeof(MyModel{}), // <2>
//	  			) // <3>
//
//		if res[0].Error != nil {
//			    return res[0].Error
//			}
//
// model := res[0].Model.(*MyModel) // <4>
//
// // modify model and store it again
// wr := p.Write(ctx,
//
//					WriteOptions{Tx:tx},
//					WriteOperation{
//						ID:res[0].ID,
//						Model: model,
//						ModelType: ModelTypeReported,
//						Version: res[0].Version},
//				)
//
//	if wr[0].Error != nil { // <5>
//		    return wr[0].Error
//	}
//
// ----
// <1> Always release the transaction to ensure it will either commit or rollback the transaction.
// <2> The `Persistence` needs to know into what type to `Unmarshal` the model into.
// <3> The transaction is passed to the `Read` operation and it will always return same amount of results as `PersistenceID`s.
// <4> The model is returned as a `any` and thus needs to be type asserted.
// <5> The `Write` will always return the same amount of results as `WriteOperation`s.
type Transactional interface {
	Persistence
	// Begin will start a new transaction.
	Begin(ctx context.Context) (*Transaction, error)
	// Release will either commit or rollback the transaction.
	//
	// If the transaction is already committed or aborted, it will not return an error so it is possible to do a
	// defer and explicit call to `Release` in the code without any errors are returned.
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

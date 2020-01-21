# How to Soft Delete

Due to our contractual obligations with the federal government, we must be able to access deleted data even several years after it’s been used in the system. For this reason, MilMove is shifting away from hard deleting data and adopting the practice to soft delete instead. Soft delete functionality has not yet been implemented throughout the entire codebase but it is expected to be the sole deletion method moving forward.

Please note that soft delete is to be treated like a hard delete in the regard that the process should never be reversed or that data can be 'un-deleted'.

## How Soft Delete Works

MilMove's implementation of soft delete takes in a model, sets a time stamp to its `DeletedAt` field, before cascading down to its children and repeating the process until there are no longer children to 'delete'.

## Prerequisites for Soft Delete

To use soft delete, a model and its children (or foreign key associations) must possess a `DeletedAt` field that corresponds to the `deleted_at` column of their table within the database.

```go
type ExampleModel struct {
    ...
    DeletedAt   *time.Time  `db:"deleted_at"`

}
```

If this has not been done, one must [create a migration](migrate-the-database.md) to make these changes.

Furthermore, any queries to fetch the model must exclude those that have been 'soft deleted'.

```go
func FetchExampleModel(ctx context.Context, db *pop.Connection, session * auth.Session, id uuid.UUID) (ExampleModel, error) {
    var exampleModel ExampleModel
    err := db.Q().Where("example_models.deleted_at is null").Eager().Find(&exampleModel, id)
    ...
}
```

## Using Soft Delete

In order to use MilMove's soft delete method, one must import the following package

```go
package models

import (
    "github.com/transcom/mymove/pkg/db/utilities"
)
```

It is recommended that any use of soft delete be wrapped in a transaction. This is to rollback the deletion should any error arise.

```go
func DeleteExampleModel(db *pop.Connection, exampleModel *ExampleModel) error {
    return db.Transaction(func(db *pop.Connection)) error {
        return utilities.SoftDestroy(db, exampleModel)
    }
}
```

package migrate

import (
	"database/sql"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func prepareCopyFromStdin(tablename string, columns []string, tx *pop.Connection) (*sql.Stmt, error) {
	// With Schema
	if strings.Contains(tablename, ".") {
		parts := strings.SplitN(tablename, ".", 2)
		preparedStmt := pq.CopyInSchema(parts[0], parts[1], columns...)
		stmt, err := tx.TX.Prepare(preparedStmt)
		if err != nil {
			return nil, errors.Wrap(err, "error preparing copy from stdin statement")
		}
		return stmt, nil
	}
	// Without Schema
	stmt, err := tx.TX.Prepare(pq.CopyIn(tablename, columns...))
	if err != nil {
		return nil, errors.Wrap(err, "error preparing copy from stdin statement")
	}
	return stmt, nil
}

package migrate

import (
	"database/sql"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func prepareCopyFromStdin(tablename string, columns []string, tx *pop.Connection) (*sql.Stmt, error) {
	//fmt.Fprintln(os.Stderr, tablename, columns)
	if strings.Contains(tablename, ".") {
		parts := strings.SplitN(tablename, ".", 2)
		stmt, err := tx.TX.Prepare(pq.CopyInSchema(parts[0], parts[1], columns...))
		if err != nil {
			return nil, errors.Wrap(err, "error preparing copy from stdin statement")
		}
		return stmt, nil
	}
	stmt, err := tx.TX.Prepare(pq.CopyIn(tablename, columns...))
	if err != nil {
		return nil, errors.Wrap(err, "error preparing copy from stdin statement")
	}
	return stmt, nil
}

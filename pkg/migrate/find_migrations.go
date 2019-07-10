package migrate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func FindMigrations(fm *pop.FileMigrator, valid map[string]struct{}, runner func(mf pop.Migration, tx *pop.Connection) error) error {

	dir := fm.Path

	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		// directory doesn't exist
		return nil
	}

	return filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			match, err := pop.ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}

			if match == nil {
				return nil
			}

			if _, ok := valid[filepath.Base(p)]; !ok {
				return errors.New(fmt.Sprintf("migration at path %q missing from manifest", p))
			}

			mf := pop.Migration{
				Path:      p,
				Version:   match.Version,
				Name:      match.Name,
				DBType:    match.DBType,
				Direction: match.Direction,
				Type:      match.Type,
				Runner:    runner,
			}
			fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		}
		return nil
	})
}

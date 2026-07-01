package migrations

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Tx(fn func(tx *sqlx.Tx) error) (err error) {
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				fmt.Printf("failed to rollback tx: %s\n", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			fmt.Printf("failed to commit tx: %s\n", commitErr)
			err = commitErr
		}
	}()
	err = fn(tx)
	return err
}

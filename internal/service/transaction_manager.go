package service

import (
	"database/sql"
	"github.com/rs/zerolog/log"
)

type TransactionManager struct {
	db *sql.DB
}

func (tm *TransactionManager) Begin() (tx *sql.Tx, err error) {
	log.Debug().Msg("Starting a new transaction")
	tx, err = tm.db.Begin()
	return tx, err
}

func (tm *TransactionManager) WithTransaction(do func(tx *sql.Tx) error) (err error) {
	tx, err := tm.Begin()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start a transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Error().Msg("Recovered from panic, rolling back transaction")
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Error().Err(err).Msg("An error occurred, rolling back transaction")
			tx.Rollback()
		} else {
			log.Debug().Msg("Committing transaction")
			err = tx.Commit()
		}
	}()

	err = do(tx)
	return err
}

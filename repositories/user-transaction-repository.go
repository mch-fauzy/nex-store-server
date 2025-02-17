package repositories

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/rs/zerolog/log"
)

func (r *Repository) UserTransactionCreate(data *models.UserTransaction) error {
	createdUserTransaction := r.PostgreSqlConn.Db.Create(data)
	err := createdUserTransaction.Error
	if err != nil {
		log.Error().Err(err).Msg("[UserTransactionCreate] Repository error creating user transaction")
		return err
	}

	return nil
}

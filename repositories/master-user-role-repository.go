package repositories

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/rs/zerolog/log"
)

func (r *Repository) MasterUserRoleFindById(primaryId models.MasterUserRolePrimaryId) (models.MasterUserRole, error) {
	var role models.MasterUserRole
	roleData := r.PostgreSqlConn.Db.First(&role, primaryId)
	err := roleData.Error
	if err != nil {
		log.Error().Err(err).Msg("[MasterUserRoleFindById] Repository error retrieving role by id")
		return models.MasterUserRole{}, err
	}

	return role, nil
}

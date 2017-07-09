package modext

import (
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

func ValidAPIUser(db boil.Executor, token string) bool {
	exists, err := models.Users(db, qm.Where("api_access=? AND api_token=?", true, token)).Exists()
	return err == nil && exists
}

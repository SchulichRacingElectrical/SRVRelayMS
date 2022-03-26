package run

import (
	"context"
	"database-ms/app/models"
)

type RunServiceI interface {
	Create(context.Context, *models.Run) error
}

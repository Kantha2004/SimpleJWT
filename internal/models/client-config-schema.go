package models

import (
	"gorm.io/datatypes"
)

// ClientConfig table
type ClientConfig struct {
	TableModel
	Settings datatypes.JSON `json:"settings" gorm:"type:jsonb"`
}

package client

import "gorm.io/gorm"

type FeatureFlag struct {
	gorm.Model
	Key     string `gorm:"uniqueIndex, primaryKey"`
	Enabled bool
}

package client

import "gorm.io/gorm"

type DB struct {
	db gorm.DB
}

func NewRDS(db gorm.DB) *DB {
	return &DB{
		db: db,
	}
}

// Create a new Flipper with a status
func (r DB) Create(key string, enabled bool) error {
	value := FeatureFlag{
		Key:     key,
		Enabled: enabled,
	}

	err := r.db.Create(&value).Error
	if err != nil {
		return err
	}

	return nil
}

// Read returns the current status of a Flipper
func (r DB) Read(key string) (bool, error) {
	var value FeatureFlag
	err := r.db.First(&value, "key=?", key).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return value.Enabled, nil
}

// Update or create a Flipper with a status
func (r DB) Update(key string, enabled bool) error {
	var value FeatureFlag
	err := r.db.First(&value, "key=?", key).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	value.Enabled = enabled

	return r.db.Save(&value).Error
}

// Delete a Flipper
func (r DB) Delete(key string) error {
	return r.db.Delete(&FeatureFlag{}, "key=?", key).Error
}

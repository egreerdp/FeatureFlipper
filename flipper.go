package featureflipper

//go:generate mockery --name=Storer --output=mocks

import (
	"errors"
	"fmt"
)

type storer interface {
	// Create a new Flipper with a status
	Create(key string, enabled bool) error
	// Read returns the current status of a Flipper
	Read(key string) (bool, error)
	// Update or create a Flipper with a status
	Update(key string, enabled bool) error
	// Delete a Flipper
	Delete(key string) error
}

// ExecutionFn is the function that will be run if a Flipper is enabled
type ExecutionFn = func() error

// FeatureFlipper executes a function if it is enabled in the flipper
type FeatureFlipper struct {
	store storer
}

func NewFeatureFlipper(store storer) *FeatureFlipper {
	return &FeatureFlipper{
		store: store,
	}
}

func (f FeatureFlipper) Create(key string, enabled bool) error {
	return f.store.Create(key, enabled)
}

func (f FeatureFlipper) Update(key string, enabled bool) error {
	return f.store.Update(key, enabled)
}

func (f FeatureFlipper) Delete(key string) error {
	return f.store.Delete(key)
}

func (f FeatureFlipper) Enabled(key string) (bool, error) {
	return f.store.Read(key)
}

// Run executes fn if key is enabled in the FlipperStore
func (f FeatureFlipper) Run(key string, fn ExecutionFn) error {
	return f.run(key, fn, true)
}

func (f FeatureFlipper) RunDisabled(key string, fn ExecutionFn) error {
	return f.run(key, fn, false)
}

func (f FeatureFlipper) run(key string, fn ExecutionFn, runWhenEnabled bool) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}

	enabled, err := f.store.Read(key)
	if err != nil {
		return err
	}

	if enabled == runWhenEnabled {
		return fn()
	}

	return fmt.Errorf("flipper disabled for key: %s", key)
}

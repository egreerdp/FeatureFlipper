package featureflipper_test

import (
	"errors"
	"testing"

	featureflipper "github.com/egreerdp/FeatureFlipper"
	flipperClient "github.com/egreerdp/FeatureFlipper/client"

	"github.com/dailypay/daily-go/pkg/ctxlogger"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type FlipperTestSuite struct {
	suite.Suite
	ff *featureflipper.FeatureFlipper
	db *gorm.DB
}

func TestFlipperTestSuite(t *testing.T) {
	suite.Run(t, new(FlipperTestSuite))
}

func (su *FlipperTestSuite) SetupTest() {
	c, err := flipperClient.NewClient(ctxlogger.NewNop())
	su.NoError(err)

	err = c.AutoMigrate(&flipperClient.FeatureFlag{})
	su.NoError(err)

	su.db = c

	s := flipperClient.NewRDS(*c)
	su.ff = featureflipper.NewFeatureFlipper(s)
}

func (su *FlipperTestSuite) TearDownTest() {
	err := su.db.Migrator().DropTable(&flipperClient.FeatureFlag{})
	su.NoError(err)

	err = su.db.AutoMigrate(&flipperClient.FeatureFlag{})
	su.NoError(err)
}

func (su *FlipperTestSuite) TestFeatureFlipper() {
	err := su.ff.Create("my-feature", true)
	su.NoError(err)

	if enabled, err := su.ff.Enabled("my-feature"); !enabled {
		su.NoError(err)
		su.FailNow(err.Error())
	}
}

func (su *FlipperTestSuite) TestFeatureFlipper_Delete() {
	err := su.ff.Create("my-feature", true)
	su.NoError(err)

	if enabled, err := su.ff.Enabled("my-feature"); !enabled {
		su.NoError(err)
		su.FailNow(err.Error())
	}

	err = su.ff.Delete("my-feature")
	su.NoError(err)

	if enabled, err := su.ff.Enabled("my-feature"); !enabled {
		su.NoError(err)
	}
}

func (su *FlipperTestSuite) TestFeatureFlipper_Update() {
	flipperKey := "my-feature"

	err := su.ff.Create(flipperKey, true)
	su.NoError(err)

	if enabled, err := su.ff.Enabled(flipperKey); !enabled {
		su.NoError(err)
		su.FailNow(err.Error())
	}

	err = su.ff.Update(flipperKey, false)
	su.NoError(err)

	if enabled, err := su.ff.Enabled(flipperKey); enabled {
		su.NoError(err)
		su.FailNow("should not be enabled")
	}
}

func (su *FlipperTestSuite) TestFeatureFlipper_RunEnabled() {
	err := su.ff.Create("my-feature", true)
	su.NoError(err)

	err = su.ff.Run("my-feature", func() error { return errors.New("I ran") })
	su.ErrorContains(err, "I ran")
}

func (su *FlipperTestSuite) TestFeatureFlipper_RunEnabled_NotRun() {
	err := su.ff.Create("my-feature", false)
	su.NoError(err)

	err = su.ff.Run("my-feature", func() error { return errors.New("I ran") })
	su.ErrorContains(err, "flipper disabled for key: my-feature")
}

func (su *FlipperTestSuite) TestFeatureFlipper_Enabled_NoErrorWithKeyNotFound() {
	err := su.ff.Run("my-feature", func() error { return errors.New("I ran") })
	su.ErrorContains(err, "flipper disabled for key: my-feature")
}

package assets

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type assetsSuite struct {
	suite.Suite
}

func TestCLISuite(t *testing.T) {
	ts := &assetsSuite{}
	suite.Run(t, ts)
}

const (
	goodPath = "notifications/templates/move_approved_template.html"
	badPath  = "bad/path/file.txt"
)

func (suite *assetsSuite) TestAsset() {
	suite.Run("golden path", func() {
		contents, err := Asset(goodPath)
		suite.NoError(err)
		suite.NotNil(contents)
		suite.True(len(contents) > 0, "Contents should have non-zero length")
	})

	suite.Run("asset missing", func() {
		contents, err := Asset(badPath)
		suite.Error(err)
		suite.Nil(contents)
	})
}

func (suite *assetsSuite) TestMustAsset() {
	suite.Run("golden path", func() {
		contents := MustAsset(goodPath)
		suite.NotNil(contents)
		suite.True(len(contents) > 0, "Contents should have non-zero length")
	})

	suite.Run("asset missing", func() {
		panicFunc := func() {
			MustAsset(badPath)
		}
		suite.Panics(panicFunc)
	})
}

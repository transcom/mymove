package cli

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckTPPSFlagsValidInput(t *testing.T) {
	v := viper.New()
	v.Set(ProcessTPPSCustomDateFile, "MILMOVE-en20250210.csv")
	v.Set(TPPSS3Bucket, "test-bucket")
	v.Set(TPPSS3Folder, "test-folder")

	err := CheckTPPSFlags(v)
	assert.NoError(t, err)
}

func TestCheckTPPSFlagsMissingProcessTPPSCustomDateFile(t *testing.T) {
	v := viper.New()
	v.Set(TPPSS3Bucket, "test-bucket")
	v.Set(TPPSS3Folder, "test-folder")

	err := CheckTPPSFlags(v)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ProcessTPPSCustomDateFile")
}

func TestCheckTPPSFlagsMissingTPPSS3Bucket(t *testing.T) {
	v := viper.New()
	v.Set(ProcessTPPSCustomDateFile, "MILMOVE-en20250210.csv")
	v.Set(TPPSS3Folder, "test-folder")

	err := CheckTPPSFlags(v)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no value for TPPSS3Bucket found")
}

func TestCheckTPPSFlagsMissingTPPSS3Folder(t *testing.T) {
	v := viper.New()
	v.Set(ProcessTPPSCustomDateFile, "MILMOVE-en20250210.csv")
	v.Set(TPPSS3Bucket, "test-bucket")

	err := CheckTPPSFlags(v)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no value for TPPSS3Folder found")
}

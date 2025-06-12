package cli

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitTPPSFlags(t *testing.T) {
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitTPPSFlags(flagSet)

	processTPPSCustomDateFile, _ := flagSet.GetString(ProcessTPPSCustomDateFile)
	assert.Equal(t, "", processTPPSCustomDateFile, "Expected ProcessTPPSCustomDateFile to have an empty default value")

	tppsS3Bucket, _ := flagSet.GetString(TPPSS3Bucket)
	assert.Equal(t, "", tppsS3Bucket, "Expected TPPSS3Bucket to have an empty default value")

	tppsS3Folder, _ := flagSet.GetString(TPPSS3Folder)
	assert.Equal(t, "", tppsS3Folder, "Expected TPPSS3Folder to have an empty default value")
}

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

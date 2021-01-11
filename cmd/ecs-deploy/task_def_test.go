package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestCreateAwsConfig(t *testing.T) {
	// This isn't a great test BUT should be a good starting point for others
	awsRegionString := "us-gov-west-1"
	awsConfig := createAwsConfig(awsRegionString)
	if *awsConfig.Region != *aws.String(awsRegionString) {
		t.Errorf("Expected AWS Config region to equal %v, instead is %v", *awsConfig.Region, *aws.String(awsRegionString))
	}
}

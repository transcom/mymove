package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestCreateAwsConfig(t *testing.T) {
	// This isn't a great test BUT should be a good starting point for others
	awsRegionString := "us-gov-west-1"
	awsConfig := createAwsConfig(awsRegionString)
	if *awsConfig.Region != *aws.String(awsRegionString) {
		t.Errorf("Expected AWS Config region to equal %v, instead is %v", *awsConfig.Region, *aws.String(awsRegionString))
	}
}

func TestRemoveSecretsWithMatchingEnvironmentVariables(t *testing.T) {
	cases := map[string]struct {
		inSecrets  []*ecs.Secret
		inEnvVars  []*ecs.KeyValuePair
		expSecrets []*ecs.Secret
	}{
		"no secrets, no env vars": {
			inSecrets:  []*ecs.Secret{},
			inEnvVars:  []*ecs.KeyValuePair{},
			expSecrets: []*ecs.Secret{},
		},
		"one secret, no env vars": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []*ecs.KeyValuePair{},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"no secrets, one env var": {
			inSecrets: []*ecs.Secret{},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []*ecs.Secret{},
		},
		"one secret, one env var, not matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"one secret, one env var, matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []*ecs.Secret{},
		},
		"two secrets, one env var, none matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting")},
			},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
		},
		"two secrets, one env var, one matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 2")},
			},
		},
		"one secret, two env vars, none matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"one secret, two env vars, one matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []*ecs.Secret{},
		},
		"two secrets, two env vars, both matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []*ecs.Secret{},
		},
		"two secrets, three env vars, two matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			expSecrets: []*ecs.Secret{},
		},
		"three secrets, two env vars, two matching": {
			inSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			inEnvVars: []*ecs.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []*ecs.Secret{
				{Name: aws.String("my setting 3")},
			},
		},
	}

	for name, tc := range cases {
		actual := removeSecretsWithMatchingEnvironmentVariables(tc.inSecrets, tc.inEnvVars)
		if !reflect.DeepEqual(actual, tc.expSecrets) {
			t.Errorf("%v: expected %v, but got %v", name, tc.expSecrets, actual)
		}
	}
}

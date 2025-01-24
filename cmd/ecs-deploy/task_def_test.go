package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func TestRemoveSecretsWithMatchingEnvironmentVariables(t *testing.T) {
	cases := map[string]struct {
		inSecrets  []ecstypes.Secret
		inEnvVars  []ecstypes.KeyValuePair
		expSecrets []ecstypes.Secret
	}{
		"no secrets, no env vars": {
			inSecrets:  []ecstypes.Secret{},
			inEnvVars:  []ecstypes.KeyValuePair{},
			expSecrets: []ecstypes.Secret{},
		},
		"one secret, no env vars": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []ecstypes.KeyValuePair{},
			expSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"no secrets, one env var": {
			inSecrets: []ecstypes.Secret{},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []ecstypes.Secret{},
		},
		"one secret, one env var, not matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"one secret, one env var, matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []ecstypes.Secret{},
		},
		"two secrets, one env var, none matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting")},
			},
			expSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
		},
		"two secrets, one env var, one matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
			},
			expSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 2")},
			},
		},
		"one secret, two env vars, none matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			expSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
		},
		"one secret, two env vars, one matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []ecstypes.Secret{},
		},
		"two secrets, two env vars, both matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []ecstypes.Secret{},
		},
		"two secrets, three env vars, two matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			expSecrets: []ecstypes.Secret{},
		},
		"three secrets, two env vars, two matching": {
			inSecrets: []ecstypes.Secret{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
				{Name: aws.String("my setting 3")},
			},
			inEnvVars: []ecstypes.KeyValuePair{
				{Name: aws.String("my setting 1")},
				{Name: aws.String("my setting 2")},
			},
			expSecrets: []ecstypes.Secret{
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

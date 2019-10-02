package iampostgres

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
)

// RDSUtilService Lightweight Interface for AWS RDS Utils
type RDSUtilService interface {
	GetToken(string, string, string, *credentials.Credentials) (string, error)
}

// RDSU type
type RDSU struct{}

// GetToken is implementation around AWS RDS Utils BuildAuthToken
func (r RDSU) GetToken(endpoint string, region string, user string, iamcreds *credentials.Credentials) (string, error) {
	authToken, err := rdsutils.BuildAuthToken(endpoint, region, user, iamcreds)
	if err != nil {
		return "", errors.New("Failed to create RDSIAM token")
	}

	return authToken, nil

}

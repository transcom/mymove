package iampostgres

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
)

// RDSUtilService Lightweight Interface for AWS RDS Utils
type RDSUtilService interface {
	GetToken(context.Context, string, string, string, aws.CredentialsProvider) (string, error)
}

// RDSU type
type RDSU struct{}

// GetToken is implementation around AWS RDS Utils BuildAuthToken
func (r RDSU) GetToken(ctx context.Context, endpoint string, region string, user string, iamcreds aws.CredentialsProvider) (string, error) {
	authToken, err := auth.BuildAuthToken(ctx, endpoint, region, user, iamcreds)
	if err != nil {
		return "", fmt.Errorf("failed to create RDSIAM token: %w", err)
	}

	return authToken, nil

}

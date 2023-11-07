package trdm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type GatewayService struct {
	httpClient  HTTPClient
	logger      *zap.Logger
	region      string
	trdmIamRole string
	gatewayURL  string
	creds       *aws.Credentials
}

func NewGatewayService(httpClient HTTPClient, logger *zap.Logger, region, trdmIamRole string, gatewayURL string, creds *aws.Credentials) *GatewayService {
	return &GatewayService{
		httpClient:  httpClient,
		logger:      logger,
		region:      region,
		trdmIamRole: trdmIamRole,
		gatewayURL:  gatewayURL,
		creds:       creds,
	}
}

const (
	// ! Not in use yet
	// TODO:
	// TrdmIamFlag is the TRDM IAM flag
	TrdmIamFlag string = "trdm-iam"
)

func (gs GatewayService) gatewayLastTableUpdate(request models.LastTableUpdateRequest) (*http.Response, error) {
	// Create the request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		gs.logger.Error("marshalling LastTableUpdate request body", zap.Error(err))
		return nil, err
	}

	// Generate a SHA256 hash for signing
	hash := GenerateSHA256Hash(requestBody)

	// Put it into a new request
	req, err := http.NewRequest("POST", gs.gatewayURL, bytes.NewBuffer(requestBody))
	if err != nil {
		gs.logger.Error("lastTableUpdate request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Sign request, this will update req in place
	err = signRequest(req, *gs.creds, string(hash), gs.region, gs.logger)
	if err != nil {
		gs.logger.Error("signing lastTableUpdate request", zap.Error(err))
		return nil, err
	}

	// Send the request
	resp, err := gs.httpClient.Do(req)
	if err != nil {
		gs.logger.Error("error sending request to API Gateway", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func GenerateSHA256Hash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func signRequest(req *http.Request, creds aws.Credentials, hash string, region string, logger *zap.Logger) error {
	signer := v4.NewSigner()

	// Provide execute-api service as we're going through the gateway for this request
	err := signer.SignHTTP(context.Background(), creds, req, hash, "execute-api", region, time.Now())
	if err != nil {
		logger.Error("error signing http request", zap.Error(err))
		return err
	}

	return nil
}

func retrieveCredentials(region string, trdmIamRole string, logger *zap.Logger) (aws.Credentials, error) {
	// We want to get the credentials from the logged in AWS
	// session rather than create directly, because the session
	// conflates the environment, shared, and container metdata
	// config within NewSession. With stscreds, we use the Secure
	// Token Service, to assume the given role (that has API
	// gateway `execute-api` permissions).
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		logger.Error("loading aws config", zap.Error(err))
		return aws.Credentials{}, err
	}

	logger.Info("assuming AWS role for API gateway execution", zap.String("role", trdmIamRole))
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, trdmIamRole)

	creds, err := provider.Retrieve(context.Background())
	if err != nil {
		logger.Error("error retrieving aws credentials", zap.Error(err))
		return aws.Credentials{}, err
	}

	return creds, nil
}

package trdm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const (
	LastTableUpdateEndpoint      string = "/api/v1/lastTableUpdate"
	GetTableEndpoint             string = "/api/v1/getTable"
	LineOfAccounting             string = "LN_OF_ACCT"
	TransportationAccountingCode string = "TRNSPRTN_ACNT"
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
	stsCreds    AssumeRoleProvider
}

func NewGatewayService(httpClient HTTPClient, logger *zap.Logger, region, trdmIamRole string, gatewayURL string, stsCreds AssumeRoleProvider) *GatewayService {
	return &GatewayService{
		httpClient:  httpClient,
		logger:      logger,
		region:      region,
		trdmIamRole: trdmIamRole,
		gatewayURL:  gatewayURL,
		stsCreds:    stsCreds,
	}
}

func (gs GatewayService) gatewayLastTableUpdate(request models.LastTableUpdateRequest) (*http.Response, error) {
	// Create the request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		gs.logger.Error("marshalling LastTableUpdate request body", zap.Error(err))
		return nil, err
	}

	// Generate a SHA256 hash for signing
	hash := generateSHA256Hash(requestBody)

	// Put it into a new request
	req, err := http.NewRequest("POST", gs.gatewayURL+LastTableUpdateEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		gs.logger.Error("lastTableUpdate request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Sign request, this will update req in place
	err = signRequest(req, gs.stsCreds, string(hash), gs.region, gs.logger)
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

func (gs GatewayService) gatewayGetTable(request models.GetTableRequest) (*http.Response, error) {
	// Create the request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		gs.logger.Error("marshalling GetTable request body", zap.Error(err))
		return nil, err
	}

	// Generate a SHA256 hash for signing
	hash := generateSHA256Hash(requestBody)

	// Put it into a new request
	req, err := http.NewRequest("POST", gs.gatewayURL+GetTableEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		gs.logger.Error("getTable request", zap.Error(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Sign request, this will update req in place
	err = signRequest(req, gs.stsCreds, string(hash), gs.region, gs.logger)
	if err != nil {
		gs.logger.Error("signing getTable request", zap.Error(err))
		return nil, err
	}

	// Send the request
	resp, err := gs.httpClient.Do(req)
	if err != nil {
		gs.logger.Error("error sending request to API Gateway", zap.Error(err))
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return resp, nil
}

func generateSHA256Hash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func signRequest(req *http.Request, stsCreds AssumeRoleProvider, hash string, region string, logger *zap.Logger) error {
	// V4 signing is used for request auth (AKA using IAM auth from Go as a client)
	signer := v4.NewSigner()

	// Generate temporary credentials
	creds, err := stsCreds.Retrieve(context.Background())
	if err != nil {
		logger.Error("error retrieving sts credentials", zap.Error(err))
		return err
	}
	// Provide execute-api service as we're going through the gateway for this request
	err = signer.SignHTTP(context.Background(), creds, req, hash, "execute-api", region, time.Now())
	if err != nil {
		logger.Error("error signing http request", zap.Error(err))
		return err
	}

	return nil
}

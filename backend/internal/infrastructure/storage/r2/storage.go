package r2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"backend/internal/usecase"
)

type storageService struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	publicURL string
}

// New returns a usecase.StorageService backed by Cloudflare R2.
// R2 is S3-compatible; the custom endpoint encodes the account ID.
func New(accountID, accessKey, secretKey, bucket, publicBaseURL string) (usecase.StorageService, error) {
	return NewWithHTTPClient(accountID, accessKey, secretKey, bucket, publicBaseURL, nil)
}

// NewWithHTTPClient is like New but accepts a custom *http.Client.
// Pass nil to use the default AWS SDK HTTP client.
// This constructor exists to allow tests to inject a mock HTTP transport.
func NewWithHTTPClient(accountID, accessKey, secretKey, bucket, publicBaseURL string, httpClient *http.Client) (usecase.StorageService, error) {
	if accountID == "" || accessKey == "" || secretKey == "" || bucket == "" || publicBaseURL == "" {
		return nil, errors.New("r2: accountID, accessKey, secretKey, bucket, and publicBaseURL are all required")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	awsOpts := []func(*config.LoadOptions) error{
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	}
	if httpClient != nil {
		awsOpts = append(awsOpts, config.WithHTTPClient(httpClient))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), awsOpts...)
	if err != nil {
		return nil, fmt.Errorf("r2: load config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &storageService{
		client:    client,
		presigner: s3.NewPresignClient(client),
		bucket:    bucket,
		publicURL: publicBaseURL,
	}, nil
}

func (s *storageService) PresignUpload(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	req, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("r2: presign put: %w", err)
	}
	return req.URL, nil
}

func (s *storageService) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2: delete object: %w", err)
	}
	return nil
}

func (s *storageService) PublicURL(key string) string {
	return s.publicURL + "/" + url.PathEscape(key)
}

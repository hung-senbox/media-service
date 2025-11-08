package uploader

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Provider struct {
	accessKey            string
	secretKey            string
	bucketName           string
	region               string
	domain               string
	cloudFrontKeyGroupID string
	cloudFrontKeyPath    string
	config               aws.Config
}

func NewS3Provider(accessKey, secretKey, bucketName, region, domain, cloudFrontKeyGroupID, cloudFrontKeyPath string) *s3Provider {
	provider := &s3Provider{
		accessKey:            accessKey,
		secretKey:            secretKey,
		bucketName:           bucketName,
		region:               region,
		domain:               domain,
		cloudFrontKeyGroupID: cloudFrontKeyGroupID,
		cloudFrontKeyPath:    cloudFrontKeyPath,
	}

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		accessKey,
		secretKey,
		"",
	))

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		panic(err)
	}

	provider.config = cfg

	return provider
}

func (p *s3Provider) SaveFileUploaded(ctx context.Context, data []byte, key string, mode UploadMode) (*string, error) {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)

	client := s3.NewFromConfig(p.config)

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(key),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
		ACL:         types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3 %w", err)
	}

	switch mode {
	case UploadPrivate:
		return p.GetFileUploaded(ctx, key, nil)
	case UploadPublic:
		duration := 100 * 365 * 24 * time.Hour
		return p.GetFileUploaded(ctx, key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}

func (p *s3Provider) SaveFileUploadedReader(ctx context.Context, r io.Reader, key string, contentType string, mode UploadMode) (*string, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	client := s3.NewFromConfig(p.config)

	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload stream to S3 %w", err)
	}

	switch mode {
	case UploadPrivate:
		return p.GetFileUploaded(ctx, key, nil)
	case UploadPublic:
		duration := 100 * 365 * 24 * time.Hour
		return p.GetFileUploaded(ctx, key, &duration)
	default:
		return nil, errors.New("invalid upload mode")
	}
}

func (p *s3Provider) GetFileUploaded(ctx context.Context, key string, duration *time.Duration) (*string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current directory")
	}

	privKeyBytes, err := os.ReadFile(dir + p.cloudFrontKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil || (block.Type != "PRIVATE KEY" && block.Type != "RSA PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	var privKey *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#1 private key: %w", err)
		}
	case "PRIVATE KEY":
		var parsedKey any
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
		}
		var ok bool
		privKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an RSA private key")
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	signer := sign.NewURLSigner(p.cloudFrontKeyGroupID, privKey)
	url := fmt.Sprintf("%s/%s", p.domain, key)

	if duration == nil {
		duration = aws.Duration(24 * time.Hour)
	}
	signedURL, err := signer.Sign(url, time.Now().Add(*duration))
	if err != nil {
		return nil, fmt.Errorf("failed to sign URL: %w", err)
	}

	return &signedURL, nil
}

func (p *s3Provider) DeleteFileUploaded(ctx context.Context, key string) error {
	client := s3.NewFromConfig(p.config)

	_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	waiter := s3.NewObjectNotExistsWaiter(client)
	err = waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(key),
	}, 5*time.Second)
	if err != nil {
		return fmt.Errorf("error waiting for object deletion: %w", err)
	}

	return nil
}

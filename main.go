// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics_test.go

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Announce struct {
	AnnounceId  string  `json:"announceId" dynamodbav:"announce_id"`
	PublishedAt string  `json:"publishedAt" dynamodbav:"published_at"`
	Title       string  `json:"title" dynamodbav:"title"`
	Description *string `json:"description" dynamodbav:"description"`
	PageUrl     *string `json:"pageUrl" dynamodbav:"page_url"`
}

type AnnounceRepository struct {
	ddbClient *dynamodb.Client
	tableName string
	logger    *zap.Logger
}

func NewAnnounceRepository(client *dynamodb.Client, tableName string, logger *zap.Logger) *AnnounceRepository {
	return &AnnounceRepository{
		ddbClient: client,
		tableName: tableName,
		logger:    logger,
	}
}

func (r *AnnounceRepository) ScanAnnounces(ctx context.Context, limit int) ([]*Announce, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}
	output, err := r.ddbClient.Scan(ctx, input)
	if err != nil {
		return nil, "", fmt.Errorf("failed to scan announces: %v", err)
	}
	announces := make([]*Announce, 0)
	if output.Count == 0 {
		return announces, "", nil
	}
	if err := attributevalue.UnmarshalListOfMaps(output.Items, &announces); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal items: %v", err)
	}
	return announces, "", nil
}

func (r *AnnounceRepository) PutAnnounce(ctx context.Context, announce *Announce) error {
	item, err := attributevalue.MarshalMap(announce)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %v", err)
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}
	output, err := r.ddbClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}
	r.logger.Info("PutAnnounce", zap.Any("output", output))
	return nil
}

type AnnounceService struct {
	logger *zap.Logger
	repo   *AnnounceRepository
}

func NewAnnounceService(logger *zap.Logger, repo *AnnounceRepository) *AnnounceService {
	return &AnnounceService{
		logger: logger,
		repo:   repo,
	}
}

func (s *AnnounceService) Scan(ctx context.Context) {
	scanResult, nextToken, err := s.repo.ScanAnnounces(ctx, 10)
	if err != nil {
		s.logger.Fatal("Failed to scan announces", zap.Error(err))
	}
	s.logger.Info("Scan result", zap.Any("announces", scanResult), zap.String("nextToken", nextToken))
}

func (s *AnnounceService) Put(ctx context.Context) {
	if err := s.repo.PutAnnounce(ctx, &Announce{
		AnnounceId:  NewAnnounceId(),
		PublishedAt: NewPublishedAt(),
		Title:       "example",
		Description: PtrStr("example"),
		PageUrl:     PtrStr("example"),
	}); err != nil {
		s.logger.Fatal("failed to put announce", zap.Error(err))
	}
}

func (s *AnnounceService) PutNil(ctx context.Context) {
	if err := s.repo.PutAnnounce(ctx, &Announce{
		AnnounceId:  NewAnnounceId(),
		PublishedAt: NewPublishedAt(),
		Title:       "example",
		Description: nil,
		PageUrl:     nil,
	}); err != nil {
		s.logger.Fatal("failed to put announce", zap.Error(err))
	}
}

func main() {
	ctx := context.Background()
	logger, _ := zap.NewDevelopment()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatal("failed to create config", zap.Error(err))
	}
	ddbClient := dynamodb.NewFromConfig(cfg)
	repo := NewAnnounceRepository(ddbClient, "GolangDynamoDBSandbox-AnnounceTable", logger)

	service := NewAnnounceService(logger, repo)

	service.Scan(ctx)
	// service.Put(ctx)
	// service.PutNil(ctx)
}

func NewAnnounceId() string {
	return uuid.New().String()
}

func NewPublishedAt() string {
	return time.Now().Format(time.RFC3339)
}

func PtrStr(s string) *string {
	return &s
}

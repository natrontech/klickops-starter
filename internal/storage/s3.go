// Package storage implements api.BlobStore against any S3-compatible
// endpoint (a klickops Bucket service, AWS S3, SeaweedFS, Garage, …).
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/natrontech/klickops-starter/internal/api"
	"github.com/natrontech/klickops-starter/internal/config"
)

type S3Store struct {
	client *s3.Client
	bucket string
}

func NewS3(ctx context.Context, cfg config.S3) (*S3Store, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build S3 config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(NormalizeEndpoint(cfg.Endpoint))
			o.UsePathStyle = true
		}
	})
	store := &S3Store{client: client, bucket: cfg.Bucket}
	store.ensureBucket(ctx)
	return store, nil
}

// NormalizeEndpoint accepts the host:port form that klickops bucket
// bindings provide and adds the scheme the SDK requires. In-cluster
// endpoints speak plain HTTP; public endpoints should be given with
// an explicit https:// prefix.
func NormalizeEndpoint(endpoint string) string {
	if strings.Contains(endpoint, "://") {
		return endpoint
	}
	return "http://" + endpoint
}

// ensureBucket creates the bucket when it does not exist yet (local dev).
// On klickops the bucket already exists and create may be denied - both
// are fine, so failures only log.
//
// ponytail: hard 5s budget, because this runs before the server listens.
// An unreachable endpoint (blocked egress, wrong host) otherwise burns
// the SDK's full retry schedule and the app never serves a request.
func (s *S3Store) ensureBucket(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &s.bucket}); err == nil {
		return
	}
	if _, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &s.bucket}); err != nil {
		slog.Warn("bucket not reachable yet", "bucket", s.bucket, "error", err)
	}
}

func (s *S3Store) ListFiles(ctx context.Context) ([]api.FileInfo, error) {
	// ponytail: single page (1000 objects); paginate when a real app needs more
	out, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: &s.bucket})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}
	files := make([]api.FileInfo, 0, len(out.Contents))
	for _, obj := range out.Contents {
		files = append(files, api.FileInfo{
			Key:          aws.ToString(obj.Key),
			Size:         aws.ToInt64(obj.Size),
			LastModified: aws.ToTime(obj.LastModified),
		})
	}
	return files, nil
}

func (s *S3Store) PutFile(ctx context.Context, key, contentType string, r io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &key,
		Body:          r,
		ContentType:   &contentType,
		ContentLength: &size,
	})
	if err != nil {
		return fmt.Errorf("failed to put object %q: %w", key, err)
	}
	return nil
}

func (s *S3Store) GetFile(ctx context.Context, key string) (io.ReadCloser, string, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &key})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			return nil, "", api.ErrNotFound
		}
		return nil, "", fmt.Errorf("failed to get object %q: %w", key, err)
	}
	contentType := aws.ToString(out.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return out.Body, contentType, nil
}

func (s *S3Store) DeleteFile(ctx context.Context, key string) error {
	if _, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.bucket, Key: &key}); err != nil {
		return fmt.Errorf("failed to delete object %q: %w", key, err)
	}
	return nil
}

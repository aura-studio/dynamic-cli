package push

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Remote interface {
	Push([]Pair)
}

type S3Remote struct {
	bucket string
	prefix string
}

func NewS3Remote(bucket string) *S3Remote {
	// Accept values like:
	// - "mirroring-lambda"
	// - "mirroring-lambda/prefix/path"
	// - "s3://mirroring-lambda"
	// - "s3://mirroring-lambda/prefix/path"
	// Normalize to bucket + optional prefix for S3 API usage.
	b := bucket
	p := ""
	b = strings.TrimPrefix(b, "s3://")
	if i := strings.IndexByte(b, '/'); i >= 0 {
		p = strings.Trim(b[i+1:], "/")
		b = b[:i]
	}
	return &S3Remote{
		bucket: b,
		prefix: p,
	}
}

func (r *S3Remote) createS3Client() (*s3.Client, error) {
	// 1. 首先使用默认配置创建一个临时客户端来查询 Bucket 所在的 Region
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load default config, %v", err)
	}

	tempClient := s3.NewFromConfig(cfg)
	output, err := tempClient.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{
		Bucket: aws.String(r.bucket),
	})

	if err != nil {
		// 如果 GetBucketLocation 失败，可能是权限问题，回退到默认配置
		log.Printf("warning: failed to get bucket location for %s, using default region: %v", r.bucket, err)
	} else {
		region := string(output.LocationConstraint)
		// AWS API 规定：如果 LocationConstraint 为空，则 Bucket 位于 us-east-1
		if region == "" {
			region = "us-east-1"
		}
		log.Printf("detected bucket %s region: %s", r.bucket, region)
		// 2. 使用检测到的 Region 重新加载配置
		cfg, err = config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
		if err != nil {
			return nil, fmt.Errorf("failed to load config with region %s, %v", region, err)
		}
	}

	return s3.NewFromConfig(cfg), nil
}

func (r *S3Remote) uploadFileToS3(remoteFilePath string, localFilePath string) error {
	client, err := r.createS3Client()
	if err != nil {
		log.Panicf("failed to create s3 client, %v", err)
	}

	// Open local file for use.
	f, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", localFilePath, err)
	}
	defer f.Close()

	// Upload file body to S3.
	key := remoteFilePath
	if r.prefix != "" {
		key = path.Join(r.prefix, key)
	}
	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func (r *S3Remote) batchUploadFilesToS3(pairs []Pair) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(pairs))
	for _, pair := range pairs {
		wg.Add(1)
		go func(pair Pair) {
			defer wg.Done()
			localFilePath := pair.LocalFilePath
			remoteFilePath := pair.RemoteFilePath

			if stat, err := os.Stat(localFilePath); err != nil {
				if os.IsNotExist(err) {
					log.Printf("%s does not exist", localFilePath)
					errChan <- err
				} else {
					log.Printf("failed to stat file, %v", err)
					errChan <- err
					return
				}
			} else if stat.Size() == 0 {
				log.Printf("%s is empty", localFilePath)
				errChan <- err
				return
			} else {
				// Log the full S3 URL for clarity
				fullKey := remoteFilePath
				if r.prefix != "" {
					fullKey = path.Join(r.prefix, remoteFilePath)
				}
				log.Printf("%s found, uploading to s3://%s/%s...", localFilePath, r.bucket, fullKey)
				if err := r.uploadFileToS3(remoteFilePath, localFilePath); err != nil {
					log.Printf("failed to upload file to s3, %v", err)
					errChan <- err
					return
				}
			}
		}(pair)
	}
	wg.Wait()

	if len(errChan) > 0 {
		return fmt.Errorf("%d errors occurred during uploading", len(errChan))
	}

	return nil
}

func (r *S3Remote) Push(tasks []Pair) {
	if err := r.batchUploadFilesToS3(tasks); err != nil {
		log.Panicf("failed to upload files to s3, %v", err)
	}
}

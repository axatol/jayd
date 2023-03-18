package miniodriver

import (
	"context"
	"net/url"
	"time"

	"github.com/axatol/jayd/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
	"github.com/rs/zerolog/log"
)

type Client struct {
	client *minio.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	creds := credentials.NewStaticV4(
		config.StorageAccessKeyID,
		config.StorageSecretAccessKey,
		"",
	)

	client, err := minio.New(config.StorageEndpoint, &minio.Options{Creds: creds, Secure: true})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, config.StorageBucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := client.MakeBucket(ctx, config.StorageBucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &Client{client}, nil
}

func (c *Client) FPut(ctx context.Context, objectName string, filePath string, tags map[string]string) error {
	res, err := c.client.FPutObject(
		ctx,
		config.StorageBucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{UserTags: tags},
	)

	if err != nil {
		return err
	}

	log.Info().
		Str("object_name", objectName).
		Str("file_path", filePath).
		Str("version_id", res.VersionID).
		Msg("put object")

	return nil
}

func (c *Client) GetPresignedURL(ctx context.Context, objectName string) (*url.URL, error) {
	presigned, err := c.client.PresignedGetObject(
		ctx,
		config.StorageBucketName,
		objectName,
		time.Minute*5,
		url.Values{},
	)

	if err != nil {
		return nil, err
	}

	log.Info().
		Str("object_name", objectName).
		Msg("got object")

	return presigned, nil
}

func (c *Client) List(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	objects := c.client.ListObjects(
		ctx,
		config.StorageBucketName,
		minio.ListObjectsOptions{Prefix: prefix},
	)

	results := []minio.ObjectInfo{}
	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}

		results = append(results, object)
	}

	log.Info().
		Str("prefix", prefix).
		Int("count", len(results)).
		Msg("listed objects")

	return results, nil
}

func (c *Client) GetTags(ctx context.Context, objectName string) (map[string]string, error) {
	tags, err := c.client.GetObjectTagging(
		ctx,
		config.StorageBucketName,
		objectName,
		minio.GetObjectTaggingOptions{},
	)

	return tags.ToMap(), err
}

func (c *Client) PutTags(ctx context.Context, objectName string, newTags map[string]string) error {
	objectTags, err := tags.MapToObjectTags(newTags)
	if err != nil {
		return err
	}

	return c.client.PutObjectTagging(
		ctx,
		config.StorageBucketName,
		objectName,
		objectTags,
		minio.PutObjectTaggingOptions{},
	)
}

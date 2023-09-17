package miniodriver

import (
	"context"
	"net/url"
	"time"

	"github.com/axatol/jayd/pkg/config"
	"github.com/axatol/jayd/pkg/config/nr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
	"github.com/rs/zerolog/log"
)

type MinioClient struct {
	client *minio.Client
}

type Tags map[string]string

var (
	c *MinioClient
)

func AssertClient(ctx context.Context) (*MinioClient, error) {
	if c != nil {
		return c, nil
	}

	creds := credentials.NewStaticV4(
		config.StorageAccessKeyID,
		config.StorageSecretAccessKey,
		"",
	)

	opts := minio.Options{Creds: creds, Secure: config.StorageSSLEnabled}
	client, err := minio.New(config.StorageEndpoint, &opts)
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

	c = &MinioClient{client}
	return c, nil
}

func (c *MinioClient) FPutObject(ctx context.Context, objectName string, filePath string, tags Tags) error {
	segment := nr.Segment(ctx, "miniodriver.FPutObject", nr.Attrs{"object_name": objectName})
	defer segment.End()

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
		Msg("put file")

	segment.AddAttribute("version_id", res.VersionID)

	return nil
}

func (c *MinioClient) RemoveObject(ctx context.Context, objectName string) error {
	defer nr.Segment(ctx, "miniodriver.RemoveObject", nr.Attrs{"object_name": objectName}).End()

	err := c.client.RemoveObject(
		ctx,
		config.StorageBucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return err
	}

	log.Info().
		Str("object_name", objectName).
		Msg("removed object")

	return nil
}

func (c *MinioClient) GetPresignedURL(ctx context.Context, objectName string) (*url.URL, error) {
	defer nr.Segment(ctx, "miniodriver.GetPresignedURL", nr.Attrs{"object_name": objectName}).End()

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

func (c *MinioClient) List(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	segment := nr.Segment(ctx, "miniodriver.List", nr.Attrs{"prefix": prefix})
	defer segment.End()

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
		Int("object_count", len(results)).
		Msg("listed objects")

	segment.AddAttribute("object_count", len(results))

	return results, nil
}

func (c *MinioClient) GetTags(ctx context.Context, objectName string) (Tags, error) {
	defer nr.Segment(ctx, "miniodriver.GetTags", nr.Attrs{"object_name": objectName}).End()

	tags, err := c.client.GetObjectTagging(
		ctx,
		config.StorageBucketName,
		objectName,
		minio.GetObjectTaggingOptions{},
	)

	return tags.ToMap(), err
}

func (c *MinioClient) PutTags(ctx context.Context, objectName string, newTags Tags) error {
	defer nr.Segment(ctx, "miniodriver.PutTags", nr.Attrs{"object_name": objectName}).End()

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

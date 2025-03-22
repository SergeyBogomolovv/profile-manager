package repo

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/SergeyBogomolovv/profile-manager/common/e"
	conf "github.com/SergeyBogomolovv/profile-manager/profile/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type imageRepo struct {
	manager *manager.Uploader
	client  *s3.Client
	bucket  string
}

func MustNewImageRepo(conf conf.S3) *imageRepo {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.Access, conf.Secret, "")),
		config.WithRegion(conf.Region),
		config.WithBaseEndpoint(conf.Endpoint),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	manager := manager.NewUploader(client)
	return &imageRepo{manager, client, conf.Bucket}
}

func (u *imageRepo) UploadAvatar(ctx context.Context, userID string, body []byte) (string, error) {
	hash := sha256.Sum256(body)
	sha256Hex := hex.EncodeToString(hash[:])
	result, err := u.manager.Upload(ctx, &s3.PutObjectInput{
		Bucket:         aws.String(u.bucket),
		Key:            aws.String(avatarKey(userID)),
		Body:           bytes.NewReader(body),
		ContentType:    aws.String("image/jpeg"),
		ContentLength:  aws.Int64(int64(len(body))),
		ChecksumSHA256: aws.String(sha256Hex),
	})
	if err != nil {
		return "", e.Wrap(err, "failed to upload avatar")
	}
	return result.Location, nil
}

func (u *imageRepo) DeleteAvatar(ctx context.Context, userID string) error {
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(avatarKey(userID)),
	})
	return e.WrapIfErr(err, "failed to delete avatar")
}

const folder = "avatars"

func avatarKey(userID string) string {
	return fmt.Sprintf("%s/%s.jpg", folder, userID)
}

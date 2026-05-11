package oss

import (
	"errors"

	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"

	"doctor-go/internal/config"
)

type Client struct {
	bucket  *aliyunoss.Bucket
	baseURL string
}

func New(cfg config.OSSConfig) (*Client, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.Bucket == "" {
		return &Client{baseURL: cfg.BaseURL}, nil
	}

	client, err := aliyunoss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, err
	}
	return &Client{bucket: bucket, baseURL: cfg.BaseURL}, nil
}

func (c *Client) IsConfigured() bool {
	return c != nil && c.bucket != nil
}

func (c *Client) SignedPutURL(objectKey string, expireSeconds int64) (string, error) {
	if !c.IsConfigured() {
		return "", errors.New("oss is not configured")
	}
	return c.bucket.SignURL(objectKey, aliyunoss.HTTPPut, expireSeconds)
}

func (c *Client) PublicURL(objectKey string) string {
	if c.baseURL == "" {
		return objectKey
	}
	return c.baseURL + "/" + objectKey
}

package main

import "time"

type Config struct {
	BackendHost              string        `env:"BACKEND_HOST"`
	BackendPort              string        `env:"BACKEND_PORT"`
	BackendDomain            string        `env:"BACKEND_DOMAIN"`
	BackendReadTimeout       time.Duration `env:"BACKEND_READ_TIMEOUT"`
	BackendReadHeaderTimeout time.Duration `env:"BACKEND_READ_HEADER_TIMEOUT"`
	BackendWriteTimeout      time.Duration `env:"BACKEND_WRITE_TIMEOUT"`
	BackendIdleTimeout       time.Duration `env:"BACKEND_IDLE_TIMEOUT"`

	AllowOriginCors string `env:"CORS_ALLOW_ORIGIN"`

	StorageHost     string `env:"STORAGE_HOST"`
	StoragePort     string `env:"STORAGE_PORT"`
	StorageName     string `env:"STORAGE_NAME"`
	StorageUser     string `env:"STORAGE_USER"`
	StoragePassword string `env:"STORAGE_PASSWORD"`
	StorageSSLMode  string `env:"STORAGE_SSL_MODE"`

	ObjectStoragePort            string `env:"OBJECT_STORAGE_PORT"`
	ObjectStorageAccessKeyID     string `env:"OBJECT_STORAGE_ACCESS_KEY_ID"`
	ObjectStorageSecretAccessKey string `env:"OBJECT_STORAGE_SECRET_ACCESS_KEY"`
	ObjectStorageEndpoint        string `env:"OBJECT_STORAGE_ENDPOINT"`
	ObjectStorageBucket          string `env:"OBJECT_STORAGE_BUCKET"`
	ObjectStorageAccessHost      string `env:"OBJECT_STORAGE_ACCESS_HOST"`

	JWKSURI string `env:"JWKS_URI"`

	TokensIssuer string `env:"TOKENS_ISSUER"`
	TokensKeyID  string `env:"TOKENS_KEY_ID"`
}

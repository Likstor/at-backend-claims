package main

type Config struct {
	StorageHost     string `env:"STORAGE_HOST"`
	StoragePort     string `env:"STORAGE_PORT"`
	StorageName     string `env:"STORAGE_NAME"`
	StorageUser     string `env:"STORAGE_USER"`
	StoragePassword string `env:"STORAGE_PASSWORD"`
	StorageSSLMode  string `env:"STORAGE_SSL_MODE"`

	MigrationsTable string `env:"MIGRATIONS_TABLE"`
	MigrationsPath  string `env:"MIGRATIONS_PATH"`
}

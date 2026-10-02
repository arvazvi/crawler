package config

import "os"

type MongoConfig struct {
	URI string
	Database string
	Collection string
}

type Config struct {
	DB *MongoConfig
}

func Resolve() *Config {
	return &Config{
		DB: &MongoConfig{
			URI: os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
			Collection: os.Getenv("MONGO_COLL"),
		},
	}
}
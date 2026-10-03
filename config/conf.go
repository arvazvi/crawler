package config

import "os"

type MongoConfig struct {
	URI string
	Database string
	Collection string
}

type AIConfig struct {
	ApiKey string
}


type Config struct {
	AI *AIConfig
	DB *MongoConfig
}

func Resolve() *Config {
	return &Config{
		AI: &AIConfig{
			ApiKey: os.Getenv("OPEN_AI_KEY"),
		},
		DB: &MongoConfig{
			URI: os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DB"),
			Collection: os.Getenv("MONGO_COLL"),
		},
	}
}
package local

import (
	"os"

	"github.com/rafaelsanzio/go-flashscore/pkg/config/key"
)

var Values = map[key.Key]string{
	key.MongoDBName:     getDefaultOrEnvVar("flashscore", "MONGO_DATABASE"),
	key.MongoDBPassword: getDefaultOrEnvVar("root", "MONGO_PASWORD"),
	key.MongoDBUsername: getDefaultOrEnvVar("root", "MONGO_USERNAME"),
	key.MongoURI:        getDefaultOrEnvVar("mongodb://root:root@mongo:27017/?connect=direct", "MONGO_URI"),
	key.Region:          "us-east-1",
	key.AppPort:         getDefaultOrEnvVar("5000", "APP_PORT"),
	key.AppAuthPort:     getDefaultOrEnvVar("5001", "APP_AUTH_PORT"),
	key.SecretApiKey:    getDefaultOrEnvVar("flashscore", "SECRET_API_KEY"),

	key.RedisPort: getDefaultOrEnvVar("6379", "REDIS_PORT"),

	key.KafkaAddress1: getDefaultOrEnvVar("kafka-1:19092", "KAFKA_ADDRESS_1"),
	key.KafkaAddress2: getDefaultOrEnvVar("kafka-2:29092", "KAFKA_ADDRESS_2"),
	key.KafkaAddress3: getDefaultOrEnvVar("kafka-3:39092", "KAFKA_ADDRESS_3"),
}

// Some of the db fields are set via env var in the makefile, so this optionally uses those to prevent test failures in jenkins
func getDefaultOrEnvVar(dfault, envVar string) string {
	val := os.Getenv(envVar)
	if val != "" {
		return val
	}
	return dfault
}

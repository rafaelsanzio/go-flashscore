package key

type Key struct {
	Name   string
	Secure bool
	Provider
}

type Provider string

var (
	ProviderStore  = Provider("store")
	ProviderEnvVar = Provider("env")
)

var (
	MongoDBName     = Key{Name: "MONGO_DATABASE", Secure: false, Provider: ProviderStore}
	MongoDBPassword = Key{Name: "MONGO_PASSWORD", Secure: true, Provider: ProviderStore}
	MongoDBUsername = Key{Name: "MONGO_USERNAME", Secure: false, Provider: ProviderStore}
	MongoURI        = Key{Name: "MONGO_URI", Secure: false, Provider: ProviderStore}
	Region          = Key{Name: "REGION", Secure: false, Provider: ProviderStore}
	AppPort         = Key{Name: "APP_PORT", Secure: false, Provider: ProviderStore}
	AppAuthPort     = Key{Name: "APP_AUTH_PORT", Secure: false, Provider: ProviderStore}
	SecretApiKey    = Key{Name: "SECRET_API_KEY", Secure: false, Provider: ProviderStore}

	KafkaAddress1 = Key{Name: "KAFKA_ADDRESS_1", Secure: false, Provider: ProviderStore}
	KafkaAddress2 = Key{Name: "KAFKA_ADDRESS_2", Secure: false, Provider: ProviderStore}
	KafkaAddress3 = Key{Name: "KAFKA_ADDRESS_3", Secure: false, Provider: ProviderStore}

/* 	MongoDBName     = Key{Name: "MONGO_DATABASE", Secure: false, Provider: ProviderEnvVar}
MongoDBPassword = Key{Name: "MONGO_PASSWORD", Secure: true, Provider: ProviderEnvVar}
MongoDBUsername = Key{Name: "MONGO_USERNAME", Secure: false, Provider: ProviderEnvVar}
MongoURI        = Key{Name: "MONGO_URI", Secure: false, Provider: ProviderEnvVar}
Region          = Key{Name: "REGION", Secure: false, Provider: ProviderEnvVar}
AppPort         = Key{Name: "APP_PORT", Secure: false, Provider: ProviderEnvVar} */
)

package app

const (
	EnvironmentVariablesFile = ".env"

	ServerHostKey          = "SERVER_HOST"
	ServerHostDefaultValue = ""

	ServerPortKey          = "SERVER_PORT"
	ServerPortDefaultValue = "8080"

	ServerPrivateKeyPathKey      = "SERVER_PRIVATE_KEY_FILE"
	ServerPrivateKeyDefaultValue = ""

	ServerCertificatePathKey      = "SERVER_CERTIFICATE_FILE"
	ServerCertificateDefaultValue = " "

	CertificateDurationKey          = "CERTIFICATE_DURATION"
	CertificateDurationDefaultValue = "8766h"

	PrivateKeyPathKey          = "PRIVATE_KEY_FILE"
	PrivateKeyPathDefaultValue = "./rsa-priv-key.pem"

	CertificatePathKey          = "CERTIFICATE_FILE"
	CertificatePathDefaultValue = "./rsa-certificate.pem"

	MongoDBDatabaseNameKey          = "MONGODB_DATABASE_NAME"
	MongoDBDatabaseNameDefaultValue = "project-jano"

	MongoDBUriKey          = "MONGODB_URI"
	MongoDBUriDefaultValue = "mongodb://${USERNAME}:${PASSWORD}@localhost:27017/?authSource="

	MongoDBUsernameKey          = "MONGODB_USERNAME"
	MongoDBUsernameDefaultValue = "jano"

	MongoDBPasswordKey          = "MONGODB_PASSWORD"
	MongoDBPasswordDefaultValue = "jano"

	MongoDBUsernamePlaceholder = "${USERNAME}"
	MongoDBPasswordPlaceholder = "${PASSWORD}"

	MongoDBMaxPoolSizeKey          = "MONGODB_MAX_POOL_SIZE"
	MongoDBMaxPoolSizeDefaultValue = "50"

	ServerReadTimeoutKey      = "SERVER_READ_TIMEOUT"
	ServerWriteTimeoutKey     = "SERVER_WRITE_TIMEOUT"
	ServerTimeoutDefaultValue = "120"

	EnvironmentKey         = "ENVIRONMENT"
	DevelopmentEnvironment = "DEVELOPMENT"

	APIAuthUsernameKey          = "API_AUTH_USERNAME"
	APIAuthUsernameDefaultValue = ""

	APIAuthPasswordKey          = "API_AUTH_PASSWORD"
	APIAuthPasswordDefaultValue = ""
)

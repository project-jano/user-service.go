package app

import (
	"fmt"
	"github.com/project-jano/user-service.go/go/api"
	"github.com/project-jano/user-service.go/go/helpers"
	"github.com/project-jano/user-service.go/go/security"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	APIConfiguration      api.Configuration
	ServerAddr            string
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	DatabaseConfiguration DatabaseConfiguration

	ServerPrivateKeyFilePath  string
	ServerCertificateFilePath string
}

type DatabaseConfiguration struct {
	Uri          string
	DatabaseName string
	MinPoolSize  uint64
	MaxPoolSize  uint64
}

func LoadAppConfiguration() (*Configuration, error) {

	viper.SetConfigFile(EnvironmentVariablesFile)
	_ = viper.ReadInConfig()

	// Parse Server host and port
	serverHost := helpers.GetEnvVar(ServerHostKey, ServerHostDefaultValue)
	if len(serverHost) > 0 && net.ParseIP(serverHost) == nil {
		return nil, fmt.Errorf("invalid %s value: '%v'", ServerHostKey, serverHost)
	}
	serverPort := helpers.GetEnvVar(ServerPortKey, ServerPortDefaultValue)
	if _, err := strconv.ParseUint(serverPort, 10, 32); err != nil {
		return nil, fmt.Errorf("invalid %s value: '%v'. %v", ServerPortKey, serverPort, err)
	}

	serverPrivateKeyFilePath := helpers.GetEnvVar(ServerPrivateKeyPathKey, ServerPrivateKeyDefaultValue)
	serverCertificateFilePath := helpers.GetEnvVar(ServerCertificatePathKey, ServerCertificateDefaultValue)

	// Parse and load private key and certificates configuration
	privateKeyFilePath := helpers.GetEnvVar(PrivateKeyPathKey, PrivateKeyPathDefaultValue)
	privateKeyPEM := helpers.ReadFile(privateKeyFilePath)

	certificateFilePath := helpers.GetEnvVar(CertificatePathKey, CertificatePathDefaultValue)
	certificatePEM := helpers.ReadFile(certificateFilePath)

	clientCertificateDuration, _ := time.ParseDuration(helpers.GetEnvVar(CertificateDurationKey, CertificateDurationDefaultValue))

	certificate, privateKey, err := security.LoadCertificateAndKey(certificatePEM, privateKeyPEM)
	if err != nil {
		return nil, err
	}

	// Parse and load database configuration
	databaseName := helpers.GetEnvVar(MongoDBDatabaseNameKey, MongoDBDatabaseNameDefaultValue)
	databaseURI := helpers.GetEnvVar(MongoDBUriKey, MongoDBUriDefaultValue+databaseName)
	databaseUsername := helpers.GetEnvVar(MongoDBUsernameKey, MongoDBUsernameDefaultValue)
	databasePassword := helpers.GetEnvVar(MongoDBPasswordKey, MongoDBPasswordDefaultValue)
	databaseURI = strings.Replace(databaseURI, MongoDBUsernamePlaceholder, databaseUsername, 1)
	databaseURI = strings.Replace(databaseURI, MongoDBPasswordPlaceholder, databasePassword, 1)

	databaseMaxPoolSize, poolSizeError := strconv.Atoi(helpers.GetEnvVar(MongoDBMaxPoolSizeKey, MongoDBMaxPoolSizeDefaultValue))
	if poolSizeError != nil || databaseMaxPoolSize < 1 {
		return nil, fmt.Errorf("invalid %s value: '%v'. %v", MongoDBMaxPoolSizeKey, databaseMaxPoolSize, poolSizeError)
	}

	// Parse Server configuration
	serverReadTimeout, serverReadTimeoutError := strconv.Atoi(helpers.GetEnvVar(ServerReadTimeoutKey, ServerTimeoutDefaultValue))
	if serverReadTimeoutError != nil || serverReadTimeout < 1 {
		return nil, fmt.Errorf("invalid %s value: '%v'. %v", ServerReadTimeoutKey, databaseMaxPoolSize, poolSizeError)
	}

	serverWriteTimeout, serverWriteTimeoutError := strconv.Atoi(helpers.GetEnvVar(ServerWriteTimeoutKey, ServerTimeoutDefaultValue))
	if serverWriteTimeoutError != nil || serverWriteTimeout < 1 {
		return nil, fmt.Errorf("invalid %s value: '%v'. %v", ServerWriteTimeoutKey, databaseMaxPoolSize, poolSizeError)
	}

	// Parse API Basic Auth params
	apiAuthUsername := helpers.GetEnvVar(APIAuthUsernameKey, APIAuthUsernameDefaultValue)
	apiAuthPassword := helpers.GetEnvVar(APIAuthPasswordKey, APIAuthPasswordDefaultValue)

	// Return application configuration
	return &Configuration{
		APIConfiguration: api.Configuration{
			AuthUsername:              apiAuthUsername,
			AuthPassword:              apiAuthPassword,
			AuthEnabled:               len(apiAuthUsername) > 0 && len(apiAuthPassword) > 0,
			TraceCallsEnabled:         helpers.GetEnvVar(EnvironmentKey, DevelopmentEnvironment) == DevelopmentEnvironment,
			CertificatePEM:            certificatePEM,
			Certificate:               certificate,
			PrivateKey:                privateKey,
			ClientCertificateDuration: clientCertificateDuration,
			DatabaseName:              databaseName,
		},
		ServerAddr:                serverHost + ":" + serverPort,
		ServerReadTimeout:         time.Duration(serverReadTimeout) * time.Second,
		ServerWriteTimeout:        time.Duration(serverWriteTimeout) * time.Second,
		ServerPrivateKeyFilePath:  serverPrivateKeyFilePath,
		ServerCertificateFilePath: serverCertificateFilePath,
		DatabaseConfiguration: DatabaseConfiguration{
			DatabaseName: databaseName,
			Uri:          databaseURI,
			MinPoolSize:  uint64(1),
			MaxPoolSize:  uint64(databaseMaxPoolSize),
		},
	}, nil
}

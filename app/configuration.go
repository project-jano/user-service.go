package app

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/project-jano/user-service.go/helpers"
	"github.com/project-jano/user-service.go/security"
	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
	API      APIConfiguration
}

type ServerConfiguration struct {
	Addr                string
	ReadTimeout         time.Duration
	ServerWriteTimeout  time.Duration
	PrivateKeyFilePath  string
	CertificateFilePath string
}

type DatabaseConfiguration struct {
	Uri          string
	DatabaseName string
	MinPoolSize  uint64
	MaxPoolSize  uint64
}

const errorFormat = "invalid %s value: '%v'. %v"

func LoadAppConfiguration() (*Configuration, error) {

	viper.SetConfigFile(EnvironmentVariablesFile)
	_ = viper.ReadInConfig()

	// Parse Server host and port
	serverHost := helpers.GetEnvVar(ServerHostKey, ServerHostDefaultValue)
	if len(serverHost) > 0 && net.ParseIP(serverHost) == nil {
		return nil, fmt.Errorf(errorFormat, ServerHostKey, serverHost, nil)
	}
	serverPort := helpers.GetEnvVar(ServerPortKey, ServerPortDefaultValue)
	if _, err := strconv.ParseUint(serverPort, 10, 32); err != nil {
		return nil, fmt.Errorf(errorFormat, ServerPortKey, serverPort, err)
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
	databaseName, databaseConfiguration, dbErr := getDatabaseConfiguration()
	if dbErr != nil {
		return nil, dbErr
	}

	// Parse Server configuration
	serverReadTimeout, serverWriteTimeout, serverConfigErr := getServerConfiguration()
	if serverConfigErr != nil {
		return nil, serverConfigErr
	}

	// Parse API Basic Auth params
	apiAuthUsername, apiAuthPassword := getAPIAuthenticationParams()

	// Return application configuration
	return &Configuration{
		API: APIConfiguration{
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
		Server: ServerConfiguration{
			Addr:                serverHost + ":" + serverPort,
			ReadTimeout:         time.Duration(serverReadTimeout) * time.Second,
			ServerWriteTimeout:  time.Duration(serverWriteTimeout) * time.Second,
			PrivateKeyFilePath:  serverPrivateKeyFilePath,
			CertificateFilePath: serverCertificateFilePath,
		},
		Database: *databaseConfiguration,
	}, nil
}

func getDatabaseConfiguration() (string, *DatabaseConfiguration, error) {
	databaseName := helpers.GetEnvVar(MongoDBDatabaseNameKey, MongoDBDatabaseNameDefaultValue)
	databaseURI := helpers.GetEnvVar(MongoDBUriKey, MongoDBUriDefaultValue+databaseName)
	databaseUsername := helpers.GetEnvVar(MongoDBUsernameKey, MongoDBUsernameDefaultValue)
	databasePassword := helpers.GetEnvVar(MongoDBPasswordKey, MongoDBPasswordDefaultValue)
	databaseURI = strings.Replace(databaseURI, MongoDBUsernamePlaceholder, databaseUsername, 1)
	databaseURI = strings.Replace(databaseURI, MongoDBPasswordPlaceholder, databasePassword, 1)

	databaseMaxPoolSize, poolSizeError := strconv.Atoi(helpers.GetEnvVar(MongoDBMaxPoolSizeKey, MongoDBMaxPoolSizeDefaultValue))
	if poolSizeError != nil || databaseMaxPoolSize < 1 {
		return "", nil, poolSizeError
	}

	databaseConfiguration := &DatabaseConfiguration{
		DatabaseName: databaseName,
		Uri:          databaseURI,
		MinPoolSize:  uint64(1),
		MaxPoolSize:  uint64(databaseMaxPoolSize),
	}

	return databaseName, databaseConfiguration, nil
}

func getServerConfiguration() (int, int, error) {
	serverReadTimeout, serverReadTimeoutError := strconv.Atoi(helpers.GetEnvVar(ServerReadTimeoutKey, ServerTimeoutDefaultValue))
	if serverReadTimeoutError != nil || serverReadTimeout < 1 {
		return 0, 0, fmt.Errorf(errorFormat, ServerReadTimeoutKey, serverReadTimeout, serverReadTimeoutError)
	}

	serverWriteTimeout, serverWriteTimeoutError := strconv.Atoi(helpers.GetEnvVar(ServerWriteTimeoutKey, ServerTimeoutDefaultValue))
	if serverWriteTimeoutError != nil || serverWriteTimeout < 1 {
		return 0, 0, fmt.Errorf(errorFormat, ServerWriteTimeoutKey, serverWriteTimeout, serverWriteTimeoutError)
	}
	return serverReadTimeout, serverWriteTimeout, nil
}

func getAPIAuthenticationParams() (string, string) {
	apiAuthUsername := helpers.GetEnvVar(APIAuthUsernameKey, APIAuthUsernameDefaultValue)
	apiAuthPassword := helpers.GetEnvVar(APIAuthPasswordKey, APIAuthPasswordDefaultValue)
	return apiAuthUsername, apiAuthPassword
}

package api

const (
	RSACipher         = "RSA"
	CertificatePrefix = "-----BEGIN CERTIFICATE-----"
	CertificateSufix  = "-----END CERTIFICATE-----"

	QueryParamUserId             = "userId"
	QueryParamDeviceId           = "deviceId"
	QueryParamDeviceIds          = "deviceIds"
	QueryParamKeyId              = "keyId"
	QueryParamShowCompleteOutput = "completeOutput"
	QueryParamShowSplittedOutput = "splittedOutput"

	QueryValueAll = "all"
	DefaultKeyId  = "default"

	DBFieldUserId       = "userId"
	DBFieldCertificates = "certificates"

	ContentType        = "Content-Type"
	DefaultContentType = "application/json; charset=UTF-8"
)

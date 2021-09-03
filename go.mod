module github.com/whiteducksoftware/azure-arm-action

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.20
	github.com/Azure/go-autorest/autorest/adal v0.9.15
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.3 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/google/uuid v1.3.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
)

replace (
	// Temporary fix until https://github.com/Azure/go-autorest/pull/653 is merged.
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.3 => ./libs/go-autorest/azure/cli
)
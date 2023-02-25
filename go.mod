module github.com/whiteducksoftware/azure-arm-action

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v57.1.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.20
	github.com/caarlos0/env/v6 v6.7.1
	github.com/google/uuid v1.3.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.8.1
	github.com/whiteducksoftware/golang-utilities/azure/auth v0.1.0-alpha3
	github.com/whiteducksoftware/golang-utilities/azure/resources v0.1.0-alpha5
	github.com/whiteducksoftware/golang-utilities/github/actions v0.1.0-alpha6

)

require (
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

// Temporary fix until https://github.com/Azure/go-autorest/pull/653 is merged.
replace github.com/Azure/go-autorest/autorest/azure/cli v0.4.3 => ./libs/@azure/go-autorest/azure/cli

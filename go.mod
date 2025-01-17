module github.com/mariotoffia/godeviceshadow

require github.com/stretchr/testify v1.10.0

// AWS
require (
	github.com/aws/aws-sdk-go-v2 v1.33.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.24
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.39.1
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.53 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.8 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/config v1.29.0
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

go 1.23

toolchain go1.23.4

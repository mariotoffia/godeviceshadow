module github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.36.1
	github.com/aws/aws-sdk-go-v2/config v1.29.6
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.15
	github.com/mariotoffia/godeviceshadow v0.0.9
	github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence v0.1.2
	github.com/stretchr/testify v1.10.0
)

// replace github.com/mariotoffia/godeviceshadow => ../..

// AWS
require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.28
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.39.5
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.59 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.14 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

module github.com/mariotoffia/godeviceshadow/examples

go 1.24

// replace github.com/mariotoffia/godeviceshadow => ..

// replace github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence => ../persistence/dynamodbpersistence/

// replace github.com/mariotoffia/godeviceshadow/notify/selectlang => ..

// Inter module dependencies
require (
	github.com/mariotoffia/godeviceshadow v0.0.9
	github.com/mariotoffia/godeviceshadow/notify/selectlang v0.1.1
	github.com/mariotoffia/godeviceshadow/persistence/dynamodbpersistence v0.1.2
)

// Test
require github.com/stretchr/testify v1.10.0

require (
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aws/aws-sdk-go-v2 v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.54 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.28 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.39.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.9 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

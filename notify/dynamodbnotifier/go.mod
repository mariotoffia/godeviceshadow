module github.com/mariotoffia/godeviceshadow/notify/dynamodbnotifier

go 1.23.5

require github.com/mariotoffia/godeviceshadow v0.0.3

// replace github.com/mariotoffia/godeviceshadow => ../..

// AWS
require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.15.28
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.39.5
)

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.24.15 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)

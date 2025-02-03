package dynamodbnotifier

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	streamTypes "github.com/aws/aws-sdk-go-v2/service/dynamodbstreams/types"
)

// PollAsDynamoDBEvent polls the stream for new records and converts them into an events.DynamoDBEvent.
//
// If the function returns an error, the first return parameter will be empty. If any error did occur
// during record decoding, it will be placed at the corresponding index in the Error slice. All other
// entries that error is `nil` will be valid events.DynamoDBEventRecord.
func (s *DynamoDBStream) PollAsDynamoDBEvent(ctx context.Context) (events.DynamoDBEvent, error) {
	records, err := s.Poll(ctx)

	if err != nil {
		return events.DynamoDBEvent{}, err
	}

	res := events.DynamoDBEvent{
		Records: make([]events.DynamoDBEventRecord, len(records)),
	}

	for i, r := range records {
		if r.Dynamodb == nil {
			// Best effort...
			continue
		}

		res.Records[i] = convertDDBRecord(r)
	}

	return res, nil
}

// convertAttributeValue converts a DynamoDB attribute value (from AWS SDK v2)
// to an events.DynamoDBAttributeValue.
func convertAttributeValue(av streamTypes.AttributeValue) events.DynamoDBAttributeValue {
	switch v := av.(type) {
	case *streamTypes.AttributeValueMemberS:
		return events.NewStringAttribute(v.Value)
	case *streamTypes.AttributeValueMemberN:
		return events.NewNumberAttribute(v.Value)
	case *streamTypes.AttributeValueMemberBOOL:
		return events.NewBooleanAttribute(v.Value)
	case *streamTypes.AttributeValueMemberB:
		return events.NewBinaryAttribute(v.Value)
	case *streamTypes.AttributeValueMemberSS:
		return events.NewStringSetAttribute(v.Value)
	case *streamTypes.AttributeValueMemberNS:
		return events.NewNumberSetAttribute(v.Value)
	case *streamTypes.AttributeValueMemberBS:
		return events.NewBinarySetAttribute(v.Value)
	case *streamTypes.AttributeValueMemberM:
		return events.NewMapAttribute(convertAttributeMap(v.Value))
	case *streamTypes.AttributeValueMemberL:
		list := make([]events.DynamoDBAttributeValue, 0, len(v.Value))

		for _, item := range v.Value {
			list = append(list, convertAttributeValue(item))
		}

		return events.NewListAttribute(list)
	case *streamTypes.AttributeValueMemberNULL:
		return events.NewNullAttribute()
	default:
		// This should never happen.
		return events.DynamoDBAttributeValue{}
	}
}

// convertAttributeMap converts a map of attribute values.
func convertAttributeMap(in map[string]streamTypes.AttributeValue) map[string]events.DynamoDBAttributeValue {
	out := make(map[string]events.DynamoDBAttributeValue, len(in))

	for key, value := range in {
		out[key] = convertAttributeValue(value)
	}
	return out
}

// convertStreamRecord converts a DynamoDB streams record to an events.DynamoDBStreamRecord.
func convertStreamRecord(sr *streamTypes.StreamRecord) events.DynamoDBStreamRecord {
	var approxCreation events.SecondsEpochTime

	if sr.ApproximateCreationDateTime != nil {
		// Use Unix time (seconds) as a float.
		approxCreation = events.SecondsEpochTime{Time: *sr.ApproximateCreationDateTime}
	}

	return events.DynamoDBStreamRecord{
		ApproximateCreationDateTime: approxCreation,
		Keys:                        convertAttributeMap(sr.Keys),
		NewImage:                    convertAttributeMap(sr.NewImage),
		OldImage:                    convertAttributeMap(sr.OldImage),
		SequenceNumber:              aws.ToString(sr.SequenceNumber),
		SizeBytes:                   aws.ToInt64(sr.SizeBytes),
		StreamViewType:              string(sr.StreamViewType),
	}
}

// convertDDBRecord converts a single streamTypes.Record into an events.DynamoDBEventRecord.
func convertDDBRecord(r streamTypes.Record) events.DynamoDBEventRecord {
	return events.DynamoDBEventRecord{
		EventID:      aws.ToString(r.EventID),
		EventName:    string(r.EventName),
		EventVersion: aws.ToString(r.EventVersion),
		EventSource:  aws.ToString(r.EventSource),
		AWSRegion:    aws.ToString(r.AwsRegion),
		Change:       convertStreamRecord(r.Dynamodb),
	}
}

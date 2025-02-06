package dynamodbnotifier

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/mariotoffia/godeviceshadow/model"
)

func processRecord(record events.DynamoDBEventRecord, tr model.TypeRegistry) (*PersistenceObject, *PersistenceObject, error) {
	var (
		oldImage, newImage *PersistenceObject
		err                error
	)

	if len(record.Change.OldImage) > 0 {
		oldImage, err = processImage(convertAttributes(record.Change.OldImage), tr)
		if err != nil {
			return nil, nil, fmt.Errorf("old image error: %v", err)
		}
	}

	// Process NewImage
	if len(record.Change.NewImage) > 0 {
		newImage, err = processImage(convertAttributes(record.Change.NewImage), tr)
		if err != nil {
			return nil, nil, fmt.Errorf("New Image Error: %v\n", err)
		}
	}

	return oldImage, newImage, nil
}

// processImage converts a DynamoDB image to serialized JSON
func processImage(image map[string]types.AttributeValue, tr model.TypeRegistry) (*PersistenceObject, error) {
	// Extract PK and SK from the image
	pk, sk, err := extractKeys(image)
	if err != nil {
		return nil, err
	}

	desired := image["Desired"]
	reported := image["Reported"]

	if _, ok := desired.(*types.AttributeValueMemberNULL); ok {
		desired = nil
	}

	if _, ok := reported.(*types.AttributeValueMemberNULL); ok {
		reported = nil
	}

	// Unlink since we're doing this manually
	image["Desired"] = nil
	image["Reported"] = nil

	var po PersistenceObject

	// Unmarshal DynamoDB image into the model
	if err := attributevalue.UnmarshalMap(image, &po); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	po.Meta = map[string]any{"pk": pk, "sk": sk}

	typeEntry, ok := ResolveType(tr, pk, sk)

	var dm, rm any

	if typeEntry.Model.Kind() == reflect.Ptr {
		dm = reflect.New(typeEntry.Model.Elem()).Interface()
		rm = reflect.New(typeEntry.Model.Elem()).Interface()
	} else {
		dm = reflect.New(typeEntry.Model).Interface()
		rm = reflect.New(typeEntry.Model).Interface()
	}

	if desired != nil {
		if err := attributevalue.UnmarshalMap(desired.(*types.AttributeValueMemberM).Value, dm); err != nil {
			return nil, fmt.Errorf("unmarshal desired error: %v", err)
		}

		po.Desired = dm
	}

	if reported != nil {
		if err := attributevalue.UnmarshalMap(reported.(*types.AttributeValueMemberM).Value, rm); err != nil {
			return nil, fmt.Errorf("unmarshal reported error: %v", err)
		}

		po.Reported = rm
	}

	if !ok {
		return nil, fmt.Errorf("no model resolved for PK=%s, SK=%s", pk, sk)
	}

	return &po, nil
}

// extractKeys retrieves PK and SK from a DynamoDB image
func extractKeys(image map[string]types.AttributeValue) (pk, sk string, err error) {
	// Adjust attribute names if your table uses different key names
	if err = attributevalue.Unmarshal(image["PK"], &pk); err != nil {
		return "", "", fmt.Errorf("failed to extract PK: %v", err)
	}
	if err = attributevalue.Unmarshal(image["SK"], &sk); err != nil {
		return "", "", fmt.Errorf("failed to extract SK: %v", err)
	}
	return pk, sk, nil
}

// convertAttributes converts a map of string -> events.DynamoDBAttributeValue
// to the equivalent map of string -> types.AttributeValue.
func convertAttributes(m map[string]events.DynamoDBAttributeValue) map[string]types.AttributeValue {
	result := make(map[string]types.AttributeValue, len(m))

	for k, v := range m {
		result[k] = eventsAttrToSDKAttr(v)
	}

	return result
}

// eventsAttrToSDKAttr converts a single events.DynamoDBAttributeValue to
// an AWS SDK for Go v2 types.AttributeValue, handling all data types.
func eventsAttrToSDKAttr(eav events.DynamoDBAttributeValue) types.AttributeValue {
	switch eav.DataType() {
	case events.DataTypeString:
		return &types.AttributeValueMemberS{Value: eav.String()}
	case events.DataTypeBoolean:
		return &types.AttributeValueMemberBOOL{Value: eav.Boolean()}
	case events.DataTypeNumber:
		return &types.AttributeValueMemberN{Value: eav.Number()}
	case events.DataTypeBinary:
		return &types.AttributeValueMemberB{Value: eav.Binary()}
	case events.DataTypeStringSet:
		return &types.AttributeValueMemberSS{Value: eav.StringSet()}
	case events.DataTypeNumberSet:
		return &types.AttributeValueMemberNS{Value: eav.NumberSet()}
	case events.DataTypeBinarySet:
		return &types.AttributeValueMemberBS{Value: eav.BinarySet()}
	case events.DataTypeList:
		list := make([]types.AttributeValue, 0, len(eav.List()))
		for _, item := range eav.List() {
			list = append(list, eventsAttrToSDKAttr(item))
		}

		return &types.AttributeValueMemberL{Value: list}
	case events.DataTypeMap:
		subMap := make(map[string]types.AttributeValue, len(eav.Map()))
		for mk, mv := range eav.Map() {
			subMap[mk] = eventsAttrToSDKAttr(mv)
		}

		return &types.AttributeValueMemberM{Value: subMap}
	case events.DataTypeNull:
		return &types.AttributeValueMemberNULL{Value: true}
	default:
		// Fallback for an unexpected or unknown type:
		return &types.AttributeValueMemberNULL{Value: true}
	}
}

func ResolveType(tr model.TypeRegistry, pk, sk string) (model.TypeEntry, bool) {
	pk = strings.TrimPrefix(pk, "DS#")

	if len(sk) > 4 && sk[3] == '#' {
		sk = sk[4:]
	}

	if resolver, ok := tr.(model.TypeRegistryResolver); ok {
		if t, ok := resolver.ResolveByID(pk, sk); ok {
			return t, true
		}
	}

	// Fallback
	return tr.Get(fmt.Sprintf("%s#%s", pk, sk))
}

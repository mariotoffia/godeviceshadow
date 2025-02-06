package dynamodbnotifier

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mariotoffia/godeviceshadow/model/persistencemodel"
)

// PersistenceObject is copied from the _dynamodbpersistence_ package.
type PersistenceObject struct {
	Version int64 `json:"version"`
	// TimeStamp is a unix64 nano timestamp with UTC time
	TimeStamp int64 `json:"timestamp"`
	// ClientToken is a unique token for the client that initiated the request
	ClientToken string `json:"clientToken,omitempty"`
	// Desired is the desired model (if such is present). Depending on how the persistence is
	// configured, this may be stored separately from the reported model.
	Desired any `json:"desired,omitempty"`
	// Reported is the reported model (if such is present). Depending on how the persistence is
	// configured, this may be stored separately from the desired model.
	Reported any `json:"reported,omitempty"`
	// Meta is a map of meta data that is not stored in the persistence object but is used by
	// the `Processor`.
	//
	// Wellknown keys are:
	// - `record`: (*events.DynamoDBEventRecord) id the record produced this object
	//
	// When old and new images are both valid, they will contain the same `record` pointer.
	Meta map[string]any `json:"-" dynamodbav:"-"`
}

type DynamoDbEventType string

const (
	DynamoDbEventTypeUnknown DynamoDbEventType = ""
	DynamoDbEventTypeInsert  DynamoDbEventType = "INSERT"
	DynamoDbEventTypeModify  DynamoDbEventType = "MODIFY"
	DynamoDbEventTypeRemove  DynamoDbEventType = "REMOVE"
)

func (p *PersistenceObject) Record() *events.DynamoDBEventRecord {
	if p == nil || p.Meta == nil {
		return nil
	}

	if v, ok := p.Meta["record"]; ok {
		return v.(*events.DynamoDBEventRecord)
	}

	return nil
}

func (p *PersistenceObject) EventType() DynamoDbEventType {
	if rec := p.Record(); rec != nil {
		return DynamoDbEventType(rec.EventName)
	}

	return DynamoDbEventTypeUnknown
}

func (p *PersistenceObject) ID() persistencemodel.ID {
	if p == nil || p.Meta == nil {
		return persistencemodel.ID{}
	}

	if pk, ok := p.Meta["pk"].(string); ok {
		if sk, ok := p.Meta["sk"].(string); ok {
			pk = strings.TrimPrefix(pk, "DS#")

			if len(sk) > 4 && sk[3] == '#' {
				sk = sk[4:]
			}

			return persistencemodel.ID{ID: pk, Name: sk}
		}
	}

	return persistencemodel.ID{}
}

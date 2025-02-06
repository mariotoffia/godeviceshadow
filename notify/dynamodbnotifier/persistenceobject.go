package dynamodbnotifier

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
}

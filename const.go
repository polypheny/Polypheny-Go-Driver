package polypheny

// Connection/API versions
const (
	majorApiVersion  = 2
	minorApiVersion  = 0
	transportVersion = "plain-v1@polypheny.com\n"
)

// Connection Status
const (
	statusDisconnected       = int32(0)
	statusServerConnected    = int32(1)
	statusPolyphenyConnected = int32(2)
)

// The delimiter between query language name and the actual query
const QueryDelimiter = ":"

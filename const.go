package polypheny

// Connection/API versions
const (
	majorApiVersion  = 2
	minorApiVersion  = 0
	transportVersion = "plain-v1@polypheny.com\n"
)

// Connection Status
const (
	statusDisconnected       = 0
	statusServerConnected    = 1
	statusPolyphenyConnected = 2
)

// The delimiter between query language name and the actual query
const QueryDelimiter = ":"

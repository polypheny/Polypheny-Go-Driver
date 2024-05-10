package polypheny

// Connection/API versions
const (
	majorApiVersion  = 1
	minorApiVersion  = 2
	transportVersion = "plain-v1@polypheny.com\n"
)

// Connection Status
const (
	statusDisconnected       = int32(0)
	statusServerConnected    = int32(1)
	statusPolyphenyConnected = int32(2)
)

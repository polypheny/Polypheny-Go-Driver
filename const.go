package polypheny

import (
	prism "github.com/polypheny/Polypheny-Go-Driver/prism"
)

// Connection/API versions
const (
	majorApiVersion  = prism.MajorVersion
	minorApiVersion  = prism.MinorVersion
	transportVersion = "plain-v1@polypheny.com\n"
)

// Connection Status
const (
	statusDisconnected       = int32(0)
	statusServerConnected    = int32(1)
	statusPolyphenyConnected = int32(2)
)

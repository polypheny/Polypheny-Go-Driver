package polypheny

import (
	testing "testing"
	protos "polypheny.com/protos"
)

func TestConnect(t *testing.T) {
	client := handleConnectRequest("localhost:20590")
	if client.isConnected != true {
		t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
	}
}

func TestDuplicateConnection(t *testing.T) {
	client := handleConnectRequest("localhost:20590")
		request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID, // use the same uuid as the previous client
        }
	_, err := client.client.Connect(client.ctx, &request)
	if err == nil {
		t.Fatal("Expected an error when sending a connect request with a duplicated UUID")
	}
}

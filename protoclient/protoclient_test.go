package protoclient

import (
	"testing"
	protos "polypheny.com/protos"
)

func TestConnect(t *testing.T) {
	client := Connect("localhost:20590")
	if client.isConnected != true {
		t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
	}
}

func TestDuplicateUUID(t *testing.T) {
	client := Connect("localhost:20590")
	request := protos.ConnectionRequest{
                MajorApiVersion: MajorApiVersion,
                MinorApiVersion: MinorApiVersion,
                ClientUuid: client.clientUUID,
        }
	_, err := client.client.Connect(client.ctx, &request)
	if err == nil {
		t.Fatal("Expected an error when sending a connect request with a duplicated UUID")
	}
}

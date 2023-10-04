package protoclient

import (
	"testing"
)

func TestConnect(t *testing.T) {
	client := Connect("localhost:20590")
	if client.isConnected != true {
		t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
	}
}

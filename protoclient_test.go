package polypheny

import (
	testing "testing"
	protos "polypheny.com/protos"
	context "context"
	time "time"
	grpc "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
)

func TestProtoConnect(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "", true)
	if client.isConnected != true {
		t.Fatal("Protoclient.isConnected is not set to true, but no error was caught")
	}
}

func TestProtoDuplicateConnection(t *testing.T) {
	username := "pa"
	password := ""
	is_auto_commit := true
	client := handleConnectRequest("localhost:20590", username, password, is_auto_commit)
	request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID, // use the same uuid as the previous client
		Username: &username,
		Password: &password,
        }
	_, err := client.client.Connect(client.ctx, &request)
	if err == nil {
		t.Fatal("Expected an error when sending a connect request with a duplicated UUID")
	}
}

func TestProtoDisconnect(t *testing.T) {
	client := handleConnectRequest("localhost:20590", "pa", "")
	client.handleDisconnectRequest()
	if client.isConnected != false {
		t.Fatal("Protoclient.isConnected is not false after disconnecting, however no error was caught")
	}
	// now we should be able to connect to the server with the previous uuid
	username := "pa"
        password := ""
	request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID, // use the same uuid as the previous client
		Username: &username,
                Password: &password,
        }
	// we cannot use client.ctx anymore, since it has been canceled
	// same with the client.Connect
	// therefore, we need to rebuild these objects
	conn, err := grpc.Dial("localhost:20590", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
        if err != nil {
                t.Fatalf("did not connect: %v", err)
        }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	c := protos.NewProtoInterfaceClient(conn)
	_, err = c.Connect(ctx, &request)
	if err != nil {
		//t.Fatal("Failed to re-connect with the same uuid")
		t.Fatal(err)
	}
}

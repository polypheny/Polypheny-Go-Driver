package protoclient

import (
	protos "polypheny.com/protos"
	uuid "github.com/google/uuid"
        context "context"
        grpc "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
        "google.golang.org/grpc/metadata"
        log "log"
        time "time"
)

const (
        MajorApiVersion = 2
        MinorApiVersion = 0
)

type ProtoClient struct {
        address string
        clientUUID string
        connection *grpc.ClientConn
        client protos.ProtoInterfaceClient
        ctx context.Context
        cancel context.CancelFunc
        isConnected bool
        responseStreamUnprepared protos.ProtoInterface_ExecuteUnparameterizedStatementClient
}

type PolyphenyKeyValuePair struct {
	key interface{}
	value interface{}
}

func newClient(address string) *ProtoClient {
	clientUUID := uuid.New().String()
        conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
                log.Fatalf("did not connect: %v", err)
        }
        c := protos.NewProtoInterfaceClient(conn)
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        ctx = metadata.AppendToOutgoingContext(ctx, "clientUUID", clientUUID)
        client := ProtoClient{
                address: address,
                clientUUID: clientUUID,
                connection: conn,
                client: c,
                ctx: ctx,
                cancel: cancel,
                isConnected: false,
                responseStreamUnprepared: nil,
        }
	return &client
}

func Connect(address string) *ProtoClient {
	client := newClient(address)
        request := protos.ConnectionRequest{
                MajorApiVersion: MajorApiVersion,
                MinorApiVersion: MinorApiVersion,
                ClientUuid: client.clientUUID,
        }
	response, err := client.client.Connect(client.ctx, &request)
        if err != nil {
                log.Fatalf("could not connect: %v", err)
        }
        if response.GetIsCompatible() != true {
                log.Fatalf("could not connect")
        }
	client.isConnected = true
        return client
}

func (c *ProtoClient) Close() {
        request := protos.DisconnectRequest{}
        _, err := c.client.Disconnect(c.ctx, &request)
        if err != nil {
                log.Fatalf("could not disconnect: %v", err)
        }
        c.cancel()
        c.connection.Close()
        c.isConnected = false
}

func (c *ProtoClient) Commit() {
	request := protos.CommitRequest{}
	_, err := c.client.CommitTransaction(c.ctx, &request)
	if err != nil {
                log.Fatalf("could not commit: %v", err)
        }
}

func (c *ProtoClient) ExecuteUnprepared(statement string, language string) bool {
        request := protos.ExecuteUnparameterizedStatementRequest{
                LanguageName: language,
                Statement: statement,
        }
        response, err := c.client.ExecuteUnparameterizedStatement(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        c.responseStreamUnprepared = response
        return true
}

func convertValues(raw protos.ProtoValue) interface{} {
        switch t := raw.GetType(); t {
        case protos.ProtoValue_UNSPECIFIED:
                return nil
        case protos.ProtoValue_BOOLEAN:
                return raw.GetBoolean().GetBoolean()
        case protos.ProtoValue_INTEGER:
                return raw.GetInteger().GetInteger()
        case protos.ProtoValue_BIGINT:
                return raw.GetLong().GetLong()
        case protos.ProtoValue_DOUBLE:
                return raw.GetDouble().GetDouble()
        case protos.ProtoValue_FLOAT:
                return raw.GetFloat().GetFloat()
        case protos.ProtoValue_VARCHAR:
                return raw.GetString_().GetString_()
        case protos.ProtoValue_BINARY:
                return raw.GetBinary().GetBinary()
        case protos.ProtoValue_VARBINARY:
                return raw.GetBinary().GetBinary()
        case protos.ProtoValue_NULL:
                return nil
        case protos.ProtoValue_ROW_ID:
                return raw.GetRowId().GetRowId()
        default:
                log.Fatalf("This is likely a bug: %T %v", raw, raw)
                return nil
        }
        return nil
}

func (c *ProtoClient) FetchResult() [][]interface{} {
	// the first is nil
	result, err := c.responseStreamUnprepared.Recv()
        if err != nil {
                log.Fatalf("%v", err)
        }
	// now the second
        result, err = c.responseStreamUnprepared.Recv()
        if err != nil {
                log.Fatalf("%v", err)
        }
        rawdata := result.GetResult()
        frame := rawdata.GetFrame()
	var values [][]interface{} // return values
	if rawdata.GetScalar() == 0 {
                if len(frame.GetRelationalFrame().GetRows()) != 0 {
                        rows := frame.GetRelationalFrame().GetRows()
                        var currentRow []interface{}
                        for _, v := range rows {
                                currentRow = []interface{}{}
                                for _, z := range v.GetValues() {
                                        currentRow = append(currentRow, convertValues(*z))
                                }
                                values = append(values, currentRow)
                        }
                        return values
                } else if len(frame.GetDocumentFrame().GetDocuments()) != 0{
			documents := frame.GetDocumentFrame().GetDocuments()
			var kv PolyphenyKeyValuePair
			var currentDocument []interface{}
			for _, entries := range documents {
				currentDocument = []interface{}{}
				for _, v := range entries.GetEntries() {
					kv.key = convertValues( *(v.GetKey()) )
					kv.value = convertValues( *(v.GetValue()) )
					currentDocument = append(currentDocument, kv)
				}
				values = append(values, currentDocument)
			}
			return values
                } else {
                        return nil//frame.GetGraphFrame()
                }
        } else {
                return nil
        }
}

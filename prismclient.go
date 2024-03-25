package polypheny

import (
	binary "encoding/binary"
	log "log"
	net "net"
	prism "polypheny/protos"

	proto "google.golang.org/protobuf/proto"
)

const (
	majorApiVersion = 2
	minorApiVersion = 0
)

const (
	statusDisconnected       = 0
	statusServerConnected    = 1
	statusPolyphenyConnected = 2
)

// Let's see if this is enough, the currecnt idea is hold all the low-level rrequests and responses here
// the fecthed results and other things should be returned to the caller ---- the polyclient.
type prismClient struct {
	address     string   // addr:port
	username    string   // username is stored, but password is not
	conn        net.Conn // the Conn struct returned by Dial
	isConnected int
}

type documentKeyValuePair struct {
	key   interface{}
	value interface{}
}

// TODO: maybe re-org these source files by seperating all the types?
type PolyphenyVersionResponse struct {
	dbmsName     string
	versionName  string
	majorVersion int32
	minorVersion int32
}

type DatabaseEntryResponse struct {
	databaseName         string
	ownerName            string
	defaultNamespaceName string
}

func newConnection(address string, username string) *prismClient { // TODO: is there a better way to pass password?
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	client := prismClient{
		address:     address,
		username:    username,
		conn:        conn,
		isConnected: statusServerConnected,
	}
	return &client
}

func (c *prismClient) serialize(m proto.Message) []byte {
	result, err := proto.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (c *prismClient) send(serialized []byte) {
	length := make([]byte, 8)
	binary.LittleEndian.PutUint32(length, uint32(len(serialized)))
	c.conn.Write(length)
	c.conn.Write(serialized)
}

func (c *prismClient) recv() []byte {
	length := make([]byte, 8)
	c.conn.Read(length)
	recvLength := binary.LittleEndian.Uint64(length)
	buf := make([]byte, recvLength)
	c.conn.Read(buf)
	return buf
}

func (c *prismClient) close() {
	err := c.conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	c.isConnected = statusDisconnected
}

func (c *prismClient) helperSendAndRecv(m proto.Message) *prism.Response {
	c.send(c.serialize(m))
	buf := c.recv()
	var response prism.Response
	proto.Unmarshal(buf, &response)
	return &response
}

func handleConnectRequest(address string, username string, password string) *prismClient {
	client := newConnection(address, username)
	request := prism.Request{
		Type: &prism.Request_ConnectionRequest{
			ConnectionRequest: &prism.ConnectionRequest{
				MajorApiVersion: majorApiVersion,
				MinorApiVersion: minorApiVersion,
				Username:        &username,
				Password:        &password,
			},
		},
	}
	client.send(client.serialize(&request))
	buf := client.recv()
	var response prism.Response
	proto.Unmarshal(buf, &response)
	if response.GetConnectionResponse().IsCompatible {
		client.isConnected = statusPolyphenyConnected
	}
	return client
}

func (c *prismClient) handleDisconnectRequest() {
	request := prism.Request{
		Type: &prism.Request_DisconnectRequest{
			DisconnectRequest: &prism.DisconnectRequest{},
		},
	}
	c.send(c.serialize(&request))
	c.recv()
	c.isConnected = statusServerConnected
	c.close()
}

func convertProtoValue(raw *prism.ProtoValue) interface{} {
	if raw.GetBoolean() != nil {
		return raw.GetBoolean().GetBoolean()
	} else if raw.GetInteger() != nil {
		return raw.GetInteger().GetInteger()
	} else if raw.GetLong() != nil {
		return raw.GetLong().GetLong()
	} else if raw.GetBigDecimal() != nil {
		// TODO: add support to big decimals
		return nil
	} else if raw.GetFloat() != nil {
		return raw.GetFloat().GetFloat()
	} else if raw.GetDouble() != nil {
		return raw.GetDouble().GetDouble()
	} else if raw.GetString_() != nil {
		return raw.GetString_().GetString_()
	} else {
		return nil
	}
}

func (c *prismClient) handleExecuteUnparameterizedStatementRequest(language string, statement string) [][]interface{} {
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName: language,
				Statement:    statement,
			},
		},
	}
	c.send(c.serialize(&request))
	buf := c.recv()
	var response prism.Response
	proto.Unmarshal(buf, &response)
	requestID := response.GetStatementResponse().GetStatementId()
	buf = c.recv() // this is the query result
	proto.Unmarshal(buf, &response)
	if requestID != response.GetStatementResponse().GetStatementId() {
		return nil
	}
	if response.GetStatementResponse().GetResult() == nil {
		return nil
	}
	if response.GetStatementResponse().GetResult().GetFrame() == nil {
		return nil
	}

	frame := response.GetStatementResponse().GetResult().GetFrame()
	var values [][]interface{}
	if frame.GetRelationalFrame() != nil {
		relationalData := frame.GetRelationalFrame()
		rows := relationalData.GetRows()
		var currentRow []interface{}
		for _, irow := range rows {
			currentRow = []interface{}{}
			for _, ivalue := range irow.GetValues() {
				currentRow = append(currentRow, convertProtoValue(ivalue))
			}
			values = append(values, currentRow)
		}
		return values
	} else if frame.GetDocumentFrame() != nil {
		documentData := frame.GetDocumentFrame().GetDocuments()
		var kv documentKeyValuePair
		var currentDocument []interface{}
		for _, entries := range documentData {
			currentDocument = []interface{}{}
			for _, v := range entries.GetEntries() {
				kv.key = convertProtoValue(v.GetKey())
				kv.value = convertProtoValue(v.GetValue())
				currentDocument = append(currentDocument, kv)
			}
			values = append(values, currentDocument)
		}
		return values
	} else {
		// graph is currently not supported
		return nil
	}

}

func (c *prismClient) handleFetchRequest(statementId int32) [][]interface{} {
	request := prism.Request{
		Type: &prism.Request_FetchRequest{
			FetchRequest: &prism.FetchRequest{
				StatementId: statementId,
			},
		},
	}
	c.send(c.serialize(&request))
	buf := c.recv()
	var response prism.Response
	proto.Unmarshal(buf, &response)
	if response.GetStatementResponse().GetResult() == nil {
		return nil
	}
	if response.GetStatementResponse().GetResult().GetFrame() == nil {
		return nil
	}

	frame := response.GetStatementResponse().GetResult().GetFrame()
	var values [][]interface{}
	if frame.GetRelationalFrame() != nil {
		relationalData := frame.GetRelationalFrame()
		rows := relationalData.GetRows()
		var currentRow []interface{}
		for _, irow := range rows {
			currentRow = []interface{}{}
			for _, ivalue := range irow.GetValues() {
				currentRow = append(currentRow, convertProtoValue(ivalue))
			}
			values = append(values, currentRow)
		}
		return values
	} else if frame.GetDocumentFrame() != nil {
		documentData := frame.GetDocumentFrame().GetDocuments()
		var kv documentKeyValuePair
		var currentDocument []interface{}
		for _, entries := range documentData {
			currentDocument = []interface{}{}
			for _, v := range entries.GetEntries() {
				kv.key = convertProtoValue(v.GetKey())
				kv.value = convertProtoValue(v.GetValue())
				currentDocument = append(currentDocument, kv)
			}
			values = append(values, currentDocument)
		}
		return values
	} else {
		// graph is currently not supported
		return nil
	}
}

func (c *prismClient) handleCommitRequest() {
	request := prism.Request{
		Type: &prism.Request_CommitRequest{
			CommitRequest: &prism.CommitRequest{},
		},
	}
	c.send(c.serialize(&request))
	c.recv()
}

func (c *prismClient) handleDbmsVersionRequest() PolyphenyVersionResponse {
	request := prism.Request{
		Type: &prism.Request_DbmsVersionRequest{
			DbmsVersionRequest: &prism.DbmsVersionRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	result := PolyphenyVersionResponse{
		dbmsName:     response.GetDbmsVersionResponse().GetDbmsName(),
		versionName:  response.GetDbmsVersionResponse().GetVersionName(),
		majorVersion: response.GetDbmsVersionResponse().GetMajorVersion(),
		minorVersion: response.GetDbmsVersionResponse().GetMinorVersion(),
	}
	return result
}

func (c *prismClient) handleLanguageRequest() []string {
	request := prism.Request{
		Type: &prism.Request_LanguageRequest{
			LanguageRequest: &prism.LanguageRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetLanguageResponse().GetLanguageNames()
}

func (c *prismClient) handleDatabaseRequest() []DatabaseEntryResponse {
	request := prism.Request{
		Type: &prism.Request_DatabasesRequest{
			DatabasesRequest: &prism.DatabasesRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []DatabaseEntryResponse
	for _, entry := range response.GetDatabasesResponse().GetDatabases() {
		result = append(result, DatabaseEntryResponse{
			databaseName:         entry.GetDatabaseName(),
			ownerName:            entry.GetOwnerName(),
			defaultNamespaceName: entry.GetDefaultNamespaceName(),
		})
	}
	return result
}

func (c *prismClient) handleTableTypeRequest() []string {
	request := prism.Request{
		Type: &prism.Request_TableTypesRequest{
			TableTypesRequest: &prism.TableTypesRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []string
	for _, typeName := range response.GetTableTypesResponse().GetTableTypes() {
		result = append(result, typeName.GetTableType())
	}
	return result
}

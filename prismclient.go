package polypheny

import (
	binary "encoding/binary"
	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
	log "log"
	net "net"

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

type UnparameterizedStatementRequest struct {
	language      string
	statement     string
	fetchSize     *int32
	namespaceName *string
}

type IndexedParametersRequest struct {
	parameters []interface{}
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

type TypeResponse struct {
	typeName        string
	precision       int32
	literalPrefix   string
	literalSuffix   string
	isCaseSensitive bool
	isSearchable    int32
	isAutoIncrement bool
	minScale        int32
	maxScale        int32
	radix           int32
}

type UserDefinedTypesResponse struct {
	typeName   string
	metaValues []string
}

type ProceduresResponse struct {
	trivialName     string
	inputParameters string
	desc            string
	returnType      int32
	uniqleName      string
}

type ClientInfoPropertyMetaResponse struct {
	key          string
	defaultValue string
	maxLength    int32
	desc         string
}

type FunctionsResponse struct {
	name             string
	syntax           string
	functionCategory string
	isTableFunction  bool
}

type NamespaceResponse struct {
	namespaceName   string
	isCaseSensitive bool
	namespaceType   string
}

type ParameterMetaResponse struct {
	precision     int32
	scale         int32
	typeName      string
	parameterName string
	name          string
}

func newConnection(address string, username string) *prismClient {
	// This function is used to establish a low-level tcp connection, not a higher level polypheny connection
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

func (c *prismClient) handleConnectionPropertiesUpdateRequest(isAutoCommit *bool, namespaceName *string) {
	request := prism.Request{
		Type: &prism.Request_ConnectionPropertiesUpdateRequest{
			ConnectionPropertiesUpdateRequest: &prism.ConnectionPropertiesUpdateRequest{
				ConnectionProperties: &prism.ConnectionProperties{
					IsAutoCommit:  isAutoCommit,
					NamespaceName: namespaceName,
				},
			},
		},
	}
	c.helperSendAndRecv(&request)
}

func (c *prismClient) handleConnectionCheckRequest() {
	request := prism.Request{
		Type: &prism.Request_ConnectionCheckRequest{
			ConnectionCheckRequest: &prism.ConnectionCheckRequest{},
		},
	}
	c.helperSendAndRecv(&request)
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

func (c *prismClient) handleCloseStatementRequest(statementId int32) {
	request := prism.Request{
		Type: &prism.Request_CloseStatementRequest{
			CloseStatementRequest: &prism.CloseStatementRequest{
				StatementId: statementId,
			},
		},
	}
	c.send(c.serialize(&request))
	c.recv()
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

func makeProtoValue(value interface{}) *prism.ProtoValue {
	var result prism.ProtoValue
	switch value.(type) {
	case bool:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Boolean{
				Boolean: &prism.ProtoBoolean{
					Boolean: value.(bool),
				},
			},
		}
	case int32:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Integer{
				Integer: &prism.ProtoInteger{
					Integer: value.(int32),
				},
			},
		}
	case int64:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Long{
				Long: &prism.ProtoLong{
					Long: value.(int64),
				},
			},
		}
	case float64:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Double{
				Double: &prism.ProtoDouble{
					Double: value.(float64),
				},
			},
		}
	case float32:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Float{
				Float: &prism.ProtoFloat{
					Float: value.(float32),
				},
			},
		}
	case string:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_String_{
				String_: &prism.ProtoString{
					String_: value.(string),
				},
			},
		}
	default:
		log.Fatalf("Lack of support to %T %v", value, value)
	}
	return &result
}

func (c *prismClient) handleExecuteUnparameterizedStatementRequest(query UnparameterizedStatementRequest) (int64, []string, [][]interface{}) {
	request := prism.Request{
		Type: &prism.Request_ExecuteUnparameterizedStatementRequest{
			ExecuteUnparameterizedStatementRequest: &prism.ExecuteUnparameterizedStatementRequest{
				LanguageName:  query.language,
				Statement:     query.statement,
				FetchSize:     query.fetchSize,
				NamespaceName: query.namespaceName,
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
		return 0, nil, nil
	}
	if response.GetStatementResponse().GetResult() == nil {
		return 0, nil, nil
	}
	affectedRows := response.GetStatementResponse().GetResult().GetScalar()
	if response.GetStatementResponse().GetResult().GetFrame() == nil {
		return affectedRows, nil, nil
	}

	frame := response.GetStatementResponse().GetResult().GetFrame()
	var values [][]interface{}
	var columns []string
	if frame.GetRelationalFrame() != nil {
		relationalData := frame.GetRelationalFrame()
		columnResponse := relationalData.GetColumnMeta()
		columns = make([]string, len(columnResponse))
		for i, v := range columnResponse {
			columns[i] = v.GetColumnName()
		}
		rows := relationalData.GetRows()
		var currentRow []interface{}
		for _, irow := range rows {
			currentRow = []interface{}{}
			for _, ivalue := range irow.GetValues() {
				currentRow = append(currentRow, convertProtoValue(ivalue))
			}
			values = append(values, currentRow)
		}
		return affectedRows, columns, values
	} else if frame.GetDocumentFrame() != nil {
		documentData := frame.GetDocumentFrame().GetDocuments()
		var kv documentKeyValuePair
		//var currentDocument []interface{}
		for _, entries := range documentData {
			//currentDocument = []interface{}{}
			for _, v := range entries.GetEntries() {
				kv.key = convertProtoValue(v.GetKey())
				kv.value = convertProtoValue(v.GetValue())
				//currentDocument = append(currentDocument, kv)
				currentRow := make([]interface{}, 2)
				currentRow[0] = kv.key
				currentRow[1] = kv.value
				values = append(values, currentRow)
			}
		}
		columns = make([]string, 2)
		columns[0] = "key"
		columns[1] = "value"
		return affectedRows, columns, values
	} else {
		// graph is currently not supported
		return 0, nil, nil
	}

}

func (c *prismClient) handleExecuteUnparameterizedStatementBatchRequest(queries []UnparameterizedStatementRequest) [][][]interface{} {
	var result [][][]interface{}
	for _, query := range queries {
		_, _, response := c.handleExecuteUnparameterizedStatementRequest(query)
		result = append(result, response)
	}
	return result
}

func (c *prismClient) handlePrepareIndexedStatementRequest(language string, statement string, namespace *string) map[int32][]ParameterMetaResponse {
	request := prism.Request{
		Type: &prism.Request_PrepareIndexedStatementRequest{
			PrepareIndexedStatementRequest: &prism.PrepareStatementRequest{
				LanguageName:  language,
				Statement:     statement,
				NamespaceName: namespace,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	result := make(map[int32][]ParameterMetaResponse)
	var parameterMetas []ParameterMetaResponse
	for _, parameter := range response.GetPreparedStatementSignature().GetParameterMetas() {
		parameterMetas = append(parameterMetas, ParameterMetaResponse{
			precision:     parameter.GetPrecision(),
			scale:         parameter.GetScale(),
			typeName:      parameter.GetTypeName(),
			parameterName: parameter.GetParameterName(),
			name:          parameter.GetName(),
		})
	}
	result[response.GetPreparedStatementSignature().GetStatementId()] = parameterMetas
	return result
}

func (c *prismClient) handlePrepareNamedStatementRequest(language string, statement string, namespace *string) map[int32][]ParameterMetaResponse {
	request := prism.Request{
		Type: &prism.Request_PrepareNamedStatementRequest{
			PrepareNamedStatementRequest: &prism.PrepareStatementRequest{
				LanguageName:  language,
				Statement:     statement,
				NamespaceName: namespace,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	result := make(map[int32][]ParameterMetaResponse)
	var parameterMetas []ParameterMetaResponse
	for _, parameter := range response.GetPreparedStatementSignature().GetParameterMetas() {
		parameterMetas = append(parameterMetas, ParameterMetaResponse{
			precision:     parameter.GetPrecision(),
			scale:         parameter.GetScale(),
			typeName:      parameter.GetTypeName(),
			parameterName: parameter.GetParameterName(),
			name:          parameter.GetName(),
		})
	}
	result[response.GetPreparedStatementSignature().GetStatementId()] = parameterMetas
	return result
}

func (c *prismClient) handleExecuteIndexedStatementRequest(statementId int32, parameters IndexedParametersRequest, fetchSize *int32) [][]interface{} {
	polyvalues := make([]*prism.ProtoValue, len(parameters.parameters))
	for i, v := range parameters.parameters {
		polyvalues[i] = makeProtoValue(v)
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteIndexedStatementRequest{
			ExecuteIndexedStatementRequest: &prism.ExecuteIndexedStatementRequest{
				StatementId: statementId,
				Parameters: &prism.IndexedParameters{
					Parameters: polyvalues,
				},
				FetchSize: fetchSize,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	if response.GetStatementResult() == nil {
		return nil
	}
	if response.GetStatementResult().GetFrame() == nil {
		return nil
	}

	frame := response.GetStatementResult().GetFrame()
	var values [][]interface{}
	// TODO: can we use a separate function to handle the following?
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

func (c *prismClient) handleExecuteIndexedStatementBatchRequest(statementId int32, parameters []IndexedParametersRequest) [][][]interface{} {
	var result [][][]interface{}
	for _, parparameter := range parameters {
		result = append(result, c.handleExecuteIndexedStatementRequest(statementId, parparameter, nil))
	}
	return result
}

func (c *prismClient) handleExecuteNamedStatementRequest(statementId int32, parameters map[string]interface{}, fetchSize *int32) [][]interface{} {
	polyvalues := make(map[string]*prism.ProtoValue)
	for k, v := range parameters {
		polyvalues[k] = makeProtoValue(v)
	}
	request := prism.Request{
		Type: &prism.Request_ExecuteNamedStatementRequest{
			ExecuteNamedStatementRequest: &prism.ExecuteNamedStatementRequest{
				StatementId: statementId,
				Parameters: &prism.NamedParameters{
					Parameters: polyvalues,
				},
				FetchSize: fetchSize,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	if response.GetStatementResult() == nil {
		return nil
	}
	if response.GetStatementResult().GetFrame() == nil {
		return nil
	}

	frame := response.GetStatementResult().GetFrame()
	var values [][]interface{}
	// TODO: can we use a separate function to handle the following?
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

func (c *prismClient) handleRollbackRequest() {
	request := prism.Request{
		Type: &prism.Request_RollbackRequest{
			RollbackRequest: &prism.RollbackRequest{},
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

func (c *prismClient) handleTypesRequest() []TypeResponse {
	request := prism.Request{
		Type: &prism.Request_TypesRequest{
			TypesRequest: &prism.TypesRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []TypeResponse
	for _, t := range response.GetTypesResponse().GetTypes() {
		result = append(result, TypeResponse{
			typeName:        t.GetTypeName(),
			precision:       t.GetPrecision(),
			literalPrefix:   t.GetLiteralPrefix(),
			literalSuffix:   t.GetLiteralSuffix(),
			isCaseSensitive: t.GetIsCaseSensitive(),
			isSearchable:    t.GetIsSearchable(),
			isAutoIncrement: t.GetIsAutoIncrement(),
			minScale:        t.GetMinScale(),
			maxScale:        t.GetMaxScale(),
			radix:           t.GetRadix(),
		})
	}
	return result
}

func (c *prismClient) handleUserDefinedTypesRequest() []UserDefinedTypesResponse {
	request := prism.Request{
		Type: &prism.Request_UserDefinedTypesRequest{
			UserDefinedTypesRequest: &prism.UserDefinedTypesRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []UserDefinedTypesResponse
	for _, t := range response.GetUserDefinedTypesResponse().GetUserDefinedTypes() {
		var metaValues []string
		for _, metaValue := range t.GetValueMetas() {
			metaValues = append(metaValues, metaValue.GetValueName())
		}
		result = append(result, UserDefinedTypesResponse{
			typeName:   t.GetTypeName(),
			metaValues: metaValues,
		})
	}
	return result
}

func (c *prismClient) handleSqlStringFunctionsRequest() string {
	request := prism.Request{
		Type: &prism.Request_SqlStringFunctionsRequest{
			SqlStringFunctionsRequest: &prism.SqlStringFunctionsRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetSqlStringFunctionsResponse().GetString_()
}

func (c *prismClient) handleSqlSystemFunctionsRequest() string {
	request := prism.Request{
		Type: &prism.Request_SqlSystemFunctionsRequest{
			SqlSystemFunctionsRequest: &prism.SqlSystemFunctionsRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetSqlSystemFunctionsResponse().GetString_()
}

func (c *prismClient) handleSqlTimeDateFunctionsRequest() string {
	request := prism.Request{
		Type: &prism.Request_SqlTimeDateFunctionsRequest{
			SqlTimeDateFunctionsRequest: &prism.SqlTimeDateFunctionsRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetSqlTimeDateFunctionsResponse().GetString_()
}

func (c *prismClient) handleSqlNumericFunctionsRequest() string {
	request := prism.Request{
		Type: &prism.Request_SqlNumericFunctionsRequest{
			SqlNumericFunctionsRequest: &prism.SqlNumericFunctionsRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetSqlNumericFunctionsResponse().GetString_()
}

func (c *prismClient) handleSqlKeywordsRequest() string {
	request := prism.Request{
		Type: &prism.Request_SqlKeywordsRequest{
			SqlKeywordsRequest: &prism.SqlKeywordsRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetSqlKeywordsResponse().GetString_()
}

func (c *prismClient) handleProceduresRequest(language string, procedureNamePattern *string) []ProceduresResponse {
	request := prism.Request{
		Type: &prism.Request_ProceduresRequest{
			ProceduresRequest: &prism.ProceduresRequest{
				Language:             language,
				ProcedureNamePattern: procedureNamePattern,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []ProceduresResponse
	for _, procedure := range response.GetProceduresResponse().GetProcedures() {
		result = append(result, ProceduresResponse{
			trivialName:     procedure.GetTrivialName(),
			inputParameters: procedure.GetInputParameters(),
			desc:            procedure.GetDescription(),
			returnType:      int32(procedure.GetReturnType()),
			uniqleName:      procedure.GetUniqueName(),
		})
	}
	return result
}

func (c *prismClient) handleClientInfoPropertiesRequest() map[string]string {
	request := prism.Request{
		Type: &prism.Request_ClientInfoPropertiesRequest{
			ClientInfoPropertiesRequest: &prism.ClientInfoPropertiesRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	return response.GetClientInfoPropertiesResponse().GetProperties()
}

func (c *prismClient) handleClientInfoPropertyMetaRequest() []ClientInfoPropertyMetaResponse {
	request := prism.Request{
		Type: &prism.Request_ClientInfoPropertyMetaRequest{
			ClientInfoPropertyMetaRequest: &prism.ClientInfoPropertyMetaRequest{},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []ClientInfoPropertyMetaResponse
	for _, entry := range response.GetClientInfoPropertyMetaResponse().GetClientInfoPropertyMetas() {
		result = append(result, ClientInfoPropertyMetaResponse{
			key:          entry.GetKey(),
			defaultValue: entry.GetDefaultValue(),
			maxLength:    entry.GetMaxlength(),
			desc:         entry.GetDescription(),
		})
	}
	return result
}

func (c *prismClient) handleFunctionsRequest(language string, category string) []FunctionsResponse {
	request := prism.Request{
		Type: &prism.Request_FunctionsRequest{
			FunctionsRequest: &prism.FunctionsRequest{
				QueryLanguage:    language,
				FunctionCategory: category,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []FunctionsResponse
	for _, f := range response.GetFunctionsResponse().GetFunctions() {
		result = append(result, FunctionsResponse{
			name:             f.GetName(),
			syntax:           f.GetSyntax(),
			functionCategory: f.GetFunctionCategory(),
			isTableFunction:  f.GetIsTableFunction(),
		})
	}
	return result
}

func (c *prismClient) handleNamespaceRequest(namespaceName string) NamespaceResponse {
	request := prism.Request{
		Type: &prism.Request_NamespaceRequest{
			NamespaceRequest: &prism.NamespaceRequest{
				NamespaceName: namespaceName,
			},
		},
	}
	c.send(c.serialize(&request))
	buf := c.recv()
	var response prism.Namespace
	proto.Unmarshal(buf, &response)
	return NamespaceResponse{
		namespaceName:   response.GetNamespaceName(),
		isCaseSensitive: response.GetIsCaseSensitive(),
		namespaceType:   response.GetNamespaceType(),
	}
}

func (c *prismClient) handleNamespacesResponse(namespacePattern *string, namespaceType *string) []NamespaceResponse {
	request := prism.Request{
		Type: &prism.Request_NamespacesRequest{
			NamespacesRequest: &prism.NamespacesRequest{
				NamespacePattern: namespacePattern,
				NamespaceType:    namespaceType,
			},
		},
	}
	response := c.helperSendAndRecv(&request)
	var result []NamespaceResponse
	for _, entry := range response.GetNamespacesResponse().GetNamespaces() {
		result = append(result, NamespaceResponse{
			namespaceName:   entry.GetNamespaceName(),
			isCaseSensitive: entry.GetIsCaseSensitive(),
			namespaceType:   entry.GetNamespaceType(),
		})
	}
	return result
}

func (c *prismClient) handleEntitiesRequest(namespaceName string, entityPattern *string) {
	request := prism.Request{
		Type: &prism.Request_EntitiesRequest{
			EntitiesRequest: &prism.EntitiesRequest{
				NamespaceName: namespaceName,
				EntityPattern: entityPattern,
			},
		},
	}
	c.helperSendAndRecv(&request)
}

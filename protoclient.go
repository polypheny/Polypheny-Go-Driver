package polypheny

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
        majorApiVersion = 2
        minorApiVersion = 0
)

func getProtoAPIVersion() [2]int {
	return [2]int{majorApiVersion, minorApiVersion}
}

// Better to not store the password of a connection
type protoClient struct {
        address string
        clientUUID string
	username string
        connection *grpc.ClientConn
        client protos.ProtoInterfaceClient
        ctx context.Context
        cancel context.CancelFunc
        isConnected bool
        responseStreamUnprepared protos.ProtoInterface_ExecuteUnparameterizedStatementClient
}

type documentKeyValuePair struct {
	key interface{}
	value interface{}
}

func newProtoClient(address string, username string) *protoClient {
	clientUUID := uuid.New().String()
        conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
                log.Fatalf("did not connect: %v", err)
        }
        c := protos.NewProtoInterfaceClient(conn)
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        ctx = metadata.AppendToOutgoingContext(ctx, "clientUUID", clientUUID)
        client := protoClient{
                address: address,
                clientUUID: clientUUID,
		username: username,
                connection: conn,
                client: c,
                ctx: ctx,
                cancel: cancel,
                isConnected: false,
                responseStreamUnprepared: nil,
        }
	return &client
}

func createConnectionProperties(properties ...interface{}) protos.ConnectionProperties {
	connectionProperties := protos.ConnectionProperties{}
        // TODO: switch to a different way to handle ConnectionProperties
        for _, arg := range properties {
                switch arg.(type) {
                case string:
                        namespace_name := arg.(string)
                        connectionProperties.NamespaceName = &namespace_name
                case bool:
                        is_auto_commit := arg.(bool)
                        connectionProperties.IsAutoCommit = &is_auto_commit
                }
        }
	return connectionProperties
}

func handleConnectRequest(address string, username string, password string,  properties ...interface{}) *protoClient {
	client := newProtoClient(address, username)
	connectionProperties := createConnectionProperties(properties)
        request := protos.ConnectionRequest{
                MajorApiVersion: majorApiVersion,
                MinorApiVersion: minorApiVersion,
                ClientUuid: client.clientUUID,
		Username: &username,
		Password: &password,
		ConnectionProperties: &connectionProperties,
        }
	response, err := client.client.Connect(client.ctx, &request)
        if err != nil {
                log.Fatalf("Failed to connect: %v", err)
        }
        if response.GetIsCompatible() != true {
                log.Fatalf("The API version is not compatible with server")
        }
	client.isConnected = true
        return client
}

func (c *protoClient) handleUpdateConnectionProperties(properties ...interface{}) {
	connectionProperties := createConnectionProperties(properties)
	request := protos.ConnectionPropertiesUpdateRequest{
		ConnectionProperties: &connectionProperties,
	}
	_, err := c.client.UpdateConnectionProperties(c.ctx, &request)
	if err != nil {
                log.Fatalf("could not update: %v", err)
        }
}

func (c *protoClient) handleDisconnectRequest() {
        request := protos.DisconnectRequest{}
        _, err := c.client.Disconnect(c.ctx, &request)
        if err != nil {
                log.Fatalf("could not disconnect: %v", err)
        }
        c.cancel()
        c.connection.Close()
        c.isConnected = false
}

func (c *protoClient) handleGetClientInfoProperties() map[string]string {
	request := protos.ClientInfoPropertiesRequest{}
	resp, err := c.client.GetClientInfoProperties(c.ctx, &request)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return resp.GetProperties()
}

func (c *protoClient) handleSetClientInfoProperties(properties map[string] string) {
	request := protos.ClientInfoProperties{
		Properties: properties,
	}
	_, err := c.client.SetClientInfoProperties(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
}

func (c *protoClient) handleConnectionCheckRequest() {
	request := protos.ConnectionCheckRequest{}
	_, err := c.client.CheckConnection(c.ctx, &request)
        if err != nil {
                log.Fatalf("Checking connection failed: %v", err)
        }
}

func (c *protoClient) handleCommitRequest() {
	request := protos.CommitRequest{}
	_, err := c.client.CommitTransaction(c.ctx, &request)
	if err != nil {
                log.Fatalf("could not commit: %v", err)
        }
}

func (c *protoClient) handleRollbackRequest() {
	request := protos.RollbackRequest{}
	_, err := c.client.RollbackTransaction(c.ctx, &request)
	if err != nil {
                log.Fatalf("could not rollback: %v", err)
        }
}

func (c *protoClient) handleExecuteUnprepared(statement string, language string) bool {
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

func (c *protoClient) handleFetchiStreamResult() [][]interface{} {
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
			var kv documentKeyValuePair
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
                        return nil
                }
        } else {
                return nil
        }
}

// Meta requests
func (c *protoClient) handleGetDBMSVersion() (string, string, int32, int32) {
	request := protos.DbmsVersionRequest{}
	resp, err := c.client.GetDbmsVersion(c.ctx, &request)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return resp.GetDbmsName(), resp.GetVersionName(), resp.GetMajorVersion(), resp.GetMinorVersion()
}

func (c *protoClient) handleGetSupportedLanguage() []string {
	request := protos.LanguageRequest{}
	resp, err := c.client.GetSupportedLanguages(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
	return resp.GetLanguageNames()
}

type ProtoDatabase struct {
	databaseName string
	ownerName string
	defaultNamespaceName string
}

func (c *protoClient) handleGetDatabases() []ProtoDatabase {
	request := protos.DatabasesRequest{}
	resp, err := c.client.GetDatabases(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
	var result []ProtoDatabase
	for _, v  := range resp.GetDatabases() {
		item := ProtoDatabase{
			databaseName: v.GetDatabaseName(),
			ownerName: v.GetOwnerName(),
			defaultNamespaceName: v.GetDefaultNamespaceName(),
		}
		result = append(result, item)
	}
	return result
}

func (c *protoClient) handleGetTableTypes() []string {
	request := protos.TableTypesRequest{}
	resp, err := c.client.GetTableTypes(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
	var result []string
	for _, v := range resp.GetTableTypes() {
		result = append(result, v.GetTableType())
	}
	return result
}

func (c *protoClient) handleGetTypes() []*protos.Type {
	request := protos.TypesRequest{}
	resp, err := c.client.GetTypes(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
	return resp.GetTypes()
}

type ProtoUserDefinedType struct {
	valueName string
	metaValue []string
}

func (c *protoClient) handleGetUserDefinedTypes() []ProtoUserDefinedType {
	// not implemented in polypheny
	// the commented code is for future use
	return nil
	/*request := protos.UserDefinedTypesRequest{}
	resp, err := c.client.GetUserDefinedTypes(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
	var result []ProtoUserDefinedType
	for _, v := range resp.GetUserDefinedTypes() {
		item := ProtoUserDefinedType{
			valueName: v.GetTypeName(),
			metaValue: []string{},
		}
		for _, valueMetas := range v.GetValueMetas() {
			item.metaValue = append(item.metaValue, valueMetas.GetValueName())
		}
		result = append(result, item)
	}
	return result*/
}

func (c *protoClient) handleGetClientInfoPropertyMetas() []*protos.ClientInfoPropertyMeta {
	// not implemented in polypheny
        // the commented code is for future use
        return nil
	/*request := protos.ClientInfoPropertyMetaRequest{}
	resp, err := c.client.GetClientInfoPropertyMetas(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetClientInfoPropertyMetas()*/
}

func (c *protoClient) handleSearchProcedures(language string, pattern ...string) []*protos.Procedure {
	var procedureNamePattern string
	if len(pattern) == 0 {
		procedureNamePattern = ""
	} else {
		procedureNamePattern = pattern[0]
	}
	request := protos.ProceduresRequest{
		Language: language,
		ProcedureNamePattern: &procedureNamePattern,
	}
	resp, err := c.client.SearchProcedures(c.ctx, &request)
	if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetProcedures()
}

func (c *protoClient) handleSearchFunctions(language string, functionCategory string) []*protos.Function {
	request := protos.FunctionsRequest{
                QueryLanguage: language,
                FunctionCategory: functionCategory,
        }
        resp, err := c.client.SearchFunctions(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetFunctions()
}

func (c *protoClient) handleSearchNamespaces(pattern string, types ...string) []*protos.Namespace {
	// polypheny will throw an exception when doing the rpc
	return nil
	/*var namespaceType string
	if len(types) == 0 {
                namespaceType = ""
        } else {
                namespaceType = types[0]
        }
	request := protos.NamespacesRequest {
		NamespacePattern: &pattern,
		NamespaceType: &namespaceType,
	}
	resp, err := c.client.SearchNamespaces(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetNamespaces()*/
}

func (c *protoClient) handleGetNamespace(namespace string) *protos.Namespace {
	request := protos.NamespaceRequest{
		NamespaceName: namespace,
	}
	resp, err := c.client.GetNamespace(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp
}

func (c *protoClient) handleSearchEntities(namespaceName string, pattern ...string) []*protos.Entity {
	var entityPattern  string
        if len(pattern) == 0 {
                entityPattern = ""
        } else {
                entityPattern = pattern[0]
        }
        request := protos.EntitiesRequest{
                NamespaceName: namespaceName,
                EntityPattern: &entityPattern,
        }
        resp, err := c.client.SearchEntities(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetEntities()
}

func (c *protoClient) handleGetSqlStringFunctions() string {
	request := protos.SqlStringFunctionsRequest{}
	resp, err := c.client.GetSqlStringFunctions(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetString_()
}

func (c *protoClient) handleGetSqlSystemFunctions() string {
        request := protos.SqlSystemFunctionsRequest{}
        resp, err := c.client.GetSqlSystemFunctions(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetString_()
}

func (c *protoClient) handleGetSqlTimeDateFunctions() string {
        request := protos.SqlTimeDateFunctionsRequest{}
        resp, err := c.client.GetSqlTimeDateFunctions(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetString_()
}

func (c *protoClient) handleGetSqlNumericFunctions() string {
        request := protos.SqlNumericFunctionsRequest{}
        resp, err := c.client.GetSqlNumericFunctions(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetString_()
}

func (c *protoClient) handleGetSqlKeywords() string {
        request := protos.SqlKeywordsRequest{}
        resp, err := c.client.GetSqlKeywords(c.ctx, &request)
        if err != nil {
                log.Fatalf("%v", err)
        }
        return resp.GetString_()
}

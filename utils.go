package polypheny

import (
	driver "database/sql/driver"
	fmt "fmt"
	strings "strings"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
)

// ParseQuery splits the query language and the actual query
// TODO: add namespace support
func parseQuery(query string) (string, string, error) {
	splitted := strings.Split(query, QueryDelimiter)
	if len(splitted) != 2 {
		return "", "", &ClientError{
			message: "A query should have the following format: QueryLanguage:Query",
		}
	}
	return splitted[0], splitted[1], nil
}

// convertProtoValue converts a ProtonValue to a go value
func convertProtoValue(raw *prism.ProtoValue) (any, error) {
	if raw.GetBoolean() != nil {
		return raw.GetBoolean().GetBoolean(), nil
	} else if raw.GetInteger() != nil {
		return raw.GetInteger().GetInteger(), nil
	} else if raw.GetLong() != nil {
		return raw.GetLong().GetLong(), nil
	} else if raw.GetBigDecimal() != nil {
		// TODO: add support to big decimals
		return nil, nil
	} else if raw.GetFloat() != nil {
		return raw.GetFloat().GetFloat(), nil
	} else if raw.GetDouble() != nil {
		return raw.GetDouble().GetDouble(), nil
	} else if raw.GetString_() != nil {
		return raw.GetString_().GetString_(), nil
	} else {
		return nil, &ClientError{
			message: "Failed to convert ProtoValue. This is likely a bug.",
		}
	}
}

// makeProtoValue converts go value to ProtoValue
func makeProtoValue(value any) (*prism.ProtoValue, error) {
	var result prism.ProtoValue
	switch value := value.(type) {
	case bool:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Boolean{
				Boolean: &prism.ProtoBoolean{
					Boolean: value,
				},
			},
		}
	case int32:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Integer{
				Integer: &prism.ProtoInteger{
					Integer: value,
				},
			},
		}
	case int64:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Long{
				Long: &prism.ProtoLong{
					Long: value,
				},
			},
		}
	case float64:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Double{
				Double: &prism.ProtoDouble{
					Double: value,
				},
			},
		}
	case float32:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_Float{
				Float: &prism.ProtoFloat{
					Float: value,
				},
			},
		}
	case string:
		result = prism.ProtoValue{
			Value: &prism.ProtoValue_String_{
				String_: &prism.ProtoString{
					String_: value,
				},
			},
		}
	default:
		return nil, &ClientError{
			message: fmt.Sprintf("Lack of support to %T %v", value, value),
		}
	}
	return &result, nil
}

// helperConvertValueToProto is a helper function which converts an array of driver.Value (any) to an array of ProtoValue
func helperConvertValueToProto(args []driver.Value) ([]*prism.ProtoValue, error) {
	pvs := make([]*prism.ProtoValue, len(args))
	for i, v := range args {
		pv, err := makeProtoValue(v)
		if err != nil {
			return nil, err
		}
		pvs[i] = pv
	}
	return pvs, nil
}

// helperConvertNamedvalueToProto converts an array of NamedValue values to an array of ProtoValue
func helperConvertNamedvalueToProto(args []driver.NamedValue) ([]*prism.ProtoValue, error) {
	pvs := make([]*prism.ProtoValue, len(args))
	for _, v := range args {
		pv, err := makeProtoValue(v.Value)
		if err != nil {
			return nil, err
		}
		pvs[v.Ordinal] = pv
	}
	return pvs, nil
}

func canConvertDocumentToRelational(documents []*prism.ProtoDocument) (bool, []string) {
	keys := make([]string, len(documents[0].GetEntries()))
	isFirst := true
	for _, document := range documents {
		if isFirst {
			isFirst = false
			for i, kvpair := range document.GetEntries() {
				key, _ := convertProtoValue(kvpair.GetKey())
				switch key := key.(type) {
				case string:
					keys[i] = key
				default:
					return false, nil
				}
			}
		} else {
			if len(document.GetEntries()) != len(keys) {
				return false, nil
			}
			for i, kvpair := range document.GetEntries() {
				key, _ := convertProtoValue(kvpair.GetKey())
				switch key := key.(type) {
				case string:
					if key != keys[i] {
						return false, nil
					}
				default:
					return false, nil
				}
			}
		}
	}
	return true, keys
}

// helperExtractResultFromStatementResult collects affectedRows and lastInsertId from StatementResult
// TODO: does polypheny currently have lastInsertId?
func helperExtractResultFromStatementResult(result *prism.StatementResult) (driver.Result, error) {
	rowsAffected := result.GetScalar()
	return &PolyphenyResult{
		lastInsertId: 0,
		rowsAffected: rowsAffected,
	}, nil
}

// helperExtractRowsFromResponse collects data from StatementResult and performs some transformation then return
func helperExtractRowsFromStatementResult(result *prism.StatementResult) (driver.Rows, error) {
	if result.GetFrame() == nil {
		return nil, &ClientError{
			message: "Query expects to return data, however the result is empty",
		}
	}
	frame := result.GetFrame()
	var values [][]any
	var columns []string
	if frame.GetRelationalFrame() != nil {
		relationalData := frame.GetRelationalFrame()
		columnResponse := relationalData.GetColumnMeta()
		columns = make([]string, len(columnResponse))
		for i, v := range columnResponse {
			columns[i] = v.GetColumnName()
		}
		rows := relationalData.GetRows()
		var currentRow []any
		for _, irow := range rows {
			currentRow = []any{}
			for _, ivalue := range irow.GetValues() {
				converted, err := convertProtoValue(ivalue)
				if err != nil {
					return nil, err
				}
				currentRow = append(currentRow, converted)
			}
			values = append(values, currentRow)
		}
		return &PolyphenyRows{
			columns:   columns,
			result:    values,
			readIndex: 0,
		}, nil
	} else if frame.GetDocumentFrame() != nil {
		documentData := frame.GetDocumentFrame().GetDocuments()
		canConvert, columns := canConvertDocumentToRelational(documentData)
		if canConvert {
			// if the documents have exactly the same schema
			// we will transform the query result into relational
			// the key will be the column name
			var currentRow []any
			for _, document := range documentData {
				currentRow = make([]any, len(columns))
				for _, kvpair := range document.GetEntries() {
					key, _ := convertProtoValue(kvpair.GetKey())
					for ki, k := range columns {
						if key == k {
							converted, err := convertProtoValue(kvpair.GetValue())
							if err != nil {
								return nil, err
							}
							currentRow[ki] = converted
						}
					}
				}
				values = append(values, currentRow)
			}
			return &PolyphenyRows{
				columns:   columns,
				result:    values,
				readIndex: 0,
			}, nil
		}
		// if we can't transform
		// TODO: need a better way to return the result
		for _, entries := range documentData {
			for _, v := range entries.GetEntries() {
				key, err := convertProtoValue(v.GetKey())
				if err != nil {
					return nil, err
				}
				value, err := convertProtoValue(v.GetValue())
				if err != nil {
					return nil, err
				}
				//currentDocument = append(currentDocument, kv)
				currentRow := make([]any, 2)
				currentRow[0] = key
				currentRow[1] = value
				values = append(values, currentRow)
			}
		}
		columns = make([]string, 2)
		columns[0] = "key"
		columns[1] = "value"
		return &PolyphenyRows{
			columns:   columns,
			result:    values,
			readIndex: 0,
		}, nil
	} else {
		// graph is currently not supported
		return nil, &ClientError{
			message: "Graph queries are currently not supported by Prism interface",
		}
	}
}

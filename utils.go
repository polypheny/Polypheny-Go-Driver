package polypheny

import (
	driver "database/sql/driver"
	fmt "fmt"
	"sync/atomic"

	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
)

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
	case int:
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
		pvs[v.Ordinal-1] = pv
	}
	return pvs, nil
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
		result := &PolyphenyRows{
			columns:   columns,
			result:    values,
			readIndex: atomic.Int32{},
		}
		result.readIndex.Store(0)
		return result, nil
	} else {
		return nil, &ClientError{
			message: "Query expects to return data, however the result is empty",
		}
	}
}

package polypheny

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math/big"
	"time"

	"github.com/polypheny/Polypheny-Go-Driver/prism"
)

// convertProtoValue converts a Prism value into a Go value.
func convertProtoValue(raw *prism.ProtoValue) (any, error) {
	switch value := raw.GetValue().(type) {
	case *prism.ProtoValue_Null:
		return nil, nil
	case *prism.ProtoValue_Boolean:
		return value.Boolean.GetBoolean(), nil
	case *prism.ProtoValue_Integer:
		return value.Integer.GetInteger(), nil
	case *prism.ProtoValue_Long:
		return value.Long.GetLong(), nil
	case *prism.ProtoValue_Float:
		return value.Float.GetFloat(), nil
	case *prism.ProtoValue_Double:
		return value.Double.GetDouble(), nil
	case *prism.ProtoValue_String_:
		return value.String_.GetString_(), nil
	case *prism.ProtoValue_Binary:
		return value.Binary.GetBinary(), nil
	case *prism.ProtoValue_File:
		return value.File.GetBinary(), nil
	case *prism.ProtoValue_Date:
		return time.Unix(value.Date.GetDate()*86400, 0).UTC(), nil
	case *prism.ProtoValue_Time:
		return time.UnixMilli(int64(value.Time.GetTime())).UTC(), nil
	case *prism.ProtoValue_Timestamp:
		return time.UnixMilli(value.Timestamp.GetTimestamp()).UTC(), nil
	case *prism.ProtoValue_BigDecimal:
		return decimalString(value.BigDecimal), nil
	case *prism.ProtoValue_List:
		result := make([]any, 0, len(value.List.GetValues()))
		for _, entry := range value.List.GetValues() {
			item, err := convertProtoValue(entry)
			if err != nil {
				return nil, err
			}
			result = append(result, item)
		}
		return result, nil
	case *prism.ProtoValue_Document:
		return convertDocument(value.Document)
	default:
		return nil, &ClientError{message: fmt.Sprintf("unsupported Prism value %T", value)}
	}
}

func decimalString(value *prism.ProtoBigDecimal) string {
	b := value.GetUnscaledValue()
	unscaled := new(big.Int).SetBytes(b)
	if len(b) > 0 && b[0]&0x80 != 0 {
		unscaled.Sub(unscaled, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
	}
	if scale := value.GetScale(); scale != 0 {
		return fmt.Sprintf("%se%d", unscaled, -int64(scale))
	}
	return unscaled.String()
}

func convertDocument(document *prism.ProtoDocument) (map[string]any, error) {
	fields := document.GetEntries()
	result := make(map[string]any, len(fields))
	for key, raw := range fields {
		value, err := convertProtoValue(raw)
		if err != nil {
			return nil, fmt.Errorf("document field %q: %w", key, err)
		}
		result[key] = value
	}
	return result, nil
}

func makeProtoValue(value any) (*prism.ProtoValue, error) {
	switch value := value.(type) {
	case nil:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Null{Null: &prism.ProtoNull{}}}, nil
	case bool:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Boolean{Boolean: &prism.ProtoBoolean{Boolean: value}}}, nil
	case int:
		return makeProtoValue(int64(value))
	case int32:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Integer{Integer: &prism.ProtoInteger{Integer: value}}}, nil
	case int64:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Long{Long: &prism.ProtoLong{Long: value}}}, nil
	case float32:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Float{Float: &prism.ProtoFloat{Float: value}}}, nil
	case float64:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Double{Double: &prism.ProtoDouble{Double: value}}}, nil
	case string:
		return &prism.ProtoValue{Value: &prism.ProtoValue_String_{String_: &prism.ProtoString{String_: value}}}, nil
	case []byte:
		if value == nil {
			return makeProtoValue(nil)
		}
		return &prism.ProtoValue{Value: &prism.ProtoValue_Binary{Binary: &prism.ProtoBinary{Binary: value}}}, nil
	case time.Time:
		return &prism.ProtoValue{Value: &prism.ProtoValue_Timestamp{Timestamp: &prism.ProtoTimestamp{Timestamp: value.UnixMilli()}}}, nil
	default:
		return nil, &ClientError{message: fmt.Sprintf("unsupported SQL argument type %T", value)}
	}
}

func namedValues(args []driver.Value) []driver.NamedValue {
	values := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		values[i] = driver.NamedValue{Ordinal: i + 1, Value: arg}
	}
	return values
}

func helperConvertNamedvalueToProto(args []driver.NamedValue) ([]*prism.ProtoValue, error) {
	values := make([]*prism.ProtoValue, len(args))
	for _, argument := range args {
		if argument.Name != "" {
			return nil, &ClientError{message: "SQL named arguments are not supported; use positional ? placeholders"}
		}
		if argument.Ordinal < 1 || argument.Ordinal > len(values) || values[argument.Ordinal-1] != nil {
			return nil, &ClientError{message: "invalid SQL argument ordinal"}
		}
		value, err := makeProtoValue(argument.Value)
		if err != nil {
			return nil, err
		}
		values[argument.Ordinal-1] = value
	}
	return values, nil
}

func helperExtractResultFromStatementResult(result *prism.StatementResult) (driver.Result, error) {
	if result == nil || result.GetFrame() != nil {
		return nil, &ClientError{message: "Prism server did not return an update count"}
	}
	return &PolyphenyResult{lastInsertId: -1, rowsAffected: result.GetScalar()}, nil
}

// readFrames fetches one page at a time and passes it to the result decoder.
func (conn *PolyphenyConn) readFrames(ctx context.Context, statementID int32, frame *prism.Frame, consume func(*prism.Frame) error) error {
	for frame != nil {
		if err := consume(frame); err != nil {
			return err
		}
		if frame.GetIsLast() {
			return nil
		}
		var err error
		frame, err = conn.fetch(ctx, statementID)
		if err != nil {
			return err
		}
	}
	return &ClientError{message: "Prism server did not return query rows"}
}

func (conn *PolyphenyConn) readRows(ctx context.Context, statementID int32, frame *prism.Frame) (driver.Rows, error) {
	rows := new(PolyphenyRows)
	err := conn.readFrames(ctx, statementID, frame, func(frame *prism.Frame) error {
		relational := frame.GetRelationalFrame()
		if relational == nil {
			return &ClientError{message: "Prism server did not return a relational result"}
		}
		if rows.columns == nil {
			columns := relational.GetColumnMeta()
			rows.columns = make([]string, len(columns))
			for i, column := range columns {
				rows.columns[i] = column.GetColumnLabel()
				if rows.columns[i] == "" {
					rows.columns[i] = column.GetColumnName()
				}
			}
		}
		for _, row := range relational.GetRows() {
			values := row.GetValues()
			if len(values) != len(rows.columns) {
				return &ClientError{message: "Prism row length does not match the column count"}
			}
			result := make([]any, len(values))
			for i, raw := range values {
				value, err := convertProtoValue(raw)
				if err != nil {
					return err
				}
				// database/sql requires its standard driver.Value types.
				switch v := value.(type) {
				case int32:
					value = int64(v)
				case float32:
					value = float64(v)
				}
				if !driver.IsValue(value) {
					return &ClientError{message: fmt.Sprintf("unsupported SQL result type %T", value)}
				}
				result[i] = value
			}
			rows.result = append(rows.result, result)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rows, nil
}

package polypheny

import (
	prism "github.com/polypheny/Polypheny-Go-Driver/protos"
	"testing"
)

func TestParseQuery(t *testing.T) {
	_, _, err := parseQuery("SELECT 42")
	if err.Error() != "A query should have the following format: QueryLanguage:Query" {
		t.Error("Error parsing query")
	}
	_, _, err = parseQuery("sql:SELECT 42")
	if err != nil {
		t.Error(err)
	}
}

func TestMakeProtoValue1(t *testing.T) {
	var result *prism.ProtoValue
	var value interface{}
	value = true
	result, _ = makeProtoValue(value)
	if result.GetBoolean().GetBoolean() != true {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetBoolean().GetBoolean())
	}
	value = int32(1)
	result, _ = makeProtoValue(value)
	if result.GetInteger().GetInteger() != value.(int32) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetInteger().GetInteger())
	}
	value = int64(100000000000)
	result, _ = makeProtoValue(value)
	if result.GetLong().GetLong() != value.(int64) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetLong().GetLong())
	}
	value = "Hello, world!"
	result, _ = makeProtoValue(value)
	if result.GetString_().GetString_() != value.(string) {
		t.Fatalf("Error in making a ProtoValue, expected %v, got %v", value, result.GetString_().GetString_())
	}
}

func TestConvertProtoValue(t *testing.T) {
	var protoValue *prism.ProtoValue
	var result interface{}
	var expected interface{} = true
	protoValue, _ = makeProtoValue(expected)
	result, _ = convertProtoValue(protoValue)
	if result.(bool) != expected {
		t.Fatalf("Failed to convert, expected %v, but got %v", expected, result)
	}
}

package webapi

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFormSerializerMarshalNilValue(t *testing.T) {
	got, err := (&formSerializer{}).Marshal(map[string]interface{}{
		"name":  "alice",
		"extra": nil,
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	values, err := url.ParseQuery(string(got))
	if err != nil {
		t.Fatalf("ParseQuery() error = %v", err)
	}
	if values.Get("name") != "alice" || values.Get("extra") != "" {
		t.Errorf("Marshal() = %q; want name=alice and empty extra value", got)
	}
}

func TestFormSerializerUnmarshalRejectsInvalidDestinations(t *testing.T) {
	tests := []struct {
		name string
		dest interface{}
	}{
		{"nil", nil},
		{"non-pointer", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := (&formSerializer{}).Unmarshal([]byte("value=1"), test.dest); err == nil {
				t.Fatal("Unmarshal() error = nil; want invalid destination error")
			}
		})
	}
}

func TestLoadFromValuesRejectsInvalidNumbersAndBooleans(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		data string
	}{
		{"integer", reflect.TypeOf(struct{ Value int }{}), "not-an-integer"},
		{"unsigned overflow", reflect.TypeOf(struct{ Value uint8 }{}), "256"},
		{"float", reflect.TypeOf(struct{ Value float64 }{}), "not-a-float"},
		{"boolean", reflect.TypeOf(struct{ Value bool }{}), "not-a-boolean"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (&param{Type: test.typ}).loadFromValues(url.Values{"Value": {test.data}})
			if err == nil {
				t.Fatalf("loadFromValues(%q) error = nil; want conversion error", test.data)
			}
		})
	}
}

func TestLoadFromValuesParsesValidNumber(t *testing.T) {
	type query struct {
		Value int
	}

	value, err := (&param{Type: reflect.TypeOf(query{})}).loadFromValues(url.Values{"Value": {"42"}})
	if err != nil {
		t.Fatalf("loadFromValues() error = %v", err)
	}
	if got := (*value).Interface().(query).Value; got != 42 {
		t.Errorf("loadFromValues() Value = %d; want 42", got)
	}
}

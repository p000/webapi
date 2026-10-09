package webapi

import (
	"net/url"
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
